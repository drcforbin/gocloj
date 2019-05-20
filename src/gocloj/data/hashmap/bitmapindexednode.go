package hashmap

import (
	"gocloj/data/atom"
	"math/bits"
)

// extra space to add when growing during assocT
const assocTGrow = 4

type mapEntryI struct {
	key interface{}
	val interface{}
}

type bitmapIndexedNode struct {
	edit   uint64
	bitmap uint32
	// if key is nil, val is a node; if not, both are atoms
	array []mapEntryI
}

func mask(hash uint32, shift uint32) uint32 {
	return (hash >> shift) & 0x01f
}

func bitpos(hash uint32, shift uint32) uint32 {
	return 1 << mask(hash, shift)
}

func createNode(edit uint64, shift uint32, key1 atom.Atom, val1 atom.Atom, key2hash uint32,
	key2 atom.Atom, val2 atom.Atom) phmNode {
	key1hash := key1.Hash()
	if key1hash == key2hash {
		return &hashCollisionNode{
			// edit?
			hash:  key1hash,
			count: 2,
			array: []mapEntry{
				mapEntry{key1, val1},
				mapEntry{key2, val2},
			},
		}
	}

	addedLeaf := false
	return emptyBin.
		assocT(edit, shift, key1hash, key1, val1, &addedLeaf).
		assocT(edit, shift, key2hash, key2, val2, &addedLeaf)
}

func cloneAndSet(array []mapEntryI, i int, entry mapEntryI) []mapEntryI {
	clone := make([]mapEntryI, len(array))
	copy(clone, array)
	clone[i] = entry
	return clone
}

func removePair(array []mapEntryI, i int) []mapEntryI {
	clone := make([]mapEntryI, len(array)-1)
	copy(clone, array[:i])
	copy(clone[i:], array[i+1:])
	return clone
}

func (bin *bitmapIndexedNode) assoc(shift uint32, hash uint32,
	key atom.Atom, val atom.Atom, addedLeaf *bool) phmNode {
	bit := bitpos(hash, shift)
	idx := bin.index(bit)

	// is it maybe already present?
	if bin.bitmap&bit != 0 {
		// if entry.key is null, entry.val is a phmNode;
		// if not, both are Atom
		entry := bin.array[idx]

		if entry.key == nil {
			n := entry.val.(phmNode).assoc(shift+5, hash, key, val, addedLeaf)
			if n == entry.val {
				return bin
			}

			// make a new bitmapIndexedNode, setting the new node in the array
			entry.val = n
			return &bitmapIndexedNode{
				bitmap: bin.bitmap,
				array:  cloneAndSet(bin.array, idx, entry),
			}
		}

		// same key?
		if key.Equals(entry.key.(atom.Atom)) {
			if val == entry.val {
				return bin
			}

			// make a new bitmapIndexedNode, setting the item in the array
			entry.val = val
			return &bitmapIndexedNode{
				bitmap: bin.bitmap,
				array:  cloneAndSet(bin.array, idx, entry),
			}
		}

		// new item, rather than a replacement
		*addedLeaf = true
		return &bitmapIndexedNode{
			bitmap: bin.bitmap,
			array: cloneAndSet(bin.array,
				idx, mapEntryI{nil,
					createNode(bin.edit, shift+5,
						entry.key.(atom.Atom), entry.val.(atom.Atom), hash, key, val)}),
		}
	} else {
		// not present

		// if we have 16 or more values, promote the node to an arrayNode
		// containing new bitmapIndexedNodes wrapping each value
		n := bits.OnesCount32(bin.bitmap)
		if n >= 16 {
			jdx := mask(hash, shift)
			newNode := &arrayNode{count: n + 1}
			newNode.array[jdx] = emptyBin.assoc(shift+5, hash, key, val, addedLeaf)

			j := 0
			for i := 0; i < 32; i++ {
				if (bin.bitmap>>uint(i))&1 != 0 {
					// if entry.key is null, entry.val is a phmNode;
					// if not, both are Atom
					entry := bin.array[j]
					if entry.key == nil {
						newNode.array[i] = entry.val.(phmNode)
					} else {
						key := entry.key.(atom.Atom)
						newNode.array[i] = emptyBin.assoc(shift+5,
							key.Hash(), key,
							entry.val.(atom.Atom), addedLeaf)
					}
					j++
				}
			}

			return newNode
		} else {
			// insert at idx
			newArray := make([]mapEntryI, n+1)
			copy(newArray, bin.array[:idx])
			newArray[idx] = mapEntryI{key, val}
			copy(newArray[idx+1:], bin.array[idx:])
			*addedLeaf = true

			return &bitmapIndexedNode{
				bitmap: bin.bitmap | bit,
				array:  newArray,
			}
		}
	}
}

func (bin *bitmapIndexedNode) without(shift uint32, hash uint32, key atom.Atom, removedLeaf *bool) phmNode {
	bit := bitpos(hash, shift)
	if bin.bitmap&bit == 0 {
		return bin
	}
	idx := bin.index(bit)
	entry := bin.array[idx]
	if entry.key == nil {
		n := entry.val.(phmNode).without(shift+5, hash, key, removedLeaf)
		if n == entry.val {
			return bin
		}
		if n != nil {
			entry.val = n
			return &bitmapIndexedNode{
				bitmap: bin.bitmap,
				array:  cloneAndSet(bin.array, idx, entry),
			}
		}
		if bin.bitmap == bit {
			return nil
		}
		return &bitmapIndexedNode{
			bitmap: bin.bitmap ^ bit,
			array:  removePair(bin.array, idx),
		}
	}
	if key.Equals(entry.key.(atom.Atom)) {
		if bin.bitmap == bit {
			return nil
		}
		return &bitmapIndexedNode{
			bitmap: bin.bitmap ^ bit,
			array:  removePair(bin.array, idx),
		}
	}
	return bin
}

func (bin *bitmapIndexedNode) find(shift uint32, hash uint32, key atom.Atom) *mapEntry {
	bit := bitpos(hash, shift)
	if bin.bitmap&bit != 0 {
		idx := bin.index(bit)

		// if entry.key is null, entry.val is a phmNode;
		// if not, both are Atom
		entry := bin.array[idx]
		if entry.key == nil {
			return entry.val.(phmNode).find(shift+5, hash, key)
		} else {
			k := entry.key.(atom.Atom)
			v := entry.val.(atom.Atom)

			if key.Equals(k) {
				return &mapEntry{key: k, val: v}
			}
		}
	}

	return nil
}

func (bin *bitmapIndexedNode) assocT(edit uint64, shift uint32, hash uint32, key atom.Atom, val atom.Atom, addedLeaf *bool) phmNode {
	bit := bitpos(hash, shift)
	idx := bin.index(bit)

	// is it maybe already present?
	if bin.bitmap&bit != 0 {
		// if entry.key is null, entry.val is a phmNode;
		// if not, both are Atom
		entry := bin.array[idx]

		if entry.key == nil {
			n := entry.val.(phmNode).assocT(edit, shift+5, hash, key, val, addedLeaf)
			if n == entry.val {
				return bin
			}

			entry.val = n
			return bin.editAndSet(edit, idx, entry)
		}

		// same key?
		if key.Equals(entry.key.(atom.Atom)) {
			if val == entry.val {
				return bin
			} else {
				entry.val = val
				return bin.editAndSet(edit, idx, entry)
			}
		}

		// new item, rather than a replacement
		*addedLeaf = true
		return bin.editAndSet(edit, idx, mapEntryI{
			nil,
			createNode(edit, shift+5,
				entry.key.(atom.Atom), entry.val.(atom.Atom), hash, key, val),
		})
	} else {
		// not present

		n := bits.OnesCount32(bin.bitmap)

		// does array have space for this item?
		if n < len(bin.array) {
			editable := bin.ensureEditable(edit)

			// insert at idx
			copy(editable.array[idx+1:], editable.array[idx:])
			editable.array[idx] = mapEntryI{key, val}

			editable.bitmap |= bit
			*addedLeaf = true

			return editable
		}

		// if we have 16 or more values, promote the node to an arrayNode
		// containing a new bitmapIndexedNode with the val in it
		if n >= 16 {
			jdx := mask(hash, shift)
			newNode := &arrayNode{edit: edit, count: n + 1}
			newNode.array[jdx] = emptyBin.assocT(edit, shift+5, hash, key, val, addedLeaf)

			j := 0
			for i := 0; i < 32; i++ {
				if (bin.bitmap>>uint(i))&1 != 0 {
					// if entry.key is null, entry.val is a phmNode;
					// if not, both are Atom
					entry := bin.array[j]
					if entry.key == nil {
						newNode.array[i] = entry.val.(phmNode)
					} else {
						key := entry.key.(atom.Atom)
						newNode.array[i] = emptyBin.assocT(edit, shift+5,
							key.Hash(), key,
							entry.val.(atom.Atom), addedLeaf)
					}
					j++
				}
			}

			return newNode
		} else {
			// allocate array with extra spaces (assume the
			// transient will grow; the space will be dropped
			// when the transient is made permanent and modified),
			// and insert at idx
			newArray := make([]mapEntryI, n+assocTGrow)
			copy(newArray, bin.array[:idx])
			newArray[idx] = mapEntryI{key, val}
			copy(newArray[idx+1:], bin.array[idx:])

			*addedLeaf = true

			// ensureEditable does extra allocations, which we'll
			// skip here since we just filled one
			// editable := bin.ensureEditable(edit)
			editable := bin
			if bin.edit != edit {
				editable = &bitmapIndexedNode{
					edit:   edit,
					bitmap: bin.bitmap,
				}
			}

			editable.array = newArray
			editable.bitmap |= bit
			return editable
		}
	}
}

func (bin *bitmapIndexedNode) withoutT(edit uint64, shift uint32, hash uint32, key atom.Atom, removedLeaf *bool) phmNode {
	bit := bitpos(hash, shift)
	if bin.bitmap&bit == 0 {
		return bin
	}
	idx := bin.index(bit)
	entry := bin.array[idx]
	if entry.key == nil {
		n := entry.val.(phmNode).withoutT(edit, shift+5, hash, key, removedLeaf)
		if n == entry.val {
			return bin
		}
		if n != nil {
			entry.val = n
			return bin.editAndSet(edit, idx, entry)
		}
		if bin.bitmap == bit {
			return nil
		}
		return bin.editAndRemovePair(edit, bit, idx)
	}
	if key.Equals(entry.key.(atom.Atom)) {
		*removedLeaf = true
		// TODO: collapse
		return bin.editAndRemovePair(edit, bit, idx)
	}
	return bin
}

func (bin *bitmapIndexedNode) iterator(handler iterHandler) atom.SeqIterator {
	return &bitmapIndexedNodeIterator{
		handler: handler,
		array:   bin.array,
	}
}

func (bin *bitmapIndexedNode) index(bit uint32) int {
	return bits.OnesCount32(bin.bitmap & (bit - 1))
}

func (bin *bitmapIndexedNode) ensureEditable(edit uint64) *bitmapIndexedNode {
	if bin.edit == edit {
		return bin
	}
	n := bits.OnesCount32(bin.bitmap)
	var count int
	if n >= 0 {
		// make room for next assoc
		count = n + 1
	} else {
		// make room for first two assocs
		count = 2
	}
	newArray := make([]mapEntryI, count)
	copy(newArray, bin.array)
	return &bitmapIndexedNode{
		edit:   edit,
		bitmap: bin.bitmap,
		array:  newArray,
	}
	return bin
}

func (bin *bitmapIndexedNode) editAndSet(edit uint64, i int, entry mapEntryI) *bitmapIndexedNode {
	editable := bin.ensureEditable(edit)
	editable.array[i] = entry
	return editable
}

func (bin *bitmapIndexedNode) editAndRemovePair(edit uint64, bit uint32, i int) *bitmapIndexedNode {
	if bin.bitmap == bit {
		return nil
	}

	editable := bin.ensureEditable(edit)
	editable.bitmap ^= bit
	copy(editable.array[i:], editable.array[i+1:])
	editable.array[len(editable.array)-1] = mapEntryI{}
	return editable
}
