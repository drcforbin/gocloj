package gocloj

import (
	"gocloj/data/atom"
	"gocloj/data/hashmap"
	"gocloj/log"
	"math/big"
)

var parserLogger = log.Get("parser")

// special value indicating that we have not
// yet found a token
var noToken = &atom.Const{Name: "notoken"}

// TODO: macros

type AtomIterator interface {
	Next() bool
	Value() (val atom.Atom)
	Err() error
}

type parseFrame struct {
	seq atom.Atom
	// used when building a atom.List
	tail *atom.ListNode
	// used when building a map
	key atom.Atom

	quoteNext bool
}

func newParseFrame(seq atom.Atom) parseFrame {
	return parseFrame{
		seq: seq,
	}
}

func (f *parseFrame) push(val atom.Atom) {
	switch s := f.seq.(type) {
	case *atom.List:
		node := &atom.ListNode{Value: val}
		if f.tail != nil {
			// add to end
			f.tail.Next = node
		} else {
			// first item
			s.Head = node
		}
		f.tail = node
	case *atom.Vec:
		s.Items = append(s.Items, val)
	case *hashmap.PersistentHashMap:
		if f.key == nil {
			f.key = val
		} else {
			key := f.key
			f.key = nil
			f.seq = s.Assoc(key, val)
		}
	default:
		// TODO: log
	}
}

type Parser struct {
	tz        Tokenizer
	stack     []parseFrame
	quoteNext bool

	value atom.Atom

	err error
}

func NewParser(tz Tokenizer) *Parser {
	p := &Parser{
		tz:    tz,
		stack: []parseFrame{},
		value: noToken,
	}
	return p
}

func (p *Parser) topFrame() (frame *parseFrame) {
	if len(p.stack) > 0 {
		frame = &p.stack[len(p.stack)-1]
	}

	return
}

func (p *Parser) pushAtom(atom atom.Atom) (err error) {
	if frame := p.topFrame(); frame != nil {
		// wrap in quote if needed
		if frame.quoteNext {
			atom = quoteAtom(atom)
			frame.quoteNext = false
		}

		frame.push(atom)
	} else {
		if p.quoteNext {
			atom = quoteAtom(atom)
			p.quoteNext = false
		}

		p.value = atom
	}

	return
}

func (p *Parser) parseError(msg string) error {
	pos := p.tz.Position()
	return NewError(msg, pos.Line, pos.File)
}

// TODO: rename quoteVal
func quoteAtom(val atom.Atom) atom.Atom {
	list := &atom.List{
		Head: &atom.ListNode{
			Value: &atom.SymName{Name: "quote"},
			Next: &atom.ListNode{
				Value: val,
			},
		},
	}
	return list
}

func (p *Parser) endFrame() (err error) {
	// TODO: error on map with leftover keys!

	switch len(p.stack) {
	case 0:
		err = p.parseError("unhandled empty stack on RParen")

	case 1:
		frame := p.topFrame()
		seq := frame.seq
		if p.quoteNext {
			seq = quoteAtom(seq)
			p.quoteNext = false
		} else if frame.quoteNext {
			// TODO: error
		}
		p.value = seq

		// clear stack
		p.stack = p.stack[:0]

	default: // > 1
		frame := p.topFrame()
		seq := frame.seq

		// pop it
		p.stack = p.stack[:len(p.stack)-1]

		frame = p.topFrame()
		// prefix quote if needed
		if frame.quoteNext {
			seq = quoteAtom(seq)
			frame.quoteNext = false
		}

		frame.push(seq)
	}

	return
}

func (p *Parser) handleToken(t Token) (err error) {
	switch t.t {
	case TokenLParen:
		frame := newParseFrame(atom.NewList())
		p.stack = append(p.stack, frame)
	case TokenLBracket:
		frame := newParseFrame(atom.NewVec())
		p.stack = append(p.stack, frame)
	case TokenLCurly:
		frame := newParseFrame(hashmap.NewPersistentHashMap())
		p.stack = append(p.stack, frame)
	case TokenRParen, TokenRBracket, TokenRCurly:
		err = p.endFrame()
	case TokenSymbol:
		val := &atom.SymName{Name: t.s}
		err = p.pushAtom(val)
	case TokenString:
		val := &atom.Str{Val: t.s}
		err = p.pushAtom(val)
	case TokenKeyword:
		val := &atom.Keyword{Name: t.s}
		err = p.pushAtom(val)
	case TokenChar:
		// we know it's a single rune from the
		// tokenizer; this is safe
		val := &atom.Char{Val: []rune(t.s)[0]}
		err = p.pushAtom(val)
	case TokenNil:
		val := atom.Nil
		err = p.pushAtom(val)
	case TokenNum:
		// note: only handling base 10
		i := &big.Int{}
		i.SetString(t.s, 10)

		val := &atom.Num{Val: i}

		err = p.pushAtom(val)
	case TokenQuote:
		if frame := p.topFrame(); frame != nil {
			frame.quoteNext = true
		} else {
			p.quoteNext = true
		}
	default:
		parserLogger.Warnf("unexpected token in parse %d", t.t)
	}

	return
}

func (p *Parser) Next() bool {
	if p.err != nil {
		return false
	}

	for p.tz.Next() {
		if p.err = p.handleToken(p.tz.Value()); p.err != nil {
			return false
		}

		// TODO: stop on EOF?
		if p.value != noToken {
			return true
		}
	}

	return false
}

func (p *Parser) Value() (val atom.Atom) {
	if p.err == nil {
		val = p.value
		p.value = noToken
	}
	return
}

func (p *Parser) Err() error {
	return p.err
}
