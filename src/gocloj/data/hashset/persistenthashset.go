package hashset

import (
	"gocloj/data/atom"
	"gocloj/data/hashmap"
	"strings"
)

// PersistentHashSet is a persistent hash set based on the
// persistent has set. Changes made to this set will result
// in the creation of a new set.
type PersistentHashSet struct {
	impl hashmap.PersistentMap
}

func NewPersistentHashSet() PersistentSet {
	return &PersistentHashSet{
		impl: hashmap.NewPersistentHashMap(),
	}
}

func (s *PersistentHashSet) String() string {
	var builder strings.Builder

	builder.WriteString("#{")

	cnt := 0
	it := s.impl.Iterator()
	for it.Next() {
		if cnt != 0 {
			builder.WriteString(" ")
		}

		// TODO: proper map entry atom type
		pair := it.Value()
		builder.WriteString(pair.(*atom.Vec).Items[0].String())

		cnt++
	}

	builder.WriteString("}")

	return builder.String()
}

func (s *PersistentHashSet) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (s *PersistentHashSet) Hash() uint32 {
	return atom.SeqHash(s.impl)
}

// Returns whether this Atom is equivalent to a given atom.
func (s *PersistentHashSet) Equals(a atom.Atom) bool {
	if val, ok := a.(atom.Seq); ok {
		// TODO: this is not correct
		// iteration order should be assumed to be indeterminate
		return atom.SeqEquals(s.impl, val)
	}

	return false
}

func (s *PersistentHashSet) Iterator() atom.SeqIterator {
	// TODO...set iterator
	return nil
}

func (s *PersistentHashSet) Length() int {
	return s.impl.Length()
}

func (s *PersistentHashSet) Get(val atom.Atom) atom.Atom {
	return s.impl.Get(val)
}

// Returns a transient (but independent) version of the set. The
// transient set can be used to efficiently modify a set, without
// affecting the original. edit should be a nonzero number that
// differentiates the transient from other transients created from
// this persistent set (each transient created from a persistent set
// should have its own edit value).
func (s *PersistentHashSet) AsTransient(edit uint64) TransientSet {
	return &TransientHashSet{
		impl: s.impl.AsTransient(edit),
	}
}

// Returns a new PersistentHashSet, adding a value to the set. atom.Nil may be
// used as a value.
func (s *PersistentHashSet) Assoc(val atom.Atom) PersistentSet {
	newImpl := s.impl.Assoc(val, val)
	if newImpl == s.impl {
		return s
	}

	return &PersistentHashSet{
		impl: newImpl,
	}
}

// Returns a new PersistentHashSet, removing a value from the set.
// atom.Nil may be used as a value.
func (s *PersistentHashSet) Without(val atom.Atom) PersistentSet {
	newImpl := s.impl.Without(val)
	if newImpl == s.impl {
		return s
	}

	return &PersistentHashSet{
		impl: newImpl,
	}
}
