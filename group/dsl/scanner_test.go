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
			expStr: `and   or   not  "tag1:somepath.to.field"  (   )`,
			expected: []expectedAtScan{
				{Tok: AND, Lit: "and", Err: nil},
				{Tok: WS, Lit: "   ", Err: nil},
				{Tok: OR, Lit: "or", Err: nil},
				{Tok: WS, Lit: "   ", Err: nil},
				{Tok: NOT, Lit: "not", Err: nil},
				{Tok: WS, Lit: "  ", Err: nil},
				{Tok: TAG, Lit: "tag1", Err: nil},
				{Tok: FIELD_PATH, Lit: "somepath.to.field", Err: nil},
				{Tok: WS, Lit: "  ", Err: nil},
				{Tok: OPPAR, Lit: "(", Err: nil},
				{Tok: WS, Lit: "   ", Err: nil},
				{Tok: CLPAR, Lit: ")", Err: nil},
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
			expStr: `"invalidTag `,
			expected: []expectedAtScan{
				{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("fail to scan tag: expected ':' but found EOF"),
				},
			},
			message: "invalid tag token",
		},
		{
			expStr: `"tag \: \" \\:path\"1\""`,
			expected: []expectedAtScan{
				{
					Tok: TAG,
					Lit: "tag : \" \\",
					Err: nil,
				},
				{
					Tok: FIELD_PATH,
					Lit: "path\"1\"",
					Err: nil,
				},
				{Tok: EOF, Lit: "", Err: nil},
			},
			message: "valid scaped tag",
		},
		{
			expStr: `"tag \s"`,
			expected: []expectedAtScan{
				{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("fail to scan tag: invalid escaped char s"),
				},
				{Tok: EOF, Lit: "", Err: nil},
			},
			message: "invalid scaped tag",
		},

		{
			expStr: `123`,
			expected: []expectedAtScan{
				{
					Tok: ILLEGAL,
					Lit: "",
					Err: fmt.Errorf("illegal char was found 1"),
				},
			},
			message: "invalid operator",
		},
	}

	for _, tc := range tests {
		scanner := NewScanner(strings.NewReader(tc.expStr))
		count := 0
		for {
			tok, lit, err := scanner.Scan()
			if count >= len(tc.expected) {
				t.Fail()
				break
			}
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
