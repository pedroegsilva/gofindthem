package finder

import (
	"strings"

	"github.com/pedroegsilva/gofindthem/dsl"
)

// Match holds the maching position
// and with substring that was found
type Match struct {
	Position int
	Term     string
}

// exprWrapper store the expression as a string and
// the SolverOrder to later solve the expression.
type exprWrapper struct {
	exprString string
	solverOrd  dsl.SolverOrder
}

// Finder stores the needed information to find the terms and solve the
// feeded expressions
type Finder struct {
	expressions       []exprWrapper
	keywords          map[string]struct{}
	regexes           map[string]struct{}
	subEng            SubstringEngine
	rgxEng            RegexEngine
	updatedSubMachine bool
	updatedRgxMachine bool
	caseSensitive     bool
}

// NewFinder retruns a new instace of Finder
// Setting the engine and if the search will be case sensitive or not.
func NewFinder(subEng SubstringEngine, rgxEng RegexEngine, caseSensitive bool) (finder *Finder) {
	return &Finder{
		expressions:       make([]exprWrapper, 0),
		keywords:          make(map[string]struct{}),
		regexes:           make(map[string]struct{}),
		subEng:            subEng,
		rgxEng:            rgxEng,
		updatedSubMachine: false,
		updatedRgxMachine: false,
		caseSensitive:     caseSensitive,
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
		finder.updatedSubMachine = false
	}

	for rgx := range p.GetRegexes() {
		finder.regexes[rgx] = struct{}{}
		finder.updatedRgxMachine = false
	}

	return nil
}

// ProcessText uses all the unique terms to create the substring engine.
// Searches for matching terms and solves the feeded expressions.
// and returns a map with the expression string as key and its evaluation as value
func (finder *Finder) ProcessText(text string) (evalResp map[string]bool, err error) {
	if !finder.caseSensitive {
		text = strings.ToLower(text)
	}

	solverMap := make(map[string]dsl.PatternResult)

	if len(finder.keywords) > 0 {
		if !finder.updatedSubMachine {
			err = finder.subEng.BuildEngine(finder.keywords, finder.caseSensitive)
			if err != nil {
				return
			}
			finder.updatedSubMachine = true
		}

		keyMaches, err := finder.subEng.FindSubstrings(text)
		if err != nil {
			return nil, err
		}
		finder.addMatchesToSolverMap(keyMaches, solverMap)
	}

	if len(finder.regexes) > 0 {
		if !finder.updatedRgxMachine {
			err = finder.rgxEng.BuildEngine(finder.regexes, finder.caseSensitive)
			if err != nil {
				return
			}
			finder.updatedRgxMachine = true
		}

		rgxMaches, err := finder.rgxEng.FindRegexes(text)
		if err != nil {
			return nil, err
		}
		finder.addMatchesToSolverMap(rgxMaches, solverMap)
	}

	evalResp, err = finder.solveExpressions(solverMap)

	return
}

func (finder *Finder) addMatchesToSolverMap(matches []*Match, solverMap map[string]dsl.PatternResult) {
	for _, match := range matches {
		term := match.Term
		// if the engine returns the substring that was actually matched
		// we turn the key to lower to avoid inconsistency
		if !finder.caseSensitive {
			term = strings.ToLower(term)
		}

		if pattRes, ok := solverMap[term]; ok {
			pattRes.Val = true
			pattRes.SortedMatchPos = append(pattRes.SortedMatchPos, match.Position)
			solverMap[term] = pattRes
		} else {
			solverMap[term] = dsl.PatternResult{
				Val:            true,
				SortedMatchPos: []int{match.Position},
			}
		}
	}
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
	if !finder.updatedSubMachine {
		err = finder.subEng.BuildEngine(finder.keywords, finder.caseSensitive)
		if err != nil {
			return
		}
		finder.updatedSubMachine = true
	}

	if !finder.updatedRgxMachine {
		err = finder.rgxEng.BuildEngine(finder.regexes, finder.caseSensitive)
		if err != nil {
			return
		}
		finder.updatedSubMachine = true
	}
	return
}

// GetKeywords returns all unique terms found on the expressions
func (finder *Finder) GetKeywords() map[string]struct{} {
	return finder.keywords
}
