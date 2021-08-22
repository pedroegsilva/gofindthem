package dsl

import (
	"fmt"
	"strings"
)

type ExprType int

const (
	// Special tokens
	UNSET_EXPR ExprType = iota
	AND_EXPR
	OR_EXPR
	NOT_EXPR
	UNIT_EXPR
)

// Expression can be a literal or a function composed by one or two other expressions.
type Expression struct {
	LExpr     *Expression
	RExpr     *Expression
	Type      ExprType
	Literal   string
	Val       bool
	Evaluated bool
}

type PatternResult struct {
	Val            bool
	SortedMatchPos []int
}

func (exp *Expression) getTypeName() string {
	switch exp.Type {
	case AND_EXPR:
		return "AND"
	case OR_EXPR:
		return "OR"
	case NOT_EXPR:
		return "NOT"
	case UNIT_EXPR:
		return "UNIT"
	}
	return "NONE"
}

func (exp *Expression) Solve(patterResByKeyword map[string]PatternResult, completeMap bool, enableCache bool) (bool, error) {
	switch exp.Type {
	case UNIT_EXPR:
		if !exp.Evaluated {
			if resp, ok := patterResByKeyword[exp.Literal]; ok {
				exp.Val = resp.Val
				if enableCache {
					exp.Evaluated = true
				}
			} else {
				if completeMap {
					return false, fmt.Errorf(fmt.Sprintf("could not find key %s on map.", exp.Literal))
				} else {
					exp.Val = false
					if enableCache {
						exp.Evaluated = true
					}
				}
			}
		}
		return exp.Val, nil
	case AND_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, fmt.Errorf(fmt.Sprintf("And statment do not have rigth or left expression: %v", exp))
		}
		if !exp.Evaluated {
			lval, err := exp.LExpr.Solve(patterResByKeyword, completeMap, enableCache)
			if err != nil {
				return false, err
			}
			rval, err := exp.RExpr.Solve(patterResByKeyword, completeMap, enableCache)
			if err != nil {
				return false, err
			}
			exp.Val = lval && rval
			if enableCache {
				exp.Evaluated = true
			}
		}
		return exp.Val, nil
	case OR_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, fmt.Errorf(fmt.Sprintf("OR statment do not have rigth or left expression: %v", exp))
		}
		if !exp.Evaluated {
			lval, err := exp.LExpr.Solve(patterResByKeyword, completeMap, enableCache)
			if err != nil {
				return false, err
			}
			rval, err := exp.RExpr.Solve(patterResByKeyword, completeMap, enableCache)
			if err != nil {
				return false, err
			}
			exp.Val = lval || rval
			if enableCache {
				exp.Evaluated = true
			}
		}
		return exp.Val, nil
	case NOT_EXPR:
		if exp.RExpr == nil {
			return false, fmt.Errorf(fmt.Sprintf("NOT statment do not have expression: %v", exp))
		}
		if !exp.Evaluated {
			lval, err := exp.RExpr.Solve(patterResByKeyword, completeMap, enableCache)
			if err != nil {
				return false, err
			}
			exp.Val = !lval
			if enableCache {
				exp.Evaluated = true
			}
		}
		return exp.Val, nil
	default:
		return false, fmt.Errorf(fmt.Sprintf("Unable to process expression type %d", exp.Type))
	}
}

func (exp *Expression) PrettyPrint() string {
	return exp.prettyPrint(0)
}

func (exp *Expression) prettyPrint(lvl int) (pprint string) {
	tabs := "    "
	onLVL := strings.Repeat(tabs, lvl)
	if exp.Type == UNIT_EXPR {
		return fmt.Sprintf("%s%s\n", onLVL, exp.Literal)
	}
	pprint = fmt.Sprintf("%s%s\n", onLVL, exp.getTypeName())
	if exp.LExpr != nil {
		pprint += exp.LExpr.prettyPrint(lvl + 1)
	}

	if exp.RExpr != nil {
		pprint += exp.RExpr.prettyPrint(lvl + 1)
	}

	return
}
