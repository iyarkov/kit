package schema

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"microserver.rockyrunstream.com/foundation/support"
	"time"
)

type Change struct {
	Version  string
	Commands []string
	Function func(ctx context.Context, tx *sql.Tx) error
	Timeout  time.Duration
}

type HistoryRecord struct {
	Id        int32
	CreatedAt time.Time
	Version   string
}

func Update(ctx context.Context, db *sql.DB, changeset []Change) (string, string, error) {
	if !assertChangeset(ctx, changeset) {
		return "", "", ErrInvalidChangeset
	}

	if err := ensureSchemaTable(ctx, db); err != nil {
		return "", "", err
	}

	history, err := LoadHistory(ctx, db)
	if err != nil {
		return "", "", err
	}
	dbVersion := ""
	versionMap := make(map[string]bool, 0)
	if len(history) > 0 {
		dbVersion = history[0].Version
		for _, r := range history {
			versionMap[r.Version] = true
		}
	}
	log := zerolog.Ctx(ctx)
	log.Debug().Msgf("Current DB version %s", dbVersion)

	expectedVersion := ""
	changed := false
	for _, change := range changeset {
		expectedVersion = change.Version
		if versionMap[change.Version] {
			log.Debug().Msgf("Skip changeset %s", change.Version)
			continue
		}
		if err = applyChange(ctx, db, change); err != nil {
			return "", "", err
		}
		changed = true
	}
	if changed {
		log.Info().Msgf("Database upgraded from %s to %s", dbVersion, expectedVersion)
	} else {
		log.Info().Msgf("Database upgraded no required, DB version: %s, expected version: %s ", dbVersion, expectedVersion)
	}
	return dbVersion, expectedVersion, nil
}

func LoadDbVersion(ctx context.Context, db *sql.DB) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	row := db.QueryRowContext(ctx, queryLastVersion)
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

func LoadHistory(ctx context.Context, db *sql.DB) ([]HistoryRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	rows, err := db.QueryContext(ctx, queryLoadHistory)
	defer support.CloseWithWarning(ctx, rows, "failed to close rows")
	if err != nil {
		return nil, fmt.Errorf("LoadHistory query failed: %w", err)
	}

	result := make([]HistoryRecord, 0)
	record := HistoryRecord{}
	for rows.Next() {
		err = rows.Scan(&record.Id, &record.CreatedAt, &record.Version)
		if err != nil {
			return nil, fmt.Errorf("LoadHistory scan error: %w", err)
		}
		result = append(result, record)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("after scan error: %w", rows.Err())
	}
	return result, nil
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

func ensureSchemaTable(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Check the table exists
	row := db.QueryRowContext(ctx, queryTableExists, defaultSchema)
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

	_, err := db.ExecContext(ctx, queryCreateTable)
	if err != nil {
		return fmt.Errorf("ensure table, create table query failed: %w", err)
	}

	log := zerolog.Ctx(ctx)
	log.Info().Msg("schema table has been created")
	return nil
}

func applyChange(ctx context.Context, db *sql.DB, change Change) error {
	// get log
	log := zerolog.Ctx(ctx)
	log.Info().Msgf("Upgrading DB schema to %s", change.Version)

	// Define timeout
	var timeout time.Duration
	if change.Timeout == 0 {
		timeout = defaultTimeout
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
		_, err = tx.ExecContext(ctx, queryInsertVersion, now, change.Version)
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
