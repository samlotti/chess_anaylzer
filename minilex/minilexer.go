package minilex

import (
	"fmt"
	"regexp"
	"strings"
)

type TokenId int

const UserTokeId = 100

const (
	TKEof TokenId = -1
)

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
	Id      TokenId
	Literal string
	Line    int
	Pos     int
}

func (t *Token) String() string {
	return fmt.Sprintf("%d:%d = %d %s", t.Line, t.Pos, t.Id, t.Literal)
}

func (t *Token) Is(id TokenId) bool {
	return t.Id == id
}

func (t *Token) AssertIs(id TokenId) error {
	if t.Id != id {
		return fmt.Errorf("expected %d, found %d : %s at %d:%d", id, t.Id, t.Literal, t.Line, t.Pos)
	}
	return nil
}

type MiniLexerOptions struct {
	whiteSpaceChars string
}

func NewMiniLexOptions() *MiniLexerOptions {
	r := &MiniLexerOptions{
		whiteSpaceChars: "\n\t\v\f\r ",
	}
	return r
}

// RemoveAsWhiteSpace - By default \n \t \f \v \r and space are whitespace.
// this api can remove one from the list
// ex:  mo.RemoveAsWhiteSpace("\n")
func (mo *MiniLexerOptions) RemoveAsWhiteSpace(chr string) {
	mo.whiteSpaceChars = strings.ReplaceAll(mo.whiteSpaceChars, chr, "")
}

// SetWhiteSpace
// Set your own whitespace
func (mo *MiniLexerOptions) SetWhiteSpace(whiteSpaceChars string) {
	mo.whiteSpaceChars = whiteSpaceChars
}

func (mo *MiniLexerOptions) GetWhiteSpaceCharacters() string {
	return mo.whiteSpaceChars
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
		return fmt.Errorf("Id must be >= UserTokenId, %d", UserTokeId)
	}
	matcher, err := NewMatcher(id, pattern)
	if err != nil {
		return err
	}

	m.matchers = append(m.matchers, matcher)
	return nil
}

// PeekToken - Returns the next token but will not advance.
func (m *MiniLexer) PeekToken() (*Token, error) {
	m.text = m.AdvanceSpaces(m.text)

	if len(m.text) == 0 {
		return &Token{
			Id:      TKEof,
			Literal: "",
			Line:    m.line,
			Pos:     m.pos,
		}, nil
	}

	var curr *Token = nil
	for _, matcher := range m.matchers {
		r := matcher._compiled.Find([]byte(m.text))
		if len(r) > 0 {
			if curr == nil {
				curr = &Token{
					Id:      matcher.Id,
					Literal: string(r),
					Line:    m.line,
					Pos:     m.pos,
				}
			} else {
				if len(curr.Literal) < len(r) {
					curr.Id = matcher.Id
					curr.Literal = string(r)
				}
			}
		}
	}
	if curr == nil {
		return nil, fmt.Errorf("invalid character found: %c", m.text[0])
	}
	return curr, nil
}

// NextToken - Returns the token for the beginning of the string
func (m *MiniLexer) NextToken() (*Token, error) {

	curr, err := m.PeekToken()
	if err != nil {
		return nil, err
	}

	m.advanceInput(curr)
	return curr, nil
}

// advanceInput
// Move the input passed the token
func (m *MiniLexer) advanceInput(tk *Token) {

	for _, c := range tk.Literal {
		m.pos += 1
		if c == '\n' {
			m.pos = 0
			m.line += 1
		}
	}

	m.text = m.text[len(tk.Literal):len(m.text)]
	//m.Pos += len(tk.Literal)
}

// AdvanceSpaces
// Move the input passed any spaces
func (m *MiniLexer) AdvanceSpaces(text string) string {

	r := strings.TrimLeftFunc(text, m.IsWhiteSpace)
	return r
}

func (m *MiniLexer) IsWhiteSpace(r rune) bool {
	var rchar string = string(r)
	var isW = strings.Contains(m.mo.whiteSpaceChars, rchar)
	if isW {
		m.pos += 1
		if rchar == "\n" {
			m.pos = 0
			m.line += 1
		}
	}
	return isW
}

func (m *MiniLexer) ReadAllTokens() ([]*Token, error) {

	var tkList []*Token = nil

	for {
		tk, err := m.NextToken()
		if err != nil {
			return nil, err
		}

		if tk.Id == TKEof {
			break
		}

		tkList = append(tkList, tk)
	}
	return tkList, nil
}

func (m *MiniLexer) IsEOF() bool {
	m.text = m.AdvanceSpaces(m.text)
	return len(m.text) == 0
}
