# gofindthem
Gofindthem is a go library that implements a domain-specific language(DSL), which is similar to a logical expression. It also uses a substring matching engine (implementations of Aho-Corasick) and a regex matching engine to enable more complex searches with easy-to-read expressions. It supports multiple expressions searches, making it an efficient way to "classify" documents according to the expressions that were matched.

## Brief History
This project was conceived in 2019 while I was searching for a way to match multiple rules on millions of documents in a more efficient way than matching multiple regex.

In my researches, I found this [article](https://medium.freecodecamp.org/regex-was-taking-5-days-flashtext-does-it-in-15-minutes-55f04411025f) written by Vikash Singh. He described a problem similar to mine, where the use of a set of regex needed to be applied on many documents, and this was quite slow. In his case, he needed to replace the found terms for another, and I needed to use the found terms to check if a set of rules were met.

In the Vikash Singh article I learned about the Aho-Corasick algorithm, which solved the "matching of multiple terms efficiently", but it didn't solve my problem completely.
The problem was that the rules that were applied needed to be managed by a team of analysts, so I couldn't just hard code the set of rules since I would have to change them constantly.
The other problem was that the team used to use regex to classify those documents. Because of that, I needed a syntax that was easy enough to convince them to change.

The idea was to create a DSL that would have the operators "AND", "OR" and "NOT", which are the same as the logical operations, and a new operator "INORD", that would check if the terms were found in the same order that they were specified.
This would allow searches that used the regex `foo.*bar` to be replaced with `INORD("foo" and "bar")` and the combination of regex 
`foo.*bar` and `bar.*foo` to be replaced as `"foo" and "bar"`. This does not replace all regex, but it would cover most use cases that I had back then.
For those cases that only regex would solve, I added a way to represent a regex with the syntax `R"foo.*bar"`, making each kind of term use its respective engine to find its matches and reducing the need to use regex for everything.

This repository is the golang implementation of this idea.

The scanner and parser from the DSL are heavily influenced by this [article](https://blog.gopheracademy.com/advent-2014/parsers-lexers/) written by Ben Johnson,
from which is heavily influenced by the [InfluxQL parser](https://github.com/influxdb/influxdb/tree/master/influxql).

## Usage/Examples

There are 2 libraries on this repository, the DSL and the Finder.

### Finder
The finder is used to manage multiple expressions. It will use the DSL to extract the terms and regex from each expression and use them to process the text with the appropriate engine.

You will need to create the Finder object. The Finder needs a `SubstringEngine` (interface can be found at `/finder/substringEngine.go`) and a `RegexEngine` (interface can be found at `/finder/regexEngine.go`)
and if the search will be case sensitive or not.
There are 3 "implementations" of `SubstringEngine` that uses the libraries from 
https://github.com/cloudflare/ahocorasick, 
https://github.com/anknown/ahocorasick and 
https://github.com/petar-dambovaliev/aho-corasick and a regexp implementation for the `RegexEngine`. 
But any other library can be used as long as it "implements" the `SubstringEngine` or `RegexEngine` interface.
```go
    subEng := &finder.CloudflareForkEngine{}
    rgxEng := &finder.RegexpEngine{}
    caseSensitive := true
    findthem := finder.NewFinder(subEng, rgxEng, caseSensitive)
```

Then you need to add the expressions that need to be solved.
```go
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
```

And finally you can check which expressions match on each text. 
```go
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
```

The full example can be found at `/examples/finder/main.go`

### DSL
#### Definition
The DSL uses 5 operators (AND, OR, NOT, R, INORD), terms (defined by "") and parentheses to form expressions. A valid expression can be:

- A single term. Eg: `"some term"`
- The result of an operation. `R"term 1"` 
- An expression enclosed by parentheses `("term 1" AND "term 2")`

Each operator functions as the following:

- **AND** - Uses the expression before and after it to solve them as a logical `AND` operator. 
    > (valid expression) AND (valid expression) eg: `"term 1" AND "term 2"` 

- **OR** - Uses the expression before and after it to solve them as a logical `OR` operator.
    > \<valid expression\> OR \<valid expression\> eg: `"term 1" OR "term 2"` 

- **NOT** - Uses the expression after it to solve them as a logical `NOT` operator.
    > NOT \<valid expression\> eg: `NOT "term 1"`

- **R** - Defines the next term as a regex. 
    > **WARNING** - To scape regex operator characters please use `\\` instead of `\` since the single reverse bar is used to scape the char `"` on the DSL.
    >
    > R \<term\> eg: `R"term1\\."`

- **INORD** - Needs to be followed by an expression enclosed in parentheses. This operator will check if there is a set of terms on the document that satisfy the same order of the enclosed terms. It will still solve the logical expressions but it will return false if the terms were not found in the defined order. _Note that the OR operator enclosed on the `INORD` operator will consider that at least one of the terms must be found and in order. For example INORD("a" and "b") or INORD("a" and "c") is equivalent to INORD("a" and ("b" or "c"))_
    > **WARNING** - The NOT and INORD operators are not permitted on the expression 
    > that is enclosed by the INORD operator. Because the expression `INORD(NOT "a" and "b")` doesn't make sense and another INORD would be redundant.
    >
    > INORD(\<valid expression\>) eg: `INORD("term 1" AND ("term 2" or "term 3"))` 

#### Package
To use this package as a stand-alone, you will need to create the parser object. The parser needs a reader with the expression that will be parsed and if it will be case-sensitive.
```go
    caseSensitive := false
    p := dsl.NewParser(strings.NewReader(`INORD("foo" and "bar" and ("dolor" or "accumsan"))`), caseSensitive)
```

Then you can parse the expression

```go 
    expression, err := p.Parse()
    if err != nil {
        log.Fatal(err)
    }
```

Once parsed you can extract which terms there were on the expression.
```go
    keywords := p.GetKeywords()
    fmt.Printf("keywords:\n%v\n", keywords)
```

And which Regexes there were on the expression.
```go
    regexes := p.GetRegexes()
    fmt.Printf("regexes:\n%v\n", regexes)
```

You can also get a pretty format to see the created Abstract Syntax Tree (AST).
```go
    fmt.Printf("pretty format:\n%s\n", expression.PrettyFormat())
```

There are two ways to solve the expression.

Recursively:
```go 
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
```
Iteratively:
```go
    solverArr := expression.CreateSolverOrder()
    responseIter, err := solverArr.Solve(matches, false)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("iterative eval ", responseIter)
```
The Iterative solution needs to create an array with the order in which the expressions need to be solved.
It is faster than the recursive if you need to solve the expression more than 8 times (the gain in performance is around 13% from the benchmark results)

The solvers also need to know if the map of matches is complete or not. If it is complete it will have the term as a key even if it was a no match.
The incomplete option will assume that if a key is not present the term was not found.
If an incomplete map is provided and the key is not found an error will be returned.

```go
    // should return an error
    _, err = expression.Solve(matches, true)
    if err != nil {

        log.Fatal(err)
    }
}

```
The complete example can be found at `/examples/dsl/main.go`
## Run Locally
This project uses Bazel to build and test the code. 
You can run this project using go as well.

### What is Bazel?
"Bazel is an open-source build and test tool similar to Make, Maven, and Gradle. It uses a human-readable, high-level build language. Bazel supports projects in multiple languages and builds outputs for multiple platforms. Bazel supports large codebases across multiple repositories and large numbers of users." 
\- from https://docs.bazel.build/versions/4.2.0/bazel-overview.html.

To install Bazel go to https://docs.bazel.build/versions/main/install.html.

Now with Bazel installed.

Clone the project.
```bash
  git clone https://github.com/pedroegsilva/gofindthem.git
```

Go to the project directory

```bash
  cd gofindthem
```

You can run the examples with the following commands.

```
  bazel run //examples/finder:finder
  bazel run //examples/dsl:dsl
```
To run all the tests use the following.

```
bazel test //...
```

To run all the benchmarks use the following.

```
bazel run //benchmarks:benchmarks_test -- -test.bench=. -test.benchmem
```
  