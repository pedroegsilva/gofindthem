package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/pedroegsilva/gofindthem/dsl"
)

func main() {
	caseSensitive := false
	p := dsl.NewParser(strings.NewReader(`"lorem ipsum" AND ("dolor" or "accumsan")`), caseSensitive)
	expression, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	keywords := p.GetKeywords()
	fmt.Printf("keywords:\n%v\n", keywords)

	fmt.Printf("pretty format:\n%s\n", expression.PrettyFormat())

	matches := map[string]dsl.PatternResult{
		"lorem ipsum": dsl.PatternResult{
			Val:            true,
			SortedMatchPos: []int{1, 3, 5},
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