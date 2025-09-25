package hashmap

import (
	"gocloj/data/atom"
)

type hashCollisionNode struct {
	edit  uint64
	hash  uint32
	count int
	array []mapEntry
}

func cloneAndSetAtomVal(array []mapEntry, i int, a atom.Atom) []mapEntry {
	clone := make([]mapEntry, len(array))
	copy(clone, array)
	clone[i].val = a
	return clone
}

func removeMapEntry(array []mapEntry, i int) []mapEntry {
	clone := make([]mapEntry, len(array)-1)
	copy(clone, array[:i])
	copy(clone[i:], array[i+1:])
	return clone
}

func (hcn *hashCollisionNode) assoc(shift uint32, hash uint32, key atom.Atom, val atom.Atom, addedLeaf *bool) phmNode {
	// check the hash; if same, we can add it to a hcn
	if hash == hcn.hash {
		idx := hcn.findIndex(key)

		// if we found the key, replace its val
		if idx != -1 {
			if hcn.array[idx].val.Equals(val) {
				return hcn
			}

			return &hashCollisionNode{
				edit:  hcn.edit,
				hash:  hash,
				count: hcn.count,
				array: cloneAndSetAtomVal(hcn.array, idx, val),
			}
		}

		newArray := make([]mapEntry, hcn.count+1)
		copy(newArray, hcn.array)
		newArray[hcn.count] = mapEntry{key, val}
		*addedLeaf = true

		return &hashCollisionNode{
			hash:  hash,
			count: hcn.count + 1,
			array: newArray,
		}
	}

	// nest it in a bitmap node
	bin := &bitmapIndexedNode{
		bitmap: bitpos(hcn.hash, shift),
		array:  []mapEntryI{mapEntryI{nil, hcn}},
	}
	return bin.assoc(shift, hash, key, val, addedLeaf)
}

func (hcn *hashCollisionNode) without(shift uint32, hash uint32, key atom.Atom, removedLeaf *bool) phmNode {
	idx := hcn.findIndex(key)
	if idx == -1 {
		return hcn
	} else if idx == 1 {
		return nil
	} else {
		return &hashCollisionNode{
			hash:  hash,
			count: hcn.count - 1,
			array: removeMapEntry(hcn.array, idx),
		}
	}
}

func (hcn *hashCollisionNode) find(shift uint32, hash uint32, key atom.Atom) *mapEntry {
	idx := hcn.findIndex(key)
	if idx != -1 {
		entry := hcn.array[idx]
		return &entry
	}
	return nil
}

func (hcn *hashCollisionNode) assocT(edit uint64, shift uint32, hash uint32, key atom.Atom, val atom.Atom, addedLeaf *bool) phmNode {
	// check the hash; if same, we can add it to a hcn
	if hash == hcn.hash {
		idx := hcn.findIndex(key)

		// if we found the key, replace its val
		if idx != -1 {
			entry := hcn.array[idx]
			if entry.val.Equals(val) {
				return hcn
			}
			entry.val = val
			return hcn.editAndSet(edit, idx, entry)
		}

		if len(hcn.array) > hcn.count {
			*addedLeaf = true
			editable := hcn.editAndSet(edit, hcn.count, mapEntry{key, val})
			editable.count++
			return editable
		}

		newArray := make([]mapEntry, hcn.count+1)
		copy(newArray, hcn.array)
		newArray[hcn.count] = mapEntry{key, val}
		*addedLeaf = true

		return hcn.ensureEditableWithArray(edit, hcn.count+1, newArray)
	}

	// nest it in a bitmap node with an extra space
	bin := &bitmapIndexedNode{
		edit:   edit,
		bitmap: bitpos(hcn.hash, shift),
		array:  []mapEntryI{mapEntryI{nil, hcn}, mapEntryI{}},
	}
	return bin.assocT(edit, shift, hash, key, val, addedLeaf)
}

func (hcn *hashCollisionNode) withoutT(edit uint64, shift uint32, hash uint32, key atom.Atom, removedLeaf *bool) phmNode {
	idx := hcn.findIndex(key)
	if idx == -1 {
		return hcn
	}

	*removedLeaf = true
	if hcn.count == 1 {
		return nil
	}

	editable := hcn.ensureEditable(edit)
	editable.array[idx] = editable.array[len(editable.array)-1]
	editable.array[len(editable.array)-1] = mapEntry{}
	editable.count--
	return editable
}

func (hcn *hashCollisionNode) iterator(handler iterHandler) atom.SeqIterator {
	return &hashCollisionNodeIterator{
		handler: handler,
		array:   hcn.array,
	}
}

func (hcn *hashCollisionNode) ensureEditable(edit uint64) *hashCollisionNode {
	if hcn.edit == edit {
		return hcn
	}
	newArray := make([]mapEntry, hcn.count+1) // make room for next assoc
	copy(newArray, hcn.array)
	return &hashCollisionNode{
		edit:  edit,
		hash:  hcn.hash,
		count: hcn.count,
		array: newArray,
	}
}

func (hcn *hashCollisionNode) ensureEditableWithArray(edit uint64, count int, array []mapEntry) *hashCollisionNode {
	if hcn.edit == edit {
		hcn.count = count
		hcn.array = array
		return hcn
	}
	return &hashCollisionNode{
		edit:  edit,
		hash:  hcn.hash,
		count: count,
		array: array,
	}
}

func (hcn *hashCollisionNode) editAndSet(edit uint64, i int, entry mapEntry) *hashCollisionNode {
	editable := hcn.ensureEditable(edit)
	editable.array[i] = entry
	return editable
}

func (hcn *hashCollisionNode) findIndex(key atom.Atom) int {
	for i := 0; i < hcn.count; i++ {
		if key.Equals(hcn.array[i].key) {
			return i
		}
	}
	return -1
}
