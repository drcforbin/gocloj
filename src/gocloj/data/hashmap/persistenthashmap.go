package hashmap

import (
	"gocloj/data/atom"
	"math"
	"strings"
)

// PersistentHashMap is, well, a persistent hash map. Changes
// made to this map will result in the creation of a new map,
// reusing whatever nodes from the previous map that can be.
type PersistentHashMap struct {
	count int
	root  phmNode

	hasNil bool
	nilVal atom.Atom
}

// Empty bitmapIndexedNode, used to create subsequent nodes.
// This saves creating a new node every time; as changes are
// made to this node, copies are returned rather than modifying
// this node itself.
var emptyBin = &bitmapIndexedNode{
	edit:  math.MaxUint64, // max to avoid collisions
	array: []mapEntryI{},
}

func NewPersistentHashMap() PersistentMap {
	return &PersistentHashMap{}
}

func (p *PersistentHashMap) String() string {
	var builder strings.Builder

	builder.WriteString("{")

	cnt := 0
	it := p.iteratorCore(func(key atom.Atom, val atom.Atom) atom.Atom {
		if cnt != 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(key.String())
		builder.WriteString(" ")
		builder.WriteString(val.String())

		cnt++

		return atom.Nil
	})

	// walk collection for side effects, we don't care about return
	for it.Next() {
		it.Value()
	}

	builder.WriteString("}")

	return builder.String()
}

func (p *PersistentHashMap) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (p *PersistentHashMap) Hash() uint32 {
	return atom.SeqHash(p)
}

// Returns whether this Atom is equivalent to a given atom.
func (p *PersistentHashMap) Equals(a atom.Atom) bool {
	if val, ok := a.(atom.Seq); ok {
		// TODO: this is not correct
		// iteration order should be assumed to be indeterminate
		return atom.SeqEquals(p, val)
	}

	return false
}

func (p *PersistentHashMap) iteratorCore(handler iterHandler) atom.SeqIterator {
	var it atom.SeqIterator
	if p.root == nil {
		it = &emptyMapIterator{}
	} else {
		it = p.root.iterator(handler)
	}

	if p.hasNil {
		it = &rootIterator{
			handler:  handler,
			nilVal:   p.nilVal,
			nestedIt: it,
		}
	}

	return it
}

func (p *PersistentHashMap) Iterator() atom.SeqIterator {
	return p.iteratorCore(iterMakePairs)
}

func (p *PersistentHashMap) Length() int {
	return p.count
}

func (p *PersistentHashMap) Get(key atom.Atom) atom.Atom {
	// TODO: handle Nil?
	if p.root != nil {
		hash := key.Hash()
		entry := p.root.find(0, hash, key)
		if entry != nil {
			return entry.val
		}
	}

	return atom.Nil
}

// Returns a transient (but independent) version of the map. The
// transient map can be used to efficiently modify a map, without
// affecting the original. edit should be a nonzero number that
// differentiates the transient from other transients created from
// this persistent map (each transient created from a persistent map
// should have its own edit value).
func (p *PersistentHashMap) AsTransient(edit uint64) TransientMap {
	return &TransientHashMap{
		edit:   edit,
		count:  p.count,
		root:   p.root,
		hasNil: p.hasNil,
		nilVal: p.nilVal,
	}
}

// Returns a new PersistentHashMap, adding or replacing a value in the map.
// atom.Nil may be used as a key.
func (p *PersistentHashMap) Assoc(key atom.Atom, val atom.Atom) PersistentMap {
	if key.Equals(atom.Nil) {
		// don't change it if val is the same
		if p.hasNil {
			if p.nilVal.Equals(val) {
				return p
			}
		}

		// add nil value
		return &PersistentHashMap{
			count:  p.count + 1,
			root:   p.root,
			hasNil: true,
			nilVal: val,
		}
	}

	// if we don't have a root, add one
	newroot := p.root
	if newroot == nil {
		newroot = emptyBin
	}

	// add the new item
	addedLeaf := false
	newroot = newroot.assoc(0, key.Hash(), key, val, &addedLeaf)
	if newroot == p.root {
		return p
	}

	var count = p.count
	if addedLeaf {
		count++
	}

	return &PersistentHashMap{
		count:  count,
		root:   newroot,
		hasNil: p.hasNil,
		nilVal: p.nilVal,
	}
}

// Returns a new PersistentHashMap, removing a value from the map.
// atom.Nil may be used as a key.
func (p *PersistentHashMap) Without(key atom.Atom) PersistentMap {
	if key.Equals(atom.Nil) {
		if p.hasNil {
			return &PersistentHashMap{
				count:  p.count - 1,
				root:   p.root,
				hasNil: false,
			}
		} else {
			return p
		}
	} else if p.root == nil {
		return p
	} else {
		// ok to ignore here; without call will return a new
		// root if a leaf was removed; this flag is redundant
		// todo: remove removedLeaf
		removedLeaf := false
		newroot := p.root.without(0, key.Hash(), key, &removedLeaf)
		if newroot == p.root {
			return p
		}
		return &PersistentHashMap{
			count:  p.count - 1,
			root:   newroot,
			hasNil: p.hasNil,
			nilVal: p.nilVal,
		}
	}
}
