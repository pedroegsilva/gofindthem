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

	matches := map[string][]int{
		"foo":   {0, 2, 5},
		"bar":   {3},
		"dolor": {1, 7},
	}

	responseRecursive, err := expression.Solve(matches)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("recursive eval ", responseRecursive)

	solverArr := expression.CreateSolverOrder()
	responseIter, err := solverArr.Solve(matches)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("iterative eval ", responseIter)

	// should return an error
	_, err = expression.Solve(matches)
	if err != nil {
		log.Fatal(err)
	}
}
