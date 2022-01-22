package main

import (
	"fmt"
	"log"

	"github.com/pedroegsilva/gofindthem/finder"
)

func main() {
	texts := []string{
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Curabitur porta lobortis nulla volutpat sagittis. 
Nulla ac sapien sodales, pulvinar elit ut, lobortis purus.
Suspendisse id luctus quam FOO.`,
		`FOO Lorem ipsum Nullam non purus eu leo accumsan cursus a quis erat. 
Etiam dictum enim eu commodo semper.
Mauris feugiat vitae eros et facilisis.
Donec facilisis mattis dignissim.`,
	}

	subEng := &finder.CloudflareForkEngine{}
	rgxEng := &finder.RegexpEngine{}
	caseSensitive := true
	findthem := finder.NewFinder(subEng, rgxEng, caseSensitive)

	if err := findthem.AddExpressionWithTag(`r"Lorem" and "ipsum"`, "test"); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpressionWithTag(`("Nullam" and not "volutpat")`, "test2"); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpressionWithTag(`"lorem ipsum" AND ("dolor" or "accumsan")`, "test"); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpression(`"purus.\nSuspendisse"`); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpression(`inord("Lorem" and "FOO")`); err != nil {
		log.Fatal(err)
	}

	for i, text := range texts {
		resp, err := findthem.ProcessText(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("----------Text %d case sensitive-----------\n", i)
		for _, expRes := range resp {
			fmt.Printf("exp %d: [%s]%s\n", expRes.ExpresionIndex, expRes.Tag, expRes.ExpresionStr)
		}
	}

	subEng2 := &finder.CloudflareForkEngine{}
	rgxEng2 := &finder.RegexpEngine{}
	findthem2 := finder.NewFinder(subEng2, rgxEng2, !caseSensitive)

	if err := findthem2.AddExpression(`"Lorem Ipsum" AND ("doLor" or "accumsan")`); err != nil {
		log.Fatal(err)
	}

	if err := findthem2.AddExpression(`R"Lorem.*Ipsum" AND (r"doLor" or r"accumsan")`); err != nil {
		log.Fatal(err)
	}

	for i, text := range texts {
		resp, err := findthem2.ProcessText(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("----------Text %d case insensitive-----------\n", i)
		for _, expRes := range resp {
			fmt.Printf("exp %d: [%s]%s\n", expRes.ExpresionIndex, expRes.Tag, expRes.ExpresionStr)
		}
	}
}
