package hlc

import (
	"testing"
)

type staticClock struct {
	value int32
}

func (s staticClock) now() int32 {
	return s.value
}

func TestSmoke(t *testing.T) {
	clock := New(1)

	t1 := clock.Next()
	t2 := clock.Next()
	t3 := clock.Next()

	if !t1.Before(t2) {
		t.Error("T1 must be smaller then T2")
	}
	if t1.After(t2) {
		t.Error("T1 must be smaller then T2")
	}
	if !t2.Before(t3) {
		t.Error("T2 must be smaller then T3")
	}
	if t2.After(t3) {
		t.Error("T2 must be smaller then T3")
	}
	if !t1.Before(t3) {
		t.Error("T1 must be smaller then T3")
	}
	if t1.After(t3) {
		t.Error("T1 must be smaller then T3")
	}
}

func TestNow(t *testing.T) {
	clock := hybridLogicalClock{
		ts:    10,
		count: 5,
		node:  3,
	}
	ts := clock.now()
	assertEquals(t, int32(10), ts.Time)
	assertEquals(t, int32(5), ts.Count)
	assertEquals(t, int32(3), ts.Node)
}

func assertEquals(t *testing.T, expected any, actual any) {
	if expected != actual {
		t.Errorf("Failed, expected %v actual %v", expected, actual)
	}
}

func TestNext(t *testing.T) {
	type Test struct {
		name     string
		clock    hybridLogicalClock
		expected Stamp
	}
	tests := []Test{
		{
			name: "Wall Ahead",
			clock: hybridLogicalClock{
				ts:    999,
				count: 55,
				node:  4,
				wallClock: staticClock{
					value: 1000,
				},
			},
			expected: Stamp{
				Time:  1000,
				Count: 0,
				Node:  4,
			},
		},
		{
			name: "Clock Ahead",
			clock: hybridLogicalClock{
				ts:    1000,
				count: 55,
				node:  4,
				wallClock: staticClock{
					value: 900,
				},
			},
			expected: Stamp{
				Time:  1000,
				Count: 56,
				Node:  4,
			},
		},
		{
			name: "Timestamps Equals",
			clock: hybridLogicalClock{
				ts:    1000,
				count: 55,
				node:  4,
				wallClock: staticClock{
					value: 1000,
				},
			},
			expected: Stamp{
				Time:  1000,
				Count: 56,
				Node:  4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.clock.Next()
			assertEquals(t, test.expected.Time, actual.Time)
			assertEquals(t, test.expected.Count, actual.Count)
			assertEquals(t, test.expected.Node, actual.Node)
		})
	}
}

func TestUpdate(t *testing.T) {
	type Test struct {
		name     string
		clock    hybridLogicalClock
		input    Stamp
		expected Stamp
	}

	tests := []Test{
		{
			name: "Wall Ahead",
			clock: hybridLogicalClock{
				ts:    1,
				count: 55,
				wallClock: staticClock{
					value: 20,
				},
			},
			input: Stamp{
				Time:  1,
				Count: 34,
			},
			expected: Stamp{
				Time:  20,
				Count: 0,
			},
		},
		{
			name: "Clock Ahead",
			clock: hybridLogicalClock{
				ts:    20,
				count: 55,
				wallClock: staticClock{
					value: 11,
				},
			},
			input: Stamp{
				Time:  10,
				Count: 34,
			},
			expected: Stamp{
				Time:  20,
				Count: 56,
			},
		},
		{
			name: "Input Ahead",
			clock: hybridLogicalClock{
				ts:    10,
				count: 55,
				wallClock: staticClock{
					value: 11,
				},
			},
			input: Stamp{
				Time:  20,
				Count: 34,
			},
			expected: Stamp{
				Time:  20,
				Count: 35,
			},
		},
		{
			name: "TS Equal, Clock Count Ahead",
			clock: hybridLogicalClock{
				ts:    10,
				count: 55,
				wallClock: staticClock{
					value: 10,
				},
			},
			input: Stamp{
				Time:  10,
				Count: 34,
			},
			expected: Stamp{
				Time:  10,
				Count: 56,
			},
		},
		{
			name: "TS Equal, Input Count Ahead",
			clock: hybridLogicalClock{
				ts:    10,
				count: 55,
				wallClock: staticClock{
					value: 10,
				},
			},
			input: Stamp{
				Time:  10,
				Count: 73,
			},
			expected: Stamp{
				Time:  10,
				Count: 74,
			},
		},
		{
			name: "All equal",
			clock: hybridLogicalClock{
				ts:    10,
				count: 55,
				node:  1,
				wallClock: staticClock{
					value: 10,
				},
			},
			input: Stamp{
				Time:  10,
				Count: 55,
				Node:  3,
			},
			expected: Stamp{
				Time:  10,
				Count: 56,
				Node:  1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.clock.Update(test.input)
			assertEquals(t, test.expected.Time, actual.Time)
			assertEquals(t, test.expected.Count, actual.Count)
			assertEquals(t, test.expected.Node, actual.Node)
		})
	}
}

func TestCompare(t *testing.T) {
	type Test struct {
		name  string
		t1    Stamp
		t2    Stamp
		after bool
	}

	var tests = []Test{
		{
			name: "Equal",
			t1: Stamp{
				Time:  1,
				Count: 1,
				Node:  1,
			},
			t2: Stamp{
				Time:  1,
				Count: 1,
				Node:  1,
			},
		},
		{
			name: "a > b, ts",
			t1: Stamp{
				Time:  2,
				Count: 1,
				Node:  1,
			},
			t2: Stamp{
				Time:  1,
				Count: 1,
				Node:  1,
			},
			after: true,
		},
		{
			name: "a < b, ts",
			t1: Stamp{
				Time:  2,
				Count: 1,
				Node:  1,
			},
			t2: Stamp{
				Time:  3,
				Count: 1,
				Node:  1,
			},
		},
		{
			name: "a < b, ts, ignore counts",
			t1: Stamp{
				Time:  2,
				Count: 10,
				Node:  1,
			},
			t2: Stamp{
				Time:  3,
				Count: 1,
				Node:  1,
			},
		},
		{
			name: "a > b, count",
			t1: Stamp{
				Time:  2,
				Count: 2,
				Node:  1,
			},
			t2: Stamp{
				Time:  2,
				Count: 1,
				Node:  1,
			},
			after: true,
		},
		{
			name: "a < b, count",
			t1: Stamp{
				Time:  2,
				Count: 1,
				Node:  1,
			},
			t2: Stamp{
				Time:  2,
				Count: 2,
				Node:  1,
			},
		},
		{
			name: "a < b, count, ignore  node",
			t1: Stamp{
				Time:  2,
				Count: 1,
				Node:  10,
			},
			t2: Stamp{
				Time:  2,
				Count: 2,
				Node:  1,
			},
		},
		{
			name: "a > b,  node",
			t1: Stamp{
				Time:  2,
				Count: 1,
				Node:  2,
			},
			t2: Stamp{
				Time:  2,
				Count: 1,
				Node:  1,
			},
			after: true,
		},
		{
			name: "a < b,  node",
			t1: Stamp{
				Time:  2,
				Count: 1,
				Node:  1,
			},
			t2: Stamp{
				Time:  2,
				Count: 1,
				Node:  2,
			},
		},
	}

	for _, test := range tests {
		actual := test.t1.After(test.t2)
		if actual != test.after {
			t.Errorf("Test %v failed, expected %v actual %v", test.name, test.after, actual)
		}
	}
}
