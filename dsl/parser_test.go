package dsl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		expStr           string
		expectedExp      Expression
		expectedKeywords map[string]struct{}
		expectedErr      error
		caseSense        bool
		message          string
	}{
		{
			expStr: `"1"`,
			expectedExp: Expression{
				Type:    UNIT_EXPR,
				Literal: "1",
			},
			expectedKeywords: map[string]struct{}{
				"1": struct{}{},
			},
			expectedErr: nil,
			caseSense:   true,
			message:     "single word",
		},
		{
			expStr: `("1")`,
			expectedExp: Expression{
				Type:    UNIT_EXPR,
				Literal: "1",
			},
			expectedKeywords: map[string]struct{}{
				"1": struct{}{},
			},
			expectedErr: nil,
			caseSense:   true,
			message:     "single word parentheses",
		},
		{
			expStr: `"1" and "2"`,
			expectedExp: Expression{
				Type: AND_EXPR,
				LExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "1",
				},
				RExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "2",
				},
			},
			expectedKeywords: map[string]struct{}{
				"1": struct{}{},
				"2": struct{}{},
			},
			expectedErr: nil,
			caseSense:   true,
			message:     "simple and",
		},
		{
			expStr: `("1" and "2")`,
			expectedExp: Expression{
				Type: AND_EXPR,
				LExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "1",
				},
				RExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "2",
				},
			},
			expectedKeywords: map[string]struct{}{
				"1": struct{}{},
				"2": struct{}{},
			},
			expectedErr: nil,
			caseSense:   true,
			message:     "simple and parentheses",
		},
		{
			expStr: `"1" or "2"`,
			expectedExp: Expression{
				Type: OR_EXPR,
				LExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "1",
				},
				RExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "2",
				},
			},
			expectedKeywords: map[string]struct{}{
				"1": struct{}{},
				"2": struct{}{},
			},
			expectedErr: nil,
			caseSense:   true,
			message:     "simple or",
		},
		{
			expStr: `not "1"`,
			expectedExp: Expression{
				Type: NOT_EXPR,
				RExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "1",
				},
			},
			expectedKeywords: map[string]struct{}{
				"1": struct{}{},
			},
			expectedErr: nil,
			caseSense:   true,
			message:     "simple not",
		},
		{
			expStr: `"1" and "2" or not "3"`,
			expectedExp: Expression{
				Type: OR_EXPR,
				LExpr: &Expression{
					Type: AND_EXPR,
					LExpr: &Expression{
						Type:    UNIT_EXPR,
						Literal: "1",
					},
					RExpr: &Expression{
						Type:    UNIT_EXPR,
						Literal: "2",
					},
				},
				RExpr: &Expression{
					Type: NOT_EXPR,
					RExpr: &Expression{
						Type:    UNIT_EXPR,
						Literal: "3",
					},
				},
			},
			expectedKeywords: map[string]struct{}{
				"1": struct{}{},
				"2": struct{}{},
				"3": struct{}{},
			},
			expectedErr: nil,
			caseSense:   true,
			message:     "multiple function no parentheses",
		},
		{
			expStr: `"1" and ("2" or not "3")`,
			expectedExp: Expression{
				Type: AND_EXPR,
				LExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "1",
				},
				RExpr: &Expression{
					Type: OR_EXPR,
					LExpr: &Expression{
						Type:    UNIT_EXPR,
						Literal: "2",
					},
					RExpr: &Expression{
						Type: NOT_EXPR,
						RExpr: &Expression{
							Type:    UNIT_EXPR,
							Literal: "3",
						},
					},
				},
			},
			expectedKeywords: map[string]struct{}{
				"1": struct{}{},
				"2": struct{}{},
				"3": struct{}{},
			},
			expectedErr: nil,
			caseSense:   true,
			message:     "multiple function with parentheses",
		},
		{
			expStr:           ``,
			expectedExp:      Expression{},
			expectedKeywords: map[string]struct{}{},
			expectedErr:      fmt.Errorf("invalid expression: unexpected EOF found"),
			caseSense:        true,
			message:          "empty expression",
		},
		{
			expStr:           `(("1")`,
			expectedExp:      Expression{},
			expectedKeywords: map[string]struct{}{},
			expectedErr:      fmt.Errorf("invalid expression: Unexpected '('"),
			caseSense:        true,
			message:          "invalid open parentheses",
		},
		{
			expStr:           `("1"))`,
			expectedExp:      Expression{},
			expectedKeywords: map[string]struct{}{},
			expectedErr:      fmt.Errorf("invalid expression: unexpected EOF found. Extra closing parentheses: 1"),
			caseSense:        true,
			message:          "invalid close parentheses",
		},
		{
			expStr:           `and`,
			expectedExp:      Expression{},
			expectedKeywords: map[string]struct{}{},
			expectedErr:      fmt.Errorf("invalid expression: no left expression was found for AND"),
			caseSense:        true,
			message:          "invalid expression empty dual exp",
		},
		{
			expStr:           ` "1" and `,
			expectedExp:      Expression{},
			expectedKeywords: map[string]struct{}{},
			expectedErr:      fmt.Errorf("invalid expression: incomplete expression AND"),
			caseSense:        true,
			message:          "invalid expression incomplete dual exp",
		},
		{
			expStr:           `or`,
			expectedExp:      Expression{},
			expectedKeywords: map[string]struct{}{},
			expectedErr:      fmt.Errorf("invalid expression: no left expression was found for OR"),
			caseSense:        true,
			message:          "invalid expression empty dual exp",
		},
		{
			expStr:           ` "1" or `,
			expectedExp:      Expression{},
			expectedKeywords: map[string]struct{}{},
			expectedErr:      fmt.Errorf("invalid expression: incomplete expression OR"),
			caseSense:        true,
			message:          "invalid expression incomplete dual exp",
		},
		{
			expStr:           `not`,
			expectedExp:      Expression{},
			expectedKeywords: map[string]struct{}{},
			expectedErr:      fmt.Errorf("invalid expression: Unexpected token 'EOF' after NOT"),
			caseSense:        true,
			message:          "invalid expression incomplete dual exp",
		},
		{
			expStr: `"CaSe In sensItIVe" and "SomeThing"`,
			expectedExp: Expression{
				Type: AND_EXPR,
				LExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "case in sensitive",
				},
				RExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "something",
				},
			},
			expectedKeywords: map[string]struct{}{
				"case in sensitive": struct{}{},
				"something":         struct{}{},
			},
			expectedErr: nil,
			caseSense:   false,
			message:     "case insensitive test",
		},
	}

	for _, tc := range tests {
		p := NewParser(strings.NewReader(tc.expStr), tc.caseSense)
		exp, err := p.Parse()
		assert.Equal(tc.expectedErr, err, tc.message)
		if err == nil {
			assert.Equal(tc.expectedExp, *exp, tc.message)
			assert.Equal(tc.expectedKeywords, p.GetKeywords(), tc.message)
		}
	}
}
