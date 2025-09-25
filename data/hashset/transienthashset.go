package hashset

import (
	"gocloj/data/atom"
	"gocloj/data/hashmap"
)

// TransientHashSet is a transient version of PersistentHashSet.
// Calling PersistentHashSet.AsTransient will result in creation
// of one of these, and it can be converted back to a PersistentHashSet
// by calling TransientHashSet.AsPersistent
type TransientHashSet struct {
	impl hashmap.TransientMap
}

// Returns a persistent version of the map. After this is called,
// the transient map should not be used again.
func (t *TransientHashSet) AsPersistent() PersistentSet {
	// TODO: ensureEditable
	return &PersistentHashSet{
		impl: t.impl.AsPersistent(),
	}
}

func (t *TransientHashSet) String() string {
	// TODO:
	return "THM"
}

func (t *TransientHashSet) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (t *TransientHashSet) Hash() uint32 {
	return t.impl.Hash()
}

// Returns whether this Atom is equivalent to a given atom.
func (t *TransientHashSet) Equals(atom atom.Atom) bool {
	return t.impl.Equals(atom)
}

func (t *TransientHashSet) Iterator() atom.SeqIterator {
	// TODO:
	return nil
}

func (t *TransientHashSet) Length() int {
	return t.impl.Length()
}

func (t *TransientHashSet) Get(val atom.Atom) atom.Atom {
	return t.impl.Get(val)
}

func (t *TransientHashSet) Assoc(val atom.Atom) TransientSet {
	newImpl := t.impl.Assoc(val, val)
	if newImpl == t.impl {
		return t
	}

	return &TransientHashSet{
		impl: newImpl,
	}
}

func (t *TransientHashSet) Without(val atom.Atom) TransientSet {
	newImpl := t.impl.Without(val)
	if newImpl == t.impl {
		return t
	}

	return &TransientHashSet{
		impl: newImpl,
	}
}
