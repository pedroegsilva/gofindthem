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
		assert.Nil(err, tc.message+" iter")
		respInt, err := exp.Solve(tc.solverMap, false)
		assert.Nil(err, tc.message+" iter")
		assert.Equal(tc.expectedResp, respInt, tc.message+" iter")
	}
}

func TestSolverInter(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range solverTestCases {
		exp, err := NewParser(strings.NewReader(tc.expStr), true).Parse()
		arr := exp.CreateSolverOrder()
		assert.Nil(err, tc.message+" iter")
		respInt, err := arr.Solve(tc.solverMap, false)
		assert.Nil(err, tc.message+" iter")
		assert.Equal(tc.expectedResp, respInt, tc.message+" iter")
	}
}

// TODO(pedro.silva) create unit test for CreateSolverOrder
func TestCreateSolverOrder(t *testing.T) {
	assert := assert.New(t)

	unit1 := &Expression{
		Type:    UNIT_EXPR,
		Literal: "sharpest",
	}
	unit2 := &Expression{
		Type:    UNIT_EXPR,
		Literal: "words",
	}
	unit3 := &Expression{
		Type:    UNIT_EXPR,
		Literal: "no one",
	}

	expr1 := &Expression{
		Type:  OR_EXPR,
		LExpr: unit2,
		RExpr: unit3,
	}

	//"sharpest" and "words" or "no one"
	exp1 := &Expression{
		Type:  AND_EXPR,
		LExpr: unit1,
		RExpr: expr1,
	}

	notExp := &Expression{
		Type:  NOT_EXPR,
		RExpr: unit1,
	}
	// Not "sharpest" and "words" or "no one"
	exp2 := &Expression{
		Type:  AND_EXPR,
		LExpr: notExp,
		RExpr: expr1,
	}

	tests := []struct {
		exp         *Expression
		expectedArr SolverOrder
		message     string
	}{
		{
			exp:         exp1,
			expectedArr: SolverOrder{exp1, unit1, expr1, unit2, unit3},
			message:     `"sharpest" and "words" or "no one"`,
		},
		{
			exp:         exp2,
			expectedArr: SolverOrder{exp2, notExp, unit1, expr1, unit2, unit3},
			message:     `Not "sharpest" and "words" or "no one"`,
		},
	}
	for _, tc := range tests {
		firstExp := tc.exp
		arr := tc.exp.CreateSolverOrder()
		assert.Equal(tc.expectedArr, arr, tc.message)
		// Asserts that the expression was not changed by the createSolver order
		assert.Equal(firstExp, tc.exp, tc.message)
	}
}

var solverTestCases = []struct {
	expStr       string
	solverMap    map[string]PatternResult
	expectedResp bool
	message      string
}{
	{
		expStr: `"1"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "single word true",
	},
	{
		expStr:       `"1"`,
		solverMap:    map[string]PatternResult{},
		expectedResp: false,
		message:      "single word false",
	},

	// and tests
	{
		expStr: `"1" and "2" and "3"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"2": PatternResult{Val: true},
			"3": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "and multi term true",
	},
	{
		expStr: `"1" and "2" and "3"`,
		solverMap: map[string]PatternResult{
			"2": PatternResult{Val: true},
			"3": PatternResult{Val: true},
		},
		expectedResp: false,
		message:      "and multi term false 1",
	},
	{
		expStr: `"1" and "2" and "3"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"3": PatternResult{Val: true},
		},
		expectedResp: false,
		message:      "and multi term false 2",
	},
	{
		expStr: `"1" and "2" and "3"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"2": PatternResult{Val: true},
		},
		expectedResp: false,
		message:      "and multi term false 3",
	},

	// or tests
	{
		expStr:       `"1" or "2" or "3"`,
		solverMap:    map[string]PatternResult{},
		expectedResp: false,
		message:      "or multi term false",
	},
	{
		expStr: `"1" or "2" or "3"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "or multi term true 1",
	},
	{
		expStr: `"1" or "2" or "3"`,
		solverMap: map[string]PatternResult{
			"2": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "or multi term true 2",
	},
	{
		expStr: `"1" or "2" or "3"`,
		solverMap: map[string]PatternResult{
			"3": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "or multi term true 3",
	},

	// not tests
	{
		expStr:       `not "1"`,
		solverMap:    map[string]PatternResult{},
		expectedResp: true,
		message:      "not true",
	},
	{
		expStr: `not "1"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
		},
		expectedResp: false,
		message:      "not false",
	},
	{
		expStr: `not "1" or not "2"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"2": PatternResult{Val: true},
		},
		expectedResp: false,
		message:      "not multi false",
	},
	{
		expStr: `not ("1" or "2")`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
		},
		expectedResp: false,
		message:      "not multi false 2",
	},
	{
		expStr: `"1" and not "2" or "3"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"3": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "not multi true 1",
	},
	{
		expStr: ` not "2" and "1" or "3"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"3": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "not multi true 2",
	},
	{
		expStr: `"1" and "3" or not "2"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"3": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "not multi true 3",
	},

	// parentheses tests
	{
		expStr: `("1" or "2") and not ("1" and "2")`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"2": PatternResult{Val: true},
		},
		expectedResp: false,
		message:      "parentheses xor 1",
	},
	{
		expStr:       `("1" or "2") and not ("1" and "2")`,
		solverMap:    map[string]PatternResult{},
		expectedResp: false,
		message:      "parentheses xor 2",
	},
	{
		expStr: `("1" or "2") and not ("1" and "2")`,
		solverMap: map[string]PatternResult{
			"2": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "parentheses xor 3",
	},
	{
		expStr: `("1" or "2") and not ("1" and "2")`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "parentheses xor 4",
	},
	{
		expStr:       `(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
		solverMap:    map[string]PatternResult{},
		expectedResp: false,
		message:      "parentheses 1",
	},
	{
		expStr: `(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
		solverMap: map[string]PatternResult{
			"4": PatternResult{Val: true},
			"6": PatternResult{Val: true},
			"8": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "parentheses 2",
	},
	{
		expStr: `(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
		solverMap: map[string]PatternResult{
			"4": PatternResult{Val: true},
			"6": PatternResult{Val: true},
		},
		expectedResp: false,
		message:      "parentheses 3",
	},
	{
		expStr: `(("1" and "2" and "3") or ("4" and not "5")) and ("6" or "7") and "8"`,
		solverMap: map[string]PatternResult{
			"1": PatternResult{Val: true},
			"2": PatternResult{Val: true},
			"3": PatternResult{Val: true},
			"6": PatternResult{Val: true},
			"8": PatternResult{Val: true},
		},
		expectedResp: true,
		message:      "parentheses 4",
	},
}
