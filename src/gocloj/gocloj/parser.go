package gocloj

import (
	"gocloj/data"
	"gocloj/log"
	"math/big"
)

var parserLogger = log.Get("parser")

// special value indicating that we have not
// yet found a token
var noToken = &data.Const{Name: "notoken"}

// TODO: macros

type AtomIterator interface {
	Next() bool
	Value() (val data.Atom)
	Err() error
}

type parseFrame struct {
	seq data.Atom
	// used when building a data.List
	tail *data.ListNode

	quoteNext bool
}

func newParseFrame(seq data.Atom) parseFrame {
	return parseFrame{
		seq: seq,
	}
}

func (f *parseFrame) push(val data.Atom) {
	switch s := f.seq.(type) {
	case *data.List:
		node := &data.ListNode{Value: val}
		if f.tail != nil {
			// add to end
			f.tail.Next = node
		} else {
			// first item
			s.Head = node
		}
		f.tail = node
	case *data.Vec:
		s.Items = append(s.Items, val)
	default:
		// TODO: log
	}
}

type Parser struct {
	tz        Tokenizer
	stack     []parseFrame
	quoteNext bool

	value data.Atom

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

func (p *Parser) pushAtom(atom data.Atom) (err error) {
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

func quoteAtom(atom data.Atom) data.Atom {
	list := &data.List{
		Head: &data.ListNode{
			Value: &data.SymName{Name: "quote"},
			Next: &data.ListNode{
				Value: atom,
			},
		},
	}
	return list
}

func (p *Parser) endFrame() (err error) {
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
		frame := newParseFrame(data.NewList())
		p.stack = append(p.stack, frame)
	case TokenLBracket:
		frame := newParseFrame(data.NewVec())
		p.stack = append(p.stack, frame)
	case TokenLCurly:
		frame := newParseFrame(data.NewHashMap())
		p.stack = append(p.stack, frame)
	case TokenRParen, TokenRBracket, TokenRCurly:
		err = p.endFrame()
	case TokenSymbol:
		val := &data.SymName{Name: t.s}
		err = p.pushAtom(val)
	case TokenString:
		val := &data.Str{Val: t.s}
		err = p.pushAtom(val)
	case TokenChar:
		// we know it's a single rune from the
		// tokenizer; this is safe
		val := &data.Char{Val: []rune(t.s)[0]}
		err = p.pushAtom(val)
	case TokenNil:
		val := data.Nil
		err = p.pushAtom(val)
	case TokenNum:
		// note: only handling base 10
		i := &big.Int{}
		i.SetString(t.s, 10)

		val := &data.Num{Val: i}

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

func (p *Parser) Value() (val data.Atom) {
	if p.err == nil {
		val = p.value
		p.value = noToken
	}
	return
}

func (p *Parser) Err() error {
	return p.err
}
