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
	TAG        // "tag"
	FIELD_PATH // "tag:fieldpath"

	// Misc characters
	QUOTATION // "
	OPPAR     // (
	CLPAR     // )

	// Operators
	AND // 'and' or 'AND'
	OR  // 'or' or 'OR'
	NOT // 'not' or 'NOT'
)

// getName returns a readable name for the Token
func (tok Token) getName() string {
	switch tok {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case TAG:
		return "TAG"
	case FIELD_PATH:
		return "FIELD_PATH"
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
	// If we see a '"' consume as a TAG.
	// If we see a '(' or ')' returns OPPAR or CLPAR respectively
	switch {
	case isWhitespace(ch):
		s.unread()
		return s.scanWhitespace()
	case ch == '"':
		s.unread()
		return s.scanTag()
	case ch == ':':
		s.unread()
		return s.scanFieldPath()
	case isLetter(ch):
		s.unread()
		return s.scanOperators()
	case ch == '(':
		return OPPAR, "(", nil
	case ch == ')':
		return CLPAR, ")", nil
	case ch == eof:
		return EOF, "", nil
	}

	return ILLEGAL, "", fmt.Errorf("illegal char was found %c", ch)
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
	default:
		return ILLEGAL, "", fmt.Errorf("failed to scan operator: unexpected operator '%s' found", lit)
	}

	return
}

// scanTag scans the tag and scape needed characters
// If a invalid scape is used an error will be returned and if EOF is found
// before a '"' returns an error as well.
func (s *Scanner) scanTag() (tok Token, lit string, err error) {
	ch := s.read()
	if ch != '"' {
		return ILLEGAL, "", fmt.Errorf("fail to scan tag: expected \" but found %c", ch)
	}
	var buf bytes.Buffer

Loop:
	for {
		ch := s.read()
		switch ch {
		case eof:
			return ILLEGAL, "", fmt.Errorf("fail to scan tag: expected ':' but found EOF")
		case '\\':
			scapedCh := s.read()
			switch scapedCh {
			case '\\', '"', ':':
				_, _ = buf.WriteRune(scapedCh)
			default:
				return ILLEGAL, "", fmt.Errorf("fail to scan tag: invalid escaped char %c", scapedCh)
			}
		case ':':
			s.unread()
			fallthrough
		case '"':
			break Loop
		default:
			_, _ = buf.WriteRune(ch)
		}
	}
	lit = strings.Trim(buf.String(), " ")
	tok = TAG
	return
}

// scanFieldPath scans the tag and scape needed characters
// If a invalid scape is used an error will be returned and if EOF is found
// before a '"' returns an error as well.
func (s *Scanner) scanFieldPath() (tok Token, lit string, err error) {
	ch := s.read()
	if ch != ':' {
		return ILLEGAL, "", fmt.Errorf("fail to scan field: expected ':' but found %c", ch)
	}
	var buf bytes.Buffer
Loop:
	for {
		ch := s.read()
		switch ch {
		case eof:
			return ILLEGAL, "", fmt.Errorf("fail to scan field: expected '\"' but found EOF")
		case '\\':
			scapedCh := s.read()
			switch scapedCh {
			case '\\', '"':
				_, _ = buf.WriteRune(scapedCh)
			default:
				return ILLEGAL, "", fmt.Errorf("fail to scan field: invalid escaped char %c", scapedCh)
			}
		case '"':
			break Loop
		default:
			_, _ = buf.WriteRune(ch)
		}
	}
	lit = strings.Trim(buf.String(), " ")
	tok = FIELD_PATH
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
