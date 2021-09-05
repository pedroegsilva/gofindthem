package dsl

import (
	"fmt"
	"io"
	"strings"
)

// Parser parser struct that holds needed information to
// parse the expression.
type Parser struct {
	s   *Scanner
	buf struct {
		tok       Token  // last read token
		lit       string // last read literal
		unscanned bool   // if it was unscanned
	}
	keywords     map[string]struct{}
	parCount     int
	casesesitive bool
	inord        bool
}

// NewParser returns a new instance of Parser.
// If casesessitive is not set all terms are changed to lowercase
func NewParser(r io.Reader, casesesitive bool) *Parser {
	return &Parser{
		s:            NewScanner(r),
		keywords:     make(map[string]struct{}),
		parCount:     0,
		casesesitive: casesesitive,
	}
}

// GetKeywords returns the set of UNIT terms (Keywords) that where
// found on the parser
func (p *Parser) GetKeywords() map[string]struct{} {
	return p.keywords
}

// Parse parses the expression and returns the root node
// of the parsed expression.
func (p *Parser) Parse() (expr *Expression, err error) {
	return p.parse()
}

// parse parses the expression and returns the root node
// of the parsed expression.
func (p *Parser) parse() (*Expression, error) {
	exp := &Expression{Inord: p.inord}
	for {
		tok, lit, err := p.scanIgnoreWhitespace()
		if err != nil {
			return exp, err
		}
		switch tok {
		case OPPAR:
			parlvl := p.parCount
			p.parCount++
			newExp, err := p.parse()
			if err != nil {
				return exp, err
			}
			if p.parCount != parlvl {
				return exp, fmt.Errorf("invalid expression: Unexpected '('")
			}

			if exp.LExpr == nil {
				exp.LExpr = newExp
			} else {
				exp.RExpr = newExp
			}

		case KEYWORD:
			if !p.casesesitive {
				lit = strings.ToLower(lit)
			}
			keyExp := &Expression{
				Type:    UNIT_EXPR,
				Literal: lit,
				Inord:   p.inord,
			}
			if exp.LExpr == nil {
				exp.LExpr = keyExp
			} else {
				exp.RExpr = keyExp
			}

			p.keywords[lit] = struct{}{}

		case AND:
			exp, err = p.handleDualOp(exp, AND_EXPR)
			if err != nil {
				return exp, err
			}

		case OR:
			exp, err = p.handleDualOp(exp, OR_EXPR)
			if err != nil {
				return exp, err
			}

		case NOT:
			if p.inord {
				return exp, fmt.Errorf("invalid expression: INORD operator must not contain NOT operator")
			}
			nextTok, nextLit, err := p.scanIgnoreWhitespace()
			if err != nil {
				return exp, err
			}

			notExp := &Expression{
				Type: NOT_EXPR,
			}

			switch nextTok {
			case KEYWORD:
				if !p.casesesitive {
					nextLit = strings.ToLower(nextLit)
				}
				notExp.RExpr = &Expression{
					Type:    UNIT_EXPR,
					Literal: nextLit,
				}
				p.keywords[nextLit] = struct{}{}

			case OPPAR:
				p.unscan()
				newExp, err := p.parse()
				if err != nil {
					return exp, err
				}
				notExp.RExpr = newExp
			default:
				return exp, fmt.Errorf("invalid expression: Unexpected token '%s' after NOT", nextTok.getName())
			}

			if exp.LExpr == nil {
				exp.LExpr = notExp
			} else {
				exp.RExpr = notExp
			}

		case INORD:
			if p.inord {
				return exp, fmt.Errorf("invalid expression: INORD operator must not contain INORD operator")
			}

			nextTok, _, err := p.scanIgnoreWhitespace()
			if err != nil {
				return exp, err
			}

			inordExp := &Expression{
				Type: INORD_EXPR,
			}

			if nextTok != OPPAR {
				return exp, fmt.Errorf("invalid expression: Unexpected token '%s' after INORD", nextTok.getName())
			}
			p.unscan()
			p.inord = true
			newExp, err := p.parse()
			if err != nil {
				return exp, err
			}

			p.inord = false
			inordExp.RExpr = newExp

			if exp.LExpr == nil {
				exp.LExpr = inordExp
			} else {
				exp.RExpr = inordExp
			}

		case CLPAR:
			p.parCount--
			fallthrough
		case EOF:
			if p.parCount < 0 {
				return exp, fmt.Errorf("invalid expression: unexpected EOF found. Extra closing parentheses: %d", p.parCount*-1)
			}

			finalExp := exp
			if exp.Type == UNSET_EXPR {
				if exp.RExpr != nil {
					finalExp = exp.RExpr
				} else if exp.LExpr != nil {
					finalExp = exp.LExpr
				} else {
					return nil, fmt.Errorf("invalid expression: unexpected EOF found")
				}
			}
			switch finalExp.Type {
			case AND_EXPR, OR_EXPR:
				if finalExp.RExpr == nil {
					return nil, fmt.Errorf("invalid expression: incomplete expression %s", finalExp.Type.GetName())
				}
			}

			return finalExp, nil

		default:
			return exp, fmt.Errorf("invalid expression: Unexpected operator was found (%d = '%s')", tok, lit)
		}
	}
}

// handleDualOp adds the needed information to the current expression and returns the next
// expression, that can be the same or another expression.
func (p *Parser) handleDualOp(exp *Expression, expType ExprType) (*Expression, error) {
	if exp.LExpr == nil {
		return exp, fmt.Errorf("invalid expression: no left expression was found for %s", expType.GetName())
	}
	if exp.RExpr == nil {
		exp.Type = expType
		return exp, nil
	}

	nextTok, _, err := p.scanIgnoreWhitespace()
	if err != nil {
		return exp, err
	}

	exp = &Expression{
		Type:  expType,
		LExpr: exp,
		Inord: p.inord,
	}

	p.unscan()
	if nextTok == OPPAR {
		newExp, err := p.parse()
		if err != nil {
			return exp, err
		}
		exp.RExpr = newExp
	}

	return exp, nil
}

// scan scans the next token and stores it on a buffer to
// make unscanning on token possible
func (p *Parser) scan() (tok Token, lit string, err error) {
	// If we have a token on the buffer, then return it.
	if p.buf.unscanned {
		p.buf.unscanned = false
		return p.buf.tok, p.buf.lit, nil
	}

	// Otherwise read the next token from the scanner.
	tok, lit, err = p.s.Scan()
	if err != nil {
		return
	}

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan sets the unscanned flag to assign the scan to
// use the buffered information.
func (p *Parser) unscan() { p.buf.unscanned = true }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string, err error) {
	tok, lit, err = p.scan()
	if err != nil {
		return
	}
	if tok == WS {
		tok, lit, err = p.scan()
	}
	return
}
