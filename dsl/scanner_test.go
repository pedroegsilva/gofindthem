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
			`and   or   not  "keyword 1"  (   )`,
			[]expectedAtScan{
				expectedAtScan{Tok: AND, Lit: "and", Err: nil},
				expectedAtScan{Tok: WS, Lit: "   ", Err: nil},
				expectedAtScan{Tok: OR, Lit: "or", Err: nil},
				expectedAtScan{Tok: WS, Lit: "   ", Err: nil},
				expectedAtScan{Tok: NOT, Lit: "not", Err: nil},
				expectedAtScan{Tok: WS, Lit: "  ", Err: nil},
				expectedAtScan{Tok: KEYWORD, Lit: "keyword 1", Err: nil},
				expectedAtScan{Tok: WS, Lit: "  ", Err: nil},
				expectedAtScan{Tok: OPPAR, Lit: "(", Err: nil},
				expectedAtScan{Tok: WS, Lit: "   ", Err: nil},
				expectedAtScan{Tok: CLPAR, Lit: ")", Err: nil},
				expectedAtScan{Tok: EOF, Lit: "", Err: nil},
			},
			"all tokens",
		},
		{
			`invalidOne`,
			[]expectedAtScan{
				expectedAtScan{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("fail to scan operator: unexpected operator 'invalidOne' found"),
				},
			},
			"invalid operator token",
		},
		{
			`"invalidKeyword `,
			[]expectedAtScan{
				expectedAtScan{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("fail to scan keyword: expected \" but found EOF"),
				},
			},
			"invalid keyword token",
		},
		{
			`"keyword \n \r \t \\ \" "`,
			[]expectedAtScan{
				expectedAtScan{
					Tok: KEYWORD,
					Lit: "keyword \n \r \t \\ \" ",
					Err: nil,
				},
				expectedAtScan{Tok: EOF, Lit: "", Err: nil},
			},
			"valid scaped keyword",
		},
		{
			`"keyword \s"`,
			[]expectedAtScan{
				expectedAtScan{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("fail to scan keyword: invalid escaped char s"),
				},
				expectedAtScan{Tok: EOF, Lit: "", Err: nil},
			},
			"invalid scaped keyword",
		},

		{
			`123`,
			[]expectedAtScan{
				expectedAtScan{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("Illegal char was found 1"),
				},
			},
			"invalid scaped keyword",
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
