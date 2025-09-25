package hashmap

import (
	"gocloj/data/atom"
)

func iterMakePairs(key atom.Atom, val atom.Atom) atom.Atom {
	vec := atom.NewVec()
	vec.Items = append(vec.Items, key, val)
	return vec
}

type emptyMapIterator struct{}

func (it *emptyMapIterator) Next() bool {
	return false
}

func (it *emptyMapIterator) Value() atom.Atom {
	// TODO: panic
	return atom.Nil
}

type rootIterator struct {
	handler iterHandler

	seen   bool
	nilVal atom.Atom

	nestedIt atom.SeqIterator
}

func (it *rootIterator) Next() bool {
	if !it.seen {
		it.seen = true
		return true
	}

	if it.nestedIt != nil {
		if it.nestedIt.Next() {
			return true
		} else {
			it.nestedIt = nil
			it.nilVal = atom.Nil
		}
	}

	return false
}

// TODO: make sure to test iterating with a nullvalue

func (it *rootIterator) Value() atom.Atom {
	if it.nestedIt != nil {
		return it.nestedIt.Value()
	} else {
		return it.handler(atom.Nil, it.nilVal)
	}
}

type arrayNodeIterator struct {
	handler iterHandler

	idx   int
	array []phmNode

	nestedIt atom.SeqIterator
}

func (it *arrayNodeIterator) Next() bool {
	for {
		if it.nestedIt != nil {
			if it.nestedIt.Next() {
				return true
			} else {
				it.nestedIt = nil
			}
		}

		if it.idx < len(it.array) {
			node := it.array[it.idx]
			it.idx++

			if node != nil {
				iter := node.iterator(it.handler)
				if iter != nil && iter.Next() {
					it.nestedIt = iter
					return true
				}
			}
		} else {
			return false
		}
	}
}

func (it *arrayNodeIterator) Value() atom.Atom {
	if it.nestedIt != nil {
		return it.nestedIt.Value()
	} else {
		return atom.Nil
	}
}

type bitmapIndexedNodeIterator struct {
	handler iterHandler

	idx   int
	array []mapEntryI

	nestedIt atom.SeqIterator
	currIdx  int
}

func (it *bitmapIndexedNodeIterator) Next() bool {
	if it.nestedIt != nil {
		if it.nestedIt.Next() {
			return true
		} else {
			it.nestedIt = nil
		}
	}

	for it.idx < len(it.array) {
		idx := it.idx
		it.idx++

		entry := it.array[idx]
		if entry.key == nil {
			iter := entry.val.(phmNode).iterator(it.handler)
			if iter != nil && iter.Next() {
				it.nestedIt = iter
				it.currIdx = -1
				return true
			}
		} else {
			it.nestedIt = nil
			it.currIdx = idx
			return true
		}
	}

	it.currIdx = -1
	it.nestedIt = nil
	return false
}

func (it *bitmapIndexedNodeIterator) Value() atom.Atom {
	if it.nestedIt != nil {
		return it.nestedIt.Value()
	} else {
		if it.currIdx != -1 {
			entry := it.array[it.currIdx]
			return it.handler(entry.key.(atom.Atom), entry.val.(atom.Atom))
		} else {
			return atom.Nil
		}
	}
}

type bitmapIndexedNode2Iterator struct {
	handler iterHandler

	idx      int
	array    [32]mapEntryI
	arrayLen int

	nestedIt atom.SeqIterator
	currIdx  int
}

func (it *bitmapIndexedNode2Iterator) Next() bool {
	if it.nestedIt != nil {
		if it.nestedIt.Next() {
			return true
		} else {
			it.nestedIt = nil
		}
	}

	for it.idx < it.arrayLen {
		idx := it.idx
		it.idx++

		entry := it.array[idx]
		if entry.key == nil {
			iter := entry.val.(phmNode).iterator(it.handler)
			if iter != nil && iter.Next() {
				it.nestedIt = iter
				it.currIdx = -1
				return true
			}
		} else {
			it.nestedIt = nil
			it.currIdx = idx
			return true
		}
	}

	it.currIdx = -1
	it.nestedIt = nil
	return false
}

func (it *bitmapIndexedNode2Iterator) Value() atom.Atom {
	if it.nestedIt != nil {
		return it.nestedIt.Value()
	} else {
		if it.currIdx != -1 {
			entry := it.array[it.currIdx]
			return it.handler(entry.key.(atom.Atom), entry.val.(atom.Atom))
		} else {
			return atom.Nil
		}
	}
}

type hashCollisionNodeIterator struct {
	handler iterHandler

	idx     int
	array   []mapEntry
	currIdx int
}

func (it *hashCollisionNodeIterator) Next() bool {
	if it.idx < len(it.array) {
		it.currIdx = it.idx
		it.idx++

		return true
	}

	it.currIdx = -1
	return false
}

func (it *hashCollisionNodeIterator) Value() atom.Atom {
	if it.currIdx != -1 {
		entry := it.array[it.currIdx]
		return it.handler(entry.key, entry.val)
	} else {
		return atom.Nil
	}
}
