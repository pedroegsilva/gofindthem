package finder

import (
	"fmt"
	"regexp"
)

type RegexEngine interface {
	// BuildEngine receive the unique terms that need
	// to be searched to create the engine support structures
	BuildEngine(regexes map[string]struct{}, caseSensitive bool) (err error)
	// FindRegexes receive the text and searchs for the feeded
	// regexes
	FindRegexes(text string) (result chan *Match, err error)
}

type RegexpEngine struct {
	compiledRegexes []*regexp.Regexp
}

func (re *RegexpEngine) BuildEngine(regexes map[string]struct{}, caseSensitive bool) (err error) {
	if len(re.compiledRegexes) > 0 {
		re.compiledRegexes = re.compiledRegexes[:0]
	}

	for rgx := range regexes {
		r, err := regexp.Compile(rgx)
		if err != nil {
			return err
		}
		re.compiledRegexes = append(re.compiledRegexes, r)
	}
	return
}

func (re *RegexpEngine) FindRegexes(text string) (matches chan *Match, err error) {
	matches = make(chan *Match, 10)
	go func() {
		defer close(matches)
		for _, rgx := range re.compiledRegexes {
			positions := rgx.FindAllStringIndex(text, -1)
			for _, pos := range positions {
				matches <- &Match{
					Term:     fmt.Sprintf("%v", rgx),
					Position: pos[0],
				}
			}
		}
	}()
	return
}
