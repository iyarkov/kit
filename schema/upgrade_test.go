package schema

import (
	"context"
	"database/sql"
	"testing"
)

func TestAssertChangesetEmpty(t *testing.T) {
	if assertChangeset(context.Background(), []Change{}) {
		t.Error("Empty changeset must be invalid")
	}
}

func TestAssertChangesetNoCommand(t *testing.T) {
	if assertChangeset(context.Background(), []Change{
		{
			Version: "1.2.3",
		},
	}) {
		t.Error("Change must have either Command or a Function")
	}
}

func TestAssertChangesetCommand(t *testing.T) {
	if !assertChangeset(context.Background(), []Change{
		{
			Version:  "1.2.3",
			Commands: []string{"command"},
		},
	}) {
		t.Error("Change with a command must be valid")
	}
}

func TestAssertChangesetFunction(t *testing.T) {
	if !assertChangeset(context.Background(), []Change{
		{
			Version: "1.2.3",
			Function: func(ctx context.Context, tx *sql.Tx) error {
				return nil
			},
		},
	}) {
		t.Error("Change with a command must be valid")
	}
}

func TestAssertChangesetDuplicateId(t *testing.T) {
	if assertChangeset(context.Background(), []Change{
		{
			Version:  "1.2.3",
			Commands: []string{"command 1"},
		},
		{
			Version:  "1.2.3",
			Commands: []string{"command 2"},
		},
	}) {
		t.Error("Change with duplicated versions must be invalid")
	}
}

func TestAssertChangesetNoDuplicateId(t *testing.T) {
	if !assertChangeset(context.Background(), []Change{
		{
			Version:  "1.2.3",
			Commands: []string{"command 1"},
		},
		{
			Version:  "1.2.4",
			Commands: []string{"command 2"},
		},
	}) {
		t.Error("Change without duplicated versions must be valid")
	}
}
