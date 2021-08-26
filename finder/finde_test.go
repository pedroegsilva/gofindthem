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
			NewFinder(&EmptyEngine{}, false),
			[]string{`"a" and "B"`, `not "C"`},
			expectedAddExpression{
				[]exprWrapper{
					exprWrapper{
						`"a" and "B"`,
						dsl.SolverOrder{
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
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "a",
							},
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "b",
							},
						},
					},
					exprWrapper{
						`not "C"`,
						dsl.SolverOrder{
							&dsl.Expression{
								Type: dsl.NOT_EXPR,
								RExpr: &dsl.Expression{
									Type:    dsl.UNIT_EXPR,
									Literal: "c",
								},
							},
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "c",
							},
						},
					},
				},
				map[string]struct{}{
					"a": struct{}{},
					"b": struct{}{},
					"c": struct{}{},
				},
				[]error{nil, nil},
			},
			"success test case insensitive",
		},
		{
			NewFinder(&EmptyEngine{}, true),
			[]string{`"A"`},
			expectedAddExpression{
				[]exprWrapper{
					exprWrapper{
						`"A"`,
						dsl.SolverOrder{
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "A",
							},
						},
					},
				},
				map[string]struct{}{
					"A": struct{}{},
				},
				[]error{nil, nil},
			},
			"success test case sensitive",
		},
		{
			NewFinder(&EmptyEngine{}, true),
			[]string{`"A"`, `invalid`},
			expectedAddExpression{
				[]exprWrapper{
					exprWrapper{
						`"A"`,
						dsl.SolverOrder{
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "A",
							},
						},
					},
				},
				map[string]struct{}{
					"A": struct{}{},
				},
				[]error{nil, fmt.Errorf("fail to scan operator: unexpected operator 'invalid' found")},
			},
			"adding invalid expression",
		},
		{
			NewFinder(&EmptyEngine{}, true),
			[]string{``},
			expectedAddExpression{
				[]exprWrapper{},
				map[string]struct{}{},
				[]error{fmt.Errorf("invalid expression: unexpected EOF found")},
			},
			"adding empty expression",
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
func (sem *SubstringEngineMock) FindSubstrings(text string) (matches chan *Match, err error) {
	args := sem.Called(text)
	return args.Get(0).(chan *Match), args.Error(1)
}

type FindSubstrMockRet struct {
	matches chan *Match
	err     error
}

func TestProcessText(t *testing.T) {
	assert := assert.New(t)
	text := `text`

	mock1 := new(SubstringEngineMock)
	chan1 := make(chan *Match, 1)
	chan1 <- &Match{1, "sharpest"}
	close(chan1)

	emptyChan := make(chan *Match, 1)
	close(emptyChan)

	mock2 := new(SubstringEngineMock)
	mock3 := new(SubstringEngineMock)

	tests := []struct {
		finder                *Finder
		mock                  *SubstringEngineMock
		buildEngMockRet       error
		buildEngExpecterInput error
		findSubstrMockRet     FindSubstrMockRet
		expectedEvalResp      map[string]bool
		expectedErr           error
		message               string
	}{
		{
			finder: &Finder{
				expressions: []exprWrapper{
					exprWrapper{
						`"sharpest"`,
						dsl.SolverOrder{
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "sharpest",
							},
						},
					},
				},
				keywords:       map[string]struct{}{"sharpest": struct{}{}},
				subEng:         mock1,
				updatedMachine: false,
			},
			mock:              mock1,
			buildEngMockRet:   nil,
			findSubstrMockRet: FindSubstrMockRet{chan1, nil},
			expectedEvalResp:  map[string]bool{`"sharpest"`: true},
			expectedErr:       nil,
			message:           "success with build",
		},
		{
			finder: &Finder{
				expressions: []exprWrapper{
					exprWrapper{
						`"sharpest"`,
						dsl.SolverOrder{
							&dsl.Expression{
								Type:    dsl.UNIT_EXPR,
								Literal: "sharpest",
							},
						},
					},
				},
				keywords:       map[string]struct{}{"sharpest": struct{}{}},
				subEng:         mock1,
				updatedMachine: true,
			},
			mock:              mock1,
			buildEngMockRet:   nil,
			findSubstrMockRet: FindSubstrMockRet{emptyChan, nil},
			expectedEvalResp:  map[string]bool{`"sharpest"`: false},
			expectedErr:       nil,
			message:           "success without build",
		},
		{
			finder:            NewFinder(mock2, true),
			mock:              mock2,
			buildEngMockRet:   fmt.Errorf("error building engine"),
			findSubstrMockRet: FindSubstrMockRet{emptyChan, nil},
			expectedEvalResp:  map[string]bool{},
			expectedErr:       fmt.Errorf("error building engine"),
			message:           "build engine error",
		},
		{
			finder:            NewFinder(mock3, true),
			mock:              mock3,
			buildEngMockRet:   nil,
			findSubstrMockRet: FindSubstrMockRet{emptyChan, fmt.Errorf("error on find")},
			expectedEvalResp:  map[string]bool{},
			expectedErr:       fmt.Errorf("error on find"),
			message:           "find substrings error",
		},
	}

	for _, tc := range tests {
		if !tc.finder.updatedMachine {
			tc.mock.On(
				"BuildEngine",
				tc.finder.keywords,
			).Return(tc.buildEngMockRet)
		}

		if tc.buildEngMockRet == nil {
			tc.mock.On(
				"FindSubstrings",
				text,
			).Return(tc.findSubstrMockRet.matches, tc.findSubstrMockRet.err)
		}

		eval, err := tc.finder.ProcessText(text)
		assert.Equal(tc.expectedErr, err, tc.message)
		if err != nil {
			continue
		}
		assert.Equal(tc.expectedEvalResp, eval, tc.message)
	}
}

func TestCreateSolverMap(t *testing.T) {
	assert := assert.New(t)

	chan1 := make(chan *Match, 6)
	chan1 <- &Match{1, "sharpest"}
	chan1 <- &Match{2, "sharpest"}
	chan1 <- &Match{3, "sharpest"}
	chan1 <- &Match{7, "words"}
	chan1 <- &Match{9, "Showman"}
	chan1 <- &Match{10, "showman"}
	close(chan1)

	chan2 := make(chan *Match, 6)
	chan2 <- &Match{1, "sharpest"}
	chan2 <- &Match{2, "sharpest"}
	chan2 <- &Match{3, "sharpest"}
	chan2 <- &Match{7, "words"}
	chan2 <- &Match{9, "Showman"}
	chan2 <- &Match{10, "showman"}
	close(chan2)

	tests := []struct {
		finder            *Finder
		inputChan         chan *Match
		expectedSolverMap map[string]dsl.PatternResult
		message           string
	}{
		{
			finder:    NewFinder(&EmptyEngine{}, false),
			inputChan: chan1,
			expectedSolverMap: map[string]dsl.PatternResult{
				"sharpest": dsl.PatternResult{
					Val:            true,
					SortedMatchPos: []int{1, 2, 3},
				},
				"words": dsl.PatternResult{
					Val:            true,
					SortedMatchPos: []int{7},
				},
				"showman": dsl.PatternResult{
					Val:            true,
					SortedMatchPos: []int{9, 10},
				},
			},
			message: "create with caseinsesitive",
		},
		{
			finder:    NewFinder(&EmptyEngine{}, true),
			inputChan: chan2,
			expectedSolverMap: map[string]dsl.PatternResult{
				"sharpest": dsl.PatternResult{
					Val:            true,
					SortedMatchPos: []int{1, 2, 3},
				},
				"words": dsl.PatternResult{
					Val:            true,
					SortedMatchPos: []int{7},
				},
				"showman": dsl.PatternResult{
					Val:            true,
					SortedMatchPos: []int{10},
				},
				"Showman": dsl.PatternResult{
					Val:            true,
					SortedMatchPos: []int{9},
				},
			},
			message: "create with casesesitive",
		},
	}

	for _, tc := range tests {
		solverMap := tc.finder.createSolverMap(tc.inputChan)
		assert.Equal(tc.expectedSolverMap, solverMap, tc.message)
	}
}
