package sql

import (
	"github.com/iyarkov/kit/support"
	"reflect"
	"strings"
	"testing"
)

func TestValidateSequenceEmpty(t *testing.T) {
	expected := Schema{
		sequencesMap: map[string]bool{},
	}
	db := Schema{
		sequencesMap: map[string]bool{},
	}
	errors := validateSequences(expected, db, true)
	if len(errors) != 0 {
		t.Errorf("No errors expected, got [%v]", errors)
	}
}

func TestValidateSequenceSame(t *testing.T) {
	expected := Schema{
		sequencesMap: map[string]bool{
			"seq1": true,
			"seq2": true,
		},
	}
	db := Schema{
		sequencesMap: map[string]bool{
			"seq1": true,
			"seq2": true,
		},
	}
	errors := validateSequences(expected, db, true)
	if len(errors) != 0 {
		t.Errorf("No errors expected, got %v", errors)
	}
}

func TestValidateSequenceMissedSequence(t *testing.T) {
	expected := Schema{
		sequencesMap: map[string]bool{
			"seq1": true,
			"seq2": true,
		},
	}
	db := Schema{
		sequencesMap: map[string]bool{
			"seq1": true,
		},
	}
	errors := validateSequences(expected, db, true)
	expectedErrors := []string{"sequence seq2 is missing"}
	if !reflect.DeepEqual(errors, expectedErrors) {
		t.Errorf("Expecting %s, got %s", expectedErrors, errors)
	}
}

func TestValidateSequenceExtraSequenceStrict(t *testing.T) {
	expected := Schema{
		sequencesMap: map[string]bool{
			"seq1": true,
			"seq2": true,
		},
	}
	db := Schema{
		sequencesMap: map[string]bool{
			"seq1": true,
			"seq2": true,
			"seq3": true,
		},
	}
	errors := validateSequences(expected, db, true)
	expectedErrors := []string{"Unexpected sequence: seq3"}
	if !reflect.DeepEqual(errors, expectedErrors) {
		t.Errorf("Expecting %s, got %s", expectedErrors, errors)
	}
}

func TestValidateSequenceExtraSequenceNonStrict(t *testing.T) {
	expected := Schema{
		sequencesMap: map[string]bool{
			"seq1": true,
			"seq2": true,
		},
	}
	db := Schema{
		sequencesMap: map[string]bool{
			"seq1": true,
			"seq2": true,
			"seq3": true,
		},
	}
	errors := validateSequences(expected, db, false)
	if len(errors) != 0 {
		t.Errorf("No errors expected, got %v", errors)
	}
}

func TestSetNamesEmptySchema(t *testing.T) {
	// Test it does not blow up
	normalize(&Schema{})
}

type testCaseSpec struct {
	name     string
	expected Schema
	actual   Schema
	strict   bool
	errors   []string
}

func TestTablesValidation(t *testing.T) {
	suite := []testCaseSpec{
		{
			name: "Missed table",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{},
			},
			errors: []string{"table Table_A is missing"},
		},
		{
			name: "Extra Table, non strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			errors: []string{},
		},
		{
			name: "Extra Table, strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			errors: []string{"Unexpected table: Table_B"},
			strict: true,
		},
		{
			name: "Tables Match",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			errors: []string{},
			strict: true,
		},
		{
			name: "Missed column",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Columns: map[string]Column{
							"Column_A": {
								Type:       "varchar",
								CharLength: 255,
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			errors: []string{
				"column Table_A.Column_A is missing",
			},
		},
		{
			name: "Extra column",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {
								Type:       "varchar",
								CharLength: 255,
							},
						},
					},
				},
			},
			errors: []string{},
		},
		{
			name: "Extra column - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {
								Type:       "varchar",
								CharLength: 255,
							},
						},
					},
				},
			},
			strict: true,
			errors: []string{
				"Unexpected column: Table_B.Column_A",
			},
		},
		{
			name: "columnsMap - matches",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {},
						},
					},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {},
						},
					},
				},
			},
			errors: []string{},
		},
		{
			name: "columnsMap - type matches",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {
								Type:         "varchar",
								CharLength:   34,
								NumPrecision: 21,
								NotNull:      false,
								IsUnique:     false,
							},
						},
					},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {
								Type:         "varchar",
								CharLength:   1,
								NumPrecision: 2,
								NotNull:      true,
								IsUnique:     true,
							},
						},
					},
				},
			},
			errors: []string{},
		},
		{
			name: "columnsMap - type does not match",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {
								Type:         "varchar",
								CharLength:   34,
								NumPrecision: 21,
								NotNull:      false,
								IsUnique:     false,
							},
						},
					},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {
								Type:         "int",
								CharLength:   1,
								NumPrecision: 2,
								NotNull:      true,
								IsUnique:     true,
							},
						},
					},
				},
			},
			errors: []string{
				"invalid column type: Table_B.Column_A, expected varchar, actual int",
			},
		},
		{
			name: "columnsMap - extra properties do not match, strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {
								Type:         "varchar",
								CharLength:   34,
								NumPrecision: 21,
								NotNull:      false,
								IsUnique:     false,
							},
						},
					},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Columns: map[string]Column{
							"Column_A": {
								Type:         "varchar",
								CharLength:   1,
								NumPrecision: 2,
								NotNull:      true,
								IsUnique:     true,
							},
						},
					},
				},
			},
			strict: true,
			errors: []string{
				"invalid column char length: Table_B.Column_A, expected 34, actual 1",
				"invalid column num precision: Table_B.Column_A, expected 21, actual 2",
				"invalid column is nullable: Table_B.Column_A, expected false, actual true",
				"invalid column is unique: Table_B.Column_A, expected false, actual true",
			},
		},
		{
			name: "Missed index",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								columnsMap: map[string]bool{
									"Column_A": true,
									"Column_B": true,
								},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			errors: []string{},
		},
		{
			name: "Missed index - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								columnsMap: map[string]bool{
									"Column_A": true,
									"Column_B": true,
								},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"index Table_A.Index_A is missing",
			},
		},
		{
			name: "Extra index - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {
						Indexes: map[string]Index{
							"Index_A": {
								columnsMap: map[string]bool{
									"Column_A": true,
									"Column_B": true,
								},
								IsUnique: true,
							},
						},
					},
				},
			},
			strict: true,
			errors: []string{
				"Unexpected index: Table_B.Index_A",
			},
		},
		{
			name: "Index matches - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								columnsMap: map[string]bool{
									"Column_A": true,
									"Column_B": true,
								},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								columnsMap: map[string]bool{
									"Column_A": true,
									"Column_B": true,
								},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{},
		},
		{
			name: "Index missed column - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								Columns:  []string{"Column_A", "Column_B"},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								Columns:  []string{"Column_A"},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"invalid index  Table_A.Index_A, missing column: Column_B",
			},
		},
		{
			name: "Index extra column - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								Columns:  []string{"Column_A"},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								Columns:  []string{"Column_A", "Column_B"},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"invalid index  Table_A.Index_A, extra column: Column_B",
			},
		},
		{
			name: "Index Unique does not match - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								columnsMap: map[string]bool{
									"Column_A": true,
									"Column_B": true,
								},
								IsUnique: true,
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						Indexes: map[string]Index{
							"Index_A": {
								columnsMap: map[string]bool{
									"Column_A": true,
									"Column_B": true,
								},
								IsUnique: false,
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"invalid index IsUnique: Table_A.Index_A, expected true, actual false",
			},
		},
		{
			name: "Missed FK",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			errors: []string{},
		},
		{
			name: "Missed FK - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"foreign keys Table_A.FK_A is missing",
			},
		},
		{
			name: "Missed FK - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"foreign keys Table_A.FK_A is missing",
			},
		},
		{
			name: "Extra FK - strict",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"Unexpected foreign keys: Table_A.FK_A",
			},
		},
		{
			name: "FK matches",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{},
		},
		{
			name: "FK Missed column",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"invalid fk: Table_A.FK_A, missed column: Column_B => Column_BB",
			},
		},
		{
			name: "FK extra column",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"invalid fk: Table_A.FK_A, extra column Column_A => Column_AA",
			},
		},
		{
			name: "FK columns mismatch",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"invalid fk: Table_A.FK_A, wrong column mapping, expected: Column_A => Column_AA, actual Column_A => Column_AA_AA",
			},
		},
		{
			name: "FK table mismatch",
			expected: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_BB",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			actual: Schema{
				Tables: map[string]Table{
					"Table_A": {
						ForeignKeys: map[string]ForeignKey{
							"FK_A": {
								ForeignTable: "Table_B",
								Columns: map[string]string{
									"Column_A": "Column_AA",
									"Column_B": "Column_BB",
								},
							},
						},
					},
					"Table_B": {},
				},
			},
			strict: true,
			errors: []string{
				"invalid fk foreign table: Table_A.FK_A, expected Table_BB, actual Table_B",
			},
		},
	}
	for _, testCase := range suite {
		t.Run(testCase.name, func(t *testing.T) {
			normalize(&testCase.expected)
			normalize(&testCase.actual)
			actualErrors := validateSchema(testCase.expected, testCase.actual, testCase.strict)
			if !support.EqualsStr(actualErrors, testCase.errors) {
				t.Errorf("result does not match, expecting: %v, actual: %v", strings.Join(testCase.errors, ","), strings.Join(actualErrors, ","))
			}
		})
	}
}
