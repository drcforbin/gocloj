package atom

import (
	"strings"
)

type Str struct {
	Val string
}

func (s Str) String() string {
	var builder strings.Builder
	builder.WriteString("\"")
	builder.WriteString(s.Val)
	builder.WriteString("\"")
	return builder.String()
}

func (s Str) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (s Str) Hash() uint32 {
	return hashString(s.Val)
}

// Returns whether this Atom is equivalent to a given atom.
func (s Str) Equals(atom Atom) bool {
	if val, ok := atom.(*Str); ok {
		return s.Val == val.Val
	}

	return false
}

type strIterator struct {
	str []rune
	idx int
}

func (s Str) Iterator() SeqIterator {
	return &strIterator{str: []rune(s.Val), idx: -1}
}

func (it *strIterator) Next() bool {
	if it != nil && it.str != nil && len(it.str) > 0 &&
		it.idx < len(it.str)-1 {
		it.idx++
		return true
	}

	return false
}

func (it *strIterator) Value() Atom {
	if it != nil && it.str != nil && len(it.str) > 0 &&
		it.idx < len(it.str) {
		return &Char{Val: it.str[it.idx]}
	}

	return nil
}
