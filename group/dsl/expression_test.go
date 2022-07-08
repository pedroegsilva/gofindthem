package dsl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolver(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range solverTestCases {
		exp, err := NewParser(strings.NewReader(tc.expStr)).Parse()
		assert.Nil(err, tc.message)
		respInt, err := exp.Solve(tc.matchedExpByFieldByTag)
		assert.Nil(err, tc.message)
		assert.Equal(tc.expectedResp, respInt, tc.message)
	}
}

var solverTestCases = []struct {
	expStr                 string
	matchedExpByFieldByTag map[string]map[string]map[string]struct{}
	expectedResp           bool
	message                string
}{
	{
		expStr: `"tag1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
		},
		expectedResp: true,
		message:      "single tag true",
	},
	{
		expStr:                 `"tag1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{},
		expectedResp:           false,
		message:                "single tag false",
	},

	// and tests
	{
		expStr: `"tag1" and "tag2" and "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag2": nil,
			"tag3": nil,
		},
		expectedResp: true,
		message:      "and multi tags true",
	},
	{
		expStr: `"tag1" and "tag2" and "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag2": nil,
			"tag3": nil,
		},
		expectedResp: false,
		message:      "and multi tags false 1",
	},
	{
		expStr: `"tag1" and "tag2" and "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag3": nil,
		},
		expectedResp: false,
		message:      "and multi tags false 2",
	},
	{
		expStr: `"tag1" and "tag2" and "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag2": nil,
		},
		expectedResp: false,
		message:      "and multi tags false 3",
	},

	// or tests
	{
		expStr:                 `"tag1" or "tag2" or "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{},
		expectedResp:           false,
		message:                "or multi tag false",
	},
	{
		expStr: `"tag1" or "tag2" or "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
		},
		expectedResp: true,
		message:      "or multi tag true 1",
	},
	{
		expStr: `"tag1" or "tag2" or "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag2": nil,
		},
		expectedResp: true,
		message:      "or multi tag true 2",
	},
	{
		expStr: `"tag1" or "tag2" or "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag3": nil,
		},
		expectedResp: true,
		message:      "or multi tag true 3",
	},

	// not tests
	{
		expStr:                 `not "tag1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{},
		expectedResp:           true,
		message:                "not true",
	},
	{
		expStr: `not "tag1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
		},
		expectedResp: false,
		message:      "not false",
	},
	{
		expStr: `not "tag1" or not "tag2"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag2": nil,
		},
		expectedResp: false,
		message:      "not multi false",
	},
	{
		expStr: `not ("tag1" or "tag2") or "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag3": nil,
		},
		expectedResp: true,
		message:      "not multi true",
	},
	{
		expStr: `"tag1" and not "tag2" or "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag3": nil,
		},
		expectedResp: true,
		message:      "not multi true 1",
	},
	{
		expStr: ` not "tag2" and "tag1" or "tag3"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag3": nil,
		},
		expectedResp: true,
		message:      "not multi true 2",
	},
	{
		expStr: `"tag1" and "tag3" or not "tag2"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag3": nil,
		},
		expectedResp: true,
		message:      "not multi true 3",
	},
	// parentheses tests
	{
		expStr: `not ("tag1" and "tag2") and ("tag1" or "tag2")`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
			"tag2": nil,
		},
		expectedResp: false,
		message:      "parentheses xor 1",
	},
	{
		expStr:                 `("tag1" or "tag2") and not ("tag1" and "tag2")`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{},
		expectedResp:           false,
		message:                "parentheses xor 2",
	},
	{
		expStr: `not ("tag1" and "tag2") and ("tag1" or "tag2")`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag2": nil,
		},
		expectedResp: true,
		message:      "parentheses xor 3",
	},
	{
		expStr: `(("tag1" and "tag2" and "tag3") or ("tag4" and not "tag5")) and ("tag6" or "tag7") and "tag8"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag4": nil,
			"tag6": nil,
			"tag8": nil,
		},
		expectedResp: true,
		message:      "parentheses 1",
	},
	{
		expStr: `(("tag1" and "tag2" and "tag3") or ("tag4" and not "tag5")) and ("tag6" or "tag7") and "tag8"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag4": nil,
			"tag6": nil,
		},
		expectedResp: false,
		message:      "parentheses 2",
	},
	// field tests
	{
		expStr: `"tag1:field1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": {"field1": nil},
		},
		expectedResp: true,
		message:      "single tag with field true 1",
	},
	{
		expStr: `"tag1:field1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": {"field3": nil, "field2": nil, "field1": nil},
		},
		expectedResp: true,
		message:      "single tag with field true 2",
	},
	{
		expStr: `"tag1:field1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": nil,
		},
		expectedResp: false,
		message:      "single tag with field false 1",
	},
	{
		expStr: `"tag1:field1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": {"field2": nil, "field3": nil},
		},
		expectedResp: false,
		message:      "single tag with field false 2",
	},
	{
		expStr: `"tag1:field1"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": {"field1.field2.field3.index(2)": nil},
		},
		expectedResp: true,
		message:      "single tag with field partial field true 1",
	},
	{
		expStr: `"tag1:field"`,
		matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
			"tag1": {"field2": nil},
		},
		expectedResp: true,
		message:      "single tag with field partial field true 2",
	},
}
