package hashset

import (
	"gocloj/data/atom"
)

type Set interface {
	// there has to be a better way than this.
	// sure wish we had generics or something
	// to make this interface smaller
	atom.Atom
	atom.Seq

	Length() int

	Get(key atom.Atom) atom.Atom
}

type PersistentSet interface {
	Set

	AsTransient(edit uint64) TransientSet

	Assoc(key atom.Atom) PersistentSet
	Without(key atom.Atom) PersistentSet
}

type TransientSet interface {
	Set

	AsPersistent() PersistentSet

	Assoc(key atom.Atom) TransientSet
	Without(key atom.Atom) TransientSet
}
