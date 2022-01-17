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
	Inord   bool
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
func (exp *Expression) Solve(sortedMatchesByKeyword map[string][]int) (bool, error) {
	eval, _, err := exp.solve(sortedMatchesByKeyword)
	return eval, err
}

//solve implements Solve
func (exp *Expression) solve(sortedMatchesByKeyword map[string][]int) (bool, []int, error) {
	switch exp.Type {
	case UNIT_EXPR:
		if sortedMatches, ok := sortedMatchesByKeyword[exp.Literal]; ok {
			return true, sortedMatches, nil
		}
		return false, nil, nil

	case AND_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, nil, fmt.Errorf("AND statment do not have rigth or left expression: %v", exp)
		}
		lval, lpos, err := exp.LExpr.solve(sortedMatchesByKeyword)
		if err != nil {
			return false, nil, err
		}
		rval, rpos, err := exp.RExpr.solve(sortedMatchesByKeyword)
		if err != nil {
			return false, nil, err
		}

		var pos []int
		if exp.Inord && len(lpos) > 0 && len(rpos) > 0 {
			idx := getLowestIdxGTVal(rpos, lpos[0])
			if idx >= 0 {
				pos = rpos[idx:]
			}
		}

		return lval && rval, pos, nil

	case OR_EXPR:

		if exp.LExpr == nil || exp.RExpr == nil {
			return false, nil, fmt.Errorf("OR statment do not have rigth or left expression: %v", exp)
		}
		lval, lpos, err := exp.LExpr.solve(sortedMatchesByKeyword)
		if err != nil {
			return false, nil, err
		}
		rval, rpos, err := exp.RExpr.solve(sortedMatchesByKeyword)
		if err != nil {
			return false, nil, err
		}

		var pos []int
		if exp.Inord {
			pos = mergeArraysSorted(lpos, rpos)
		}

		return lval || rval, pos, nil

	case NOT_EXPR:
		if exp.RExpr == nil {
			return false, nil, fmt.Errorf("NOT statement do not have expression: %v", exp)
		}
		rval, _, err := exp.RExpr.solve(sortedMatchesByKeyword)
		if err != nil {
			return false, nil, err
		}

		return !rval, nil, nil

	case INORD_EXPR:
		if exp.RExpr == nil {
			return false, nil, fmt.Errorf("INORD statement do not have expression: %v", exp)
		}
		rval, rpos, err := exp.RExpr.solve(sortedMatchesByKeyword)
		if err != nil {
			return false, nil, err
		}
		return rval && len(rpos) > 0, nil, nil

	default:
		return false, nil, fmt.Errorf("unable to process expression type %d", exp.Type)
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

type valAndPos struct {
	Val bool
	Pos []int
}

// Solve solves the expresion iteratively. It has the option to use a complete map of
// PatternResult or a incomplete map. If the complete map option is used the map must have
// all the terms needed to solve de expression or it will return an error.
// If the incomplete map is used, missing keys will be considered as a no match on the
// document.
func (so SolverOrder) Solve(sortedMatchesByKeyword map[string][]int) (bool, error) {
	values := make(map[*Expression]valAndPos)
	for i := len(so) - 1; i >= 0; i-- {
		exp := so[i]
		if exp == nil {
			return false, fmt.Errorf("malformed solver order - solver order should not have nil values")
		}
		switch exp.Type {
		case UNIT_EXPR:
			if sortedMatches, ok := sortedMatchesByKeyword[exp.Literal]; ok {
				values[exp] = valAndPos{
					Val: true,
					Pos: sortedMatches,
				}
			} else {
				values[exp] = valAndPos{Val: false}
			}

		case AND_EXPR:
			l, lOk := values[exp.LExpr]
			r, rOk := values[exp.RExpr]
			if !lOk || !rOk {
				return false, fmt.Errorf("AND statement do not have right or left expression: %v", exp)
			}
			vap := valAndPos{Val: l.Val && r.Val}

			if exp.Inord {
				lpos := l.Pos
				rpos := r.Pos
				if exp.Inord && len(lpos) > 0 && len(rpos) > 0 {
					idx := getLowestIdxGTVal(rpos, lpos[0])
					if idx >= 0 {
						vap.Pos = rpos[idx:]
					}
				}
			}
			values[exp] = vap

		case OR_EXPR:
			l, lOk := values[exp.LExpr]
			r, rOk := values[exp.RExpr]
			if !lOk || !rOk {
				return false, fmt.Errorf("OR statement do not have right or left expression: %v", exp)
			}
			vap := valAndPos{Val: l.Val || r.Val}
			if exp.Inord {
				vap.Pos = mergeArraysSorted(l.Pos, r.Pos)
			}
			values[exp] = vap

		case NOT_EXPR:
			r, rOk := values[exp.RExpr]
			if !rOk {
				return false, fmt.Errorf("NOT statement do not have expression: %v", exp)
			}
			values[exp] = valAndPos{Val: !r.Val}

		case INORD_EXPR:
			r, rOk := values[exp.RExpr]
			if !rOk {
				return false, fmt.Errorf("INORD statement do not have expression: %v", exp)
			}
			values[exp] = valAndPos{Val: r.Val && len(r.Pos) > 0}
		default:
			return false, fmt.Errorf("unable to process expression type %d", exp.Type)
		}
	}
	return values[so[0]].Val, nil
}

// CreateSolverOrder traverses the expression tree in Preorder and
// stores the expressions on SolverOrder
func (exp *Expression) CreateSolverOrder() *SolverOrder {
	solverOrder := new(SolverOrder)
	cpExp := exp
	createSolverOrder(cpExp, solverOrder)
	return solverOrder
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
