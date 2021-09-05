package dsl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	KEYWORD // "keyword"

	// Misc characters
	QUOTATION // "
	OPPAR     // (
	CLPAR     // )

	// Operators
	AND   // 'and' or 'AND'
	OR    // 'or' or 'OR'
	NOT   // 'not' or 'NOT'
	INORD // 'inord' or 'INORD'
)

// getName retuns a readable name for the Token
func (tok Token) getName() string {
	switch tok {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case KEYWORD:
		return "KEYWORD"
	case QUOTATION:
		return "QUOTATION"
	case OPPAR:
		return "OPPAR"
	case CLPAR:
		return "CLPAR"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case INORD:
		return "INORD"
	default:
		return "UNEXPECTED"
	}
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string, err error) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an operator.
	// If we see a '"' consume as a KEYWORD.
	// If we see a '(' or ')' returns OPPAR or CLPAR respectively
	switch {
	case isWhitespace(ch):
		s.unread()
		return s.scanWhitespace()
	case isLetter(ch):
		s.unread()
		return s.scanOperators()
	case ch == '"':
		s.unread()
		return s.scanKeyword()
	case ch == '(':
		return OPPAR, "(", nil
	case ch == ')':
		return CLPAR, ")", nil
	case ch == eof:
		return EOF, "", nil
	}

	return ILLEGAL, "", fmt.Errorf("Illegal char was found %c", ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string, err error) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String(), nil
}

// scanOperators consumes the current rune and all contiguous operator runes.
func (s *Scanner) scanOperators() (tok Token, lit string, err error) {
	// Create a buffer and read the current character into it.
	ch := s.read()
	if !isLetter(ch) {
		return ILLEGAL, "", fmt.Errorf("fail to scan operator: expected letter but found %c", ch)
	}
	var buf bytes.Buffer

	buf.WriteRune(ch)

	// Read every subsequent operator character into the buffer.
	// Non-operator characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a operator then return that operator.
	// Otherwise return an error.
	lit = buf.String()
	switch strings.ToUpper(lit) {
	case "AND":
		tok = AND
	case "OR":
		tok = OR
	case "NOT":
		tok = NOT
	case "INORD":
		tok = INORD
	default:
		return ILLEGAL, "", fmt.Errorf("failed to scan operator: unexpected operator '%s' found", lit)
	}

	return
}

// scanKeyword scans the keyword and scape needed characters
// If a invalid scape is used an error will be returned and if EOF is found
// before a '"' returns an error as well.
func (s *Scanner) scanKeyword() (tok Token, lit string, err error) {
	ch := s.read()
	if ch != '"' {
		return ILLEGAL, "", fmt.Errorf("fail to scan keyword: expected \" but found %c", ch)
	}
	var buf bytes.Buffer

	endloop := false
	for {
		ch := s.read()
		switch ch {
		case eof:
			return ILLEGAL, "", fmt.Errorf("fail to scan keyword: expected \" but found EOF")
		case '\\':
			scapedCh := s.read()
			switch scapedCh {
			case '\\':
				_, _ = buf.WriteRune(scapedCh)
			case 'n':
				_, _ = buf.WriteRune('\n')
			case 'r':
				_, _ = buf.WriteRune('\r')
			case 't':
				_, _ = buf.WriteRune('\t')
			case '"':
				_, _ = buf.WriteRune(scapedCh)
			default:
				return ILLEGAL, "", fmt.Errorf("fail to scan keyword: invalid escaped char %c", scapedCh)
			}
		case '"':
			endloop = true
		default:
			_, _ = buf.WriteRune(ch)
		}
		if endloop {
			break
		}
	}
	lit = buf.String()
	tok = KEYWORD
	return
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
