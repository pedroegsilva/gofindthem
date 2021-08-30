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
Suspendisse id luctus quam.`,
		`Lorem ipsum Nullam non purus eu leo accumsan cursus a quis erat. 
Etiam dictum enim eu commodo semper.
Mauris feugiat vitae eros et facilisis.
Donec facilisis mattis dignissim.`,
	}

	subEng := &finder.PetarDambovalievEngine{}
	caseSensitive := true
	findthem := finder.NewFinder(subEng, caseSensitive)

	if err := findthem.AddExpression(`"Lorem" and "ipsum"`); err != nil {
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

	findthem2 := finder.NewFinder(subEng, !caseSensitive)

	if err := findthem2.AddExpression(`"lorem ipsum" AND ("dolor" or "accumsan")`); err != nil {
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
