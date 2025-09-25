package hashmap

import (
	"gocloj/data/atom"
)

type Map interface {
	// there has to be a better way than this.
	// sure wish we had generics or something
	// to make this interface smaller
	atom.Atom
	atom.Seq

	Length() int

	Get(key atom.Atom) atom.Atom
}

type PersistentMap interface {
	Map

	AsTransient(edit uint64) TransientMap

	Assoc(key atom.Atom, val atom.Atom) PersistentMap
	Without(key atom.Atom) PersistentMap
}

type TransientMap interface {
	Map

	AsPersistent() PersistentMap

	Assoc(key atom.Atom, val atom.Atom) TransientMap
	Without(key atom.Atom) TransientMap
}

type iterHandler func(key atom.Atom, val atom.Atom) atom.Atom

type phmNode interface {
	assoc(shift uint32, hash uint32, key atom.Atom, val atom.Atom, addedLeaf *bool) phmNode
	without(shift uint32, hash uint32, key atom.Atom, removedLeaf *bool) phmNode
	find(shift uint32, hash uint32, key atom.Atom) *mapEntry

	// transient funcs
	assocT(edit uint64, shift uint32, hash uint32, key atom.Atom, val atom.Atom, addedLeaf *bool) phmNode
	withoutT(edit uint64, shift uint32, hash uint32, key atom.Atom, removedLeaf *bool) phmNode

	iterator(handler iterHandler) atom.SeqIterator
}
