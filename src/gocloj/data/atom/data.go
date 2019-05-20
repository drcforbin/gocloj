package atom

import (
	"gocloj/log"
	"math/big"
	"math/bits"
	"strings"
)

var dataLogger = log.Get("data")

type Const struct {
	Name string
}

func (c Const) String() string {
	return c.Name
}

func (c Const) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (c Const) Hash() uint32 {
	return hashString(c.Name)
}

// Returns whether this Atom is equivalent to a given atom.
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

// Returns a hash value for this Atom.
func (c Char) Hash() uint32 {
	return uint32(c.Val)
}

// Returns whether this Atom is equivalent to a given atom.
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

// Returns a hash value for this Atom.
func (n Num) Hash() uint32 {
	// not calling n.Val.Bytes to save allocation
	const wordBytes = bits.UintSize / 8
	var bytes [wordBytes]byte

	// NOTE: THIS DISCARDS ALL BUT LAST WORD!

	valWords := n.Val.Bits()
	for i := 0; i < len(valWords); i++ {
		valWord := valWords[i]
		for j := 0; j < wordBytes; j++ {
			bytes[j] = byte(valWord)
			valWord >>= 8
		}
	}

	return MurmurHash3(bytes[:], 0)
}

// Returns whether this Atom is equivalent to a given atom.
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

// Returns a hash value for this Atom.
func (s SymName) Hash() uint32 {
	return hashString(s.Name)
}

// Returns whether this Atom is equivalent to a given atom.
func (s SymName) Equals(atom Atom) bool {
	if val, ok := atom.(*SymName); ok {
		return s.Name == val.Name
	}

	return false
}

type Keyword struct {
	Name string
}

func (k Keyword) String() string {
	if k.Name != "" {
		return k.Name
	} else {
		return "nil"
	}
}

func (k *Keyword) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (k Keyword) Hash() uint32 {
	return hashString(k.Name)
}

// Returns whether this Atom is equivalent to a given atom.
func (k Keyword) Equals(atom Atom) bool {
	if val, ok := atom.(*Keyword); ok {
		return k.Name == val.Name
	}

	return false
}
