package hashmap

import (
	"fmt"
	"gocloj/data/atom"
	"math/big"
	"math/rand"
	"strings"
)

func randomPairs(count int) []mapEntry {
	// constant seed for random
	rand.Seed(3289417)
	pairs := []mapEntry{}
	for i := 0; i < count; i++ {
		pairs = append(pairs, mapEntry{
			key: &atom.Num{Val: big.NewInt(rand.Int63())},
			val: &atom.Num{Val: big.NewInt(rand.Int63())},
		})
	}
	return pairs
}

func walkArrayNode(parentName string, parentPort string, cnt int, n *arrayNode) {
	nodeName := fmt.Sprintf("%s_%d", parentName, cnt)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("{arrayNode|%d|{", n.count))
	for i, _ := range n.array {
		if i != 0 {
			builder.WriteString("|")
		}
		builder.WriteString(fmt.Sprintf("<f%d> %d", i, i))
	}
	builder.WriteString("}}")

	fmt.Printf("%s [label=\"%s\"]; %s -> %s\n",
		nodeName, builder.String(), parentPort, nodeName)

	for i, child := range n.array {
		parent := fmt.Sprintf("%s:<f%d>", nodeName, i)
		walkPhmNode(nodeName, parent, i, child)
	}
}

func walkBitmapIndexedNode(parentName string, parentPort string, cnt int, n *bitmapIndexedNode) {
	nodeName := fmt.Sprintf("%s_%d", parentName, cnt)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("{bitmapIndexedNode|0x%08X|{", n.bitmap))
	for i, _ := range n.array {
		if i != 0 {
			builder.WriteString("|")
		}
		builder.WriteString(fmt.Sprintf("<f%d> %d", i, i))
	}
	builder.WriteString("}}")

	fmt.Printf("%s [label=\"%s\"]; %s -> %s\n",
		nodeName, builder.String(), parentPort, nodeName)

	for i, entry := range n.array {
		parent := fmt.Sprintf("%s:<f%d>", nodeName, i)

		if entry.key == nil {
			if entry.val != nil {
				walkPhmNode(nodeName, parent, i, entry.val.(phmNode))
			}
		} else {
			k := entry.key.(atom.Atom)
			v := entry.val.(atom.Atom)

			entryName := fmt.Sprintf("%s_e_%d", nodeName, i)
			fmt.Printf("%s [label=\"{%s|%s}\"]; %s -> %s\n",
				entryName, k.String(), v.String(), parent, entryName)
		}
	}
}

func walkHashCollisionNode(parentName string, parentPort string, cnt int, n *hashCollisionNode) {
	nodeName := fmt.Sprintf("%s_%d", parentName, cnt)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("{hashCollisionNode|0x%08X|%d|{",
		n.hash, n.count))
	for i, _ := range n.array {
		if i != 0 {
			builder.WriteString("|")
		}
		builder.WriteString(fmt.Sprintf("<f%d>", i))
	}
	builder.WriteString("}}")

	fmt.Printf("%s [label=\"%s\"]; %s -> %s\n",
		nodeName, builder.String(), parentPort, nodeName)

	for i, entry := range n.array {
		parent := fmt.Sprintf("%s:<f%d>", nodeName, i)

		k := entry.key
		v := entry.val

		entryName := fmt.Sprintf("%s_e_%d", nodeName, i)
		fmt.Printf("%s [label=\"{%s|%s}\"]; %s -> %s\n",
			entryName, k.String(), v.String(), parent, entryName)
	}
}

func walkPhmNode(parentName string, parentPort string, cnt int, node phmNode) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *arrayNode:
		walkArrayNode(parentName, parentPort, cnt, n)
	case *bitmapIndexedNode:
		walkBitmapIndexedNode(parentName, parentPort, cnt, n)
	case *hashCollisionNode:
		walkHashCollisionNode(parentName, parentPort, cnt, n)
	default:
		fmt.Println("UNKNOWN NODE TYPE!")
	}
}

func dumpPhm(m *PersistentHashMap) {
	fmt.Println("digraph {")
	fmt.Println("rankdir=LR")
	fmt.Println("node [shape=record, rankdir=LR]")
	fmt.Printf("map [label=\"map|%d\"]\n",
		m.count)

	if m.hasNil {
		fmt.Printf("nilVal [label=\"%s\"]; map -> nilVal [label=\"nilVal\"]\n",
			m.nilVal.String())
	}

	walkPhmNode("map", "map", 0, m.root)
	fmt.Println("}")
}

func Dump() {
	fillPersistentMap := func(m PersistentMap, pairs []mapEntry) PersistentMap {
		for _, pair := range pairs {
			m = m.Assoc(pair.key, pair.val)
		}
		return m
	}

	pairs := randomPairs(10)
	m := NewPersistentHashMap()
	m = fillPersistentMap(m, pairs)

	dumpPhm(m.(*PersistentHashMap))
}
