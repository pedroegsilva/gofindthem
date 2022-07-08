package dsl

import (
	"fmt"
	"io"
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
	parCount int
	fields   map[string]struct{}
	tags     map[string]struct{}
}

// NewParser returns a new instance of Parser.
// If case sensitive is not set all terms are changed to lowercase
func NewParser(r io.Reader) *Parser {
	return &Parser{
		s:        NewScanner(r),
		parCount: 0,
		fields:   make(map[string]struct{}),
		tags:     make(map[string]struct{}),
	}
}

// Parse parses the expression and returns the root node
// of the parsed expression.
func (p *Parser) Parse() (expr *Expression, err error) {
	return p.parse()
}

// parse implementation of Parse()
func (p *Parser) parse() (*Expression, error) {
	exp := &Expression{}
	for {
		tok, lit, err := p.scanIgnoreWhitespace()
		if err != nil {
			return exp, err
		}
		switch tok {
		case OPPAR:
			newExp, err := p.handleOpenPar()
			if err != nil {
				return exp, err
			}

			if exp.LExpr == nil {
				exp.LExpr = newExp
			} else {
				exp.RExpr = newExp
			}

		case TAG:
			p.unscan()
			tag, err := p.parseTagInfo()
			if err != nil {
				return exp, err
			}

			keyExp := &Expression{
				Type: UNIT_EXPR,
				Tag:  tag,
			}
			if exp.LExpr == nil {
				exp.LExpr = keyExp
			} else {
				exp.RExpr = keyExp
			}
			p.tags[tag.Name] = struct{}{}
			if tag.FieldPath != "" {
				p.fields[tag.FieldPath] = struct{}{}
			}

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
			nextTok, _, err := p.scanIgnoreWhitespace()
			if err != nil {
				return exp, err
			}

			notExp := &Expression{
				Type: NOT_EXPR,
			}

			switch nextTok {
			case TAG:
				p.unscan()
				tag, err := p.parseTagInfo()
				if err != nil {
					return exp, err
				}
				notExp.RExpr = &Expression{
					Type: UNIT_EXPR,
					Tag:  tag,
				}
				p.tags[tag.Name] = struct{}{}
				if tag.FieldPath != "" {
					p.fields[tag.FieldPath] = struct{}{}
				}

			case OPPAR:
				newExp, err := p.handleOpenPar()
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

	exp = &Expression{
		Type:  expType,
		LExpr: exp,
	}

	nextTok, _, err := p.scanIgnoreWhitespace()
	if err != nil {
		return exp, err
	}

	if nextTok == OPPAR {
		newExp, err := p.handleOpenPar()
		if err != nil {
			return exp, err
		}
		exp.RExpr = newExp
	} else {
		p.unscan()
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

// handleOpenPar gets the expression that is inside the parentheses
func (p *Parser) handleOpenPar() (*Expression, error) {
	parlvl := p.parCount
	p.parCount++
	newExp, err := p.parse()
	if err != nil {
		return newExp, err
	}
	if p.parCount != parlvl {
		return newExp, fmt.Errorf("invalid expression: Unexpected '('")
	}
	return newExp, nil
}

// handleOpenPar gets the expression that is inside the parentheses
func (p *Parser) parseTagInfo() (TagInfo, error) {
	tagInfo := TagInfo{}
	tok, lit, err := p.scanIgnoreWhitespace()
	if err != nil {
		return tagInfo, err
	}

	if tok != TAG {
		return tagInfo, fmt.Errorf("invalid expression: Expecting TAG but found %s", tok.getName())
	}

	if lit == "" {
		return tagInfo, fmt.Errorf("invalid expression: Found empty TAG")
	}

	tagInfo.Name = lit

	nextTok, nextLit, err := p.scanIgnoreWhitespace()
	if err != nil {
		return tagInfo, err
	}

	if nextTok != FIELD_PATH {
		p.unscan()
		return tagInfo, nil
	}

	tagInfo.FieldPath = nextLit
	return tagInfo, nil
}

// GetFields returns the list of unique fields that were found on the expression
func (p *Parser) GetFields() (fields []string) {
	for field := range p.fields {
		fields = append(fields, field)
	}
	return fields
}

// GetTags returns the list of unique tags that were found on the expression
func (p *Parser) GetTags() (tags []string) {
	for tag := range p.tags {
		tags = append(tags, tag)
	}
	return tags
}
