package main

import (
	"fmt"
	"log"

	"github.com/pedroegsilva/gofindthem/finder"
)

func main() {
	texts := []string{
		`lore ipsum`,
		`test`,
		`something`,
	}

	subEng := &finder.PetarDambovalievEngine{}
	rgxEng := &finder.RegexpEngine{}
	caseSensitive := true
	findthem := finder.NewFinder(subEng, rgxEng, caseSensitive)

	if err := findthem.AddExpressionWithTag(`"test"`, "test"); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpressionWithTag(`"something"`, "test2"); err != nil {
		log.Fatal(err)
	}

	for i, text := range texts {
		resp, err := findthem.ProcessText(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("----------Text %d case sensitive-----------\n", i)
		for _, expRes := range resp {
			fmt.Printf("exp: [%s]%s | %v\n", expRes.Tag, expRes.ExpresionStr, expRes.Evaluation)
		}
	}

	subEng2 := &finder.PetarDambovalievEngine{}
	rgxEng2 := &finder.RegexpEngine{}
	findthem2 := finder.NewFinder(subEng2, rgxEng2, !caseSensitive)

	if err := findthem2.AddExpressionWithTag(`"test"`, "test"); err != nil {
		log.Fatal(err)
	}

	if err := findthem2.AddExpressionWithTag(`"something"`, "test2"); err != nil {
		log.Fatal(err)
	}

	for i, text := range texts {
		resp, err := findthem2.ProcessText(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("----------Text %d case insensitive-----------\n", i)
		for _, expRes := range resp {
			fmt.Printf("exp: [%s]%s | %v\n", expRes.Tag, expRes.ExpresionStr, expRes.Evaluation)
		}
	}
}
