package main

import (
	"fmt"

	"github.com/pedroegsilva/gofindthem/finder"
	gfinder "github.com/pedroegsilva/gofindthem/group/finder"
)

func main() {
	gofindthemRules := map[string][]string{
		"tag1": {
			`"string1"`,
			`"string2"`,
		},
		"tag2": {
			`"string3"`,
			`"string4"`,
		},
		"tag3": {
			`"string5"`,
			`"string6"`,
		},
		"tag4": {
			`"string7"`,
			`"string8"`,
		},
	}

	rules := map[string][]string{
		"rule1": {`"tag1" or "tag2"`},
		"rule2": {`"tag3:Field3.SomeField1" or "tag4"`},
		"rule3": {`"tag3:Field3" or "tag4"`},
	}

	gft, err := finder.NewFinderWithExpressions(
		&finder.CloudflareForkEngine{},
		&finder.RegexpEngine{},
		false,
		gofindthemRules,
	)

	if err != nil {
		panic(err)
	}

	gftg, err := gfinder.NewFinderWithRules(gft, rules)
	if err != nil {
		panic(err)
	}

	someObject := struct {
		Field1 string
		Field2 int
		Field3 struct {
			SomeField1 string
			SomeField2 []string
		}
	}{
		Field1: "some pretty text with string1",
		Field2: 42,
		Field3: struct {
			SomeField1 string
			SomeField2 []string
		}{
			SomeField1: "some pretty text with string5",
			SomeField2: []string{"some pretty text with string5", "some pretty text with string2", "some pretty text with string3"},
		},
	}

	matchedExpByFieldByTag, err := gftg.TagObject(someObject, gftg.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}

	for tag, expressionsByField := range matchedExpByFieldByTag {
		fmt.Println("Tag: ", tag)
		for field, exprs := range expressionsByField {
			fmt.Println("    Field: ", field)
			for exp := range exprs {
				fmt.Println("        Expressions: ", exp)
			}
		}
	}

	res, err := gftg.ProcessObject(someObject, gftg.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessObject: ", res)

	fmt.Println("-----------------------------")
	arr := []struct {
		FieldN string
		FieldX string
	}{
		{FieldN: "some pretty text with string5"},
		{FieldN: "some pretty text with string2"},
		{FieldN: "some pretty text with string3"},
	}

	matchedExpByFieldByTag2, err := gftg.TagObject(arr, nil, nil)
	if err != nil {
		panic(err)
	}
	for tag, expressionsByField := range matchedExpByFieldByTag2 {
		fmt.Println("Tag: ", tag)
		for field, exprs := range expressionsByField {
			fmt.Println("    Field: ", field)
			for exp := range exprs {
				fmt.Println("        Expressions: ", exp)
			}
		}
	}

	res2, err := gftg.ProcessObject(arr, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessObject2: ", res2)

	fmt.Println("-----------------------------")
	rawJson := `
	{
		"Field1": "some pretty text with string1",
		"Field2": 42,
		"Field3":
		{
			"SomeField1": "some pretty text with string5",
			"SomeField2":
			[
				"some pretty text with string5",
				"some pretty text with string2",
				"some pretty text with string3"
			]
		}
	}
	`

	matchedExpByFieldByTag3, err := gftg.TagJson(rawJson, gftg.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	for tag, expressionsByField := range matchedExpByFieldByTag3 {
		fmt.Println("Tag: ", tag)
		for field, exprs := range expressionsByField {
			fmt.Println("    Field: ", field)
			for exp := range exprs {
				fmt.Println("        Expressions: ", exp)
			}
		}
	}
	res3, err := gftg.ProcessJson(rawJson, gftg.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessJson: ", res3)
}
