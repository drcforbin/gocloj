package hashmap

import (
	"gocloj/data/atom"
)

// TransientHashMap is a transient version of PersistentHashMap.
// Calling PersistentHashMap.AsTransient will result in creation
// of one of these, and it can be converted back to a PersistentHashMap
// by calling TransientHashMap.AsPersistent
type TransientHashMap struct {
	edit     uint64
	count    int
	root     phmNode
	hasNil   bool
	nilVal   atom.Atom
	leafFlag bool
}

// Returns a persistent version of the map. After this is called,
// the transient map should not be used again.
func (t *TransientHashMap) AsPersistent() PersistentMap {
	t.edit = 0

	// TODO: ensureEditable
	return &PersistentHashMap{
		count:  t.count,
		root:   t.root,
		hasNil: t.hasNil,
		nilVal: t.nilVal,
	}
}

func (t *TransientHashMap) String() string {
	// TODO:
	return "THM"
}

func (t *TransientHashMap) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (t *TransientHashMap) Hash() uint32 {
	// TODO:
	return 0
}

// Returns whether this Atom is equivalent to a given atom.
func (t *TransientHashMap) Equals(atom atom.Atom) bool {
	// TODO:
	return false
}

func (t *TransientHashMap) Iterator() atom.SeqIterator {
	// TODO:
	return nil
}

func (t *TransientHashMap) Length() int {
	return t.count
}

func (t *TransientHashMap) Get(key atom.Atom) atom.Atom {
	// TODO: remove
	if t.root != nil {
		hash := key.Hash()
		entry := t.root.find(0, hash, key)
		if entry != nil {
			return entry.val
		}
	}

	return atom.Nil
}

func (t *TransientHashMap) Assoc(key atom.Atom, val atom.Atom) TransientMap {
	t.ensureEditable()

	if key.Equals(atom.Nil) {
		t.nilVal = val
		if !t.hasNil {
			t.count++
			t.hasNil = true
		}
	} else {
		// if we don't have a root, add one
		newroot := t.root
		if newroot == nil {
			newroot = emptyBin
		}

		// add the new item
		addedLeaf := false
		newroot = newroot.assocT(t.edit, 0, key.Hash(), key, val, &addedLeaf)
		if newroot != t.root {
			t.root = newroot
		}

		if addedLeaf {
			t.count++
		}
	}
	return t
}

func (t *TransientHashMap) Without(key atom.Atom) TransientMap {
	t.ensureEditable()

	if key.Equals(atom.Nil) {
		if !t.hasNil {
			return t
		}
		t.hasNil = false
		t.nilVal = nil
		t.count--
		return t
	}
	if t.root == nil {
		return t
	}

	removedLeaf := false
	n := t.root.withoutT(t.edit, 0, key.Hash(), key, &removedLeaf)
	if n != t.root {
		t.root = n
	}

	if removedLeaf {
		t.count--
	}
	return t
}

func (t *TransientHashMap) ensureEditable() {
	// TODO: error?
	// if t.edit == 0 { panic? }
}
