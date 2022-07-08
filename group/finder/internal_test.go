package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isValidateFieldPath(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		fieldPath    string
		includePaths []string
		excludePaths []string
	}
	tests := []struct {
		args     args
		expected bool
		message  string
	}{
		{
			args: args{
				fieldPath:    "field1.inner1.inner2",
				includePaths: []string{},
				excludePaths: []string{},
			},
			expected: true,
			message:  "empty includes and excludes",
		},
		{
			args: args{
				fieldPath:    "field1.inner1.inner2",
				includePaths: []string{"field1.inner1.inner2"},
				excludePaths: []string{},
			},
			expected: true,
			message:  "include exact match",
		},
		{
			args: args{
				fieldPath:    "field1.inner1.inner2",
				includePaths: []string{"field2"},
				excludePaths: []string{},
			},
			expected: false,
			message:  "include no match ",
		},
		{
			args: args{
				fieldPath:    "field1.inner1.inner2",
				includePaths: []string{},
				excludePaths: []string{"field1.inner1.inner2"},
			},
			expected: false,
			message:  "exclude exact match",
		},
		{
			args: args{
				fieldPath:    "field1.inner1.inner2",
				includePaths: []string{"field1.inner1.inner2"},
				excludePaths: []string{"field1.inner1.inner2"},
			},
			expected: false,
			message:  "exclude and include exact match",
		},
		{
			args: args{
				fieldPath:    "field1.inner1.inner2",
				includePaths: []string{"field1.inner1"},
				excludePaths: []string{},
			},
			expected: true,
			message:  "include partial match",
		},
		{
			args: args{
				fieldPath:    "field1.inner1.inner2",
				includePaths: []string{},
				excludePaths: []string{"field1.inner1"},
			},
			expected: false,
			message:  "exclude partial match",
		},
		{
			args: args{
				fieldPath:    "field1.inner1.inner2",
				includePaths: []string{"field1.inner1"},
				excludePaths: []string{"field1.inner1"},
			},
			expected: false,
			message:  "exclude and include partial match",
		},
	}
	for _, tc := range tests {
		res := isValidateFieldPath(tc.args.fieldPath, tc.args.includePaths, tc.args.excludePaths)
		assert.Equal(tc.expected, res, tc.message)
	}
}
