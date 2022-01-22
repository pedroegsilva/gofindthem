package finderbenchmarks

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	akahocorasick "github.com/anknown/ahocorasick"
	cfahocorasick "github.com/cloudflare/ahocorasick"
	"github.com/pedroegsilva/gofindthem/dsl"
	"github.com/pedroegsilva/gofindthem/finder"
	pdahocorasick "github.com/petar-dambovaliev/aho-corasick"
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

	exp100, solverMapComp100, solverMapPart100 = createRandExpressionAndSolverMap(100)
	exp10000, solverMapComp10000, solverMapPart10000 = createRandExpressionAndSolverMap(10000)
	exps10 = createExpressions(10)
	exps100 = createExpressions(100)
	exps1000 = createExpressions(1000)
	knownTerms := make([]string, 10000)
	count = 0
	for word := range solverMapComp10000 {
		knownTerms[count] = word
		count++
	}
	randText100000 = createText(100000, knownTerms)
}

var alphabet = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

var (
	words                                  [466550]string
	exp100, exp10000                       string
	solverMapComp100, solverMapPart100     map[string]dsl.PatternResult
	solverMapComp10000, solverMapPart10000 map[string]dsl.PatternResult
	exps10                                 []string
	exps100                                []string
	exps1000                               []string
	randText100000                         string
)

const (
	EN_WORDS_FILE = "./files/words.txt"
)

// 100 terms

func BenchmarkParser100(b *testing.B) {
	BMParser(exp100, b)
}

func BenchmarkSolverCompleteMap100(b *testing.B) {
	BMSolver(exp100, solverMapComp100, true, b)
}

func BenchmarkSolverPartialMap100(b *testing.B) {
	BMSolver(exp100, solverMapPart100, false, b)
}

func BenchmarkSolverIterCompleteMap100(b *testing.B) {
	BMSolverIter(exp100, solverMapComp100, true, b)
}

func BenchmarkSolverIterPartialMap100(b *testing.B) {
	BMSolverIter(exp100, solverMapPart100, false, b)
}

func BenchmarkAhocorasickCloudFlareBuild100(b *testing.B) {
	BMCloudFlareBuild(exp100, b)
}

func BenchmarkAhocorasickAnknownBuild100(b *testing.B) {
	BMAnknownBuild(exp100, b)
}

func BenchmarkAhocorasickPetarDambovalievBuild100(b *testing.B) {
	BMPetarDambovalievBuild(exp100, b)
}

func BenchmarkAhocorasickCloudFlareSearch100(b *testing.B) {
	BMCloudFlareSearch([]string{exp100}, b)
}

func BenchmarkDslWithCloudFlare100(b *testing.B) {
	BMDslSearch([]string{exp100}, &finder.CloudflareEngine{}, b)
}

func BenchmarkAhocorasickAnknownSearch100(b *testing.B) {
	BMAnknownSearch([]string{exp100}, b)
}

func BenchmarkDslWithAnknown100(b *testing.B) {
	BMDslSearch([]string{exp100}, &finder.AnknownEngine{}, b)
}

func BenchmarkAhocorasickPetarDambovalievSearch100(b *testing.B) {
	BMPetarDambovalievSearch([]string{exp100}, b)
}

func BenchmarkDslWithPetarDambovaliev100(b *testing.B) {
	BMDslSearch([]string{exp100}, &finder.PetarDambovalievEngine{}, b)
}

// 10000 terms

func BenchmarkParser10000(b *testing.B) {
	BMParser(exp10000, b)
}

func BenchmarkSolverCompleteMap10000(b *testing.B) {
	BMSolver(exp10000, solverMapComp10000, true, b)
}

func BenchmarkSolverPartialMap10000(b *testing.B) {
	BMSolver(exp10000, solverMapPart10000, false, b)
}

func BenchmarkSolverIterCompleteMap10000(b *testing.B) {
	BMSolverIter(exp10000, solverMapComp10000, true, b)
}

func BenchmarkSolverIterPartialMap10000(b *testing.B) {
	BMSolverIter(exp10000, solverMapPart10000, false, b)
}

func BenchmarkAhocorasickCloudFlareBuild10000(b *testing.B) {
	BMCloudFlareBuild(exp10000, b)
}

func BenchmarkAhocorasickAnknownBuild10000(b *testing.B) {
	BMAnknownBuild(exp10000, b)
}

func BenchmarkAhocorasickPetarDambovalievBuild10000(b *testing.B) {
	BMPetarDambovalievBuild(exp10000, b)
}

func BenchmarkAhocorasickCloudFlareSearch10000(b *testing.B) {
	BMCloudFlareSearch([]string{exp10000}, b)
}

func BenchmarkDslWithCloudFlare10000(b *testing.B) {
	BMDslSearch([]string{exp10000}, &finder.CloudflareEngine{}, b)
}

func BenchmarkAhocorasickAnknownSearch10000(b *testing.B) {
	BMAnknownSearch([]string{exp10000}, b)
}

func BenchmarkDslWithAnknown10000(b *testing.B) {
	BMDslSearch([]string{exp10000}, &finder.AnknownEngine{}, b)
}

func BenchmarkAhocorasickPetarDambovalievSearch10000(b *testing.B) {
	BMPetarDambovalievSearch([]string{exp10000}, b)
}

func BenchmarkDslWithPetarDambovaliev10000(b *testing.B) {
	BMDslSearch([]string{exp10000}, &finder.PetarDambovalievEngine{}, b)
}

// dsl specific
func BenchmarkDslWithEmptyEngine10Exps(b *testing.B) {
	BMDslSearch(exps10, &finder.EmptyEngine{}, b)
}

func BenchmarkOnlyCloudFlare10Exps(b *testing.B) {
	BMCloudFlareSearch(exps10, b)
}

func BenchmarkDslWithCloudFlare10Exps(b *testing.B) {
	BMDslSearch(exps10, &finder.CloudflareEngine{}, b)
}

func BenchmarkOnlyAnknown10Exps(b *testing.B) {
	BMAnknownSearch(exps10, b)
}

func BenchmarkDslWithAnknown10Exps(b *testing.B) {
	BMDslSearch(exps10, &finder.AnknownEngine{}, b)
}

func BenchmarkOnlyPetarDambovaliev10Exps(b *testing.B) {
	BMPetarDambovalievSearch(exps10, b)
}

func BenchmarkDslWithPetarDambovaliev10Exps(b *testing.B) {
	BMDslSearch(exps10, &finder.PetarDambovalievEngine{}, b)
}

func BenchmarkDslWithEmptyEngine100Exps(b *testing.B) {
	BMDslSearch(exps100, &finder.EmptyEngine{}, b)
}

func BenchmarkOnlyCloudFlare100Exps(b *testing.B) {
	BMCloudFlareSearch(exps100, b)
}

func BenchmarkDslWithCloudFlare100Exps(b *testing.B) {
	BMDslSearch(exps100, &finder.CloudflareEngine{}, b)
}

func BenchmarkOnlyAnknown100Exps(b *testing.B) {
	BMAnknownSearch(exps100, b)
}

func BenchmarkDslWithAnknown100Exps(b *testing.B) {
	BMDslSearch(exps100, &finder.AnknownEngine{}, b)
}

func BenchmarkOnlyPetarDambovaliev100Exps(b *testing.B) {
	BMPetarDambovalievSearch(exps100, b)
}

func BenchmarkDslWithPetarDambovaliev100Exps(b *testing.B) {
	BMDslSearch(exps100, &finder.PetarDambovalievEngine{}, b)
}

func BenchmarkDslWithEmptyEngine1000Exps(b *testing.B) {
	BMDslSearch(exps1000, &finder.EmptyEngine{}, b)
}

func BenchmarkOnlyCloudFlare1000Exps(b *testing.B) {
	BMCloudFlareSearch(exps1000, b)
}

func BenchmarkDslWithCloudFlare1000Exps(b *testing.B) {
	BMDslSearch(exps1000, &finder.CloudflareEngine{}, b)
}

func BenchmarkOnlyAnknown1000Exps(b *testing.B) {
	BMAnknownSearch(exps1000, b)
}

func BenchmarkDslWithAnknown1000Exps(b *testing.B) {
	BMDslSearch(exps1000, &finder.AnknownEngine{}, b)
}

func BenchmarkOnlyPetarDambovaliev1000Exps(b *testing.B) {
	BMPetarDambovalievSearch(exps1000, b)
}

func BenchmarkDslWithPetarDambovaliev1000Exps(b *testing.B) {
	BMDslSearch(exps1000, &finder.PetarDambovalievEngine{}, b)
}

func BenchmarkUseCasesDsl(b *testing.B) {
	expressions := []string{
		`"foo" and "bar"`,
	}
	BMDslSearch(expressions, &finder.PetarDambovalievEngine{}, b)
}

func BenchmarkUseCasesDslWithRegex(b *testing.B) {
	expressions := []string{
		`r"foo.*bar" and r"bar.*foo"`,
	}
	BMDslSearch(expressions, &finder.PetarDambovalievEngine{}, b)
}

func BenchmarkUseCasesDslWithInord(b *testing.B) {
	expressions := []string{
		`INORD("foo" and "bar") and INORD("bar" and "foo")`,
	}
	BMDslSearch(expressions, &finder.PetarDambovalievEngine{}, b)
}

func BenchmarkUseCasesRegexOnly(b *testing.B) {
	rgx1 := regexp.MustCompile("foo.*bar")
	rgx2 := regexp.MustCompile("bar.*foo")

	for i := 0; i < b.N; i++ {
		rgx1.FindAllStringIndex(randText100000, -1)
		rgx2.FindAllStringIndex(randText100000, -1)
	}
}

// test funcs

func BMParser(exp string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		dsl.NewParser(strings.NewReader(exp), true).Parse()
	}
}

func BMSolver(exp string, solverMap map[string]dsl.PatternResult, completeMap bool, b *testing.B) {
	p := dsl.NewParser(strings.NewReader(exp), true)
	e, _ := p.Parse()
	for i := 0; i < b.N; i++ {
		e.Solve(solverMap, completeMap)
	}
}

func BMSolverIter(exp string, solverMap map[string]dsl.PatternResult, completeMap bool, b *testing.B) {
	p := dsl.NewParser(strings.NewReader(exp), true)
	e, _ := p.Parse()
	so := e.CreateSolverOrder()
	for i := 0; i < b.N; i++ {
		so.Solve(solverMap, completeMap)
	}
}

func BMCloudFlareBuild(exp string, b *testing.B) {
	p := dsl.NewParser(strings.NewReader(exp), true)
	p.Parse()
	dict := [][]byte{}
	for key := range p.GetKeywords() {
		dict = append(dict, []byte(key))
	}

	for i := 0; i < b.N; i++ {
		cfahocorasick.NewMatcher(dict)
	}
}

func BMAnknownBuild(exp string, b *testing.B) {
	p := dsl.NewParser(strings.NewReader(exp), true)
	p.Parse()
	dict := [][]rune{}
	for key := range p.GetKeywords() {
		dict = append(dict, []rune(key))
	}

	for i := 0; i < b.N; i++ {
		m := new(akahocorasick.Machine)
		m.Build(dict)
	}
}

func BMPetarDambovalievBuild(exp string, b *testing.B) {
	p := dsl.NewParser(strings.NewReader(exp100), true)
	p.Parse()
	dict := []string{}

	for key := range p.GetKeywords() {
		dict = append(dict, key)
	}

	for i := 0; i < b.N; i++ {
		builder := pdahocorasick.NewAhoCorasickBuilder(pdahocorasick.Opts{
			AsciiCaseInsensitive: true,
			MatchOnlyWholeWords:  false,
			MatchKind:            pdahocorasick.LeftMostLongestMatch,
			DFA:                  true,
		})
		builder.Build(dict)
	}
}

func BMCloudFlareSearch(exps []string, b *testing.B) {
	findthem := finder.NewFinder(&finder.EmptyEngine{}, &finder.RegexpEngine{}, true)
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	dict := [][]byte{}
	for key := range findthem.GetKeywords() {
		dict = append(dict, []byte(key))
	}

	m := cfahocorasick.NewMatcher(dict)

	content := []byte(randText100000)
	for i := 0; i < b.N; i++ {
		m.Match(content)
	}
}

func BMAnknownSearch(exps []string, b *testing.B) {
	findthem := finder.NewFinder(&finder.EmptyEngine{}, &finder.RegexpEngine{}, true)
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	dict := [][]rune{}
	for key := range findthem.GetKeywords() {
		dict = append(dict, []rune(key))
	}

	m := new(akahocorasick.Machine)
	m.Build(dict)

	contentRune := bytes.Runes([]byte(randText100000))
	for i := 0; i < b.N; i++ {
		m.MultiPatternSearch(contentRune, false)
	}
}

func BMPetarDambovalievSearch(exps []string, b *testing.B) {
	findthem := finder.NewFinder(&finder.EmptyEngine{}, &finder.RegexpEngine{}, true)
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	dict := []string{}

	for key := range findthem.GetKeywords() {
		dict = append(dict, key)
	}

	builder := pdahocorasick.NewAhoCorasickBuilder(pdahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  true,
		MatchKind:            pdahocorasick.LeftMostLongestMatch,
		DFA:                  false,
	})
	bld := builder.Build(dict)
	for i := 0; i < b.N; i++ {
		bld.FindAll(randText100000)
	}
}

func BMDslSearch(exps []string, subEng finder.SubstringEngine, b *testing.B) {
	findthem := finder.NewFinder(subEng, &finder.RegexpEngine{}, false)
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	findthem.ForceBuild()
	for i := 0; i < b.N; i++ {
		findthem.ProcessText(randText100000)
	}
}

// aux

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func createRandExpressionAndSolverMap(
	numTerm int,
) (string, map[string]dsl.PatternResult, map[string]dsl.PatternResult) {
	solverMapComp := make(map[string]dsl.PatternResult, numTerm)
	solverMapPart := make(map[string]dsl.PatternResult, numTerm/4)
	expression := ""
	dictLen := len(words)
	for i := 1; i <= numTerm; i++ {
		idx := rand.Intn(dictLen)
		keyword := words[idx]
		if i == numTerm {
			expression = expression + fmt.Sprintf(`"%s"`, keyword)
		} else {
			expression = expression + fmt.Sprintf(`"%s" AND `, keyword)
		}
		if rand.Intn(4) == 0 {
			solverMapComp[keyword] = dsl.PatternResult{
				Val: true,
			}
			solverMapPart[keyword] = dsl.PatternResult{
				Val: true,
			}
		} else {
			solverMapComp[keyword] = dsl.PatternResult{
				Val: false,
			}
		}
	}
	expression = "INORD(" + expression + ")"
	return expression, solverMapComp, solverMapPart
}

func createExpressions(numExp int) []string {
	exps := make([]string, numExp)
	for i := 0; i < numExp; i++ {
		exp, _, _ := createRandExpressionAndSolverMap(1 + rand.Intn(10))
		exps[i] = exp
	}
	return exps
}

func createText(numTerm int, knownTerms []string) string {
	text := ""
	dictLen := len(words)
	ktermsLen := len(knownTerms)
	for i := 1; i <= numTerm; i++ {
		if ktermsLen > 0 {
			if rand.Intn(20) == 0 {
				text += knownTerms[rand.Intn(ktermsLen)] + " "
			} else {
				text += words[rand.Intn(dictLen)] + " "
			}
		} else {
			text += words[rand.Intn(dictLen)] + " "
		}
	}
	return text
}
