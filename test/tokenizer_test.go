package gocloj

import (
	"strings"
	"testing"
	"gocloj/gocloj"
)

func TestTokenKeyword(t *testing.T) {
	str := ":keyword"
	tz := gocloj.NewTokenizer(strings.NewReader(str), "internal-test")
	if tz.Err() != nil {
		t.Errorf("got error %s", tz.Err())
	}
	if !tz.Next() {
		t.Errorf("expected next")
	}

	tok := tz.Value()
	if tok.T != gocloj.TokenKeyword {
		t.Errorf("got unexpected token type %s", tok.t)
	}

	if tz.Next() {
		t.Errorf("expected end, got %s", tz.Value())
	}

	/*
		TokenLParen
		TokenRParen
		TokenLBracket
		TokenRBracket
		TokenLCurly
		TokenRCurly
		TokenSymbol
		TokenString
		TokenChar
		TokenNil
		TokenNum
		TokenQuote
		TokenKeyword
	*/
}
