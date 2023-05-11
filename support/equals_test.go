package support

import "testing"

func TestEqualsStr(t *testing.T) {
	type equalsStrTestCase struct {
		name   string
		a      []string
		b      []string
		result bool
	}

	suite := []equalsStrTestCase{
		{
			name:   "all nil",
			a:      nil,
			b:      nil,
			result: true,
		},
		{
			name:   "a not nil",
			a:      []string{},
			b:      nil,
			result: false,
		},
		{
			name:   "b not nil",
			a:      nil,
			b:      []string{},
			result: false,
		},
		{
			name:   "size does not match",
			a:      []string{"aaa", "bbb"},
			b:      []string{"bb"},
			result: false,
		},
		{
			name:   "not equal",
			a:      []string{"aaa", "bbb"},
			b:      []string{"bbb", "aaa"},
			result: false,
		},
		{
			name:   "equal",
			a:      []string{"aaa", "bbb"},
			b:      []string{"aaa", "bbb"},
			result: true,
		},
	}

	for _, test := range suite {
		t.Run(test.name, func(t *testing.T) {
			actual := EqualsStr(test.a, test.b)
			if actual != test.result {
				t.Errorf("Test %s failes, expecting %t actual %t", test.name, test.result, actual)
			}
		})
	}

}
