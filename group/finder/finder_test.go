package finder

import (
	"fmt"
	"testing"

	gofindthem "github.com/pedroegsilva/gofindthem/finder"
	"github.com/pedroegsilva/gofindthem/group/dsl"
	"github.com/stretchr/testify/assert"
)

func TestNewFinder(t *testing.T) {
	gft := gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.EmptyRgxEngine{}, false)
	assert := assert.New(t)
	tests := []struct {
		expectedFinder *GroupFinder
		message        string
	}{
		{
			expectedFinder: &GroupFinder{
				findthem:                    gft,
				expressionWrapperByExprName: make(map[string][]ExpressionWrapper),
				fields:                      make(map[string]struct{}),
				tags:                        make(map[string]struct{}),
			},
			message: "empty gftg",
		},
		{
			expectedFinder: &GroupFinder{
				findthem:                    gft,
				expressionWrapperByExprName: make(map[string][]ExpressionWrapper),
				fields:                      make(map[string]struct{}),
				tags:                        make(map[string]struct{}),
			},
			message: "gftgwith empty gftgs",
		},
	}

	for _, tc := range tests {
		gftg := NewFinder(gft)
		assert.Equal(tc.expectedFinder, gftg, tc.message)
	}
}

func TestNewFinderWithRules(t *testing.T) {
	gft := gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.EmptyRgxEngine{}, false)
	assert := assert.New(t)
	tests := []struct {
		rulesByName map[string][]string
		groupFinder *GroupFinder
		expectedErr error
		message     string
	}{
		{
			rulesByName: map[string][]string{
				"rule1": {
					`"tag1"`,
					`"tag2:field1"`,
				},
				"rule2": {
					`"tag3:field2.field3"`,
					`"tag4"`,
				},
			},
			groupFinder: &GroupFinder{
				findthem: gft,
				expressionWrapperByExprName: map[string][]ExpressionWrapper{
					"rule1": {
						{
							ExpressionString: `"tag1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag1", FieldPath: ""},
							},
						},
						{
							ExpressionString: `"tag2:field1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag2", FieldPath: "field1"},
							},
						},
					},
					"rule2": {
						{
							ExpressionString: `"tag3:field2.field3"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag3", FieldPath: "field2.field3"},
							},
						},
						{
							ExpressionString: `"tag4"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag4", FieldPath: ""},
							},
						},
					},
				},
				fields: map[string]struct{}{
					"field1":        {},
					"field2.field3": {},
				},
				tags: map[string]struct{}{
					"tag1": {},
					"tag2": {},
					"tag3": {},
					"tag4": {},
				},
			},
			expectedErr: nil,
			message:     "new gftgwith valid rules",
		},
		{
			rulesByName: map[string][]string{
				"rule1": {`"tag1`},
			},
			groupFinder: &GroupFinder{
				findthem:                    gft,
				expressionWrapperByExprName: map[string][]ExpressionWrapper{},
				fields:                      map[string]struct{}{},
				tags:                        map[string]struct{}{},
			},
			expectedErr: fmt.Errorf("fail to scan tag: expected ':' but found EOF"),
			message:     "new gftgwith invalid rules",
		},
	}

	for _, tc := range tests {
		gftg, err := NewFinderWithRules(gft, tc.rulesByName)
		assert.Equal(tc.expectedErr, err, tc.message)
		assert.Equal(tc.groupFinder, gftg, tc.message)
	}
}

func TestAddRule(t *testing.T) {
	assert := assert.New(t)
	gft := gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.EmptyRgxEngine{}, false)
	tests := []struct {
		ruleName    string
		expressions []string
		groupFinder *GroupFinder
		expectedErr error
		message     string
	}{
		{
			ruleName: "rule1",
			expressions: []string{
				`"tag1"`,
				`"tag2:field1"`,
			},
			groupFinder: &GroupFinder{
				findthem: gft,
				expressionWrapperByExprName: map[string][]ExpressionWrapper{
					"rule1": {
						{
							ExpressionString: `"tag1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag1", FieldPath: ""},
							},
						},
						{
							ExpressionString: `"tag2:field1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag2", FieldPath: "field1"},
							},
						},
					},
				},
				fields: map[string]struct{}{
					"field1": {},
				},
				tags: map[string]struct{}{
					"tag1": {},
					"tag2": {},
				},
			},
			expectedErr: nil,
			message:     "add valid expressions",
		},
		{
			ruleName: "rule1",
			expressions: []string{
				`"tag1`,
			},
			groupFinder: &GroupFinder{
				findthem:                    gft,
				expressionWrapperByExprName: map[string][]ExpressionWrapper{},
				fields:                      map[string]struct{}{},
				tags:                        map[string]struct{}{},
			},
			expectedErr: fmt.Errorf("fail to scan tag: expected ':' but found EOF"),
			message:     "add invalid expression",
		},
	}

	for _, tc := range tests {
		gftg := NewFinder(gft)
		err := gftg.AddRule(tc.ruleName, tc.expressions)
		assert.Equal(tc.expectedErr, err, tc.message)
		assert.Equal(tc.groupFinder, gftg, tc.message)
	}
}

func TestAddRules(t *testing.T) {
	assert := assert.New(t)
	gft := gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.EmptyRgxEngine{}, false)
	tests := []struct {
		rulesByName map[string][]string
		groupFinder *GroupFinder
		expectedErr error
		message     string
	}{
		{
			rulesByName: map[string][]string{
				"rule1": {
					`"tag1"`,
					`"tag2:field1"`,
				},
				"rule2": {
					`"tag3:field2.field3"`,
					`"tag4"`,
				},
			},
			groupFinder: &GroupFinder{
				findthem: gft,
				expressionWrapperByExprName: map[string][]ExpressionWrapper{
					"rule1": {
						{
							ExpressionString: `"tag1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag1", FieldPath: ""},
							},
						},
						{
							ExpressionString: `"tag2:field1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag2", FieldPath: "field1"},
							},
						},
					},
					"rule2": {
						{
							ExpressionString: `"tag3:field2.field3"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag3", FieldPath: "field2.field3"},
							},
						},
						{
							ExpressionString: `"tag4"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag4", FieldPath: ""},
							},
						},
					},
				},
				fields: map[string]struct{}{
					"field1":        {},
					"field2.field3": {},
				},
				tags: map[string]struct{}{
					"tag1": {},
					"tag2": {},
					"tag3": {},
					"tag4": {},
				},
			},
			message: "add rules with valid rules",
		},
		{
			rulesByName: map[string][]string{
				"rule1": {`"tag1`},
			},
			groupFinder: &GroupFinder{
				findthem:                    gft,
				expressionWrapperByExprName: map[string][]ExpressionWrapper{},
				fields:                      map[string]struct{}{},
				tags:                        map[string]struct{}{},
			},
			expectedErr: fmt.Errorf("fail to scan tag: expected ':' but found EOF"),
			message:     "add rules with invalid rule",
		},
	}

	for _, tc := range tests {
		gftg := NewFinder(gft)
		err := gftg.AddRules(tc.rulesByName)
		assert.Equal(tc.expectedErr, err, tc.message)
		assert.Equal(tc.groupFinder, gftg, tc.message)
	}
}

func TestTagJson(t *testing.T) {
	assert := assert.New(t)
	gft := gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.EmptyRgxEngine{}, false)
	tests := []struct {
		rawJsonStr             string
		matchedExpByFieldByTag map[string]map[string]map[string]struct{}
		expectedErr            error
		message                string
	}{
		{
			rawJsonStr: `{"strField": "some string", "intField": 42, "floatField": 42.42}`,
			matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
				"strTag": {"strField": nil},
			},
			expectedErr: nil,
			message:     "tag json",
		},
	}

	for _, tc := range tests {
		gftg := NewFinder(gft)
		res, err := gftg.TagJson(tc.rawJsonStr, nil, nil)
		assert.Equal(tc.expectedErr, err, tc.message+" expected error")
		for _, resFileds := range res {
			for _, fields := range tc.matchedExpByFieldByTag {
				assert.Equal(resFileds, fields, tc.message+" expected field element")
			}
		}
	}

}

func TestTagObject(t *testing.T) {
	assert := assert.New(t)
	gft := gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.EmptyRgxEngine{}, false)
	gft.AddExpressionWithTag(`"string"`, "strTag")
	tests := []struct {
		object                 interface{}
		matchedExpByFieldByTag map[string]map[string]map[string]struct{}
		expectedErr            error
		message                string
	}{
		{
			object: `some random string`,
			matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
				"strTag": {"": {`"string"`: struct{}{}}},
			},
			expectedErr: nil,
			message:     "tag object raw string",
		},
		{
			object: []string{`some random string`, `some random string`},
			matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
				"strTag": {
					"index(0)": {`"string"`: struct{}{}},
					"index(1)": {`"string"`: struct{}{}},
				},
			},
			expectedErr: nil,
			message:     "tag object array of string",
		},
		{
			object: struct {
				StrField   string
				StrArray   []string
				AnotherObj struct {
					Field1        int
					Field2        float32
					internalField int
				}
				internalStr string
				internalArr []string
				internalObj struct {
					Field3 float64
				}
			}{
				StrField: "some random string",
				StrArray: []string{
					"some random string 1",
					"some random string 2",
				},
				AnotherObj: struct {
					Field1        int
					Field2        float32
					internalField int
				}{
					Field1:        42,
					Field2:        42.42,
					internalField: 0,
				},
				internalStr: "some internal value",
				internalArr: []string{"some internal value 0"},
				internalObj: struct {
					Field3 float64
				}{
					Field3: 0.0,
				},
			},
			matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
				"strTag": {
					"StrField":          {`"string"`: struct{}{}},
					"StrArray.index(0)": {`"string"`: struct{}{}},
					"StrArray.index(1)": {`"string"`: struct{}{}},
				},
			},
			expectedErr: nil,
			message:     "tag object struct with internal fields",
		},
	}

	rules := map[string][]string{
		"test": {`"strTag"`},
	}
	gftg, _ := NewFinderWithRules(gft, rules)
	for _, tc := range tests {
		res, err := gftg.TagObject(tc.object, nil, nil)
		assert.Equal(tc.expectedErr, err, tc.message+" expected error")
		assert.Equal(tc.matchedExpByFieldByTag, res, tc.message+"  expected equal elements found")
	}

}

func TestTagText(t *testing.T) {
	assert := assert.New(t)
	gft := gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.EmptyRgxEngine{}, false)
	gft.AddExpressionWithTag(`"string"`, "strTag")
	tests := []struct {
		text            string
		matchedExpByTag map[string][]string
		expectedErr     error
		message         string
	}{
		{
			text: `some random string`,
			matchedExpByTag: map[string][]string{
				"strTag": {`"string"`},
			},
			expectedErr: nil,
			message:     "tag text",
		},
	}
	for _, tc := range tests {
		gftg := NewFinder(gft)
		res, err := gftg.TagText(tc.text)
		assert.Equal(tc.expectedErr, err, tc.message+" expected error")
		assert.Equal(tc.matchedExpByTag, res, tc.message+" result")
	}
}

func TestEvaluateRules(t *testing.T) {
	assert := assert.New(t)
	gft := gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.EmptyRgxEngine{}, false)
	tests := []struct {
		rulesByName               map[string][]string
		matchedExpByFieldByTag    map[string]map[string]map[string]struct{}
		expectedExpressionsByRule map[string][]string
		expectedErr               error
		message                   string
	}{
		{
			rulesByName: map[string][]string{
				"rule1": {
					`"tag1" and "tag2"`,
				},
				"rule2": {
					`"tag3:field1" and "tag4:field2.innerfield1"`,
				},
				"rule3": {
					`"tag5:field3"`,
				},
				"unmatched rule 1": {
					`"tag4:field1"`,
				},
			},
			matchedExpByFieldByTag: map[string]map[string]map[string]struct{}{
				"tag1": {"randomFiled": nil},
				"tag2": {"randomFiled2": nil},
				"tag3": {"field1": nil},
				"tag4": {"field2.innerfield1": nil},
				"tag5": {"field3.innerfield2": nil},
			},
			expectedExpressionsByRule: map[string][]string{
				"rule1": {
					`"tag1" and "tag2"`,
				},
				"rule2": {
					`"tag3:field1" and "tag4:field2.innerfield1"`,
				},
				"rule3": {
					`"tag5:field3"`,
				},
			},
			expectedErr: nil,
			message:     "evaluate rules",
		},
	}
	for _, tc := range tests {
		gftg, _ := NewFinderWithRules(gft, tc.rulesByName)
		extractorInfoByTaggerName, err := gftg.EvaluateRules(tc.matchedExpByFieldByTag)
		assert.Equal(tc.expectedErr, err, tc.message+" expected error")
		assert.Equal(tc.expectedExpressionsByRule, extractorInfoByTaggerName, tc.message+" result")
	}
}
