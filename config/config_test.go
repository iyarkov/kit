package config

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

type level3Config struct {
	StrVal   string
	BoolVal  bool
	IntVal   int64
	FloatVal float64
	PwdVal   Password
}

type level2Config struct {
	StrVal   string
	BoolVal  bool
	IntVal   int64
	FloatVal float64
	PwdVal   Password

	Bottom level3Config
}

type topConfig struct {
	StrVal   string
	BoolVal  bool
	IntVal   int64
	FloatVal float64
	PwdVal   Password

	Left  level2Config
	Right level2Config
}

func assertEqualL3(t *testing.T, expected, actual *level3Config, name string) {
	if expected.StrVal != actual.StrVal {
		t.Errorf("%s.StrVal [%s]!=[%s]", name, expected.StrVal, actual.StrVal)
	}
	if expected.BoolVal != actual.BoolVal {
		t.Errorf("%s.BoolVal [%t]!=[%t]", name, expected.BoolVal, actual.BoolVal)
	}
	if expected.IntVal != actual.IntVal {
		t.Errorf("%s.IntVal [%d]!=[%d]", name, expected.IntVal, actual.IntVal)
	}
	if expected.FloatVal != actual.FloatVal {
		t.Errorf("%s.FloatVal [%f]!=[%f]", name, expected.FloatVal, actual.FloatVal)
	}
	if expected.PwdVal != actual.PwdVal {
		t.Errorf("PwdVal [%s]!=[%s]", expected.PwdVal, actual.PwdVal)
	}
}

func assertEqualL2(t *testing.T, expected, actual *level2Config, name string) {
	if expected.StrVal != actual.StrVal {
		t.Errorf("%s.StrVal [%s]!=[%s]", name, expected.StrVal, actual.StrVal)
	}
	if expected.BoolVal != actual.BoolVal {
		t.Errorf("%s.BoolVal [%t]!=[%t]", name, expected.BoolVal, actual.BoolVal)
	}
	if expected.IntVal != actual.IntVal {
		t.Errorf("%s.IntVal [%d]!=[%d]", name, expected.IntVal, actual.IntVal)
	}
	if expected.FloatVal != actual.FloatVal {
		t.Errorf("%s.FloatVal [%f]!=[%f]", name, expected.FloatVal, actual.FloatVal)
	}
	if expected.PwdVal != actual.PwdVal {
		t.Errorf("PwdVal [%s]!=[%s]", expected.PwdVal, actual.PwdVal)
	}
	assertEqualL3(t, &expected.Bottom, &actual.Bottom, fmt.Sprintf("%s.bottom", name))
}

func assertEqual(t *testing.T, expected, actual *topConfig) {
	if expected.StrVal != actual.StrVal {
		t.Errorf("StrVal [%s]!=[%s]", expected.StrVal, actual.StrVal)
	}
	if expected.BoolVal != actual.BoolVal {
		t.Errorf("BoolVal [%t]!=[%t]", expected.BoolVal, actual.BoolVal)
	}
	if expected.IntVal != actual.IntVal {
		t.Errorf("IntVal [%d]!=[%d]", expected.IntVal, actual.IntVal)
	}
	if expected.FloatVal != actual.FloatVal {
		t.Errorf("FloatVal [%f]!=[%f]", expected.FloatVal, actual.FloatVal)
	}
	if expected.PwdVal != actual.PwdVal {
		t.Errorf("PwdVal [%s]!=[%s]", string(expected.PwdVal), string(actual.PwdVal))
	}
	assertEqualL2(t, &expected.Left, &actual.Left, "left")
	assertEqualL2(t, &expected.Right, &actual.Right, "right")
}

func TestUpdateVal(t *testing.T) {
	type spec struct {
		name        string
		expected    topConfig
		params      map[string]string
		expectedErr error
	}

	suite := []spec{
		{
			name: "String",
			expected: topConfig{
				StrVal: "fooo",
			},
			params: map[string]string{
				"StrVal": "fooo",
			},
		},
		{
			name: "Bool",
			expected: topConfig{
				BoolVal: true,
			},
			params: map[string]string{
				"BoolVal": "true",
			},
		},
		{
			name:        "BoolInvalid",
			expectedErr: strconv.ErrSyntax,
			params: map[string]string{
				"BoolVal": "sf",
			},
		},
		{
			name: "IntVal",
			expected: topConfig{
				IntVal: 123,
			},
			params: map[string]string{
				"IntVal": "123",
			},
		},
		{
			name:        "IntValInvalid",
			expectedErr: strconv.ErrSyntax,
			params: map[string]string{
				"IntVal": "12sf",
			},
		},
		{
			name: "FloatVal",
			expected: topConfig{
				FloatVal: 123.34,
			},
			params: map[string]string{
				"FloatVal": "123.34",
			},
		},
		{
			name:        "FloatValInvalid",
			expectedErr: strconv.ErrSyntax,
			params: map[string]string{
				"FloatVal": "12.sf",
			},
		},
		{
			name: "PwdVal",
			expected: topConfig{
				PwdVal: "topSecret",
			},
			params: map[string]string{
				"PwdVal": "topSecret",
			},
		},
		{
			name: "nestedLvl2",
			expected: topConfig{
				Left: level2Config{
					StrVal: "Left",
				},
				Right: level2Config{
					BoolVal: true,
				},
			},
			params: map[string]string{
				"Left.StrVal":   "Left",
				"Right.BoolVal": "true",
			},
		},
		{
			name: "nestedLvl3",
			expected: topConfig{
				Left: level2Config{
					Bottom: level3Config{
						BoolVal: true,
					},
				},
			},
			params: map[string]string{
				"Left.Bottom.BoolVal": "true",
			},
		},
		{
			name: "nestedLvl3Error",
			params: map[string]string{
				"Left.Bottom.BoolVal": "xyz",
			},
			expectedErr: strconv.ErrSyntax,
		},
	}

	for _, test := range suite {
		t.Run(test.name, func(t *testing.T) {
			value := topConfig{}
			if err := updateConfig(&value, &test.params); err != nil {
				if test.expectedErr == nil || !errors.Is(err, test.expectedErr) {
					t.Errorf("Unexpected error, exepcting %v, got %v", test.expectedErr, err)
				}
			} else {
				assertEqual(t, &test.expected, &value)
			}
		})
	}
}

func TestNormalizeToUpper(t *testing.T) {
	type spec struct {
		input    string
		expected string
	}

	suite := []spec{
		{
			input:    "abc",
			expected: "Abc",
		},
		{
			input:    "ABC",
			expected: "Abc",
		},
		{
			input:    "прИвет",
			expected: "Привет",
		},
		{
			input:    "SNAKE_FOO_BAR",
			expected: "Snake.Foo.Bar",
		},
		{
			input:    "...",
			expected: "...",
		},
		{
			input:    "___a",
			expected: "...A",
		},
	}

	for i, test := range suite {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			output := normalizeEnvKey(test.input)
			if output != test.expected {
				t.Errorf("expected [%s], actual [%s]", test.expected, output)
			}
		})
	}
}

func TestGetFlag(t *testing.T) {
	type spec struct {
		flag     string
		input    []string
		expected string
	}
	suite := []spec{
		{
			flag:     "-f",
			input:    []string{},
			expected: "",
		},
		{
			flag:     "-f",
			input:    []string{"aaa"},
			expected: "",
		},
		{
			flag:     "-f",
			input:    []string{"aaa", "bbb"},
			expected: "",
		},
		{
			flag:     "-f",
			input:    []string{"aaa", "-f"},
			expected: "",
		},
		{
			flag:     "-f",
			input:    []string{"-fabc"},
			expected: "abc",
		},
		{
			flag:     "-f",
			input:    []string{"-fabc", "-fcde"},
			expected: "abc",
		},
	}
	for i, test := range suite {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			output := getFlag(test.flag, test.input)
			if output != test.expected {
				t.Errorf("expected [%s], actual [%s]", test.expected, output)
			}
		})
	}

}
