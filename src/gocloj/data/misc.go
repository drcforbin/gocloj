package data

import (
	"errors"
	"fmt"
)

func Truthy(atom Atom) bool {
	if atom == False ||
		atom == Nil ||
		atom.IsNil() {
		return false
	}

	return true
}

func ValidBinding(binding Atom) (err error) {
	// TODO: handle vec and map
	switch v := binding.(type) {
	case *SymName:
		// pass, leaving err nil

	case *Vec:
		for _, atom := range v.Items {
			if err = ValidBinding(atom); err != nil {
				return
			}
		}

	// TODO: destructure map

	default:
		err = errors.New(fmt.Sprintf(
			"unexpected type for destructure binding %T", binding))
	}

	return
}

/*
// -1 if x <  y
//  0 if x == y
// +1 if x >  y
func Compare(atom1 Atom, atom2 Atom) (cmp int) {
// NOTE: T > everything else
// NOTE: nil < everything else
}
*/

func seqEquals(ita SeqIterator, itb SeqIterator) bool {
	// walk a's
	for ita.Next() {
		// are we out of b's?
		if !itb.Next() {
			return false
		}

		if !Equals(ita.Value(), itb.Value()) {
			return false
		}
	}

	// do we still have more b's?
	if itb.Next() {
		return false
	}

	return true
}

func Equals(a Atom, b Atom) bool {
	// is one nil but not the other?
	aisnil := a.IsNil()
	bisnil := b.IsNil()
	if aisnil != bisnil {
		return false
	}

	// are they both nil? (checking
	// bisnil is redundant with last check)
	if aisnil {
		return true
	}

	return a.Equals(b)
}
