package ai

import (
	"fmt"
	"regexp"
	"strings"
)

type TokenId int

const UserTokeId = 100

type MatcherEntry struct {
	Regex     string
	Id        TokenId
	_compiled *regexp.Regexp
}

// NewMatcher
// return the matcher struct or error
func NewMatcher(id TokenId, matchExpression string) (*MatcherEntry, error) {
	r, err := regexp.Compile(fmt.Sprintf("^(%s)", matchExpression))
	if err != nil {
		return nil, err
	}
	return &MatcherEntry{
		Regex:     matchExpression,
		Id:        id,
		_compiled: r,
	}, err
}

type Token struct {
	id      TokenId
	literal string
	line    int
	pos     int
}

type MiniLexerOptions struct {
	whiteSpaceChar string
}

func NewMiniLexOptions() *MiniLexerOptions {
	r := &MiniLexerOptions{
		whiteSpaceChar: "\n\t\v\f\r ",
	}
	return r
}

// RemoveAsWhiteSpace - By default \n \t \f \v \r and space are whitespace.
// this api can remove one from the list
// ex:  mo.RemoveAsWhiteSpace("\n")
func (mo *MiniLexerOptions) RemoveAsWhiteSpace(chr string) {
	mo.whiteSpaceChar = strings.ReplaceAll(mo.whiteSpaceChar, chr, "")
}

type MiniLexer struct {
	matchers []*MatcherEntry
	text     string
	line     int
	pos      int
	mo       *MiniLexerOptions
}

func NewMiniLexer(text string, mo *MiniLexerOptions) *MiniLexer {
	var m = &MiniLexer{}
	m.text = text
	m.line = 1
	m.pos = 0
	m.mo = mo
	return m
}

// AddPattern - Add the patterns the lexer.
// Note all patters will have ^ implicitly added.
func (m *MiniLexer) AddPattern(id TokenId, pattern string) error {
	if id < UserTokeId {
		return fmt.Errorf("id must be >= UserTokenId, %d", UserTokeId)
	}
	matcher, err := NewMatcher(id, pattern)
	if err != nil {
		return err
	}

	m.matchers = append(m.matchers, matcher)
	return nil
}

// NextToken - Returns the token for the beginning of the string
func (m *MiniLexer) NextToken() (*Token, error) {
	m.text = m.AdvanceSpaces(m.text)
	var curr *Token = nil
	for _, matcher := range m.matchers {
		r := matcher._compiled.Find([]byte(m.text))
		if len(r) > 0 {
			if curr == nil {
				curr = &Token{
					id:      matcher.Id,
					literal: string(r),
					line:    m.line,
					pos:     m.pos,
				}
			} else {
				if len(curr.literal) < len(r) {
					curr.id = matcher.Id
					curr.literal = string(r)
				}
			}
		}
	}
	if curr == nil {
		return nil, fmt.Errorf("invalid character found: %c", m.text[0])
	}
	m.advanceInput(curr)
	return curr, nil
}

func (m *MiniLexer) Lex(text string) ([]*Token, error) {
	return nil, nil
}

// advanceInput
// Move the input passed the token
func (m *MiniLexer) advanceInput(tk *Token) {

	for _, c := range tk.literal {
		m.pos += 1
		if c == '\n' {
			m.pos = 0
			m.line += 1
		}
	}

	m.text = m.text[len(tk.literal):len(m.text)]
	//m.pos += len(tk.literal)
}

// AdvanceSpaces
// Move the input passed any spaces
func (m *MiniLexer) AdvanceSpaces(text string) string {

	slen := len(m.text)
	r := strings.TrimLeftFunc(text, m.IsWhiteSpace)
	m.pos += slen - len(m.text)
	return r
}

func (m *MiniLexer) IsWhiteSpace(r rune) bool {
	var rchar string = string(r)
	m.pos += 1
	if rchar == "\n" {
		m.pos = 0
		m.line += 1
	}
	return strings.Contains(m.mo.whiteSpaceChar, rchar)
}
