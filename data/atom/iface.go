package atom

import (
	"fmt"
)

// Interface for all atoms on the lisp side.
type Atom interface {
	fmt.Stringer

	IsNil() bool
	// Returns a hash value for this Atom.
	Hash() uint32
	// Returns whether this Atom is equivalent to a given atom.
	Equals(atom Atom) bool
}

// Interface for iterating a sequence of atoms.
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
