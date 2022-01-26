package dsl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolver(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range solverTestCases {
		exp, err := NewParser(strings.NewReader(tc.expStr), true).Parse()
		assert.Nil(err, tc.message+" recursive")
		respInt, err := exp.Solve(tc.sortedMatchesByKeyword)
		assert.Nil(err, tc.message+" recursive")
		assert.Equal(tc.expectedResp, respInt, tc.message+" recursive")
	}
}

var solverTestCases = []struct {
	expStr                 string
	sortedMatchesByKeyword map[string][]int
	expectedResp           bool
	message                string
}{
	{
		expStr: `"1"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
		},
		expectedResp: true,
		message:      "single word true",
	},
	{
		expStr:                 `"1"`,
		sortedMatchesByKeyword: map[string][]int{},
		expectedResp:           false,
		message:                "single word false",
	},

	// and tests
	{
		expStr: `"1" and "2" and "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"2": nil,
			"3": nil,
		},
		expectedResp: true,
		message:      "and multi term true",
	},
	{
		expStr: `"1" and "2" and "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"2": nil,
			"3": nil,
		},
		expectedResp: false,
		message:      "and multi term false 1",
	},
	{
		expStr: `"1" and "2" and "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"3": nil,
		},
		expectedResp: false,
		message:      "and multi term false 2",
	},
	{
		expStr: `"1" and "2" and "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"2": nil,
		},
		expectedResp: false,
		message:      "and multi term false 3",
	},

	// or tests
	{
		expStr:                 `"1" or "2" or "3"`,
		sortedMatchesByKeyword: map[string][]int{},
		expectedResp:           false,
		message:                "or multi term false",
	},
	{
		expStr: `"1" or "2" or "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
		},
		expectedResp: true,
		message:      "or multi term true 1",
	},
	{
		expStr: `"1" or "2" or "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"2": nil,
		},
		expectedResp: true,
		message:      "or multi term true 2",
	},
	{
		expStr: `"1" or "2" or "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"3": nil,
		},
		expectedResp: true,
		message:      "or multi term true 3",
	},

	// not tests
	{
		expStr:                 `not "1"`,
		sortedMatchesByKeyword: map[string][]int{},
		expectedResp:           true,
		message:                "not true",
	},
	{
		expStr: `not "1"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
		},
		expectedResp: false,
		message:      "not false",
	},
	{
		expStr: `not "1" or not "2"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"2": nil,
		},
		expectedResp: false,
		message:      "not multi false",
	},
	{
		expStr: `not ("1" or "2") or "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"3": nil,
		},
		expectedResp: true,
		message:      "not multi true",
	},
	{
		expStr: `"1" and not "2" or "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"3": nil,
		},
		expectedResp: true,
		message:      "not multi true 1",
	},
	{
		expStr: ` not "2" and "1" or "3"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"3": nil,
		},
		expectedResp: true,
		message:      "not multi true 2",
	},
	{
		expStr: `"1" and "3" or not "2"`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"3": nil,
		},
		expectedResp: true,
		message:      "not multi true 3",
	},
	// parentheses tests
	{
		expStr: `not ("1" and "2") and ("1" or "2")`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"2": nil,
		},
		expectedResp: false,
		message:      "parentheses xor 1",
	},
	{
		expStr:                 `("1" or "2") and not ("1" and "2")`,
		sortedMatchesByKeyword: map[string][]int{},
		expectedResp:           false,
		message:                "parentheses xor 2",
	},
	{
		expStr: `not ("1" and "2") and ("1" or "2")`,
		sortedMatchesByKeyword: map[string][]int{
			"2": nil,
		},
		expectedResp: true,
		message:      "parentheses xor 3",
	},
	{
		expStr: `("1" or "2") and not (r"1" and "2")`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
		},
		expectedResp: true,
		message:      "parentheses xor 4 with regex",
	},
	{
		expStr: `(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
		sortedMatchesByKeyword: map[string][]int{
			"4": nil,
			"6": nil,
			"8": nil,
		},
		expectedResp: true,
		message:      "parentheses 1",
	},
	{
		expStr: `(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
		sortedMatchesByKeyword: map[string][]int{
			"4": nil,
			"6": nil,
		},
		expectedResp: false,
		message:      "parentheses 2",
	},
	{
		expStr: `(not ("1" and "2" and r"3") or ("4" and not r"5"))`,
		sortedMatchesByKeyword: map[string][]int{
			"1": nil,
			"2": nil,
			"3": nil,
			"4": nil,
		},
		expectedResp: true,
		message:      "parentheses with regex",
	},
	// inord tests
	{
		expStr: `inord("a" and "b" and "c")`,
		// acabXaXcb
		sortedMatchesByKeyword: map[string][]int{
			"a": {0, 2, 5},
			"b": {3, 8},
			"c": {1, 7},
		},
		expectedResp: true,
		message:      "inord true",
	},
	{
		expStr: `inord("a" and ("b" or "c"))`,
		// cbac
		sortedMatchesByKeyword: map[string][]int{
			"a": {2},
			"b": {1},
			"c": {0, 3},
		},
		expectedResp: true,
		message:      "inord with or",
	},
	{
		expStr: `inord("a" and "b" and "c")`,
		// bacb
		sortedMatchesByKeyword: map[string][]int{
			"a": {1},
			"b": {0, 3},
			"c": {2},
		},
		expectedResp: false,
		message:      "inord false 1",
	},
	{
		expStr: `inord("a" and "b" and "c")`,
		// bcab
		sortedMatchesByKeyword: map[string][]int{
			"a": {2},
			"b": {0, 3},
			"c": {1},
		},
		expectedResp: false,
		message:      "inord false 2",
	},
	{
		expStr: `inord(("b" or "c") and ("a" or "b"))`,
		// bcab
		sortedMatchesByKeyword: map[string][]int{
			"a": {2},
			"b": {0, 3},
			"c": {1},
		},
		expectedResp: true,
		message:      "inord multiple or with repeated key",
	},
	{
		expStr: `inord("b" and "c") and inord("a" and "b")`,
		// bcab
		sortedMatchesByKeyword: map[string][]int{
			"a": {2},
			"b": {0, 3},
			"c": {1},
		},
		expectedResp: true,
		message:      "multiple inord",
	},
	{
		expStr: `inord("b" and r"c") and inord(r"a" and "b")`,
		// bcab
		sortedMatchesByKeyword: map[string][]int{
			"a": {2},
			"b": {0, 3},
			"c": {1},
		},
		expectedResp: true,
		message:      "multiple inord with regex",
	},
}
