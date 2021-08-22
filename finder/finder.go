package finder

import (
	"fmt"
	"strings"

	"github.com/pedroegsilva/gofindthem/dsl"
)

type ExprWrapper struct {
	Expr       *dsl.Expression
	ExprString string
}

type Finder struct {
	Expressions    []ExprWrapper
	Keywords       map[string]bool
	SubEng         SubstringEngine
	updatedMachine bool
}

func NewFinder(subEng SubstringEngine) (finder *Finder) {
	return &Finder{
		Keywords:       make(map[string]bool),
		SubEng:         subEng,
		updatedMachine: false,
	}
}

func (finder *Finder) AddExpression(expression string) error {
	p := dsl.NewParser(strings.NewReader(expression))
	exp, err := p.Parse()
	if err != nil {
		return err
	}

	finder.Expressions = append(finder.Expressions, ExprWrapper{exp, expression})
	for key := range p.GetKeywords() {
		finder.Keywords[key] = true
	}
	finder.updatedMachine = false

	return nil
}

func (finder *Finder) AddExpressionInter(expression string) error {
	p := dsl.NewParser(strings.NewReader(expression))
	exp, err := p.ParseInter()
	if err != nil {
		return err
	}

	fmt.Println(exp.PrettyPrint())

	finder.Expressions = append(finder.Expressions, ExprWrapper{exp, expression})
	for key := range p.GetKeywords() {
		finder.Keywords[key] = true
	}
	finder.updatedMachine = false

	return nil
}

func (finder *Finder) ProcessText(text string, completeMap bool) (evalResp map[string]bool, err error) {

	if !finder.updatedMachine {
		err = finder.SubEng.BuildEngine(finder.Keywords)
		if err != nil {
			return
		}
		finder.updatedMachine = true
	}

	matches, err := finder.SubEng.ProcessText(text)
	if err != nil {
		return
	}

	solverMap := finder.createSolverMap(matches, completeMap)
	evalResp, err = finder.solveExpressions(solverMap, completeMap)

	return
}

func (finder *Finder) createSolverMap(matches chan *Match, completeMap bool) (solverMap map[string]dsl.PatternResult) {
	solverMap = make(map[string]dsl.PatternResult)

	if completeMap {
		for key := range finder.Keywords {
			solverMap[key] = dsl.PatternResult{
				Val: false,
			}
		}
	}

	for match := range matches {
		if pattRes, ok := solverMap[match.Substring]; ok {
			pattRes.Val = true
			pattRes.SortedMatchPos = append(pattRes.SortedMatchPos, match.Position)
			solverMap[match.Substring] = pattRes
		} else {
			solverMap[match.Substring] = dsl.PatternResult{
				Val:            true,
				SortedMatchPos: []int{match.Position},
			}
		}
	}
	return
}

func (finder *Finder) solveExpressions(solverMap map[string]dsl.PatternResult, completeMap bool) (evalResp map[string]bool, err error) {
	evalResp = make(map[string]bool)
	for _, exp := range finder.Expressions {
		res, err := exp.Expr.Solve(solverMap, completeMap, true)
		if err != nil {
			return nil, err
		}
		evalResp[exp.ExprString] = res
	}
	return
}

func (finder *Finder) ForceBuild() (err error) {
	if !finder.updatedMachine {
		err = finder.SubEng.BuildEngine(finder.Keywords)
		if err != nil {
			return
		}
		finder.updatedMachine = true
	}
	return
}
