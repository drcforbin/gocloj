package hashmap

import (
	"fmt"
	"gocloj/data/atom"
	"math/bits"
	"testing"
)

type mockHashAtom struct {
	hash uint32
}

func (mha mockHashAtom) String() string {
	return fmt.Sprintf("0x%08X", mha.hash)
}

func (mha mockHashAtom) IsNil() bool {
	return false
}

func (mha mockHashAtom) Hash() uint32 {
	return mha.hash
}

func (mha mockHashAtom) Equals(atom atom.Atom) bool {
	// no two of these are equal, so we can easily
	// get collisions
	return false
}

// TODO: test overflow to two bitmapIndexedNodes (happens when hashes
// differ, but they map to the same bucket)
func TestBitmapIndexedNodeOverflowToArrayNode(t *testing.T) {
	var node phmNode = &bitmapIndexedNode{
		edit:  0,
		array: []mapEntryI{},
	}

	// fill 16 items in (0-16)
	var i int
	for i = 0; i < 16; i++ {
		a := &mockHashAtom{hash: uint32(i)}
		// using same val for key and val, it's easier
		addedLeaf := false
		node = node.assoc(0, a.Hash(), a, a, &addedLeaf)
		if !addedLeaf {
			t.Errorf("assoc item %d did not set addedLeaf", i)
		}

		if n, ok := node.(*bitmapIndexedNode); !ok {
			t.Errorf("assoc item %d did not return bitmapIndexedNode", i)
		} else {
			count := bits.OnesCount32(n.bitmap)
			if count != i+1 {
				t.Errorf("adding node %d resulted in wrong count %d", i, count)
			}
		}
	}

	// add a 17th item. This should kick it over, causing bin to
	// become an arraynode containing 17 bitmapIndexedNodes
	i = 16
	a := &mockHashAtom{hash: uint32(i)}
	// using same val for key and val, it's easier
	addedLeaf := false
	node = node.assoc(0, a.Hash(), a, a, &addedLeaf)
	if !addedLeaf {
		t.Errorf("assoc item %d did not set addedLeaf", i)
	}

	if n, ok := node.(*arrayNode); !ok {
		t.Errorf("assoc item %d did not return arrayNode", i)
	} else {
		if n.count != 17 {
			t.Errorf("adding node %d resulted in wrong count %d", i, n.count)
		}

		for i = 0; i < 17; i++ {
			if _, ok := n.array[i].(*bitmapIndexedNode); !ok {
				t.Errorf("subnode %d was wrong type %T", i, n.array[i])
			}
		}

		for ; i < 32; i++ {
			if n.array[i] != nil {
				t.Errorf("subnode %d was wrong type %T", i, n.array[i])
			}
		}
	}
}

func TestBitmapIndexedNodeCollide(t *testing.T) {
	var node phmNode = &bitmapIndexedNode{
		edit:  0,
		array: []mapEntryI{},
	}

	// add two different items with same hash value
	a := &mockHashAtom{hash: uint32(2894734455)}
	b := &mockHashAtom{hash: uint32(2894734455)}

	addedLeaf := false
	node = node.assoc(0, a.Hash(), a, a, &addedLeaf)
	if !addedLeaf {
		t.Errorf("assoc item did not set addedLeaf")
	}
	addedLeaf = false
	node = node.assoc(0, b.Hash(), b, b, &addedLeaf)
	if !addedLeaf {
		t.Errorf("assoc item did not set addedLeaf")
	}

	if n, ok := node.(*bitmapIndexedNode); !ok {
		t.Errorf("assoc items did not return bitmapIndexedNode")
	} else {
		// we expect the bitmapIndexedNode to contain only a single child
		count := bits.OnesCount32(n.bitmap)
		if count != 1 {
			t.Errorf("adding nodes resulted in wrong count %d", count)
		}

		// and that child node should be a hcn
		if n.array[0].key != nil {
			t.Errorf("expected key for child 0 to be nil (memory leak?)")
		}
		if _, ok := n.array[0].val.(*hashCollisionNode); !ok {
			t.Errorf("items did not contain a hashCollisionNode")
		}

		// make sure the rest of the array is nil
		for i := 1; i < len(n.array); i++ {
			if n.array[i].key != nil {
				t.Errorf("expected key for child %d to be nil", i)
			}
			if n.array[i].val != nil {
				t.Errorf("expected val for child %d to be nil", i)
			}
		}
	}
}
