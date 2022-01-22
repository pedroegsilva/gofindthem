package finder

import (
	goahocorasick "github.com/anknown/ahocorasick"
	cfahocorasick "github.com/cloudflare/ahocorasick"
	forkahocorasick "github.com/pedroegsilva/ahocorasick/ahocorasick"
)

// SubstringEngine interface that finder
// uses to search for substrings on a given text
type SubstringEngine interface {
	// BuildEngine receive the unique terms that need
	// to be searched to create the engine support structures
	BuildEngine(keywords map[string]struct{}, caseSensitive bool) (err error)
	// FindSubstrings receive the text and searchs for the feeded
	// terms
	FindSubstrings(text string) (matches []*Match, err error)
}

// AnknownEngine implements SubstringEngine using the
// github.com/anknown/ahocorasick package
type AnknownEngine struct {
	AhoEngine *goahocorasick.Machine
}

// BuildEngine implements BuildEngine using the
// github.com/anknown/ahocorasick package
func (am *AnknownEngine) BuildEngine(keywords map[string]struct{}, caseSensitive bool) (err error) {
	dict := [][]rune{}
	for key := range keywords {
		dict = append(dict, []rune(key))
	}
	ahoMachine := new(goahocorasick.Machine)
	err = ahoMachine.Build(dict)
	if err != nil {
		return
	}
	am.AhoEngine = ahoMachine
	return
}

// FindSubstrings implements FindSubstrings using the
// github.com/anknown/ahocorasick package
func (am *AnknownEngine) FindSubstrings(text string) (matches []*Match, err error) {
	ms := am.AhoEngine.MultiPatternSearch([]rune(text), false)
	for _, m := range ms {
		matches = append(matches, &Match{
			Term:     string(m.Word),
			Position: m.Pos,
		})
	}
	return
}

// CloudflareEngine implements SubstringEngine using the
// github.com/cloudflare/ahocorasick package. This engine
// does not support the use of INORD operator
type CloudflareEngine struct {
	Matcher *cfahocorasick.Matcher
	Dict    []string
}

// BuildEngine implements BuildEngine using the
// github.com/cloudflare/ahocorasick package
func (cfm *CloudflareEngine) BuildEngine(keywords map[string]struct{}, caseSensitive bool) (err error) {
	dict := []string{}
	for key := range keywords {
		dict = append(dict, key)
	}
	cfm.Matcher = cfahocorasick.NewStringMatcher(dict)
	cfm.Dict = dict
	return
}

// FindSubstrings implements FindSubstrings using the
// github.com/cloudflare/ahocorasick package
func (cfm *CloudflareEngine) FindSubstrings(text string) (matches []*Match, err error) {
	ms := cfm.Matcher.Match([]byte(text))
	for _, m := range ms {
		matches = append(matches, &Match{
			Term:     cfm.Dict[m],
			Position: 0,
		})
	}
	return
}

// CloudflareEngine implements SubstringEngine using the
// github.com/pedroegsilva/ahocorasick package. This engine
// does not support the use of INORD operator
type CloudflareForkEngine struct {
	Matcher *forkahocorasick.Matcher
	Dict    []string
}

// BuildEngine implements BuildEngine using the
// github.com/cloudflare/ahocorasick package
func (cffm *CloudflareForkEngine) BuildEngine(keywords map[string]struct{}, caseSensitive bool) (err error) {
	dict := []string{}
	for key := range keywords {
		dict = append(dict, key)
	}
	cffm.Matcher = forkahocorasick.NewStringMatcher(dict)
	cffm.Dict = dict
	return
}

// FindSubstrings implements FindSubstrings using the
// github.com/cloudflare/ahocorasick package
func (cffm *CloudflareForkEngine) FindSubstrings(text string) (matches []*Match, err error) {
	ms := cffm.Matcher.MatchAll([]byte(text))
	for _, hit := range ms {
		matches = append(matches, &Match{
			Term:     cffm.Dict[hit.DictIndex],
			Position: hit.Position,
		})
	}
	return
}

// EmptyEngine implements SubstringEngine with an nop
type EmptyEngine struct {
}

// BuildEngine implements BuildEngine with an nop
func (pdm *EmptyEngine) BuildEngine(keywords map[string]struct{}, caseSensitive bool) (err error) {
	return
}

// BuildEngine implements FindSubstrings with an nop
func (pdm *EmptyEngine) FindSubstrings(text string) (matches []*Match, err error) {
	return
}
