package hashmap

import (
	gdat "gocloj/data/atom/testing"
	"math/rand"
)

func randomDupPairs(pairs []mapEntry) []mapEntry {
	// constant seed for random
	rand.Seed(923437776)
	numPairs := len(pairs)
	nums := gdat.RandomNums(len(pairs))
	newPairs := make([]mapEntry, numPairs, numPairs)
	for i := 0; i < len(pairs); i++ {
		newPairs[i] = mapEntry{
			key: pairs[i].key,
			val: nums[i],
		}
	}
	return newPairs
}
