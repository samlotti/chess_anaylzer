package ai

import (
	"fmt"
	"regexp"
)

type Token struct {
	key     string
	literal string
}

type MiniLexer struct {
	patterns map[string]*regexp.Regexp
}

func NewMiniLexer() *MiniLexer {
	var m = &MiniLexer{}
	m.patterns = make(map[string]*regexp.Regexp)
	return m
}

// AddPatterns - Add the patterns the lexer.
// Note all patters will have ^ implicitly added.
func (m *MiniLexer) AddPatterns(patterns map[string]string) error {
	for k, v := range patterns {
		r, err := regexp.Compile(fmt.Sprintf("^(%s)", v))
		if err != nil {
			return err
		}
		m.patterns[k] = r
	}
	return nil
}

// GetToken - Returns the token for the beginning of the string
func (m *MiniLexer) GetToken(text string) (*Token, error) {
	var curr *Token = nil
	for k, v := range m.patterns {
		r := v.Find([]byte(text))
		if len(r) > 0 {
			if curr == nil {
				curr = &Token{
					key:     k,
					literal: string(r),
				}
			} else {
				if len(curr.literal) < len(r) {
					curr.key = k
					curr.literal = string(r)
				}
			}
		}
	}
	if curr == nil {
		return nil, fmt.Errorf("invalid character found: %c", text[0])
	}
	return curr, nil
}

func (m *MiniLexer) Lex(text string) ([]*Token, error) {
	return nil, nil
}
