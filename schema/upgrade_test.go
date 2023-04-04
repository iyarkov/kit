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

func TestFindStartIdx(t *testing.T) {
	type spec struct {
		name           string
		currentVersion string
		changeset      []Change
		expected       int
	}
	suite := []spec{
		{
			name:           "from scratch",
			currentVersion: "",
			changeset: []Change{
				{
					Version: "1.0.0",
				},
			},
			expected: 0,
		},
		{
			name:           "upgrade required",
			currentVersion: "1.0.2",
			changeset: []Change{
				{
					Version: "1.0.1",
				},
				{
					Version: "1.0.2",
				},
				{
					Version: "1.0.3",
				},
			},
			expected: 1,
		},
		{
			name:           "upgrade not required",
			currentVersion: "1.0.3",
			changeset: []Change{
				{
					Version: "1.0.1",
				},
				{
					Version: "1.0.2",
				},
				{
					Version: "1.0.3",
				},
			},
			expected: 2,
		},
		{
			name:           "app behind the schema",
			currentVersion: "1.0.4",
			changeset: []Change{
				{
					Version: "1.0.1",
				},
				{
					Version: "1.0.2",
				},
				{
					Version: "1.0.3",
				},
			},
			expected: -1,
		},
	}

	for _, test := range suite {
		t.Run(test.name, func(t *testing.T) {
			actual := findStartIdx(test.currentVersion, test.changeset)
			if actual != test.expected {
				t.Errorf("failed expected %d, actual %d", test.expected, actual)
			}
		})
	}
}
