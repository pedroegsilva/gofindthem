package main

import (
	"fmt"
	"strings"

	"github.com/pedroegsilva/gofindthem/dsl"
)

func main() {
	solverMap := map[string]dsl.PatternResult{
		"something": dsl.PatternResult{
			Val: true,
		},
		"test": dsl.PatternResult{
			Val: true,
		},
	}
	exp, err := dsl.NewParser(strings.NewReader(`not (("something" and "test") or NOT ("another test" OR "some"))`)).Parse()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println("parse :", exp.PrettyPrint())
	expint, err := dsl.NewParser(strings.NewReader(`not (("something" and "test") or NOT ("another test" OR "some"))`)).ParseInter()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println("parse inter: ", expint.PrettyPrint())
	resp, err := expint.Solve(solverMap, false, false)
	if err != nil {
		fmt.Printf("solve err: %v\n", err)
		return
	}
	fmt.Println("resp: ", resp)
}
