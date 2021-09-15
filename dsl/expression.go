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
	INORD_EXPR
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
	case INORD_EXPR:
		return "INORD"
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
	Inord   bool
}

// PatternResult stores if the patter was matched on
// the text and the positions it was found
type PatternResult struct {
	Val            bool
	SortedMatchPos []int
}

// GetTypeName returns the type of the expression with a readable name
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
	eval, _, err := exp.solve(patterResByKeyword, completeMap)
	return eval, err
}

//solve implements Solve
func (exp *Expression) solve(
	patterResByKeyword map[string]PatternResult,
	completeMap bool,
) (bool, []int, error) {
	switch exp.Type {
	case UNIT_EXPR:
		var pos []int
		if resp, ok := patterResByKeyword[exp.Literal]; ok {
			exp.Val = resp.Val
			if exp.Inord {
				pos = resp.SortedMatchPos
			}
		} else {
			if completeMap {
				return false, pos, fmt.Errorf("could not find key %s on map.", exp.Literal)
			} else {
				exp.Val = false
			}
		}
		return exp.Val, pos, nil
	case AND_EXPR:
		var pos []int
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, pos, fmt.Errorf("And statment do not have rigth or left expression: %v", exp)
		}
		lval, lpos, err := exp.LExpr.solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, pos, err
		}
		rval, rpos, err := exp.RExpr.solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, pos, err
		}
		exp.Val = lval && rval
		if exp.Inord && len(lpos) > 0 && len(rpos) > 0 {
			idx := getLowestIdxGTVal(rpos, lpos[0])
			if idx >= 0 {
				pos = rpos[idx:]
			}
		}

		return exp.Val, pos, nil
	case OR_EXPR:
		var pos []int
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, pos, fmt.Errorf("OR statment do not have rigth or left expression: %v", exp)
		}
		lval, lpos, err := exp.LExpr.solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, pos, err
		}
		rval, rpos, err := exp.RExpr.solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, pos, err
		}
		exp.Val = lval || rval
		if exp.Inord {
			pos = mergeArraysSorted(lpos, rpos)
		}

		return exp.Val, pos, nil
	case NOT_EXPR:
		var pos []int
		if exp.RExpr == nil {
			return false, pos, fmt.Errorf("NOT statement do not have expression: %v", exp)
		}
		rval, _, err := exp.RExpr.solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, pos, err
		}
		exp.Val = !rval
		return exp.Val, pos, nil
	case INORD_EXPR:
		var pos []int
		if exp.RExpr == nil {
			return false, pos, fmt.Errorf("INORD statement do not have expression: %v", exp)
		}
		rval, rpos, err := exp.RExpr.solve(patterResByKeyword, completeMap)
		if err != nil {
			return false, pos, err
		}
		exp.Val = rval && len(rpos) > 0
		return exp.Val, pos, nil
	default:
		return false, nil, fmt.Errorf("Unable to process expression type %d", exp.Type)
	}
}

// PrettyPrint returns the expression formated on a tabbed structure
// Eg: for the expression ("a" and "b") or "c"
//    OR
//        AND
//            a
//            b
//        c
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

// positionStack stack of []int used to solve inord operations
// iteratively
type positionStack [][]int

// add implements add for positionStack
func (ps *positionStack) add(pos []int) {
	*ps = append(*ps, pos)
}

// pop implements pop for positionStack
func (ps *positionStack) pop() []int {
	var topPos []int
	n := len(*ps) - 1
	if n >= 0 {
		topPos = (*ps)[n]
		*ps = (*ps)[:n]
	}

	return topPos
}

// Solve solves the expresion iteratively. It has the option to use a complete map of
// PatternResult or a incomplete map. If the complete map option is used the map must have
// all the terms needed to solve de expression or it will return an error.
// If the incomplete map is used, missing keys will be considered as a no match on the
// document.
func (so SolverOrder) Solve(patterResByKeyword map[string]PatternResult, completeMap bool) (bool, error) {
	posStack := &positionStack{}
	for i := len(so) - 1; i >= 0; i-- {
		exp := so[i]
		if exp == nil {
			continue
		}
		switch exp.Type {
		case UNIT_EXPR:
			if resp, ok := patterResByKeyword[exp.Literal]; ok {
				exp.Val = resp.Val
				if exp.Inord {
					posStack.add(resp.SortedMatchPos)
				}
			} else {
				if completeMap {
					return false, fmt.Errorf("could not find key %s on map.", exp.Literal)
				} else {
					exp.Val = false
				}
			}
		case AND_EXPR:
			if exp.LExpr == nil || exp.RExpr == nil {
				return false, fmt.Errorf("And statement do not have right or left expression: %v", exp)
			}
			exp.Val = exp.LExpr.Val && exp.RExpr.Val
			if exp.Inord {
				lpos := posStack.pop()
				rpos := posStack.pop()
				if exp.Inord && len(lpos) > 0 && len(rpos) > 0 {
					idx := getLowestIdxGTVal(rpos, lpos[0])
					if idx >= 0 {
						posStack.add(rpos[idx:])
					}
				}
			}

		case OR_EXPR:
			if exp.LExpr == nil || exp.RExpr == nil {
				return false, fmt.Errorf("OR statement do not have right or left expression: %v", exp)
			}
			exp.Val = exp.LExpr.Val || exp.RExpr.Val
			if exp.Inord {
				lpos := posStack.pop()
				rpos := posStack.pop()
				posStack.add(mergeArraysSorted(lpos, rpos))
			}

		case NOT_EXPR:
			if exp.RExpr == nil {
				return false, fmt.Errorf("NOT statement do not have expression: %v", exp)
			}

			exp.Val = !exp.RExpr.Val

		case INORD_EXPR:
			if exp.RExpr == nil {
				return false, fmt.Errorf("INORD statement do not have expression: %v", exp)
			}
			rpos := posStack.pop()
			if len(*posStack) > 0 {
				return false, fmt.Errorf("INORD did not clear the position stack:")
			}
			exp.Val = exp.RExpr.Val && len(rpos) > 0
		default:
			return false, fmt.Errorf("Unable to process expression type %d", exp.Type)
		}
	}
	return so[0].Val, nil
}

// CreateSolverOrder traverses the expression tree in Preorder and
// stores the expressions on SolverOrder
func (exp *Expression) CreateSolverOrder() SolverOrder {
	solverOrder := new(SolverOrder)
	cpExp := exp
	createSolverOrder(cpExp, solverOrder)
	return *solverOrder
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

// getLowestIdxGTVal uses binary search to find the
// index of the lowest element that is greater than 'value'
func getLowestIdxGTVal(positions []int, value int) int {
	left := 0
	right := len(positions) - 1
	lwGrtI := -1
	for left <= right {
		half := (left + right) >> 1 // divide by 2
		if positions[half] > value {
			lwGrtI = half
			right = half - 1
		} else {
			left = half + 1
		}
	}
	return lwGrtI
}

// mergeArraysSorted merges two sorted arrays into a new sorted array
func mergeArraysSorted(lArr []int, rArr []int) []int {
	leftIdx := 0
	rightIdx := 0
	if len(lArr) == 0 {
		return rArr
	}
	if len(rArr) == 0 {
		return lArr
	}
	lSize := len(lArr)
	rSize := len(rArr)
	sumSize := lSize + rSize
	outArr := make([]int, sumSize)
	count := 0

	for count < sumSize {
		switch {
		case leftIdx == lSize:
			outArr[count] = rArr[rightIdx]
			rightIdx++
		case rightIdx == rSize:
			outArr[count] = lArr[leftIdx]
			leftIdx++
		case lArr[leftIdx] < rArr[rightIdx]:
			outArr[count] = lArr[leftIdx]
			leftIdx++
		default:
			outArr[count] = rArr[rightIdx]
			rightIdx++
		}
		count++
	}
	return outArr
}
