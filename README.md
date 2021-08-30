# gofindthem
Gofindthem is a go library that combines a domain specific language (DSL), that is like a logical expression, with a substring matching engine (implementations of Aho-Corasick at the moment).
Enabling more complex searches with easy to read expressions.
It supports multiple expressions searches, making it an efficient way to "classify" documents according to the expressions that were matched.

This project was conceived in 2019 while I was searching a way to process millions of documents in a more efficient way.
In my researches I found this post [medium flashtext](https://medium.freecodecamp.org/regex-was-taking-5-days-flashtext-does-it-in-15-minutes-55f04411025f) of Vikash Singh.
And discover the Aho-Corasick, which is an awesome algorithm, but it didn't solve my problem completely.
I needed something that could process as fast as the Aho-Corasick implementations could, but also, something that was able to do more complex searches.
The other problem was that the expressions were managed by a team of analysts that were used to use regex to classify those documents.
So I needed a syntax that was easy enough to convince them to change.
The idea was to create a DSL that would have the operators "AND", "OR" and "NOT" that were the same of the logical operations and a new operator "INORD"
that would check if the terms were found on the same order that they were specified.
This would allow searches that used the regex `foo.*bar` to be replaced with `INORD("foo" and "bar")` and the combination of regexes 
`foo.*bar` and `bar.*foo` to be replaced as `"foo" and "bar"`. This is not supposed to be a replacement for regex, but it was enough for most use cases that I had back then.
For those cases that only regex would solve, I added a way to represent a regex with the syntax `R"foo.*bar"`.
Making each kind of terms use its respective engine to find its matches and reducing the need to use regex for everything.

This repository is the golang implementation of this idea.

The scanner and parser form the DSL are heavily influenced by this post [medium parsers-lexers](https://blog.gopheracademy.com/advent-2014/parsers-lexers/) of Ben Johnson,
from which is heavily influenced by the [InfluxQL parser](https://github.com/influxdb/influxdb/tree/master/influxql).

PS: The INORD operator and Regex are not yet supported on this version.

## Usage/Examples

There are 2 libraries on this repository, the DSL and the Finder.

### Finder
First you need to create the Finder. The Finder needs a `SubstringEngine` (interface found at /finder/substringEngine.go) 
and if the search will be case sensitive or not.
there are 3 implementations of `SubstringEngine` that uses the libraries from 
https://github.com/cloudflare/ahocorasick, 
https://github.com/anknown/ahocorasick and 
https://github.com/petar-dambovaliev/aho-corasick. 
But any other library can be used as long as it implements the `SubstringEngine` interface.
```go
    subEng := &finder.PetarDambovalievEngine{}
	caseSensitive := true
	findthem := finder.NewFinder(subEng, caseSensitive)
```

Them you need to add the expressions that need to be solved.
```go
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
```

And finaly you can check which expressions match on each text. 
```go
	for i, text := range texts {
		resp, err := findthem.ProcessText(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("----------Text %d -----------\n", i)
		for exp, val := range resp {
			fmt.Printf("exp: %s | %v\n", exp, val)
		}
	}
}
```

The full example can be found at `/examples/finder/main.go`

### DSL
First you need to create the parser object.
The parser needs a reader with the expression that will be parsed and if it will be case sensitive.
```go
    caseSensitive := false
	p := dsl.NewParser(strings.NewReader(`"lorem ipsum" AND ("dolor" or "accumsan")`), caseSensitive)
```
Them you can parse the expression
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

Format a pretty print to see the Abstract Syntax Tree (AST).
```go
    fmt.Printf("pretty format:\n%s\n", expression.PrettyFormat())
```

There are two ways to solve the expression.

Recursively:
```go 
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
It is faster then the recursive if you need to solve the expression more then 8 times (the gain in performance is around 13% from the benchmark results)

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
This projects uses bazel to build and test the code. 
You can run this project using go as well.

### What is bazel?
"Bazel is an open-source build and test tool similar to Make, Maven, and Gradle. It uses a human-readable, high-level build language. Bazel supports projects in multiple languages and builds outputs for multiple platforms. Bazel supports large codebases across multiple repositories, and large numbers of users." 
\- from https://docs.bazel.build/versions/4.2.0/bazel-overview.html.

To install bazel go to https://docs.bazel.build/versions/main/install.html.

Now with bazel installed.

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

To run all the benchmark use the following.

```
bazel run //benchmarks:benchmarks_test -- -test.bench=. -test.benchmem
```
  