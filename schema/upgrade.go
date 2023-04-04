package schema

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type Change struct {
	Version  string
	Commands []string
	Function func(ctx context.Context, tx *sql.Tx) error
	Timeout  time.Duration
}

var ErrInvalidChangeset = fmt.Errorf("InvalidChangeset")

const DefaultTimeout = time.Second * 30

func Update(ctx context.Context, db *sql.DB, changeset []Change) error {
	if !assertChangeset(ctx, changeset) {
		return ErrInvalidChangeset
	}

	if err := ensureSchemaTable(ctx, db); err != nil {
		return err
	}

	currentVersion, err := getVersion(ctx, db)
	if err != nil {
		return err
	}
	log := zerolog.Ctx(ctx)
	log.Debug().Msgf("Current DB version %s", currentVersion)
	startChangesetIdx := findStartIdx(currentVersion, changeset)
	if startChangesetIdx == -1 || startChangesetIdx == len(changeset) {
		log.Info().Msgf("DB schema version %s, upgrade is not required", currentVersion)
		return nil
	}
	for i := startChangesetIdx; i < len(changeset); i++ {
		change := changeset[i]
		if err := applyChange(ctx, db, change); err != nil {
			return err
		}
	}
	log.Info().Msgf("Database upgraded from %s to %s", currentVersion, changeset[len(changeset)-1].Version)
	return nil
}

func ensureSchemaTable(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Check the table exists
	row := db.QueryRowContext(ctx, `
SELECT EXISTS (
	SELECT FROM pg_tables WHERE
schemaname = 'public' AND
tablename  = 'schema_history'
)`)
	if row.Err() != nil {
		return fmt.Errorf("ensure table, table exist query failed: %w", row.Err())
	}
	var tableExist bool
	if err := row.Scan(&tableExist); err != nil {
		return fmt.Errorf("ensure table, table exist scan failed: %w", err)
	}

	if tableExist {
		return nil
	}

	_, err := db.ExecContext(ctx, `
CREATE TABLE schema_history (
	id SERIAL,
	created_at TIMESTAMP(3) WITHOUT TIME ZONE,
	version VARCHAR(255),
	PRIMARY KEY (id)
)`)
	if err != nil {
		return fmt.Errorf("ensure table, create table query failed: %w", err)
	}

	log := zerolog.Ctx(ctx)
	log.Info().Msg("schema table has been created")
	return nil
}

func getVersion(ctx context.Context, db *sql.DB) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	row := db.QueryRowContext(ctx, `SELECT version FROM schema_history WHERE id = (SELECT MAX(id) FROM schema_history)`)
	if row.Err() != nil {
		return "", fmt.Errorf("get version, version query failed: %w", row.Err())
	}

	var version string

	if err := row.Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("get version, version scan failed: %w", err)
	}
	return version, nil
}

func assertChangeset(ctx context.Context, changeset []Change) bool {
	seen := make(map[string]bool)
	valid := true
	log := zerolog.Ctx(ctx)
	for i, change := range changeset {
		if change.Version == "" {
			log.Error().Msgf("Line %d missing Version", i)
			valid = false
			continue
		}
		if seen[change.Version] {
			log.Error().Msgf("Line %d duplicated Version %s", i, change.Version)
			valid = false
		} else {
			seen[change.Version] = true
		}
		if (change.Commands == nil || len(change.Commands) == 0) && change.Function == nil {
			log.Error().Msgf("Line %d Version %s, either Command or Function required", i, change.Version)
			valid = false
		}
	}
	if len(changeset) == 0 {
		log.Error().Msg("Changeset is empty")
		valid = false
	}
	return valid
}

func findStartIdx(currentVersion string, changeset []Change) int {
	if currentVersion == "" {
		return 0
	}
	for i, change := range changeset {
		if change.Version == currentVersion {
			return i + 1
		}
	}
	return -1
}

func applyChange(ctx context.Context, db *sql.DB, change Change) error {
	// get log
	log := zerolog.Ctx(ctx)
	log.Info().Msgf("Upgrading DB schema to %s", change.Version)

	// Define timeout
	var timeout time.Duration
	if change.Timeout == 0 {
		timeout = DefaultTimeout
	} else {
		timeout = change.Timeout
	}

	// Apply timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Start transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("execute command, begin tx failed: %w", err)
	}
	log.Debug().Msg("Transaction begin")
	defer func() {
		if err == nil {
			log.Debug().Msg("Transaction succeeded")
			return
		}
		log.Debug().Msg("Transaction failed, rolling back")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zerolog.Ctx(ctx).Warn().Err(rollbackErr).Msg("tx rollback error")

		}
		log.Debug().Msg("Transaction rolled back")
	}()

	err = func() error {
		// Apply the change
		if change.Function != nil {
			err = change.Function(ctx, tx)
			if err != nil {
				return fmt.Errorf("execute command %s, exec failed: %w", change.Version, err)
			}
		} else {
			for _, command := range change.Commands {
				_, err = tx.ExecContext(ctx, command)
				if err != nil {
					return fmt.Errorf("execute command %s, exec failed: %w", change.Version, err)
				}
			}
		}

		// Update history
		now := time.Now().UTC()
		_, err = tx.ExecContext(ctx, "INSERT INTO schema_history(created_at, version) VALUES($1, $2)", now, change.Version)
		if err != nil {
			return fmt.Errorf("execute command, update history failed: %w", err)
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("execute command, tx commit failed: %w", err)
		}
		log.Info().Msgf("DB schema upgraded, new version %s", change.Version)
		return nil
	}()
	return err
}
