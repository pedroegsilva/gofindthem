package dsl

import (
	"fmt"
	"io"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
	keywords map[string]bool
	parCount int
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r), keywords: make(map[string]bool), parCount: 0}
}

func (p *Parser) Parse() (expr *Expression, err error) {
	parcount := 0
	return p.parse(nil, &parcount)
}

func (p *Parser) ParseInter() (expr *Expression, err error) {
	return p.parseInter()
}

func (p *Parser) parse(exp *Expression, parCount *int) (expr *Expression, err error) {
	tok, lit, err := p.scanIgnoreWhitespace()
	if err != nil {
		return
	}
	switch tok {
	case OPPAR:
		*parCount++
		expression, err := p.parse(nil, parCount)
		if err != nil {
			return nil, err
		}
		return p.parse(expression, parCount)

	case CLPAR:
		if *parCount > 0 && exp != nil {
			*parCount--
			return exp, nil
		} else {
			return nil, fmt.Errorf("invalid expression: unexpected closing paranteses found. parCount: %d exp: %v", *parCount, exp)
		}

	case KEYWORD:
		if exp != nil {
			return nil, fmt.Errorf("invalid expression: unexpected keyword was received when another keyword was on queue: '%s'", exp.Literal)
		}
		p.keywords[lit] = true
		keyExpr := &Expression{
			Type:      UNIT_EXPR,
			Literal:   lit,
			Evaluated: false,
		}
		return p.parse(keyExpr, parCount)

	case EOF:
		if *parCount == 0 && exp != nil {
			return exp, nil
		} else {
			return nil, fmt.Errorf("invalid expression: unexpected EOF found. parCount: %d exp: %v", *parCount, exp)
		}

	case AND:
		if exp == nil {
			return nil, fmt.Errorf("invalid expression: expected parser to receive expresion on AND but received nil")
		}
		return p.dualExp(
			AND_EXPR,
			exp,
			parCount,
		)

	case OR:
		if exp == nil {
			return nil, fmt.Errorf("invalid expression: expected parser to receive expresion on OR but received nil")
		}
		return p.dualExp(
			OR_EXPR,
			exp,
			parCount,
		)

	case NOT:
		expression, err := p.parse(nil, parCount)
		if err != nil {
			return nil, err
		}
		notExpr := &Expression{
			RExpr:     expression,
			Type:      NOT_EXPR,
			Evaluated: false,
		}
		return p.parse(notExpr, parCount)
	default:
		return nil, fmt.Errorf("invalid expression: Unexpected operator was found (%d = '%s')", tok, lit)
	}

}

func (p *Parser) GetKeywords() map[string]bool {
	return p.keywords
}

func (p *Parser) scan() (tok Token, lit string, err error) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
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

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) dualExp(exprType ExprType, leftExp *Expression, parCount *int) (expr *Expression, err error) {
	expr = &Expression{
		Type:  exprType,
		LExpr: leftExp,
	}
	rExp, err := p.parse(nil, parCount)
	if err != nil {
		return nil, err
	}
	expr.RExpr = rExp
	return
}

func (p *Parser) parseInter() (*Expression, error) {
	exp := &Expression{
		Evaluated: false,
	}
	for {
		tok, lit, err := p.scanIgnoreWhitespace()
		if err != nil {
			return exp, err
		}
		switch tok {
		case OPPAR:
			parlvl := p.parCount
			p.parCount++
			newExp, err := p.parseInter()
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
			keyExp := &Expression{
				Type:      UNIT_EXPR,
				Literal:   lit,
				Evaluated: false,
			}
			if exp.LExpr == nil {
				exp.LExpr = keyExp
			} else {
				exp.RExpr = keyExp
			}
			p.keywords[lit] = true

		case AND:
			exp.Type = AND_EXPR
			if exp.RExpr == nil {
				continue
			}

			tok, _, err = p.scanIgnoreWhitespace()
			if err != nil {
				return exp, err
			}

			exp = &Expression{
				Type:      AND_EXPR,
				LExpr:     exp,
				Evaluated: false,
			}

			if tok == OPPAR {
				parlvl := p.parCount
				p.parCount++
				newExp, err := p.parseInter()
				if err != nil {
					return exp, err
				}
				if p.parCount != parlvl {
					return exp, fmt.Errorf("invalid expression: Unexpected '('")
				}
				exp.RExpr = newExp
			}

		case OR:
			exp.Type = OR_EXPR
			if exp.RExpr == nil {
				continue
			}

			tok, _, err = p.scanIgnoreWhitespace()
			if err != nil {
				return exp, err
			}

			exp = &Expression{
				Type:      OR_EXPR,
				LExpr:     exp,
				Evaluated: false,
			}

			if tok == OPPAR {
				parlvl := p.parCount
				p.parCount++
				newExp, err := p.parseInter()
				if err != nil {
					return exp, err
				}
				if p.parCount != parlvl {
					return exp, fmt.Errorf("invalid expression: Unexpected '('")
				}
				exp.RExpr = newExp
			}

		case NOT:
			if exp.RExpr != nil {
				return exp, fmt.Errorf("invalid expression: Unexpected NOT")
			}

			newExp, err := p.parseInter()
			if err != nil {
				return exp, err
			}

			// in case the expression starts with a not
			exp.RExpr = &Expression{
				RExpr:     newExp,
				Type:      NOT_EXPR,
				Evaluated: false,
			}
		case CLPAR:
			p.parCount--
			fallthrough
		case EOF:
			if p.parCount < 0 {
				return exp, fmt.Errorf("invalid expression: unexpected EOF found. parCount: %d", p.parCount)
			}

			if exp.Type == UNSET_EXPR {
				if exp.RExpr != nil {
					return exp.RExpr, nil
				} else if exp.LExpr != nil {
					return exp.LExpr, nil
				} else {
					return nil, fmt.Errorf("invalid expression: unexpected EOF found")
				}
			}
			return exp, nil

		default:
			return exp, fmt.Errorf("invalid expression: Unexpected operator was found (%d = '%s')", tok, lit)
		}
	}
}
