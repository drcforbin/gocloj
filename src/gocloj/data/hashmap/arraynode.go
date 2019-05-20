package hashmap

import (
	"gocloj/data/atom"
)

// arrayNode is used to store references to up to 32 subnodes, either
// bitmapIndexedNode or hashCollisionNode structures.
type arrayNode struct {
	edit  uint64
	count int

	array [32]phmNode
}

func (an *arrayNode) assoc(shift uint32, hash uint32, key atom.Atom, val atom.Atom, addedLeaf *bool) phmNode {
	idx := mask(hash, shift)
	node := an.array[idx]

	// add node with new value if not found
	if node == nil {
		newNode := &arrayNode{
			count: an.count + 1,
			array: an.array,
		}
		newNode.array[idx] = emptyBin.assoc(shift+5, hash, key, val, addedLeaf)
		return newNode
	}

	// otherwise, add the value to the node
	n := node.assoc(shift+5, hash, key, val, addedLeaf)
	if n == node {
		return an
	}

	newNode := &arrayNode{
		count: an.count,
		array: an.array,
	}
	newNode.array[idx] = n
	return newNode
}

func (an *arrayNode) without(shift uint32, hash uint32, key atom.Atom, removedLeaf *bool) phmNode {
	idx := mask(hash, shift)
	node := an.array[idx]
	if node == nil {
		return an
	}

	n := node.without(shift+5, hash, key, removedLeaf)
	if n == node {
		return an
	}
	if n == nil {
		if an.count <= 8 {
			// shrink
			return an.pack(0, int(idx))
		}

		newNode := &arrayNode{
			count: an.count - 1,
			array: an.array,
		}
		newNode.array[idx] = n
		return newNode
	} else {
		newNode := &arrayNode{
			count: an.count,
			array: an.array,
		}
		newNode.array[idx] = n
		return newNode
	}
}

func (an *arrayNode) find(shift uint32, hash uint32, key atom.Atom) *mapEntry {
	idx := mask(hash, shift)
	node := an.array[idx]
	if node != nil {
		return node.find(shift+5, hash, key)
	}
	return nil
}

func (an *arrayNode) assocT(edit uint64, shift uint32, hash uint32, key atom.Atom, val atom.Atom, addedLeaf *bool) phmNode {
	idx := mask(hash, shift)
	node := an.array[idx]

	// add node with new value if not found
	if node == nil {
		editable := an.editAndSet(edit, int(idx),
			emptyBin.assocT(edit, shift+5, hash, key, val, addedLeaf))
		editable.count++
		return editable
	}

	// otherwise, add the value to the node
	n := node.assocT(edit, shift+5, hash, key, val, addedLeaf)
	if n == node {
		return an
	}

	return an.editAndSet(edit, int(idx), n)
}

func (an *arrayNode) withoutT(edit uint64, shift uint32, hash uint32, key atom.Atom, removedLeaf *bool) phmNode {
	idx := int(mask(hash, shift))
	node := an.array[idx]
	if node == nil {
		return an
	}
	n := node.withoutT(edit, shift+5, hash, key, removedLeaf)
	if n == node {
		return an
	}
	if n == nil {
		if an.count <= 8 {
			// shrink
			return an.pack(edit, idx)
		}
		editable := an.editAndSet(edit, idx, n)
		editable.count--
		return editable
	}
	return an.editAndSet(edit, idx, n)
}

func (an *arrayNode) iterator(handler iterHandler) atom.SeqIterator {
	return &arrayNodeIterator{
		handler: handler,
		array:   an.array[:],
	}
}

func (an *arrayNode) ensureEditable(edit uint64) *arrayNode {
	if an.edit == edit {
		return an
	}
	return &arrayNode{edit: edit, count: an.count, array: an.array}
}

func (an *arrayNode) editAndSet(edit uint64, i int, n phmNode) *arrayNode {
	editable := an.ensureEditable(edit)
	editable.array[i] = n
	return editable
}

func (an *arrayNode) pack(edit uint64, idx int) phmNode {
	newArray := make([]mapEntryI, an.count-1)

	j := 1
	bitmap := uint32(0)
	for i := 0; i < idx; i++ {
		if an.array[i] != nil {
			newArray[j].val = an.array[i]
			bitmap |= 1 << uint32(i)
			j++
		}
	}

	for i := idx + 1; i < len(an.array); i++ {
		if an.array[i] != nil {
			newArray[j].val = an.array[i]
			bitmap |= 1 << uint32(i)
			j++
		}
	}

	return &bitmapIndexedNode{
		edit:   edit,
		bitmap: bitmap,
		array:  newArray,
	}
}
