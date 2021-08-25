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
		expStr      string
		expectedExp Expression
		expectedErr error
		message     string
	}{
		{
			`"1"`,
			Expression{
				Type:    UNIT_EXPR,
				Literal: "1",
			},
			nil,
			"single word",
		},
		{
			`("1")`,
			Expression{
				Type:    UNIT_EXPR,
				Literal: "1",
			},
			nil,
			"single word parentheses",
		},
		{
			`"1" and "2"`,
			Expression{
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
			nil,
			"simple and",
		},
		{
			`("1" and "2")`,
			Expression{
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
			nil,
			"simple and parentheses",
		},
		{
			`"1" or "2"`,
			Expression{
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
			nil,
			"simple or",
		},
		{
			`not "1"`,
			Expression{
				Type: NOT_EXPR,
				RExpr: &Expression{
					Type:    UNIT_EXPR,
					Literal: "1",
				},
			},
			nil,
			"simple not",
		},
		{
			`"1" and "2" or not "3"`,
			Expression{
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
			nil,
			"multiple function no parentheses",
		},
		{
			`"1" and ("2" or not "3")`,
			Expression{
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
			nil,
			"multiple function with parentheses",
		},
		{
			``,
			Expression{},
			fmt.Errorf("invalid expression: unexpected EOF found"),
			"empty expression",
		},
		{
			`(("1")`,
			Expression{},
			fmt.Errorf("invalid expression: Unexpected '('"),
			"invalid open parentheses",
		},
		{
			`("1"))`,
			Expression{},
			fmt.Errorf("invalid expression: unexpected EOF found. Extra closing parentheses: 1"),
			"invalid close parentheses",
		},
		{
			`and`,
			Expression{},
			fmt.Errorf("invalid expression: no left expression was found for AND"),
			"invalid expression empty dual exp",
		},
		{
			` "1" and `,
			Expression{},
			fmt.Errorf("invalid expression: incomplete expression AND"),
			"invalid expression incomplete dual exp",
		},
		{
			`or`,
			Expression{},
			fmt.Errorf("invalid expression: no left expression was found for OR"),
			"invalid expression empty dual exp",
		},
		{
			` "1" or `,
			Expression{},
			fmt.Errorf("invalid expression: incomplete expression OR"),
			"invalid expression incomplete dual exp",
		},
		{
			`not`,
			Expression{},
			fmt.Errorf("invalid expression: Unexpected token 'EOF' after NOT"),
			"invalid expression incomplete dual exp",
		},
	}

	for _, tc := range tests {
		exp, err := NewParser(strings.NewReader(tc.expStr)).Parse()
		assert.Equal(tc.expectedErr, err, tc.message)
		if err == nil {
			assert.Equal(tc.expectedExp, *exp, tc.message)
		}
	}
}
