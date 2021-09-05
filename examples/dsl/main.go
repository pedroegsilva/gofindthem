package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/pedroegsilva/gofindthem/dsl"
)

func main() {
	caseSensitive := false
	p := dsl.NewParser(strings.NewReader(`INORD("foo" and "bar" and (r"dolor" or "accumsan"))`), caseSensitive)
	expression, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	keywords := p.GetKeywords()
	fmt.Printf("keywords:\n%v\n", keywords)

	regexes := p.GetRegexes()
	fmt.Printf("regexes:\n%v\n", regexes)

	fmt.Printf("pretty format:\n%s\n", expression.PrettyFormat())

	matches := map[string]dsl.PatternResult{
		"foo": dsl.PatternResult{
			Val:            true,
			SortedMatchPos: []int{0, 2, 5},
		},
		"bar": dsl.PatternResult{
			Val:            true,
			SortedMatchPos: []int{3},
		},
		"dolor": dsl.PatternResult{
			Val:            true,
			SortedMatchPos: []int{1, 7},
		},
	}

	responseRecursive, err := expression.Solve(matches, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("recursive eval ", responseRecursive)

	solverArr := expression.CreateSolverOrder()
	responseIter, err := solverArr.Solve(matches, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("iterative eval ", responseIter)

	// should return an error
	_, err = expression.Solve(matches, true)
	if err != nil {
		log.Fatal(err)
	}
}
