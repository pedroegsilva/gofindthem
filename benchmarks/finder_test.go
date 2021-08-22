package finderbenchmarks

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	akahocorasick "github.com/anknown/ahocorasick"
	cfahocorasick "github.com/cloudflare/ahocorasick"
	"github.com/pedroegsilva/gofindthem/dsl"
	"github.com/pedroegsilva/gofindthem/finder"
	pdahocorasick "github.com/petar-dambovaliev/aho-corasick"
)

// bazel run //dsl:dsl_test -- -test.bench=. -test.benchmem
// bazel run //dsl:dsl_test -- -test.bench=Exps -test.benchmem

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
	randText10000 = createText(10000, knownTerms)
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
	randText5000, randText10000            string
)

const (
	EN_WORDS_FILE = "./files/words.txt"
)

// 100 terms

func BenchmarkParser100(b *testing.B) {
	BMParser(exp100, b)
}

func BenchmarkParserInter100(b *testing.B) {
	BMParserInter(exp100, b)
}

func BenchmarkSolverNoCacheCompleteMap100(b *testing.B) {
	BMSolver(exp100, solverMapComp100, true, b)
}

func BenchmarkSolverNoCachePartialMap100(b *testing.B) {
	BMSolver(exp100, solverMapPart100, false, b)
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
	BMDslSearch([]string{exp100}, &finder.CloudflareMachine{}, b)
}

func BenchmarkAhocorasickAnknownSearch100(b *testing.B) {
	BMAnknownSearch([]string{exp100}, b)
}

func BenchmarkDslWithAnknown100(b *testing.B) {
	BMDslSearch([]string{exp100}, &finder.AnknownMachine{}, b)
}

func BenchmarkAhocorasickPetarDambovalievSearch100(b *testing.B) {
	BMPetarDambovalievSearch([]string{exp100}, b)
}

func BenchmarkDslWithPetarDambovaliev100(b *testing.B) {
	BMDslSearch([]string{exp100}, &finder.PetarDambovalievMachine{}, b)
}

// 10000 terms

func BenchmarkParser10000(b *testing.B) {
	BMParser(exp10000, b)
}

func BenchmarkParserInter10000(b *testing.B) {
	BMParserInter(exp10000, b)
}

func BenchmarkSolverNoCacheCompleteMap10000(b *testing.B) {
	BMSolver(exp10000, solverMapComp10000, true, b)
}

func BenchmarkSolverNoCachePartialMap10000(b *testing.B) {
	BMSolver(exp10000, solverMapPart10000, false, b)
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
	BMDslSearch([]string{exp10000}, &finder.CloudflareMachine{}, b)
}

func BenchmarkAhocorasickAnknownSearch10000(b *testing.B) {
	BMAnknownSearch([]string{exp10000}, b)
}

func BenchmarkDslWithAnknown10000(b *testing.B) {
	BMDslSearch([]string{exp10000}, &finder.AnknownMachine{}, b)
}

func BenchmarkAhocorasickPetarDambovalievSearch10000(b *testing.B) {
	BMPetarDambovalievSearch([]string{exp10000}, b)
}

func BenchmarkDslWithPetarDambovaliev10000(b *testing.B) {
	BMDslSearch([]string{exp10000}, &finder.PetarDambovalievMachine{}, b)
}

// dsl specific
func BenchmarkDslWithEmptyMachine10Exps(b *testing.B) {
	BMDslSearch(exps10, &finder.EmptyMachine{}, b)
}

func BenchmarkOnlyCloudFlare10Exps(b *testing.B) {
	BMCloudFlareSearch(exps10, b)
}

func BenchmarkDslWithCloudFlare10Exps(b *testing.B) {
	BMDslSearch(exps10, &finder.CloudflareMachine{}, b)
}

func BenchmarkOnlyAnknown10Exps(b *testing.B) {
	BMAnknownSearch(exps10, b)
}

func BenchmarkDslWithAnknown10Exps(b *testing.B) {
	BMDslSearch(exps10, &finder.AnknownMachine{}, b)
}

func BenchmarkOnlyPetarDambovaliev10Exps(b *testing.B) {
	BMPetarDambovalievSearch(exps10, b)
}

func BenchmarkDslWithPetarDambovaliev10Exps(b *testing.B) {
	BMDslSearch(exps10, &finder.PetarDambovalievMachine{}, b)
}

func BenchmarkDslWithEmptyMachine100Exps(b *testing.B) {
	BMDslSearch(exps100, &finder.EmptyMachine{}, b)
}

func BenchmarkOnlyCloudFlare100Exps(b *testing.B) {
	BMCloudFlareSearch(exps100, b)
}

func BenchmarkDslWithCloudFlare100Exps(b *testing.B) {
	BMDslSearch(exps100, &finder.CloudflareMachine{}, b)
}

func BenchmarkOnlyAnknown100Exps(b *testing.B) {
	BMAnknownSearch(exps100, b)
}

func BenchmarkDslWithAnknown100Exps(b *testing.B) {
	BMDslSearch(exps100, &finder.AnknownMachine{}, b)
}

func BenchmarkOnlyPetarDambovaliev100Exps(b *testing.B) {
	BMPetarDambovalievSearch(exps100, b)
}

func BenchmarkDslWithPetarDambovaliev100Exps(b *testing.B) {
	BMDslSearch(exps100, &finder.PetarDambovalievMachine{}, b)
}

func BenchmarkDslWithEmptyMachine1000Exps(b *testing.B) {
	BMDslSearch(exps1000, &finder.EmptyMachine{}, b)
}

func BenchmarkOnlyCloudFlare1000Exps(b *testing.B) {
	BMCloudFlareSearch(exps1000, b)
}

func BenchmarkDslWithCloudFlare1000Exps(b *testing.B) {
	BMDslSearch(exps1000, &finder.CloudflareMachine{}, b)
}

func BenchmarkOnlyAnknown1000Exps(b *testing.B) {
	BMAnknownSearch(exps1000, b)
}

func BenchmarkDslWithAnknown1000Exps(b *testing.B) {
	BMDslSearch(exps1000, &finder.AnknownMachine{}, b)
}

func BenchmarkOnlyPetarDambovaliev1000Exps(b *testing.B) {
	BMPetarDambovalievSearch(exps1000, b)
}

func BenchmarkDslWithPetarDambovaliev1000Exps(b *testing.B) {
	BMDslSearch(exps1000, &finder.PetarDambovalievMachine{}, b)
}

// test funcs

func BMParser(exp string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		dsl.NewParser(strings.NewReader(exp)).Parse()
	}
}

func BMParserInter(exp string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		dsl.NewParser(strings.NewReader(exp)).ParseInter()
	}
}

func BMSolver(exp string, solverMap map[string]dsl.PatternResult, completeMap bool, b *testing.B) {
	p := dsl.NewParser(strings.NewReader(exp))
	e, _ := p.Parse()
	for i := 0; i < b.N; i++ {
		e.Solve(solverMap, completeMap, false)
	}
}

func BMCloudFlareBuild(exp string, b *testing.B) {
	p := dsl.NewParser(strings.NewReader(exp))
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
	p := dsl.NewParser(strings.NewReader(exp))
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
	p := dsl.NewParser(strings.NewReader(exp100))
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
	findthem := finder.NewFinder(&finder.EmptyMachine{})
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	dict := [][]byte{}
	for key := range findthem.Keywords {
		dict = append(dict, []byte(key))
	}

	m := cfahocorasick.NewMatcher(dict)

	content := []byte(randText10000)
	for i := 0; i < b.N; i++ {
		m.Match(content)
	}
}

func BMAnknownSearch(exps []string, b *testing.B) {
	findthem := finder.NewFinder(&finder.EmptyMachine{})
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	dict := [][]rune{}
	for key := range findthem.Keywords {
		dict = append(dict, []rune(key))
	}

	m := new(akahocorasick.Machine)
	m.Build(dict)

	contentRune := bytes.Runes([]byte(randText10000))
	for i := 0; i < b.N; i++ {
		m.MultiPatternSearch(contentRune, false)
	}
}

func BMPetarDambovalievSearch(exps []string, b *testing.B) {
	findthem := finder.NewFinder(&finder.EmptyMachine{})
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	dict := []string{}

	for key := range findthem.Keywords {
		dict = append(dict, key)
	}

	builder := pdahocorasick.NewAhoCorasickBuilder(pdahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  false,
		MatchKind:            pdahocorasick.LeftMostLongestMatch,
		DFA:                  true,
	})
	bld := builder.Build(dict)
	for i := 0; i < b.N; i++ {
		bld.FindAll(randText10000)
	}
}

func BMDslSearch(exps []string, subEng finder.SubstringEngine, b *testing.B) {
	findthem := finder.NewFinder(subEng)
	for _, exp := range exps {
		findthem.AddExpression(exp)
	}

	findthem.ForceBuild()
	for i := 0; i < b.N; i++ {
		findthem.ProcessText(randText10000, false)
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
