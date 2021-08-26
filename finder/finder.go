package finder

import (
	"strings"

	"github.com/pedroegsilva/gofindthem/dsl"
)

// exprWrapper store the expression as a string and
// the SolverOrder to later solve the expression.
type exprWrapper struct {
	exprString string
	solverOrd  dsl.SolverOrder
}

// Finder stores the needed information to find the terms and solve the
// feeded expressions
type Finder struct {
	expressions    []exprWrapper
	keywords       map[string]struct{}
	subEng         SubstringEngine
	updatedMachine bool
	caseSensitive  bool
}

// NewFinder retruns a new instace of Finder
// Setting the engine and if the search will be case sensitive or not.
func NewFinder(subEng SubstringEngine, caseSensitive bool) (finder *Finder) {
	return &Finder{
		expressions:    make([]exprWrapper, 0),
		keywords:       make(map[string]struct{}),
		subEng:         subEng,
		updatedMachine: false,
		caseSensitive:  caseSensitive,
	}
}

// AddExpression adds the expression to the finder. It also collect
// and store the terms that are going to be used by the substring engine
// If the expression is malformed returns an erro.
func (finder *Finder) AddExpression(expression string) error {
	p := dsl.NewParser(strings.NewReader(expression), finder.caseSensitive)
	exp, err := p.Parse()
	if err != nil {
		return err
	}

	finder.expressions = append(finder.expressions, exprWrapper{expression, exp.CreateSolverOrder()})
	for key := range p.GetKeywords() {
		finder.keywords[key] = struct{}{}
	}
	finder.updatedMachine = false

	return nil
}

// ProcessText uses all the unique terms to create the substring engine.
// Searches for matching terms and solves the feeded expressions.
// and returns a map with the expression string as key and its evaluation as value
func (finder *Finder) ProcessText(text string) (evalResp map[string]bool, err error) {
	if !finder.updatedMachine {
		err = finder.subEng.BuildEngine(finder.keywords, finder.caseSensitive)
		if err != nil {
			return
		}
		finder.updatedMachine = true
	}
	if !finder.caseSensitive {
		text = strings.ToLower(text)
	}
	matches, err := finder.subEng.FindSubstrings(text)
	if err != nil {
		return
	}

	solverMap := finder.createSolverMap(matches)
	evalResp, err = finder.solveExpressions(solverMap)

	return
}

// createSolverMap creates a map with the matching terms positions and value
func (finder *Finder) createSolverMap(matches chan *Match) (solverMap map[string]dsl.PatternResult) {
	solverMap = make(map[string]dsl.PatternResult)
	for match := range matches {
		key := match.Substring
		// if the engine returns the substring that was actualy matched
		// we turn the key to lower to avoid inconsistency
		if !finder.caseSensitive {
			key = strings.ToLower(key)
		}

		if pattRes, ok := solverMap[key]; ok {
			pattRes.Val = true
			pattRes.SortedMatchPos = append(pattRes.SortedMatchPos, match.Position)
			solverMap[key] = pattRes
		} else {
			solverMap[key] = dsl.PatternResult{
				Val:            true,
				SortedMatchPos: []int{match.Position},
			}
		}
	}
	return
}

// solveExpressions solves all feeded expressions using the values of the solverMap
func (finder *Finder) solveExpressions(solverMap map[string]dsl.PatternResult) (evalResp map[string]bool, err error) {
	evalResp = make(map[string]bool)
	for _, exp := range finder.expressions {
		res, err := exp.solverOrd.Solve(solverMap, false)
		if err != nil {
			return nil, err
		}
		evalResp[exp.exprString] = res
	}
	return
}

// ForceBuild forces the substring engine to be built if needed
func (finder *Finder) ForceBuild() (err error) {
	if !finder.updatedMachine {
		err = finder.subEng.BuildEngine(finder.keywords, finder.caseSensitive)
		if err != nil {
			return
		}
		finder.updatedMachine = true
	}
	return
}

// GetKeywords returns all unique terms found on the expressions
func (finder *Finder) GetKeywords() map[string]struct{} {
	return finder.keywords
}
