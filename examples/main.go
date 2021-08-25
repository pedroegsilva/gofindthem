package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/pedroegsilva/gofindthem/dsl"

	"github.com/pedroegsilva/gofindthem/finder"
)

func main() {
	texts := []string{
		`A domain-specific language (DSL) is a computer language specialized to a particular application domain.
This is in contrast to a general-purpose language (GPL), which is broadly applicable across domains.
There are a wide variety of DSLs, ranging from widely used languages for common domains, such as HTML for web pages, down to languages used by only one or a few pieces of software, such as MUSH soft code.
DSLs can be further subdivided by the kind of language, and include domain-specific markup languages, domain-specific modeling languages (more generally, specification languages), and domain-specific programming languages.
Special-purpose computer languages have always existed in the computer age, but the term "domain-specific language" has become more popular due to the rise of domain-specific modeling.
Simpler DSLs, particularly ones used by a single application, are sometimes informally called mini-languages.
from "https://en.wikipedia.org/wiki/Domain-specific_language"`,
		`The line between general-purpose languages and domain-specific languages is not always sharp, as a language may have specialized features for a particular domain but be applicable more broadly, or conversely may in principle be capable of broad application but in practice used primarily for a specific domain.
For example, Perl was originally developed as a text-processing and glue language, for the same domain as AWK and shell scripts, but was mostly used as a general-purpose programming language later on.
By contrast, PostScript is a Turing complete language, and in principle can be used for any task, but in practice is narrowly used as a page description language.
from "https://en.wikipedia.org/wiki/Domain-specific_language"`,
	}

	findthem := finder.NewFinder(&finder.CloudflareEngine{})

	if err := findthem.AddExpression(`"computer" and "language"`); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpression(`("language" and not "feature")`); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpression(`("HTML" OR "features")`); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpression(`"domain.\nThis"`); err != nil {
		log.Fatal(err)
	}

	if err := findthem.AddExpression(`"ween"`); err != nil {
		log.Fatal(err)
	}

	for _, text := range texts {
		resp, err := findthem.ProcessText(text)
		if err != nil {
			log.Fatal(err)
		}
		for exp, val := range resp {
			fmt.Printf("exp: %s | %v\n", exp, val)
		}
		fmt.Println("------------------------------------")
	}
	exp, err := dsl.NewParser(strings.NewReader(`(("1" and( "1" and "1")) and( "1" and "1")) and (("1" and ("1" and "1")) and ("1" and "1"))`)).Parse()
	//exp, err := dsl.NewParser(strings.NewReader(`"1" and "2" and "3"`)).Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(exp.PrettyFormat())
	exp.CreateSolverOrder()
}
