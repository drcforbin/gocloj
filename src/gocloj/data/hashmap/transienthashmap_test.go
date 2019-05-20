package hashmap

import (
	"gocloj/data/atom"
	"math/big"
	"testing"
)

func fillTransientMap(m PersistentMap, pairs []mapEntry) PersistentMap {
	t := m.AsTransient(1)
	for _, pair := range pairs {
		t = t.Assoc(pair.key, pair.val)
	}
	return t.AsPersistent()
}

func TestTransientHashMapSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		pairs := randomPairs(count)
		m := NewPersistentHashMap()
		m = fillTransientMap(m, pairs)
		if !getPairs(m, pairs) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestTransientHashMapDupSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		pairs := randomPairs(count)
		dupPairs := randomDupPairs(pairs)
		m := NewPersistentHashMap()
		m = fillTransientMap(m, pairs)
		m = fillTransientMap(m, dupPairs)
		if !getPairs(m, dupPairs) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestTransientHashMapDupImmutability(t *testing.T) {
	k1 := &atom.Num{Val: big.NewInt(8)}
	k2 := &atom.Num{Val: big.NewInt(10)}

	// start with empty map, make transient, add items
	m := NewPersistentHashMap().
		AsTransient(1).
		Assoc(k1, &atom.Num{Val: big.NewInt(99)}).
		Assoc(k2, &atom.Num{Val: big.NewInt(108983)}).
		AsPersistent()

	if m.Length() != 2 {
		t.Errorf("unexpected length for m, %d", m.Length())
	}

	if !m.Get(k1).Equals(&atom.Num{Val: big.NewInt(99)}) {
		t.Errorf("unexpected value for m k1, %s", m.Get(k1))
	}
	if !m.Get(k2).Equals(&atom.Num{Val: big.NewInt(108983)}) {
		t.Errorf("unexpected value for m k2, %s", m.Get(k2))
	}

	// start with m, make transient, dup items
	m = m.AsTransient(1).
		Assoc(k1, &atom.Num{Val: big.NewInt(9387)}).
		Assoc(k2, &atom.Num{Val: big.NewInt(3)}).
		AsPersistent()

	if m.Length() != 2 {
		t.Errorf("unexpected length for m, %d", m.Length())
	}

	if !m.Get(k1).Equals(&atom.Num{Val: big.NewInt(9387)}) {
		t.Errorf("unexpected value for m k1, %s", m.Get(k1))
	}
	if !m.Get(k2).Equals(&atom.Num{Val: big.NewInt(3)}) {
		t.Errorf("unexpected value for m k2, %s", m.Get(k2))
	}
}

func TestTransientHashMapContext(t *testing.T) {
	k0 := &atom.Num{Val: big.NewInt(222)}
	k1 := &atom.Num{Val: big.NewInt(8)}
	k2 := &atom.Num{Val: big.NewInt(10)}

	// start with empty map, add a single value (later maps
	// should share it between them)
	m0 := NewPersistentHashMap()
	m0 = m0.Assoc(k0, &atom.Num{Val: big.NewInt(1129374)})

	// make transients with different contexts
	tm1 := m0.AsTransient(1)
	tm2 := m0.AsTransient(2)

	// set differing values for each, making sure that
	// neither depends on the other
	tm1 = tm1.Assoc(k1, &atom.Num{Val: big.NewInt(99)})
	tm2 = tm2.Assoc(k1, &atom.Num{Val: big.NewInt(9387)})
	tm2 = tm2.Assoc(k2, &atom.Num{Val: big.NewInt(3)})
	tm1 = tm1.Assoc(k2, &atom.Num{Val: big.NewInt(108983)})

	// make persistent again
	m1 := tm1.AsPersistent()
	m2 := tm2.AsPersistent()

	if m0.Length() != 1 {
		t.Errorf("m0 was modified by transient operations!")
	}
	if m1.Length() != 3 || m2.Length() != 3 {
		t.Errorf("unexpected length for m1 (%d) or m2 (%d)", m1.Length(), m2.Length())
	}

	// check value and pointer equality; they should both share k0 val
	if !m0.Get(k0).Equals(&atom.Num{Val: big.NewInt(1129374)}) ||
		!m1.Get(k0).Equals(&atom.Num{Val: big.NewInt(1129374)}) ||
		!m2.Get(k0).Equals(&atom.Num{Val: big.NewInt(1129374)}) ||
		m0.Get(k0) != m1.Get(k0) || m1.Get(k0) != m2.Get(k0) {
		t.Error("expected k0 value to be shared")
	}

	if !m1.Get(k1).Equals(&atom.Num{Val: big.NewInt(99)}) {
		t.Errorf("unexpected value for m1 k1, %s", m1.Get(k1))
	}
	if !m1.Get(k2).Equals(&atom.Num{Val: big.NewInt(108983)}) {
		t.Errorf("unexpected value for m1 k2, %s", m1.Get(k2))
	}

	if !m2.Get(k1).Equals(&atom.Num{Val: big.NewInt(9387)}) {
		t.Errorf("unexpected value for m2 k1, %s", m2.Get(k1))
	}
	if !m2.Get(k2).Equals(&atom.Num{Val: big.NewInt(3)}) {
		t.Errorf("unexpected value for m2 k2, %s", m2.Get(k2))
	}
}

func BenchmarkTransientHashMap_SetRandom1000(b *testing.B) {
	pairs := randomPairs(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewPersistentHashMap()
		gm = fillTransientMap(m, pairs)
	}
}
