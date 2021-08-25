package dsl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolver(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		expStr       string
		solverMap    map[string]PatternResult
		expectedResp bool
		message      string
	}{
		{
			`"1"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			true,
			"single word true",
		},
		{
			`"1"`,
			map[string]PatternResult{},
			false,
			"single word false",
		},

		// and tests
		{
			`"1" and "2" and "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			true,
			"and multi term true",
		},
		{
			`"1" and "2" and "3"`,
			map[string]PatternResult{
				"2": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			false,
			"and multi term false 1",
		},
		{
			`"1" and "2" and "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			false,
			"and multi term false 2",
		},
		{
			`"1" and "2" and "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
			},
			false,
			"and multi term false 3",
		},

		// or tests
		{
			`"1" or "2" or "3"`,
			map[string]PatternResult{},
			false,
			"or multi term false",
		},
		{
			`"1" or "2" or "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			true,
			"or multi term true 1",
		},
		{
			`"1" or "2" or "3"`,
			map[string]PatternResult{
				"2": PatternResult{Val: true},
			},
			true,
			"or multi term true 2",
		},
		{
			`"1" or "2" or "3"`,
			map[string]PatternResult{
				"3": PatternResult{Val: true},
			},
			true,
			"or multi term true 3",
		},

		// not tests
		{
			`not "1"`,
			map[string]PatternResult{},
			true,
			"not true",
		},
		{
			`not "1"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			false,
			"not false",
		},
		{
			`not "1" or not "2"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
			},
			false,
			"not multi false",
		},
		{
			`not ("1" or "2")`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			false,
			"not multi false 2",
		},
		{
			`"1" and not "2" or "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			true,
			"not multi true 1",
		},
		{
			` not "2" and "1" or "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			true,
			"not multi true 2",
		},
		{
			`"1" and "3" or not "2"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			true,
			"not multi true 3",
		},

		// parentheses tests
		{
			`("1" or "2") and not ("1" and "2")`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
			},
			false,
			"parentheses xor 1",
		},
		{
			`("1" or "2") and not ("1" and "2")`,
			map[string]PatternResult{},
			false,
			"parentheses xor 2",
		},
		{
			`("1" or "2") and not ("1" and "2")`,
			map[string]PatternResult{
				"2": PatternResult{Val: true},
			},
			true,
			"parentheses xor 3",
		},
		{
			`("1" or "2") and not ("1" and "2")`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			true,
			"parentheses xor 4",
		},
		{
			`(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
			map[string]PatternResult{},
			false,
			"parentheses 1",
		},
		{
			`(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
			map[string]PatternResult{
				"4": PatternResult{Val: true},
				"6": PatternResult{Val: true},
				"8": PatternResult{Val: true},
			},
			true,
			"parentheses 2",
		},
		{
			`(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
			map[string]PatternResult{
				"4": PatternResult{Val: true},
				"6": PatternResult{Val: true},
			},
			false,
			"parentheses 3",
		},
		{
			`(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
				"3": PatternResult{Val: true},
				"6": PatternResult{Val: true},
				"8": PatternResult{Val: true},
			},
			true,
			"parentheses 4",
		},
	}

	for _, tc := range tests {
		exp, err := NewParser(strings.NewReader(tc.expStr)).Parse()
		assert.Nil(err, tc.message+" iter")
		respInt, err := exp.Solve(tc.solverMap, false)
		assert.Nil(err, tc.message+" iter")
		assert.Equal(tc.expectedResp, respInt, tc.message+" iter")
	}
}

func TestSolverInter(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		expStr       string
		solverMap    map[string]PatternResult
		expectedResp bool
		message      string
	}{
		{
			`"1"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			true,
			"single word true",
		},
		{
			`"1"`,
			map[string]PatternResult{},
			false,
			"single word false",
		},

		// and tests
		{
			`"1" and "2" and "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			true,
			"and multi term true",
		},
		{
			`"1" and "2" and "3"`,
			map[string]PatternResult{
				"2": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			false,
			"and multi term false 1",
		},
		{
			`"1" and "2" and "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			false,
			"and multi term false 2",
		},
		{
			`"1" and "2" and "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
			},
			false,
			"and multi term false 3",
		},

		// or tests
		{
			`"1" or "2" or "3"`,
			map[string]PatternResult{},
			false,
			"or multi term false",
		},
		{
			`"1" or "2" or "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			true,
			"or multi term true 1",
		},
		{
			`"1" or "2" or "3"`,
			map[string]PatternResult{
				"2": PatternResult{Val: true},
			},
			true,
			"or multi term true 2",
		},
		{
			`"1" or "2" or "3"`,
			map[string]PatternResult{
				"3": PatternResult{Val: true},
			},
			true,
			"or multi term true 3",
		},

		// not tests
		{
			`not "1"`,
			map[string]PatternResult{},
			true,
			"not true",
		},
		{
			`not "1"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			false,
			"not false",
		},
		{
			`not "1" or not "2"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
			},
			false,
			"not multi false",
		},
		{
			`not ("1" or "2")`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			false,
			"not multi false 2",
		},
		{
			`"1" and not "2" or "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			true,
			"not multi true 1",
		},
		{
			` not "2" and "1" or "3"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			true,
			"not multi true 2",
		},
		{
			`"1" and "3" or not "2"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"3": PatternResult{Val: true},
			},
			true,
			"not multi true 3",
		},

		// parentheses tests
		{
			`("1" or "2") and not ("1" and "2")`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
			},
			false,
			"parentheses xor 1",
		},
		{
			`("1" or "2") and not ("1" and "2")`,
			map[string]PatternResult{},
			false,
			"parentheses xor 2",
		},
		{
			`("1" or "2") and not ("1" and "2")`,
			map[string]PatternResult{
				"2": PatternResult{Val: true},
			},
			true,
			"parentheses xor 3",
		},
		{
			`("1" or "2") and not ("1" and "2")`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
			},
			true,
			"parentheses xor 4",
		},
		{
			`(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
			map[string]PatternResult{},
			false,
			"parentheses 1",
		},
		{
			`(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
			map[string]PatternResult{
				"4": PatternResult{Val: true},
				"6": PatternResult{Val: true},
				"8": PatternResult{Val: true},
			},
			true,
			"parentheses 2",
		},
		{
			`(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
			map[string]PatternResult{
				"4": PatternResult{Val: true},
				"6": PatternResult{Val: true},
			},
			false,
			"parentheses 3",
		},
		{
			`(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
			map[string]PatternResult{
				"1": PatternResult{Val: true},
				"2": PatternResult{Val: true},
				"3": PatternResult{Val: true},
				"6": PatternResult{Val: true},
				"8": PatternResult{Val: true},
			},
			true,
			"parentheses 4",
		},
	}

	for _, tc := range tests {
		exp, err := NewParser(strings.NewReader(tc.expStr)).Parse()
		arr := exp.CreateIteractive()
		assert.Nil(err, tc.message+" iter")
		respInt, err := arr.Solve(tc.solverMap, false)
		assert.Nil(err, tc.message+" iter")
		assert.Equal(tc.expectedResp, respInt, tc.message+" iter")
	}
}
