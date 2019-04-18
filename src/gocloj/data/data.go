package data

import (
	"fmt"
	"gocloj/log"
	"math/big"
	"strings"
)

var dataLogger = log.Get("data")

var Nil = &Const{Name: "nil"}
var True = &Const{Name: "true"}
var False = &Const{Name: "false"}

type Atom interface {
	fmt.Stringer
	IsNil() bool
	Hash() uint32
	Equals(atom Atom) bool
}

type Const struct {
	Name string
}

func (c Const) String() string {
	return c.Name
}

func (c Const) IsNil() bool {
	return false
}

func (c Const) Hash() uint32 {
	return hashString(c.Name)
}

func (c Const) Equals(atom Atom) bool {
	if val, ok := atom.(*Const); ok {
		return c.Name == val.Name
	}

	return false
}

type Char struct {
	Val rune
}

func (c Char) String() string {
	var builder strings.Builder
	builder.WriteString("\\")
	builder.WriteRune(c.Val)
	return builder.String()
}

func (c Char) IsNil() bool {
	return false
}

func (c Char) Hash() uint32 {
	return uint32(c.Val)
}

func (c Char) Equals(atom Atom) bool {
	if val, ok := atom.(*Char); ok {
		return c.Val == val.Val
	}

	return false
}

type Num struct {
	Val *big.Int
}

func (s Num) String() string {
	return s.Val.String()
}

func (n Num) IsNil() bool {
	return false
}

func (n Num) Hash() uint32 {
	hash := seed

	bytes := n.Val.Bytes()
	byteLen := len(bytes)

	for idx := 0; idx < byteLen; idx += 4 {
		val := uint32(0)

		chunkLen := (byteLen - idx) % 4
		switch chunkLen {
		case 1:
			val = uint32(bytes[idx+0])
		case 2:
			val = uint32(bytes[idx+0])<<8 |
				uint32(bytes[idx+1])
		case 3:
			val = uint32(bytes[idx+0])<<16 |
				uint32(bytes[idx+1])<<8 |
				uint32(bytes[idx+2])
		default: // bytes left >= 4
			val = uint32(bytes[idx+0])<<24 |
				uint32(bytes[idx+1])<<16 |
				uint32(bytes[idx+2])<<8 |
				uint32(bytes[idx+3])
		}

		hash = mixH1(hash, mixK1(val))
	}

	return fmix(hash, uint32(byteLen))
}

func (n Num) Equals(atom Atom) bool {
	if val, ok := atom.(*Num); ok {
		return n.Val.Cmp(val.Val) == 0
	}

	return false
}

type SymName struct {
	Name string
}

func (s SymName) String() string {
	if s.Name != "" {
		return s.Name
	} else {
		return "{}"
	}
}

func (s *SymName) IsNil() bool {
	return false
}

func (s SymName) Hash() uint32 {
	return hashString(s.Name)
}

func (s SymName) Equals(atom Atom) bool {
	if val, ok := atom.(*SymName); ok {
		return s.Name == val.Name
	}

	return false
}

type SeqIterator interface {
	Next() bool
	Value() Atom
}

type Seq interface {
	Iterator() SeqIterator
}

type Indexable interface {
	Length() int
	Item(idx int) Atom
}
