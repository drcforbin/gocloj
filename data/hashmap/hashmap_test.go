package hashmap

/*
import (
	"testing"
)

func fillMap(m Map, pairs []mapEntry) {
	for _, pair := range pairs {
		m.Set(pair.key, pair.val)
	}
}

func getHashMapPairs(m Map, pairs []mapEntry) bool {
	for _, pair := range pairs {
		actual := m.Get(pair.key)
		if !pair.val.Equals(actual) {
			return false
		}
	}
	return true
}

func TestHashMapSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		pairs := randomPairs(count)
		m := NewHashMap()
		fillMap(m, pairs)
		if !getHashMapPairs(m, pairs) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func TestHashMapDupSetGet(t *testing.T) {
	counts := []int{5, 100, 1000, 10000}
	for _, count := range counts {
		pairs := randomPairs(count)
		dupPairs := randomDupPairs(pairs)
		m := NewHashMap()
		fillMap(m, pairs)
		fillMap(m, dupPairs)
		if !getHashMapPairs(m, dupPairs) {
			t.Errorf("failed to retrieve matching vals for count %d", count)
		}
	}
}

func BenchmarkHashMapSetRandom1000(b *testing.B) {
	pairs := randomPairs(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewHashMap()
		fillMap(m, pairs)
	}
}

func BenchmarkHashMapGetRandom1000(b *testing.B) {
	pairs := randomPairs(1000)
	m := NewHashMap()
	fillMap(m, pairs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !getHashMapPairs(m, pairs) {
			b.Errorf("failed to retrieve matching vals")
		}
	}
}

func BenchmarkHashMapSetDupRandom1000(b *testing.B) {
	pairs := randomPairs(1000)
	dupPairs := randomDupPairs(pairs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewHashMap()
		fillMap(m, pairs)
		fillMap(m, dupPairs)
	}
}
*/
