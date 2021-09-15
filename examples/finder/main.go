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

	subEng := &finder.PetarDambovalievEngine{}
	rgxEng := &finder.RegexpEngine{}
	caseSensitive := true
	findthem := finder.NewFinder(subEng, rgxEng, caseSensitive)

	if err := findthem.AddExpression(`r"Lorem" and "ipsum"`); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpression(`("Nullam" and not "volutpat")`); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpression(`"lorem ipsum" AND ("dolor" or "accumsan")`); err != nil {
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
		for exp, val := range resp {
			fmt.Printf("exp: %s | %v\n", exp, val)
		}
	}

	subEng2 := &finder.PetarDambovalievEngine{}
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
		for exp, val := range resp {
			fmt.Printf("exp: %s | %v\n", exp, val)
		}
	}
}
