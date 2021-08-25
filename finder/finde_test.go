package finder

import (
	"testing"

	"github.com/pedroegsilva/gofindthem/dsl"

	"github.com/stretchr/testify/assert"
)

type expectedAddExpression struct {
	exprs    []exprWrapper
	keywords map[string]struct{}
	errors   []error
}

func TestAddExpression(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		finder      *Finder
		expressions []string
		expected    expectedAddExpression
		message     string
	}{
		{
			NewFinder(&EmptyEngine{}),
			[]string{`"a" and "b"`, `not "c"`},
			expectedAddExpression{
				[]exprWrapper{
					exprWrapper{
						`"a" and "b"`,
						dsl.SolverOrder{
							&dsl.Expression{
								Type: dsl.AND_EXPR,
								LExpr: &dsl.Expression{
									Type:    dsl.UNIT_EXPR,
									Literal: "a",
								},
								RExpr: &dsl.Expression{
									Type:    dsl.UNIT_EXPR,
									Literal: "b",
								},
							},
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "a",
							},
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "b",
							},
						},
					},
					exprWrapper{
						`not "c"`,
						dsl.SolverOrder{
							&dsl.Expression{
								Type: dsl.NOT_EXPR,
								RExpr: &dsl.Expression{
									Type:    dsl.UNIT_EXPR,
									Literal: "c",
								},
							},
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "c",
							},
						},
					},
				},
				map[string]struct{}{
					"a": struct{}{},
					"b": struct{}{},
					"c": struct{}{},
				},
				[]error{nil, nil},
			},
			"success test",
		},
	}

	for _, tc := range tests {
		count := 0
		for _, exp := range tc.expressions {
			err := tc.finder.AddExpression(exp)
			assert.Equal(tc.expected.errors[count], err, tc.message)
			count++
		}
		assert.Equal(tc.expected.exprs, tc.finder.expressions, tc.message)
		assert.Equal(tc.expected.keywords, tc.finder.keywords, tc.message)
	}
}
