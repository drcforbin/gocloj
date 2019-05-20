package hashset

import (
	"gocloj/data/atom"
	"math/big"
	"testing"
)

// TODO: iteration tests

var gm PersistentSet

func fillPersistentSet(m PersistentSet, vals []atom.Atom) PersistentSet {
	for _, val := range vals {
		m = m.Assoc(val)
	}
	return m
}

func getVals(m PersistentSet, vals []atom.Atom) bool {
	for _, val := range vals {
		actual := m.Get(val)
		if !val.Equals(actual) {
			return false
		}
	}
	return true
}

func TestPersistentHashSetSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		vals := randomVals(count)
		m := NewPersistentHashSet()
		m = fillPersistentSet(m, vals)
		if !getVals(m, vals) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestPersistentHashSetDupSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		vals := randomVals(count)
		dupVals := randomDupVals(vals)
		m := NewPersistentHashSet()
		m = fillPersistentSet(m, vals)
		m = fillPersistentSet(m, dupVals)
		if !getVals(m, dupVals) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestPersistentHashSetSetImmutability(t *testing.T) {
	vals := []atom.Atom{
		&atom.Num{Val: big.NewInt(8)},
		&atom.Num{Val: big.NewInt(10)},
		&atom.Num{Val: big.NewInt(13)},
		&atom.Num{Val: big.NewInt(3545)},
	}

	// TODO: test immutability in removal

	// start with empty set
	last := NewPersistentHashSet()
	sets := []PersistentSet{last}

	for i, outerVal := range vals {
		// assoc one of the elements
		s := last.Assoc(outerVal)
		sets = append(sets, s)

		// for each set so far, check/recheck its length
		for j, m := range sets {
			if m.Length() != j {
				t.Errorf("set %d contained unexpected number of values", j)
			}

			// ...and make sure it returns the right vals or Nil
			for k, val := range vals {
				if k < j {
					if !m.Get(val).Equals(val) {
						t.Errorf("unexpected value iter %d, set %d, val %d", i, j, k)
					}
				} else {
					if !m.Get(val).Equals(atom.Nil) {
						t.Errorf("expected Nil iter %d, set %d, val %d", i, j, k)
					}
				}
			}
		}

		last = s
	}
}

func TestPersistentHashSetDupImmutability(t *testing.T) {
	v1 := &atom.Num{Val: big.NewInt(8)}
	v2 := &atom.Num{Val: big.NewInt(10)}
	v3 := &atom.Num{Val: big.NewInt(9)}

	// start with empty set, add a couple
	s := NewPersistentHashSet().
		Assoc(v1).
		Assoc(v2)

	// add a new item to make sure we get a new set, then replace them
	s2 := s.Assoc(&atom.Num{Val: big.NewInt(9)}).
		Assoc(&atom.Num{Val: big.NewInt(10)}).
		Assoc(&atom.Num{Val: big.NewInt(8)})

	if s.Length() != 2 || s2.Length() != 3 {
		t.Errorf("unexpected length for s (%d) or s2 (%d)", s.Length(), s2.Length())
	}

	if !s.Get(v1).Equals(v1) {
		t.Errorf("unexpected value")
	}
	if !s.Get(v2).Equals(v2) {
		t.Errorf("unexpected value")
	}
	if !s.Get(v3).Equals(atom.Nil) {
		t.Errorf("unexpected value")
	}

	if !s2.Get(v1).Equals(v1) {
		t.Errorf("unexpected value")
	}
	if !s2.Get(v1).Equals(v1) {
		t.Errorf("unexpected value")
	}
	if !s2.Get(v3).Equals(v3) {
		t.Errorf("unexpected value")
	}
}

func TestPersistentHashSetWithout(t *testing.T) {
	v1 := &atom.Num{Val: big.NewInt(8)}
	v2 := &atom.Num{Val: big.NewInt(10)}

	// start with empty set, add a couple and replace them
	s := NewPersistentHashSet().
		Assoc(v1).
		Assoc(v2).
		Assoc(v2).
		Assoc(v1)

	// remove v1
	s2 := s.Without(v1)

	if s.Length() != 2 || s2.Length() != 1 {
		t.Errorf("unexpected length for s (%d) or s2 (%d)", s.Length(), s2.Length())
	}

	if !s.Get(v1).Equals(v1) {
		t.Errorf("unexpected value")
	}
	if !s.Get(v2).Equals(v2) {
		t.Errorf("unexpected value")
	}

	if !s2.Get(v1).Equals(atom.Nil) {
		t.Errorf("unexpected value")
	}
	if !s2.Get(v2).Equals(v2) {
		t.Errorf("unexpected value")
	}
}

func BenchmarkPersistentHashSet_SetRandom1000(b *testing.B) {
	vals := randomVals(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewPersistentHashSet()
		gm = fillPersistentSet(m, vals)
	}
}

func BenchmarkPersistentHashSet_GetRandom1000(b *testing.B) {
	vals := randomVals(1000)
	m := NewPersistentHashSet()
	m = fillTransientSet(m, vals)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !getVals(m, vals) {
			b.Errorf("failed to retrieve matching vals")
		}
	}
}

func BenchmarkPersistentHashSet_SetDupRandom1000(b *testing.B) {
	vals := randomVals(1000)
	dupVals := randomDupVals(vals)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewPersistentHashSet()
		gm = fillPersistentSet(m, vals)
		gm = fillPersistentSet(m, dupVals)
	}
}

func BenchmarkTransientHashSet_SetDupRandom1000(b *testing.B) {
	vals := randomVals(1000)
	dupVals := randomDupVals(vals)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewPersistentHashSet()
		gm = fillTransientSet(m, vals)
		gm = fillTransientSet(m, dupVals)
	}
}
