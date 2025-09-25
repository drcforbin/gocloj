package gocloj

import (
	"gocloj/data/atom"
	"testing"
)

type mockTokenizer struct {
	file string

	idx    int
	tokens []Token

	err    error
	errPos Position
}

func (m mockTokenizer) Position() (pos Position) {
	if m.err != nil {
		pos = m.errPos
	} else if m.idx < len(m.tokens) {
		pos = m.tokens[m.idx].pos
	}

	return
}

func (m *mockTokenizer) Next() bool {
	m.idx++
	return m.idx < len(m.tokens)
}

func (m mockTokenizer) Value() (tk Token) {
	if 0 <= m.idx && m.idx < len(m.tokens) {
		tk = m.tokens[m.idx]
	}

	return
}

func (m mockTokenizer) Err() error {
	return m.err
}

func newMockTokenizer(file string, tokens []Token) Tokenizer {
	// set up positions if they're missing
	for i, tk := range tokens {
		if tk.pos.File == "" || tk.pos.Line == 0 {
			tk.pos.File = file
			tk.pos.Line = i + 1
			tokens[i] = tk
		}
	}

	return &mockTokenizer{
		idx:    -1,
		file:   file,
		tokens: tokens,
	}
}

func testParseToAtom(t *testing.T, tokens []Token) atom.Atom {
	file := "test-file"
	tk := newMockTokenizer(file, tokens)

	p := NewParser(tk)

	if !p.Next() {
		t.Error("expected an atom from parser")
	}

	if p.Err() != nil {
		t.Error("did not expect error from parser")
	}

	val := p.Value()

	if p.Next() {
		t.Error("expected to be out of atoms from parser")
	}

	if p.Err() != nil {
		t.Error("did not expect error from parser")
	}

	return val
}

func TestParseKeyword(t *testing.T) {
	val := testParseToAtom(t, []Token{
		Token{t: TokenKeyword, s: ":abc"},
	},
	)

	if kw, ok := val.(*atom.Keyword); ok {
		if kw.Name != ":abc" {
			t.Error("expected keyword with proper name")
		}
	} else {
		t.Error("expected keyword")
	}
}
