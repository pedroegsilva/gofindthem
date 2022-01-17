package finder

import (
	goahocorasick "github.com/anknown/ahocorasick"
	cfahocorasick "github.com/cloudflare/ahocorasick"
	pdahocorasick "github.com/petar-dambovaliev/aho-corasick"
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
		if key == "" {
			continue
		}
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

// PetarDambovalievEngine implements SubstringEngine using the
// github.com/petar-dambovaliev/aho-corasick package
type PetarDambovalievEngine struct {
	AhoEngine pdahocorasick.AhoCorasick
}

// BuildEngine implements BuildEngine using the
// github.com/petar-dambovaliev/aho-corasick package
func (pdm *PetarDambovalievEngine) BuildEngine(keywords map[string]struct{}, caseSensitive bool) (err error) {
	dict := make([]string, 0, len(keywords))
	for key := range keywords {
		dict = append(dict, key)
	}
	builder := pdahocorasick.NewAhoCorasickBuilder(pdahocorasick.Opts{
		AsciiCaseInsensitive: !caseSensitive,
		MatchOnlyWholeWords:  false,
		MatchKind:            pdahocorasick.StandardMatch,
		DFA:                  true,
	})

	pdm.AhoEngine = builder.Build(dict)
	return
}

// FindSubstrings implements FindSubstrings using the
// github.com/petar-dambovaliev/aho-corasick package
func (pdm *PetarDambovalievEngine) FindSubstrings(text string) (matches []*Match, err error) {
	ms := pdm.AhoEngine.FindAll(text)
	for _, m := range ms {
		matches = append(matches, &Match{
			Term:     text[m.Start():m.End()],
			Position: m.Start(),
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
