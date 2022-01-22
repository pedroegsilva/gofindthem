package dsl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectedAtScan struct {
	Tok Token
	Lit string
	Err error
}

func TestScanner(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		expStr   string
		expected []expectedAtScan
		message  string
	}{
		{
			expStr: `and   or   not  "keyword 1"  (   ) inord r`,
			expected: []expectedAtScan{
				{Tok: AND, Lit: "and", Err: nil},
				{Tok: WS, Lit: "   ", Err: nil},
				{Tok: OR, Lit: "or", Err: nil},
				{Tok: WS, Lit: "   ", Err: nil},
				{Tok: NOT, Lit: "not", Err: nil},
				{Tok: WS, Lit: "  ", Err: nil},
				{Tok: KEYWORD, Lit: "keyword 1", Err: nil},
				{Tok: WS, Lit: "  ", Err: nil},
				{Tok: OPPAR, Lit: "(", Err: nil},
				{Tok: WS, Lit: "   ", Err: nil},
				{Tok: CLPAR, Lit: ")", Err: nil},
				{Tok: WS, Lit: " ", Err: nil},
				{Tok: INORD, Lit: "inord", Err: nil},
				{Tok: WS, Lit: " ", Err: nil},
				{Tok: REGEX, Lit: "r", Err: nil},
				{Tok: EOF, Lit: "", Err: nil},
			},
			message: "all tokens",
		},
		{
			expStr: `invalidOne`,
			expected: []expectedAtScan{
				{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("failed to scan operator: unexpected operator 'invalidOne' found"),
				},
			},
			message: "invalid operator token",
		},
		{
			expStr: `"invalidKeyword `,
			expected: []expectedAtScan{
				{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("fail to scan keyword: expected \" but found EOF"),
				},
			},
			message: "invalid keyword token",
		},
		{
			expStr: `"keyword \n \r \t \\ \" "`,
			expected: []expectedAtScan{
				{
					Tok: KEYWORD,
					Lit: "keyword \n \r \t \\ \" ",
					Err: nil,
				},
				{Tok: EOF, Lit: "", Err: nil},
			},
			message: "valid scaped keyword",
		},
		{
			expStr: `"keyword \s"`,
			expected: []expectedAtScan{
				{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("fail to scan keyword: invalid escaped char s"),
				},
				{Tok: EOF, Lit: "", Err: nil},
			},
			message: "invalid scaped keyword",
		},

		{
			expStr: `123`,
			expected: []expectedAtScan{
				{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("Illegal char was found 1"),
				},
			},
			message: "invalid scaped keyword",
		},
	}

	for _, tc := range tests {
		scanner := NewScanner(strings.NewReader(tc.expStr))
		count := 0
		for {
			tok, lit, err := scanner.Scan()
			expected := tc.expected[count]
			assert.Equal(expected.Err, err, tc.message)
			assert.Equal(expected.Tok, tok, tc.message)
			assert.Equal(expected.Lit, lit, tc.message)
			if err != nil {
				break
			}

			count++
			if tok == EOF {
				break
			}
		}
	}
}
