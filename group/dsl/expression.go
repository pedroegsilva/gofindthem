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

// TagInfo holds the name of the tag and the field path prefix that it need to be found at.
type TagInfo struct {
	Name      string
	FieldPath string
}

// Expression can be a TagInfo (UNIT) or a function composed by
// one or two other expressions (NOT, AND, OR).
type Expression struct {
	LExpr *Expression
	RExpr *Expression
	Type  ExprType
	Tag   TagInfo
}

// GetTypeName returns the type of the expression with a readable name
func (exp *Expression) GetTypeName() string {
	return exp.Type.GetName()
}

// Solve solves the expresion using the ginven values of fieldPathByTag.
// fieldPathByTag will hold the values of all tags that were found with a
// list of field paths that the tag was found
func (exp *Expression) Solve(
	matchedExpByFieldByTag map[string]map[string]map[string]struct{},
) (bool, error) {
	eval, err := exp.solve(matchedExpByFieldByTag)
	return eval, err
}

//solve implements Solve
func (exp *Expression) solve(matchedExpByFieldByTag map[string]map[string]map[string]struct{}) (bool, error) {
	switch exp.Type {
	case UNIT_EXPR:
		if fieldPaths, ok := matchedExpByFieldByTag[exp.Tag.Name]; ok {
			if exp.Tag.FieldPath == "" {
				return true, nil
			}

			for fieldPath, _ := range fieldPaths {
				if strings.HasPrefix(fieldPath, exp.Tag.FieldPath) {
					return true, nil
				}
			}
		}

		return false, nil

	case AND_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, fmt.Errorf("AND statement do not have right or left expression: %v", exp)
		}
		lval, err := exp.LExpr.solve(matchedExpByFieldByTag)
		if err != nil {
			return false, err
		}
		rval, err := exp.RExpr.solve(matchedExpByFieldByTag)
		if err != nil {
			return false, err
		}

		return lval && rval, nil
	case OR_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, fmt.Errorf("OR statement do not have right or left expression: %v", exp)
		}
		lval, err := exp.LExpr.solve(matchedExpByFieldByTag)
		if err != nil {
			return false, err
		}
		rval, err := exp.RExpr.solve(matchedExpByFieldByTag)
		if err != nil {
			return false, err
		}

		return lval || rval, nil
	case NOT_EXPR:
		if exp.RExpr == nil {
			return false, fmt.Errorf("NOT statement do not have expression: %v", exp)
		}
		rval, err := exp.RExpr.solve(matchedExpByFieldByTag)
		if err != nil {
			return false, err
		}
		return !rval, nil
	default:
		return false, fmt.Errorf("unable to process expression type %d", exp.Type)
	}
}

// PrettyFormat returns the expression formated on a tabbed structure
// Eg: for the expression ("a" and "b") or "c"
//    OR
//        AND
//            a
//            b
//        c
func (exp *Expression) PrettyFormat() string {
	return exp.prettyFormat(0)
}

// prettyFormat implementation of PrettyFormat()
func (exp *Expression) prettyFormat(lvl int) (pprint string) {
	tabs := "    "
	onLVL := strings.Repeat(tabs, lvl)
	if exp.Type == UNIT_EXPR {
		fieldPath := ""
		if exp.Tag.FieldPath != "" {
			fieldPath = fmt.Sprintf("[%s]", exp.Tag.FieldPath)
		}
		return fmt.Sprintf("%s%s%s\n", onLVL, exp.Tag.Name, fieldPath)
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
