package data

import (
	"errors"
	"fmt"
	"gocloj/data/atom"
)

func Truthy(a atom.Atom) bool {
	if a == atom.False ||
		a == atom.Nil ||
		a.IsNil() {
		return false
	}

	return true
}

func ValidBinding(binding atom.Atom) (err error) {
	// TODO: handle vec and map
	switch v := binding.(type) {
	case *atom.SymName, *atom.Keyword:
		// pass, leaving err nil

	case *atom.Vec:
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

func Equals(a atom.Atom, b atom.Atom) bool {
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
