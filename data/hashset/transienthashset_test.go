package hashset

import (
	"gocloj/data/atom"
	"math/big"
	"testing"
)

func fillTransientSet(m PersistentSet, vals []atom.Atom) PersistentSet {
	t := m.AsTransient(1)
	for _, val := range vals {
		t = t.Assoc(val)
	}
	return t.AsPersistent()
}

func TestTransientHashSetSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		vals := randomVals(count)
		m := NewPersistentHashSet()
		m = fillTransientSet(m, vals)
		if !getVals(m, vals) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestTransientHashSetDupSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		vals := randomVals(count)
		dupVals := randomDupVals(vals)
		m := NewPersistentHashSet()
		m = fillTransientSet(m, vals)
		m = fillTransientSet(m, dupVals)
		if !getVals(m, dupVals) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestTransientHashSetDupImmutability(t *testing.T) {
	v1 := &atom.Num{Val: big.NewInt(8)}
	v2 := &atom.Num{Val: big.NewInt(10)}
	v3 := &atom.Num{Val: big.NewInt(9)}

	// start with empty set, make transient, add items
	s := NewPersistentHashSet().
		AsTransient(1).
		Assoc(v1).
		Assoc(v2).
		AsPersistent()

	if s.Length() != 2 {
		t.Errorf("unexpected length for s, %d", s.Length())
	}

	if !s.Get(v1).Equals(v1) {
		t.Errorf("unexpected value for s v1, %s", s.Get(v1))
	}
	if !s.Get(v2).Equals(v2) {
		t.Errorf("unexpected value for s v2, %s", s.Get(v2))
	}

	// add a new item to make sure we get a new set, then replace them
	s = s.AsTransient(1).
		Assoc(&atom.Num{Val: big.NewInt(9)}).
		Assoc(&atom.Num{Val: big.NewInt(10)}).
		Assoc(&atom.Num{Val: big.NewInt(8)}).
		AsPersistent()

	if s.Length() != 3 {
		t.Errorf("unexpected length for s, %d", s.Length())
	}

	if !s.Get(v1).Equals(v1) {
		t.Errorf("unexpected value for s v1, %s", s.Get(v1))
	}
	if !s.Get(v2).Equals(v2) {
		t.Errorf("unexpected value for s v2, %s", s.Get(v2))
	}
	if !s.Get(v3).Equals(v3) {
		t.Errorf("unexpected value for s v2, %s", s.Get(v3))
	}
}

func TestTransientHashSetContext(t *testing.T) {
	v0 := &atom.Num{Val: big.NewInt(222)}

	// start with empty set, add a single value (later sets
	// should share it between them)
	s0 := NewPersistentHashSet()
	s0 = s0.Assoc(v0)

	// make transients with different contexts
	ts1 := s0.AsTransient(1)
	ts2 := s0.AsTransient(2)

	// set same v1 value on each, but differing values v2 and v3,
	// making sure that neither depends on the other by interleaving
	v1 := &atom.Num{Val: big.NewInt(8)}
	v2 := &atom.Num{Val: big.NewInt(10)}
	v3 := &atom.Num{Val: big.NewInt(99)}
	ts1 = ts1.Assoc(v1)
	ts2 = ts2.Assoc(v2)
	ts2 = ts2.Assoc(v1)
	ts1 = ts1.Assoc(v3)

	// make persistent again
	s1 := ts1.AsPersistent()
	s2 := ts2.AsPersistent()

	if s0.Length() != 1 {
		t.Errorf("s0 was modified by transient operations!")
	}
	if s1.Length() != 3 || s2.Length() != 3 {
		t.Errorf("unexpected length for s1 (%d) or s2 (%d)", s1.Length(), s2.Length())
	}

	// check value and pointer equality; all three should share v0
	if !s0.Get(v0).Equals(v0) ||
		!s1.Get(v0).Equals(v0) ||
		!s2.Get(v0).Equals(v0) ||
		s0.Get(v0) != s1.Get(v0) || s1.Get(v0) != s2.Get(v0) {
		t.Error("expected v0 value to be shared")
	}
	// check value and pointer equality; s1 and s2 should share v1
	if !s1.Get(v1).Equals(v1) ||
		!s2.Get(v1).Equals(v1) ||
		s1.Get(v1) != s2.Get(v1) {
		t.Error("expected v1 value to be shared")
	}

	if !s1.Get(v2).Equals(atom.Nil) {
		t.Errorf("unexpected value for s1 v2, %s", s1.Get(v2))
	}
	if !s1.Get(v3).Equals(v3) {
		t.Errorf("unexpected value for s1 v3, %s", s1.Get(v3))
	}

	if !s2.Get(v2).Equals(v2) {
		t.Errorf("unexpected value for s2 v2, %s", s2.Get(v2))
	}
	if !s2.Get(v3).Equals(atom.Nil) {
		t.Errorf("unexpected value for s2 v3, %s", s2.Get(v3))
	}
}

func BenchmarkTransientHashSet_SetRandom1000(b *testing.B) {
	vals := randomVals(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewPersistentHashSet()
		gm = fillTransientSet(m, vals)
	}
}
