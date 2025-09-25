package testing

import (
	"gocloj/data/atom"
	"math/big"
	"math/rand"
)

func RandomNums(count int) []atom.Atom {
	// note: NOT seeding rng; leave that to caller
	nums := []atom.Atom{}
	for i := 0; i < count; i++ {
		nums = append(nums,
			&atom.Num{Val: big.NewInt(rand.Int63())},
		)
	}
	return nums
}
