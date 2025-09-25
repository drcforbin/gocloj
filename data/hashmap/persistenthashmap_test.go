package hashmap

import (
	"gocloj/data/atom"
	"math/big"
	"testing"
)

// TODO: iteration tests

var gm PersistentMap

func fillPersistentMap(m PersistentMap, pairs []mapEntry) PersistentMap {
	for _, pair := range pairs {
		m = m.Assoc(pair.key, pair.val)
	}
	return m
}

func getPairs(m PersistentMap, pairs []mapEntry) bool {
	for _, pair := range pairs {
		actual := m.Get(pair.key)
		if !pair.val.Equals(actual) {
			return false
		}
	}
	return true
}

func TestPersistentHashMapSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		pairs := randomPairs(count)
		m := NewPersistentHashMap()
		m = fillPersistentMap(m, pairs)
		if !getPairs(m, pairs) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestPersistentHashMapDupSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		pairs := randomPairs(count)
		dupPairs := randomDupPairs(pairs)
		m := NewPersistentHashMap()
		m = fillPersistentMap(m, pairs)
		m = fillPersistentMap(m, dupPairs)
		if !getPairs(m, dupPairs) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestPersistentHashMapSetImmutability(t *testing.T) {
	pairs := []mapEntry{
		mapEntry{&atom.Num{Val: big.NewInt(8)}, &atom.Num{Val: big.NewInt(99)}},
		mapEntry{&atom.Num{Val: big.NewInt(10)}, &atom.Num{Val: big.NewInt(108983)}},
		mapEntry{&atom.Num{Val: big.NewInt(13)}, &atom.Num{Val: big.NewInt(600)}},
		mapEntry{&atom.Num{Val: big.NewInt(3545)}, &atom.Num{Val: big.NewInt(1)}},
	}

	// TODO: test immutability in removal

	// start with empty map
	last := NewPersistentHashMap()
	maps := []PersistentMap{last}

	for i, outerPair := range pairs {
		// assoc one of the elements
		m := last.Assoc(outerPair.key, outerPair.val)
		maps = append(maps, m)

		// for each map so far, check/recheck its length
		for j, m := range maps {
			if m.Length() != j {
				t.Errorf("map %d contained unexpected number of values", j)
			}

			// ...and make sure it returns the right vals or Nil
			for k, pair := range pairs {
				if k < j {
					if !m.Get(pair.key).Equals(pair.val) {
						t.Errorf("unexpected value iter %d, map %d, pair %d", i, j, k)
					}
				} else {
					if !m.Get(pair.key).Equals(atom.Nil) {
						t.Errorf("expected Nil iter %d, map %d, pair %d", i, j, k)
					}
				}
			}
		}

		last = m
	}
}

func TestPersistentHashMapDupImmutability(t *testing.T) {
	k1 := &atom.Num{Val: big.NewInt(8)}
	k2 := &atom.Num{Val: big.NewInt(10)}

	// start with empty map, add a couple
	m := NewPersistentHashMap().
		Assoc(k1, &atom.Num{Val: big.NewInt(99)}).
		Assoc(k2, &atom.Num{Val: big.NewInt(108983)})

	// replace them
	m2 := m.Assoc(k1, &atom.Num{Val: big.NewInt(9387)}).
		Assoc(k2, &atom.Num{Val: big.NewInt(3)})

	if m.Length() != 2 || m2.Length() != 2 {
		t.Errorf("unexpected length for m (%d) or m2 (%d)", m.Length(), m2.Length())
	}

	if !m.Get(k1).Equals(&atom.Num{Val: big.NewInt(99)}) {
		t.Errorf("unexpected value")
	}
	if !m.Get(k2).Equals(&atom.Num{Val: big.NewInt(108983)}) {
		t.Errorf("unexpected value")
	}

	if !m2.Get(k1).Equals(&atom.Num{Val: big.NewInt(9387)}) {
		t.Errorf("unexpected value")
	}
	if !m2.Get(k2).Equals(&atom.Num{Val: big.NewInt(3)}) {
		t.Errorf("unexpected value")
	}
}

func TestPersistentHashMapWithout(t *testing.T) {
	k1 := &atom.Num{Val: big.NewInt(8)}
	k2 := &atom.Num{Val: big.NewInt(10)}

	// start with empty map, add a couple and replace them
	m := NewPersistentHashMap().
		Assoc(k1, &atom.Num{Val: big.NewInt(99)}).
		Assoc(k2, &atom.Num{Val: big.NewInt(108983)}).
		Assoc(k1, &atom.Num{Val: big.NewInt(9387)}).
		Assoc(k2, &atom.Num{Val: big.NewInt(3)})

	// remove k1
	m2 := m.Without(k1)

	if m.Length() != 2 || m2.Length() != 1 {
		t.Errorf("unexpected length for m (%d) or m2 (%d)", m.Length(), m2.Length())
	}

	if !m.Get(k1).Equals(&atom.Num{Val: big.NewInt(9387)}) {
		t.Errorf("unexpected value")
	}
	if !m.Get(k2).Equals(&atom.Num{Val: big.NewInt(3)}) {
		t.Errorf("unexpected value")
	}

	if !m2.Get(k1).Equals(atom.Nil) {
		t.Errorf("unexpected value")
	}
	if !m2.Get(k2).Equals(&atom.Num{Val: big.NewInt(3)}) {
		t.Errorf("unexpected value")
	}
}

func BenchmarkPersistentHashMap_SetRandom1000(b *testing.B) {
	pairs := randomPairs(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewPersistentHashMap()
		gm = fillPersistentMap(m, pairs)
	}
}

func BenchmarkPersistentHashMap_GetRandom1000(b *testing.B) {
	pairs := randomPairs(1000)
	m := NewPersistentHashMap()
	m = fillTransientMap(m, pairs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !getPairs(m, pairs) {
			b.Errorf("failed to retrieve matching vals")
		}
	}
}

func BenchmarkPersistentHashMap_SetDupRandom1000(b *testing.B) {
	pairs := randomPairs(1000)
	dupPairs := randomDupPairs(pairs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewPersistentHashMap()
		gm = fillPersistentMap(m, pairs)
		gm = fillPersistentMap(m, dupPairs)
	}
}

func BenchmarkTransientHashMap_SetDupRandom1000(b *testing.B) {
	pairs := randomPairs(1000)
	dupPairs := randomDupPairs(pairs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewPersistentHashMap()
		gm = fillTransientMap(m, pairs)
		gm = fillTransientMap(m, dupPairs)
	}
}
