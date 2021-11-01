package finderbenchmarks

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/pedroegsilva/gofindthem/finder"
)

func init() {
	rand.Seed(1629074756677820700)
	wordsPath, err := filepath.Abs(EN_WORDS_FILE)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(wordsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	count := 0
	for scanner.Scan() {
		if count < 466550 {
			words[count] = scanner.Text()
		} else {
			break
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

var (
	words [466550]string
)

const (
	EN_WORDS_FILE = "./files/words.txt"
)

func createText(numTerm int) string {
	text := ""
	dictLen := len(words)
	for i := 1; i <= numTerm; i++ {
		text += words[rand.Intn(dictLen)] + " "
	}
	return text
}

func BenchmarkIncresingTextSize(b *testing.B) {
	for _, step := range []int{10, 100, 1000} {
		BMIncresingTextSize(b, step)
	}
}

func BenchmarkIncresingTermsCount(b *testing.B) {
	for _, step := range []int{10, 100, 1000} {
		BMIncresingTermsCount(b, step)
	}
}

func BMIncresingTextSize(b *testing.B, step int) {
	for i := step; i <= step*10; i = i + step {
		text := "foo " + createText(i) + " bar"
		expressions1 := []string{
			`"foo" and "bar"`,
		}
		BMDslSearch(fmt.Sprintf("TextDSL_%d", i), expressions1, &finder.PetarDambovalievEngine{}, text, b)

		expressions2 := []string{
			`r"foo.*bar" and r"bar.*foo"`,
		}
		BMDslSearch(fmt.Sprintf("TextDSLRegex_%d", i), expressions2, &finder.PetarDambovalievEngine{}, text, b)

		expressions3 := []string{
			`INORD("foo" and "bar") and INORD("bar" and "foo")`,
		}
		BMDslSearch(fmt.Sprintf("TextDSLInord_%d", i), expressions3, &finder.PetarDambovalievEngine{}, text, b)

		regexes := []string{
			"foo.*bar",
			"bar.*foo",
		}
		BMUseCasesRegexOnly(fmt.Sprintf("TextRegex_%d", i), text, regexes, b)
	}
}

func BMIncresingTermsCount(b *testing.B, step int) {
	for i := step; i <= step*10; i = i + step {
		text := createText(10)
		terms := getTerms(i)
		regexes, andExps, inordExps, rgxExps := createExpressions(terms)

		BMDslSearch(fmt.Sprintf("TermsDSL_%d", i), andExps, &finder.PetarDambovalievEngine{}, text, b)

		BMDslSearch(fmt.Sprintf("TermsDSLRegex_%d", i), rgxExps, &finder.PetarDambovalievEngine{}, text, b)

		BMDslSearch(fmt.Sprintf("TermsDSLInord_%d", i), inordExps, &finder.PetarDambovalievEngine{}, text, b)

		BMUseCasesRegexOnly(fmt.Sprintf("TermsRegex_%d", i), text, regexes, b)
	}
}

func BMUseCasesRegexOnly(name string, text string, regexestr []string, b *testing.B) {
	var regex []*regexp.Regexp
	for _, rgxstr := range regexestr {
		regex = append(regex, regexp.MustCompile(rgxstr))
	}

	b.ResetTimer()
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, r := range regex {
				r.FindAllStringIndex(text, -1)
			}
		}
	})
}

func BMDslSearch(name string, exps []string, subEng finder.SubstringEngine, text string, b *testing.B) {
	findthem := finder.NewFinder(subEng, &finder.RegexpEngine{}, false)
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	findthem.ForceBuild()
	b.ResetTimer()
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			findthem.ProcessText(text)
		}
	})
}

func getTerms(numTerm int) []string {
	terms := make([]string, numTerm)
	dictLen := len(words)
	for i := 1; i < numTerm; i++ {
		terms[i] = words[rand.Intn(dictLen)]
	}
	return terms
}

func createExpressions(terms []string) (regexes []string, andExps []string, inordExps []string, rgxExps []string) {
	for i := 0; i < len(terms); {
		firstTerm := terms[i]
		i++
		secondTerm := terms[i]
		i++
		regexes = append(regexes, fmt.Sprintf("%s.*%s", firstTerm, secondTerm))
		andExps = append(andExps, fmt.Sprintf("\"%s\" and \"%s\"", firstTerm, secondTerm))
		inordExps = append(inordExps, fmt.Sprintf("inord(\"%s\" and \"%s\")", firstTerm, secondTerm))
		rgxExps = append(rgxExps, fmt.Sprintf("r\"%s\" and r\"%s\"", firstTerm, secondTerm))
	}
	return
}
