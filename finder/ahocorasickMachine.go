package finder

import (
	goahocorasick "github.com/anknown/ahocorasick"
	cfahocorasick "github.com/cloudflare/ahocorasick"
	pdahocorasick "github.com/petar-dambovaliev/aho-corasick"
)

type Match struct {
	Position  int
	Substring string
}

type SubstringEngine interface {
	BuildEngine(keywords map[string]bool) (err error)
	ProcessText(text string) (result chan *Match, err error)
}

type AnknownMachine struct {
	AhoMachine *goahocorasick.Machine
}

func (am *AnknownMachine) BuildEngine(keywords map[string]bool) (err error) {
	dict := [][]rune{}
	for key := range keywords {
		dict = append(dict, []rune(key))
	}
	ahoMachine := new(goahocorasick.Machine)
	err = ahoMachine.Build(dict)
	if err != nil {
		return
	}
	am.AhoMachine = ahoMachine
	return
}

func (am *AnknownMachine) ProcessText(text string) (matches chan *Match, err error) {
	matches = make(chan *Match, 100)
	go func() {
		defer close(matches)
		ms := am.AhoMachine.MultiPatternSearch([]rune(text), false)
		for _, m := range ms {
			matches <- &Match{
				Substring: string(m.Word),
				Position:  m.Pos,
			}
		}
	}()
	return
}

type CloudflareMachine struct {
	Matcher *cfahocorasick.Matcher
	Dict    []string
}

func (cfm *CloudflareMachine) BuildEngine(keywords map[string]bool) (err error) {
	dict := []string{}
	for key := range keywords {
		dict = append(dict, key)
	}
	cfm.Matcher = cfahocorasick.NewStringMatcher(dict)
	cfm.Dict = dict
	return
}

func (cfm *CloudflareMachine) ProcessText(text string) (matches chan *Match, err error) {
	matches = make(chan *Match, 100)
	go func() {
		defer close(matches)
		ms := cfm.Matcher.Match([]byte(text))
		for _, m := range ms {
			matches <- &Match{
				Substring: cfm.Dict[m],
				Position:  0,
			}
		}
	}()
	return
}

type PetarDambovalievMachine struct {
	AhoMachine pdahocorasick.AhoCorasick
}

func (pdm *PetarDambovalievMachine) BuildEngine(keywords map[string]bool) (err error) {
	dict := make([]string, 0, len(keywords))
	for key := range keywords {
		dict = append(dict, key)
	}
	builder := pdahocorasick.NewAhoCorasickBuilder(pdahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  false,
		MatchKind:            pdahocorasick.StandardMatch,
		DFA:                  true,
	})

	pdm.AhoMachine = builder.Build(dict)
	return
}

func (pdm *PetarDambovalievMachine) ProcessText(text string) (matches chan *Match, err error) {
	matches = make(chan *Match, 100)
	go func() {
		defer close(matches)
		ms := pdm.AhoMachine.FindAll(text)
		for _, m := range ms {
			matches <- &Match{
				Substring: text[m.Start():m.End()],
				Position:  m.Start(),
			}
		}
	}()
	return
}

type EmptyMachine struct {
}

func (pdm *EmptyMachine) BuildEngine(keywords map[string]bool) (err error) {
	return
}

func (pdm *EmptyMachine) ProcessText(text string) (matches chan *Match, err error) {
	matches = make(chan *Match, 100)
	defer close(matches)
	return
}
