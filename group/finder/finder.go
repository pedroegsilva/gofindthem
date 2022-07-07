package finder

import (
	"encoding/json"
	"strings"

	"github.com/pedroegsilva/gofindthem/finder"
	"github.com/pedroegsilva/gofindthem/group/dsl"
)

// GroupFinder stores all values needed for the rules
type GroupFinder struct {
	findthem                    *finder.Finder
	expressionWrapperByExprName map[string][]ExpressionWrapper
	fields                      map[string]struct{}
	tags                        map[string]struct{}
}

// ExpressionWrapper store the parsed expression and the raw expressions
type ExpressionWrapper struct {
	ExpressionString string
	Expression       *dsl.Expression
}

// NewFinder returns initialized instancy of GroupFinder.
func NewFinder(findthem *finder.Finder) *GroupFinder {
	return &GroupFinder{
		findthem:                    findthem,
		expressionWrapperByExprName: make(map[string][]ExpressionWrapper),
		fields:                      make(map[string]struct{}),
		tags:                        make(map[string]struct{}),
	}
}

// NewFinderWithRules returns initialized instancy of GroupFinder with the given rules.
func NewFinderWithRules(findthem *finder.Finder, rulesByName map[string][]string) (rules *GroupFinder, err error) {
	rules = NewFinder(findthem)
	err = rules.AddRules(rulesByName)
	return
}

// AddRule adds the given expressions with the rule name to the tagger.
func (rf *GroupFinder) AddRule(ruleName string, expressions []string) error {
	for _, rawExpr := range expressions {
		p := dsl.NewParser(strings.NewReader(rawExpr))
		exp, err := p.Parse()
		if err != nil {
			return err
		}
		expWrapper := ExpressionWrapper{
			ExpressionString: rawExpr,
			Expression:       exp,
		}
		rf.expressionWrapperByExprName[ruleName] = append(rf.expressionWrapperByExprName[ruleName], expWrapper)
		for _, tag := range p.GetTags() {
			rf.tags[tag] = struct{}{}
		}
		for _, field := range p.GetFields() {
			rf.fields[field] = struct{}{}
		}
	}
	return nil
}

// AddRules adds the given expressions with the rule names (key of the map) to the tagger.
func (rf *GroupFinder) AddRules(rulesByName map[string][]string) error {
	for key, exprs := range rulesByName {
		err := rf.AddRule(key, exprs)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFieldNames returns all the unique fields that can be found on all the expressions.
func (rf *GroupFinder) GetFieldNames() (fields []string) {
	for field := range rf.fields {
		fields = append(fields, field)
	}
	return
}

// TagJson tags the fields of a data of type json. Warning at the moment golang json unmarshal,
// when provided a interface{} as the target object, consider all numbers as float64.
// So use the FloatTagger instead of IntTagger for tagging numbers.
func (rf *GroupFinder) TagJson(
	data string,
	includePaths []string,
	excludePaths []string,
) (matchedExpByFieldByTag map[string]map[string]map[string]struct{}, err error) {
	var genericObj interface{}
	err = json.Unmarshal([]byte(data), &genericObj)
	if err != nil {
		return
	}

	return rf.TagObject(genericObj, includePaths, excludePaths)
}

// TagObject tags the fields of a data of type interface.
func (rf *GroupFinder) TagObject(
	data interface{},
	includePaths []string,
	excludePaths []string,
) (matchedExpByFieldByTag map[string]map[string]map[string]struct{}, err error) {
	matchedExpByFieldByTag = make(map[string]map[string]map[string]struct{})
	err = rf.getRulesInfo(data, "", includePaths, excludePaths, matchedExpByFieldByTag)
	return
}

// TagText tags the fields of a string.
func (rf *GroupFinder) TagText(
	data string,
) (matchedExpByTag map[string][]string, err error) {
	matchedExpByTag = make(map[string][]string)
	matchedExpByFieldByTag, err := rf.TagObject(data, nil, nil)
	if err != nil {
		return
	}

	for tag, fields := range matchedExpByFieldByTag {
		for exp := range fields[""] {
			matchedExpByTag[tag] = append(matchedExpByTag[tag], exp)
		}
	}
	return
}

// EvaluateRules evaluate all rules with the given fields by tag.
func (rf *GroupFinder) EvaluateRules(
	matchedExpByFieldByTag map[string]map[string]map[string]struct{},
) (expressionsByRule map[string][]string, err error) {
	expressionsByRule = make(map[string][]string)
	for name, exprWrappers := range rf.expressionWrapperByExprName {
		for _, ew := range exprWrappers {
			eval, err := ew.Expression.Solve(matchedExpByFieldByTag)
			if err != nil {
				return nil, err
			}
			if eval {
				expressionsByRule[name] = append(expressionsByRule[name], ew.ExpressionString)
			}
		}
	}
	return
}

// ProcessJson extract all tags and evaluate all rules for the given data of type json.
// includePaths can be used to specify what fields will be used on the tagging.
// if empty array or nil is passed to 'includePaths' it will consider all fields as taggable.
// excludePaths can be used to specify what fields will be skipped on the tagging, if
// there is a conflict on a specific field the excludePath has precedence over the include paths.
// Empty array or nil can be used to not exclude any fields
func (rf *GroupFinder) ProcessJson(
	rawJson string,
	includePaths []string,
	excludePaths []string,
) (expressionsByRule map[string][]string, err error) {
	matchedExpByFieldByTag, err := rf.TagJson(rawJson, includePaths, excludePaths)
	if err != nil {
		return nil, err
	}

	return rf.EvaluateRules(matchedExpByFieldByTag)
}

// ProcessObject extract all tags and evaluate all rules for the given data of type interface.
// includePaths can be used to specify what fields will be used on the tagging.
// if empty array or nil is passed to 'includePaths' it will consider all fields as taggable.
// excludePaths can be used to specify what fields will be skipped on the tagging, if
// there is a conflict on a specific field the excludePath has precedence over the include paths.
// Empty array or nil can be used to not exclude any fields
func (rf *GroupFinder) ProcessObject(
	obj interface{},
	includePaths []string,
	excludePaths []string,
) (expressionsByRule map[string][]string, err error) {
	matchedExpByFieldByTag, err := rf.TagObject(obj, includePaths, excludePaths)
	if err != nil {
		return nil, err
	}

	return rf.EvaluateRules(matchedExpByFieldByTag)
}

// ProcessText extract all tags and evaluate all rules for the given string.
func (rf *GroupFinder) ProcessText(
	data string,
) (expressionsByRule map[string][]string, err error) {
	matchedExpByFieldByTag, err := rf.TagObject(data, nil, nil)
	if err != nil {
		return nil, err
	}
	return rf.EvaluateRules(matchedExpByFieldByTag)
}
