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
		expStr        string
		expectedExp   Expression
		expectedTags  map[string]struct{}
		expectedPaths map[string]struct{}
		expectedErr   error
		message       string
	}{
		{
			expStr: `"tag1"`,
			expectedExp: Expression{
				Type: UNIT_EXPR,
				Tag: TagInfo{
					Name:      "tag1",
					FieldPath: "",
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
			},
			expectedPaths: map[string]struct{}{},
			expectedErr:   nil,
			message:       "single tag",
		},
		{
			expStr: `"tag1:field1"`,
			expectedExp: Expression{
				Type: UNIT_EXPR,
				Tag: TagInfo{
					Name:      "tag1",
					FieldPath: "field1",
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
			},
			expectedErr: nil,
			message:     "tag with field",
		},
		{
			expStr: `("tag1:field1")`,
			expectedExp: Expression{
				Type: UNIT_EXPR,
				Tag: TagInfo{
					Name:      "tag1",
					FieldPath: "field1",
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
			},
			expectedErr: nil,
			message:     "single tag parentheses",
		},
		{
			expStr: `"tag1:field1" and "tag2:field2"`,
			expectedExp: Expression{
				Type: AND_EXPR,
				LExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag1",
						FieldPath: "field1",
					},
				},
				RExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag2",
						FieldPath: "field2",
					},
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
				"tag2": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
				"field2": {},
			},
			expectedErr: nil,
			message:     "simple and",
		},
		{
			expStr: `("tag1:field1" and "tag2:field2")`,
			expectedExp: Expression{
				Type: AND_EXPR,
				LExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag1",
						FieldPath: "field1",
					},
				},
				RExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag2",
						FieldPath: "field2",
					},
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
				"tag2": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
				"field2": {},
			},
			expectedErr: nil,
			message:     "simple and with parentheses",
		},
		{
			expStr: `"tag1:field1" or "tag2:field2"`,
			expectedExp: Expression{
				Type: OR_EXPR,
				LExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag1",
						FieldPath: "field1",
					},
				},
				RExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag2",
						FieldPath: "field2",
					},
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
				"tag2": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
				"field2": {},
			},
			expectedErr: nil,
			message:     "simple or",
		},
		{
			expStr: `not "tag1:field1"`,
			expectedExp: Expression{
				Type: NOT_EXPR,
				RExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag1",
						FieldPath: "field1",
					},
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
			},
			expectedErr: nil,
			message:     "simple not",
		},
		{
			expStr: `"tag1:field1" and "tag2:field2" or not "tag3:field3"`,
			expectedExp: Expression{
				Type: OR_EXPR,
				LExpr: &Expression{
					Type: AND_EXPR,
					LExpr: &Expression{
						Type: UNIT_EXPR,
						Tag: TagInfo{
							Name:      "tag1",
							FieldPath: "field1",
						},
					},
					RExpr: &Expression{
						Type: UNIT_EXPR,
						Tag: TagInfo{
							Name:      "tag2",
							FieldPath: "field2",
						},
					},
				},
				RExpr: &Expression{
					Type: NOT_EXPR,
					RExpr: &Expression{
						Type: UNIT_EXPR,
						Tag: TagInfo{
							Name:      "tag3",
							FieldPath: "field3",
						},
					},
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
				"tag2": {},
				"tag3": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
				"field2": {},
				"field3": {},
			},
			expectedErr: nil,
			message:     "multiple function no parentheses",
		},
		{
			expStr: `"tag1:field1" and ("tag2:field2" or not "tag3:field3")`,
			expectedExp: Expression{
				Type: AND_EXPR,
				LExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag1",
						FieldPath: "field1",
					},
				},
				RExpr: &Expression{
					Type: OR_EXPR,
					LExpr: &Expression{
						Type: UNIT_EXPR,
						Tag: TagInfo{
							Name:      "tag2",
							FieldPath: "field2",
						},
					},
					RExpr: &Expression{
						Type: NOT_EXPR,
						RExpr: &Expression{
							Type: UNIT_EXPR,
							Tag: TagInfo{
								Name:      "tag3",
								FieldPath: "field3",
							},
						},
					},
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
				"tag2": {},
				"tag3": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
				"field2": {},
				"field3": {},
			},
			expectedErr: nil,
			message:     "multiple function with parentheses",
		},
		{
			expStr: `not("tag2:field2" or "tag3:field3") and "tag1:field1"`,
			expectedExp: Expression{
				Type: AND_EXPR,
				LExpr: &Expression{
					Type: NOT_EXPR,
					RExpr: &Expression{
						Type: OR_EXPR,
						LExpr: &Expression{
							Type: UNIT_EXPR,
							Tag: TagInfo{
								Name:      "tag2",
								FieldPath: "field2",
							},
						},
						RExpr: &Expression{
							Type: UNIT_EXPR,
							Tag: TagInfo{
								Name:      "tag3",
								FieldPath: "field3",
							},
						},
					},
				},
				RExpr: &Expression{
					Type: UNIT_EXPR,
					Tag: TagInfo{
						Name:      "tag1",
						FieldPath: "field1",
					},
				},
			},
			expectedTags: map[string]struct{}{
				"tag1": {},
				"tag2": {},
				"tag3": {},
			},
			expectedPaths: map[string]struct{}{
				"field1": {},
				"field2": {},
				"field3": {},
			},
			expectedErr: nil,
			message:     "not with parentheses",
		},
		{
			expStr:        ``,
			expectedExp:   Expression{},
			expectedTags:  map[string]struct{}{},
			expectedPaths: map[string]struct{}{},
			expectedErr:   fmt.Errorf("invalid expression: unexpected EOF found"),
			message:       "empty expression",
		},
		{
			expStr:        `(("tag1")`,
			expectedExp:   Expression{},
			expectedTags:  map[string]struct{}{},
			expectedPaths: map[string]struct{}{},
			expectedErr:   fmt.Errorf("invalid expression: Unexpected '('"),
			message:       "invalid open parentheses",
		},
		{
			expStr:        `("tag1"))`,
			expectedExp:   Expression{},
			expectedTags:  map[string]struct{}{},
			expectedPaths: map[string]struct{}{},
			expectedErr:   fmt.Errorf("invalid expression: unexpected EOF found. Extra closing parentheses: 1"),
			message:       "invalid close parentheses",
		},
		{
			expStr:        `and`,
			expectedExp:   Expression{},
			expectedTags:  map[string]struct{}{},
			expectedPaths: map[string]struct{}{},
			expectedErr:   fmt.Errorf("invalid expression: no left expression was found for AND"),
			message:       "invalid expression empty dual exp",
		},
		{
			expStr:        ` "tag1" and `,
			expectedExp:   Expression{},
			expectedTags:  map[string]struct{}{},
			expectedPaths: map[string]struct{}{},
			expectedErr:   fmt.Errorf("invalid expression: incomplete expression AND"),

			message: "invalid expression incomplete dual exp",
		},
		{
			expStr:        `or`,
			expectedExp:   Expression{},
			expectedTags:  map[string]struct{}{},
			expectedPaths: map[string]struct{}{},
			expectedErr:   fmt.Errorf("invalid expression: no left expression was found for OR"),

			message: "invalid expression empty dual exp",
		},
		{
			expStr:        ` "tag1" or `,
			expectedExp:   Expression{},
			expectedTags:  map[string]struct{}{},
			expectedPaths: map[string]struct{}{},
			expectedErr:   fmt.Errorf("invalid expression: incomplete expression OR"),

			message: "invalid expression incomplete dual exp",
		},
		{
			expStr:        `not`,
			expectedExp:   Expression{},
			expectedTags:  map[string]struct{}{},
			expectedPaths: map[string]struct{}{},
			expectedErr:   fmt.Errorf("invalid expression: Unexpected token 'EOF' after NOT"),
			message:       "invalid expression incomplete dual exp",
		},
	}

	for _, tc := range tests {
		p := NewParser(strings.NewReader(tc.expStr))
		exp, err := p.Parse()
		assert.Equal(tc.expectedErr, err, tc.message)
		if err == nil {
			assert.Equal(tc.expectedExp, *exp, tc.message)
			assert.Equal(tc.expectedTags, p.tags, tc.message)
			assert.Equal(tc.expectedPaths, p.fields, tc.message)
		}
	}
}
