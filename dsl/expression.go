package dsl

import (
	"fmt"
	"strings"
)

// ExprType are Special tokens used to define the expression type
type ExprType int

const (
	UNSET_EXPR ExprType = iota
	AND_EXPR
	OR_EXPR
	NOT_EXPR
	UNIT_EXPR
)

// GetName returns a readable name for the ExprType value
func (exprType ExprType) GetName() string {
	switch exprType {
	case UNSET_EXPR:
		return "UNSET"
	case AND_EXPR:
		return "AND"
	case OR_EXPR:
		return "OR"
	case NOT_EXPR:
		return "NOT"
	case UNIT_EXPR:
		return "UNIT"
	default:
		return "UNEXPECTED"
	}
}

// Expression can be a literal (UNIT) or a function composed by
// one or two other expressions (NOT, AND, OR).
type Expression struct {
	LExpr   *Expression
	RExpr   *Expression
	Type    ExprType
	Literal string
	Val     bool
}

// PatternResult stores if the patter was matched on
// the text and the positions it was found
type PatternResult struct {
	Val            bool
	SortedMatchPos []int
}

// getTypeName returns the type of the expression with a readable name
func (exp *Expression) GetTypeName() string {
	return exp.Type.GetName()
}

// Solve solves the expresion recursively. It has the option to use a complete map of
// PatternResult or a incomplete map. If the complete map option is used the map must have
// all the terms needed to solve de expression or it will return an error.
// If the incomplete map is used, missing keys will be considered as a no match on the
// document.
func (exp *Expression) Solve(
	patterResByKeyword map[string]PatternResult,
	completeMap bool,
) (bool, error) {
	switch exp.Type {
	case UNIT_EXPR:
		if resp, ok := patterResByKeyword[exp.Literal]; ok {
			exp.Val = resp.Val
		} else {
			if completeMap {
				return false, fmt.Errorf(fmt.Sprintf("could not find key %s on map.", exp.Literal))
			} else {
				exp.Val = false
			}
		}
		return exp.Val, nil
	case AND_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, fmt.Errorf(fmt.Sprintf("And statment do not have rigth or left expression: %v", exp))
		}
		lval, err := exp.LExpr.Solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, err
		}
		rval, err := exp.RExpr.Solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, err
		}
		exp.Val = lval && rval

		return exp.Val, nil
	case OR_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, fmt.Errorf(fmt.Sprintf("OR statment do not have rigth or left expression: %v", exp))
		}
		lval, err := exp.LExpr.Solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, err
		}
		rval, err := exp.RExpr.Solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, err
		}
		exp.Val = lval || rval
		return exp.Val, nil
	case NOT_EXPR:
		if exp.RExpr == nil {
			return false, fmt.Errorf(fmt.Sprintf("NOT statment do not have expression: %v", exp))
		}
		lval, err := exp.RExpr.Solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, err
		}
		exp.Val = !lval
		return exp.Val, nil
	default:
		return false, fmt.Errorf(fmt.Sprintf("Unable to process expression type %d", exp.Type))
	}
}

// PrettyPrint returns the expression formated on a tabbed structure
// Eg: for the expression ("a" and "b") or "c"
// OR
//     AND
//         a
//         b
func (exp *Expression) PrettyFormat() string {
	return exp.prettyFormat(0)
}

func (exp *Expression) prettyFormat(lvl int) (pprint string) {
	tabs := "    "
	onLVL := strings.Repeat(tabs, lvl)
	if exp.Type == UNIT_EXPR {
		return fmt.Sprintf("%s%s\n", onLVL, exp.Literal)
	}
	pprint = fmt.Sprintf("%s%s\n", onLVL, exp.GetTypeName())
	if exp.LExpr != nil {
		pprint += exp.LExpr.prettyFormat(lvl + 1)
	}

	if exp.RExpr != nil {
		pprint += exp.RExpr.prettyFormat(lvl + 1)
	}

	return
}

// SolverOrder store the expressions Preorder
type SolverOrder []*Expression

// Solve solves the expresion iteratively. It has the option to use a complete map of
// PatternResult or a incomplete map. If the complete map option is used the map must have
// all the terms needed to solve de expression or it will return an error.
// If the incomplete map is used, missing keys will be considered as a no match on the
// document.
func (so SolverOrder) Solve(patterResByKeyword map[string]PatternResult, completeMap bool) (bool, error) {
	for i := len(so) - 1; i >= 0; i-- {
		exp := so[i]
		if exp == nil {
			continue
		}
		switch exp.Type {
		case UNIT_EXPR:
			if resp, ok := patterResByKeyword[exp.Literal]; ok {
				exp.Val = resp.Val
			} else {
				if completeMap {
					return false, fmt.Errorf(fmt.Sprintf("could not find key %s on map.", exp.Literal))
				} else {
					exp.Val = false
				}
			}
		case AND_EXPR:
			if exp.LExpr == nil || exp.RExpr == nil {
				return false, fmt.Errorf(fmt.Sprintf("And statment do not have rigth or left expression: %v", exp))
			}
			exp.Val = exp.LExpr.Val && exp.RExpr.Val

		case OR_EXPR:
			if exp.LExpr == nil || exp.RExpr == nil {
				return false, fmt.Errorf(fmt.Sprintf("OR statment do not have rigth or left expression: %v", exp))
			}
			exp.Val = exp.LExpr.Val || exp.RExpr.Val

		case NOT_EXPR:
			if exp.RExpr == nil {
				return false, fmt.Errorf(fmt.Sprintf("NOT statment do not have expression: %v", exp))
			}

			exp.Val = !exp.RExpr.Val

		default:
			return false, fmt.Errorf(fmt.Sprintf("Unable to process expression type %d", exp.Type))
		}
	}
	return so[0].Val, nil
}

// CreateSolverOrder traverses the expression tree in Preorder and
// stores the expressions on SolverOrder
func (exp *Expression) CreateSolverOrder() SolverOrder {
	test := new(SolverOrder)
	createSolverOrder(exp, test)
	return *test
}

// createSolverOrder recursion that traverses the expression
// tree in Preorder
func createSolverOrder(exp *Expression, arr *SolverOrder) {
	(*arr) = append((*arr), exp)

	if exp.LExpr != nil {
		createSolverOrder(exp.LExpr, arr)
	}

	if exp.RExpr != nil {
		createSolverOrder(exp.RExpr, arr)
	}
}
