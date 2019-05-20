package hashset

import (
	"gocloj/data/atom"
	// gdat "gocloj/data/atom/testing"
	"math/big"
	"math/rand"
)

func randomVals(count int) []atom.Atom {
	// constant seed for random
	rand.Seed(3289417)
	vals := []atom.Atom{}
	for i := 0; i < count; i++ {
		vals = append(vals, &atom.Num{Val: big.NewInt(rand.Int63())})
	}
	return vals
}

func randomDupVals(vals []atom.Atom) []atom.Atom {
	// TODO: randomize vals
	return vals
}
