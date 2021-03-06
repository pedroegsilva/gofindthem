package finder

import (
	"fmt"
	"testing"

	"github.com/pedroegsilva/gofindthem/dsl"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type expectedAddExpression struct {
	exprs    []exprWrapper
	keywords map[string]struct{}
	regexes  map[string]struct{}
	errors   []error
}

func TestAddExpression(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		finder      *Finder
		expressions []string
		expected    expectedAddExpression
		message     string
	}{
		{
			finder:      NewFinder(&EmptyEngine{}, &EmptyRgxEngine{}, false),
			expressions: []string{`"a" and r"B"`, `not "C"`},
			expected: expectedAddExpression{
				exprs: []exprWrapper{
					{
						`"a" and r"B"`,
						&dsl.Expression{
							Type: dsl.AND_EXPR,
							LExpr: &dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "a",
							},
							RExpr: &dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "b",
							},
						},
						"",
					},
					{
						`not "C"`,
						&dsl.Expression{
							Type: dsl.NOT_EXPR,
							RExpr: &dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "c",
							},
						},
						"",
					},
				},
				keywords: map[string]struct{}{
					"a": {},
					"c": {},
				},
				regexes: map[string]struct{}{
					"b": {},
				},
				errors: []error{nil, nil},
			},
			message: "success test case insensitive",
		},
		{
			finder:      NewFinder(&EmptyEngine{}, &EmptyRgxEngine{}, true),
			expressions: []string{`"A"`},
			expected: expectedAddExpression{
				exprs: []exprWrapper{
					{
						`"A"`,
						&dsl.Expression{
							Type:    dsl.UNIT_EXPR,
							Literal: "A",
						},
						"",
					},
				},
				keywords: map[string]struct{}{
					"A": {},
				},
				regexes: map[string]struct{}{},
				errors:  []error{nil},
			},
			message: "success test case sensitive",
		},
		{
			finder:      NewFinder(&EmptyEngine{}, &EmptyRgxEngine{}, true),
			expressions: []string{`"A"`, `invalid`},
			expected: expectedAddExpression{
				exprs: []exprWrapper{
					{
						`"A"`,
						&dsl.Expression{
							Type:    dsl.UNIT_EXPR,
							Literal: "A",
						},
						"",
					},
				},
				keywords: map[string]struct{}{
					"A": {},
				},
				regexes: map[string]struct{}{},
				errors:  []error{nil, fmt.Errorf("failed to scan operator: unexpected operator 'invalid' found")},
			},
			message: "adding invalid expression",
		},
		{
			finder:      NewFinder(&EmptyEngine{}, &EmptyRgxEngine{}, true),
			expressions: []string{``},
			expected: expectedAddExpression{
				exprs:    []exprWrapper{},
				keywords: map[string]struct{}{},
				regexes:  map[string]struct{}{},
				errors:   []error{fmt.Errorf("invalid expression: unexpected EOF found")},
			},
			message: "adding empty expression",
		},
	}

	for _, tc := range tests {
		count := 0
		for _, exp := range tc.expressions {
			err := tc.finder.AddExpression(exp)
			assert.Equal(tc.expected.errors[count], err, tc.message)
			count++
		}
		assert.Equal(tc.expected.exprs, tc.finder.expressions, tc.message)
		assert.Equal(tc.expected.keywords, tc.finder.keywords, tc.message)
		assert.Equal(tc.expected.regexes, tc.finder.regexes, tc.message)
	}
}

type SubstringEngineMock struct {
	mock.Mock
}

// BuildEngine implements BuildEngine with an nop
func (sem *SubstringEngineMock) BuildEngine(keywords map[string]struct{}, caseSensitive bool) (err error) {
	args := sem.Called(keywords)
	return args.Error(0)
}

// BuildEngine implements FindSubstrings with an nop
func (sem *SubstringEngineMock) FindSubstrings(text string) (matches []*Match, err error) {
	args := sem.Called(text)
	return args.Get(0).([]*Match), args.Error(1)
}

type RegexEngineMock struct {
	mock.Mock
}

// BuildEngine implements BuildEngine with an nop
func (rem *RegexEngineMock) BuildEngine(regexes map[string]struct{}, caseSensitive bool) (err error) {
	args := rem.Called(regexes)
	return args.Error(0)
}

// FindRegexes implements FindRegexes with an nop
func (rem *RegexEngineMock) FindRegexes(text string) (matches []*Match, err error) {
	args := rem.Called(text)
	return args.Get(0).([]*Match), args.Error(1)
}

type FindMockRet struct {
	matches []*Match
	err     error
}

func TestProcessText(t *testing.T) {
	assert := assert.New(t)
	text := `text`

	matches1 := []*Match{
		{1, "sharpest"},
	}

	matches2 := []*Match{
		{2, "words"},
	}

	emptyMatches := []*Match{}
	subMock1 := new(SubstringEngineMock)
	subMock2 := new(SubstringEngineMock)
	subMock3 := new(SubstringEngineMock)
	subMock4 := new(SubstringEngineMock)
	subMock5 := new(SubstringEngineMock)
	subMock6 := new(SubstringEngineMock)

	rgxMock1 := new(RegexEngineMock)
	rgxMock2 := new(RegexEngineMock)
	rgxMock3 := new(RegexEngineMock)
	rgxMock4 := new(RegexEngineMock)
	rgxMock5 := new(RegexEngineMock)
	rgxMock6 := new(RegexEngineMock)

	finderBuildErrSub := NewFinder(subMock3, rgxMock3, true)
	finderBuildErrSub.keywords = map[string]struct{}{"1": {}}

	finderBuildErrRgx := NewFinder(subMock4, rgxMock4, true)
	finderBuildErrRgx.regexes = map[string]struct{}{"1": {}}

	finderFindErrSub := NewFinder(subMock5, rgxMock5, true)
	finderFindErrSub.keywords = map[string]struct{}{"1": {}}

	finderFindErrRgx := NewFinder(subMock6, rgxMock6, true)
	finderFindErrRgx.regexes = map[string]struct{}{"1": {}}

	tests := []struct {
		finder                *Finder
		subMock               *SubstringEngineMock
		rgxMock               *RegexEngineMock
		buildSubEngMockRet    error
		buildRgxEngMockRet    error
		buildEngExpecterInput error
		findSubMockRet        FindMockRet
		findRgxMockRet        FindMockRet
		expectedExpRes        []ExpressionResult
		expectedErr           error
		message               string
	}{
		{
			finder: &Finder{
				expressions: []exprWrapper{
					{
						`"sharpest"`,
						&dsl.Expression{
							Type:    dsl.UNIT_EXPR,
							Literal: "sharpest",
						},
						"",
					},
					{
						`r"words"`,

						&dsl.Expression{
							Type:    dsl.UNIT_EXPR,
							Literal: "words",
						},
						"",
					},
				},
				keywords:          map[string]struct{}{"sharpest": {}},
				regexes:           map[string]struct{}{"words": {}},
				subEng:            subMock1,
				rgxEng:            rgxMock1,
				updatedSubMachine: false,
				updatedRgxMachine: false,
			},
			subMock:            subMock1,
			rgxMock:            rgxMock1,
			buildSubEngMockRet: nil,
			buildRgxEngMockRet: nil,
			findSubMockRet:     FindMockRet{matches1, nil},
			findRgxMockRet:     FindMockRet{matches2, nil},
			expectedExpRes: []ExpressionResult{
				{
					ExpresionStr:   `"sharpest"`,
					ExpresionIndex: 0,
				},
				{
					ExpresionStr:   `r"words"`,
					ExpresionIndex: 1,
				},
			},
			expectedErr: nil,
			message:     "success with build",
		},
		{
			finder: &Finder{
				expressions: []exprWrapper{
					{
						`"sharpest"`,
						&dsl.Expression{
							Type:    dsl.UNIT_EXPR,
							Literal: "sharpest",
						},
						"",
					},
					{
						`r"words"`,
						&dsl.Expression{
							Type:    dsl.UNIT_EXPR,
							Literal: "words",
						},
						"",
					},
				},
				keywords:          map[string]struct{}{"sharpest": {}},
				regexes:           map[string]struct{}{"words": {}},
				subEng:            subMock2,
				rgxEng:            rgxMock2,
				updatedSubMachine: true,
				updatedRgxMachine: true,
			},
			subMock:            subMock2,
			rgxMock:            rgxMock2,
			buildSubEngMockRet: nil,
			buildRgxEngMockRet: nil,
			findSubMockRet:     FindMockRet{emptyMatches, nil},
			findRgxMockRet:     FindMockRet{matches2, nil},
			expectedExpRes: []ExpressionResult{
				{
					ExpresionStr:   `r"words"`,
					ExpresionIndex: 1,
				},
			},
			expectedErr: nil,
			message:     "success without build",
		},
		{
			finder:             finderBuildErrSub,
			subMock:            subMock3,
			rgxMock:            rgxMock3,
			buildSubEngMockRet: fmt.Errorf("error building sub engine"),
			buildRgxEngMockRet: nil,
			findSubMockRet:     FindMockRet{emptyMatches, nil},
			findRgxMockRet:     FindMockRet{emptyMatches, nil},
			expectedExpRes:     []ExpressionResult{},
			expectedErr:        fmt.Errorf("error building sub engine"),
			message:            "build engine error substring",
		},
		{
			finder:             finderBuildErrRgx,
			subMock:            subMock4,
			rgxMock:            rgxMock4,
			buildSubEngMockRet: nil,
			buildRgxEngMockRet: fmt.Errorf("error building rgx engine"),
			findSubMockRet:     FindMockRet{emptyMatches, nil},
			findRgxMockRet:     FindMockRet{emptyMatches, nil},
			expectedExpRes:     []ExpressionResult{},
			expectedErr:        fmt.Errorf("error building rgx engine"),
			message:            "build engine error regexes",
		},
		{
			finder:             finderFindErrSub,
			subMock:            subMock5,
			rgxMock:            rgxMock5,
			buildSubEngMockRet: nil,
			buildRgxEngMockRet: nil,
			findSubMockRet:     FindMockRet{emptyMatches, fmt.Errorf("error on sub find")},
			findRgxMockRet:     FindMockRet{emptyMatches, nil},
			expectedExpRes:     []ExpressionResult{},
			expectedErr:        fmt.Errorf("error on sub find"),
			message:            "find substrings error",
		},
		{
			finder:             finderFindErrRgx,
			subMock:            subMock6,
			rgxMock:            rgxMock6,
			buildSubEngMockRet: nil,
			buildRgxEngMockRet: nil,
			findSubMockRet:     FindMockRet{emptyMatches, nil},
			findRgxMockRet:     FindMockRet{emptyMatches, fmt.Errorf("error on rgx find")},
			expectedExpRes:     []ExpressionResult{},
			expectedErr:        fmt.Errorf("error on rgx find"),
			message:            "find regex error",
		},
	}

	for _, tc := range tests {
		if !tc.finder.updatedSubMachine {
			tc.subMock.On(
				"BuildEngine",
				tc.finder.keywords,
			).Return(tc.buildSubEngMockRet)
		}

		if tc.buildSubEngMockRet == nil {
			tc.subMock.On(
				"FindSubstrings",
				text,
			).Return(tc.findSubMockRet.matches, tc.findSubMockRet.err)
		}

		if !tc.finder.updatedRgxMachine {
			tc.rgxMock.On(
				"BuildEngine",
				tc.finder.regexes,
			).Return(tc.buildRgxEngMockRet)
		}

		if tc.buildRgxEngMockRet == nil {
			tc.rgxMock.On(
				"FindRegexes",
				text,
			).Return(tc.findRgxMockRet.matches, tc.findRgxMockRet.err)
		}

		expRes, err := tc.finder.ProcessText(text)
		assert.Equal(tc.expectedErr, err, tc.message)
		if err != nil {
			continue
		}
		assert.Equal(tc.expectedExpRes, expRes, tc.message)
	}
}

func TestAddMatchesToSolverMap(t *testing.T) {
	assert := assert.New(t)
	matches1 := []*Match{
		{1, "sharpest"},
		{2, "sharpest"},
		{3, "sharpest"},
		{7, "words"},
		{9, "Showman"},
		{10, "showman"},
	}

	matches2 := []*Match{
		{1, "sharpest"},
		{2, "sharpest"},
		{3, "sharpest"},
		{7, "words"},
		{9, "Showman"},
		{10, "showman"},
	}

	tests := []struct {
		finder                         *Finder
		matches                        []*Match
		expectedSortedMatchesByKeyword map[string][]int
		message                        string
	}{
		{
			finder:  NewFinder(&EmptyEngine{}, &EmptyRgxEngine{}, false),
			matches: matches1,
			expectedSortedMatchesByKeyword: map[string][]int{
				"sharpest": {1, 2, 3},
				"words":    {7},
				"showman":  {9, 10},
			},
			message: "create with caseinsesitive",
		},
		{
			finder:  NewFinder(&EmptyEngine{}, &EmptyRgxEngine{}, true),
			matches: matches2,
			expectedSortedMatchesByKeyword: map[string][]int{
				"sharpest": {1, 2, 3},
				"words":    {7},
				"showman":  {10},
				"Showman":  {9},
			},
			message: "create with casesesitive",
		},
	}

	for _, tc := range tests {
		sortedMatchesByKeyword := make(map[string][]int)
		tc.finder.addMatchesToSolverMap(tc.matches, sortedMatchesByKeyword)
		assert.Equal(tc.expectedSortedMatchesByKeyword, sortedMatchesByKeyword, tc.message)
	}
}

func TestSolveExpressions(t *testing.T) {
	assert := assert.New(t)
	lexp1 := &dsl.Expression{
		Type:    dsl.UNIT_EXPR,
		Literal: "sharpest",
	}
	rexp1 := &dsl.Expression{
		Type:    dsl.UNIT_EXPR,
		Literal: "words",
	}
	lexp2 := &dsl.Expression{
		Type:    dsl.UNIT_EXPR,
		Literal: "no one",
	}
	rexp2 := &dsl.Expression{
		Type:    dsl.UNIT_EXPR,
		Literal: "Can get in the way",
	}
	finder := &Finder{
		expressions: []exprWrapper{
			{
				`"sharpest" and "words"`,
				&dsl.Expression{
					Type:  dsl.AND_EXPR,
					LExpr: lexp1,
					RExpr: rexp1,
				},
				"",
			},
			{
				`"no one" or "Can get in the way"`,
				&dsl.Expression{
					Type:  dsl.OR_EXPR,
					LExpr: lexp2,
					RExpr: rexp2,
				},
				"",
			},
		},
	}

	tests := []struct {
		finder                 *Finder
		sortedMatchesByKeyword map[string][]int
		expectedExpRes         []ExpressionResult
		expectedErr            error
		message                string
	}{
		{
			finder: finder,
			sortedMatchesByKeyword: map[string][]int{
				"sharpest": {},
				"words":    {},
			},
			expectedExpRes: []ExpressionResult{
				{
					ExpresionIndex: 0,
					ExpresionStr:   `"sharpest" and "words"`,
				},
			},
			expectedErr: nil,
			message:     "first exp true",
		},
		{
			finder: finder,
			sortedMatchesByKeyword: map[string][]int{
				"no one": {},
			},
			expectedExpRes: []ExpressionResult{
				{
					ExpresionIndex: 1,
					ExpresionStr:   `"no one" or "Can get in the way"`,
				},
			},
			expectedErr: nil,
			message:     "second exp true",
		},
		{
			finder: finder,
			sortedMatchesByKeyword: map[string][]int{
				"words": {},
			},
			expectedExpRes: []ExpressionResult{},
			expectedErr:    nil,
			message:        "both false",
		},
		{
			finder: finder,
			sortedMatchesByKeyword: map[string][]int{
				"sharpest":           {},
				"words":              {},
				"Can get in the way": {},
			},
			expectedExpRes: []ExpressionResult{
				{
					ExpresionIndex: 0,
					ExpresionStr:   `"sharpest" and "words"`,
				},
				{
					ExpresionIndex: 1,
					ExpresionStr:   `"no one" or "Can get in the way"`,
				},
			},
			expectedErr: nil,
			message:     "both true",
		},
	}

	for _, tc := range tests {
		expRes, err := tc.finder.solveExpressions(tc.sortedMatchesByKeyword)
		assert.Equal(tc.expectedErr, err, tc.message)
		if err == nil {
			assert.Equal(tc.expectedExpRes, expRes, tc.message)
		}
	}
}
