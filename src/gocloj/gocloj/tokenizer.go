package gocloj

import (
	"bufio"
	"fmt"
	"gocloj/log"
	"io"
	"io/ioutil"
	"regexp"
	// "strings"
)

var tokLogger = log.Get("tokenizer")

type TokenType int

const (
	TokenNone TokenType = iota
	TokenLParen
	TokenRParen
	TokenLBracket
	TokenRBracket
	TokenLCurly
	TokenRCurly
	TokenSymbol
	TokenString
	TokenNil
	TokenNum
	TokenQuote
	TokenEOF
)

func (t Token) String() string {
	switch t.t {
	case TokenNone:
		return "<none>"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenLCurly:
		return "{"
	case TokenRCurly:
		return "}"
	case TokenSymbol:
		return t.s
	case TokenString:
		return t.s
	case TokenNil:
		return "<nil>"
	case TokenNum:
		return t.s
	case TokenQuote:
		return "'"
	case TokenEOF:
		return "<eof>"
	default:
		return fmt.Sprintf("<unknown %d>", t.t)
	}

	return t.s
}

type Tokenizer interface {
	Position() Position

	Next() bool
	Value() Token
	Err() error
}

type Token struct {
	t TokenType
	s string

	pos Position
}

type Position struct {
	Line int
	File string
}

// symbols begin with nonumeric, can contain
//   *, +, !, -, _, ', ?, <, >, and =
//   can contain one or more non repeating :
//     symbols beginning or ending w/ : are reserved
//   / namespace indicator
//   . class namespace
// ".." - string; can be multiple lines
// 'form - quote form
// \c - char literal
//    \newline
//    \space
//    \tab
//    \formfeed
//    \backspace
//    \return
// \uNNNN - unicode literal
// nil - nil
// true, false - booleans
// ##Inf, ##-Inf, ##NaN
// keywords can and must begin with :
//  cannot contain . inCurr name part
//  can contain / for ns, which may include .
//  keyword beginning with two colons gets current ns
// (...) - list
// [...] - vec
// {...} - map, must be pairs
// , - whitespace
// todo: map namespace syntax
// #{...} - set
// todo: deftype, defrecord, and constructor calls
// ; - comment
// @ - deref, @form -> (deref form)
// ^ - metadata
//     ^Type → ^{:tag Type}
//     ^:key → ^{:key true}
//       e.g., ^:dynamic ^:private ^:doc ^:const
// ~ - unquote
// ~@ - unquote splicing
//
// #{} - set
// #" - rx pattern
// #' - #'x -> (var x)
// #(...) - (fn [args] (...)), args %, %n, %&
// #_ - ignore next form
//
// ` - for forms other than sym, list, vec, set, map, `x is 'x
//     todo
//
// skipping tagged literals, reader conditional, splicing reader conditional
//
// {...} - map
//

// special forms:
// def
// if
// do
// let
// quote
// var
// fn
// loop
// recur
// throw, try, monitor-enter, monitor-exit, new, set!
//
// - destructuring
//
// & - rest

var tokRegex = regexp.MustCompile(
	// ignore whitespace
	`[\s,]*` +
		// match ~@ explicitly
		"(~@|" +
		"[\\[\\]{}()~^@`'&]|" + // match special chars
		// match strings, requiring escapes to be followed by char
		`"(?:\\.|[^\\"])*"?|` +
		// match comments
		`;.*|` +
		// capture token chars
		`[^\s\[\]{}()'"` + "`,;]*)",
)
var numRegex = regexp.MustCompile(`^-?[0-9]+$`)
var strRegex = regexp.MustCompile(`^"(?:\\.|[^\\"])*"?$`)

type tokenizer struct {
	reader *bufio.Reader

	buffer   string
	parsePos int

	pos Position

	tk  Token
	err error
}

func NewTokenizer(r io.Reader, file string) Tokenizer {
	return &tokenizer{
		reader:   bufio.NewReader(r),
		parsePos: -1,
		pos:      Position{Line: 1, File: file},
	}
}

func (tz *tokenizer) Position() Position {
	return tz.pos
}

func (tz *tokenizer) beginRead() {
	if tz.parsePos == -1 {
		// TODO: ignoring err
		b, _ := ioutil.ReadAll(tz.reader)
		tz.buffer = string(b)
		tz.parsePos = 0
	}

}

func (tz *tokenizer) handleToken(s string) (tk Token) {
	tk.s = s

	switch s {
	case "(":
		tk.t = TokenLParen
	case ")":
		tk.t = TokenRParen
	case "[":
		tk.t = TokenLBracket
	case "]":
		tk.t = TokenRBracket
	case "{":
		tk.t = TokenLCurly
	case "}":
		tk.t = TokenRCurly
	case "'":
		tk.t = TokenQuote
	case "nil":
		tk.t = TokenNil
	case "~", "^", "@", "`", "~@":
		// TODO
	default:
		if numRegex.MatchString(s) {
			tk.t = TokenNum
		} else if strRegex.MatchString(s) {
			tk.t = TokenString
		} else {
			tk.t = TokenSymbol
		}
	}

	return
}

func (tz *tokenizer) Next() bool {
	tz.beginRead()

	for tz.parsePos < len(tz.buffer) {
		if match := tokRegex.FindStringSubmatch(
			tz.buffer[tz.parsePos:]); match != nil {
			// inc by whole matched string
			tz.parsePos += len(match[0])

			// check on our capture
			capture := match[1]
			if capture != "" && capture[0] != ';' {
				tk := tz.handleToken(capture)
				// return true if we got a token
				if tk.t != TokenNone {
					tz.tk = tk
					return true
				}
			}
		}
	}

	if tz.parsePos >= len(tz.buffer) {
		tz.tk = Token{t: TokenEOF}
	}
	/*
		// return 'next' if we have one
		if tz.tkNext.t != TokenNone {
			// copy to curr and reset
			tz.tkCurr = tz.tkNext
			tz.tkNext.t = TokenNone
			return true
		} else {
			// otherwise, hunt for tokens

			for tz.handleRune != nil {
				// get next rune
				if r, ok := tz.readRune(); ok {
					tk := tz.handleRune(tz, r)

					// return true if we got a token
					if tk.t != TokenNone {
						tz.tkCurr = tk
						return true
					}
				}
			}
		}
	*/

	return false
}

func (tz *tokenizer) Value() Token {
	return tz.tk
}

func (tz *tokenizer) Err() error {
	return tz.err
}
