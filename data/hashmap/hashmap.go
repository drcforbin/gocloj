package hashmap

/*
import (
	"gocloj/data/atom"
	"strings"
)

type hashMap struct {
	table [][]mapEntry
	size  int
}

func NewHashMap() Map {
	return &hashMap{
		// prealloc to len 2
		table: make([][]mapEntry, 2),
	}
}

func (m hashMap) String() string {
	var builder strings.Builder

	builder.WriteString("{")
	cnt := 0
	for _, bucket := range m.table {
		if bucket != nil {
			for _, pair := range bucket {
				if cnt != 0 {
					builder.WriteString(" ")
				}
				builder.WriteString(pair.key.String())
				builder.WriteString(" ")
				builder.WriteString(pair.val.String())

				cnt++
			}
		}
	}
	builder.WriteString("}")

	return builder.String()
}

func (m *hashMap) IsNil() bool {
	return false
}

	// Returns a hash value for this Atom.
func (m hashMap) Hash() uint32 {
	return atom.SeqHash(m)
}

// Returns whether this Atom is equivalent to a given atom.
func (m hashMap) Equals(a atom.Atom) bool {
	if val, ok := a.(atom.Seq); ok {
		// TODO: this is not correct
		// iteration order is indeterminate
		return atom.SeqEquals(m, val)
	}

	return false
}

func (m hashMap) Length() int {
	if !m.IsNil() {
		cnt := 0
		for _, bucket := range m.table {
			if bucket != nil {
				cnt += len(bucket)
			}
		}
		return cnt
	}

	return 0
}

func (m hashMap) Get(key atom.Atom) atom.Atom {
	if !m.IsNil() {
		bucketIdx := key.Hash() % uint32(len(m.table))
		bucket := m.table[bucketIdx]
		for _, pair := range bucket {
			if key.Equals(pair.key) {
				return pair.val
			}
		}
	}

	return atom.Nil
}

func (m *hashMap) Set(key atom.Atom, val atom.Atom) {
	if !m.IsNil() {
		pair := mapEntry{
			key: key,
			val: val,
		}

		bucketIdx := key.Hash() % uint32(len(m.table))
		bucket := m.table[bucketIdx]
		if bucket != nil {
			for i, currPair := range bucket {
				if key.Equals(currPair.key) {
					bucket[i] = pair
					return
				}
			}
		}

		// not found; add
		if bucket == nil {
			bucket = []mapEntry{pair}
		} else {
			bucket = append(bucket, pair)
		}
		m.table[bucketIdx] = bucket
		m.size++

		// resize if load factor >= 0.7
		if float32(m.size)/float32(len(m.table)) >= 0.7 {
			newLen := len(m.table) * 2
			newTable := make([][]mapEntry, newLen)
			for _, bucket := range m.table {
				if bucket != nil {
					for _, pair := range bucket {
						bucketIdx := pair.key.Hash() % uint32(newLen)
						bucket := newTable[bucketIdx]
						if bucket == nil {
							bucket = []mapEntry{pair}
						} else {
							bucket = append(bucket, pair)
						}
						newTable[bucketIdx] = bucket
					}
				}
			}

			m.table = newTable
		}
	}
}

type hashMapIterator struct {
	table              [][]mapEntry
	bucketIdx, pairIdx int
}

func (m hashMap) Iterator() atom.SeqIterator {
	return &hashMapIterator{
		table:     m.table,
		bucketIdx: 0,
		pairIdx:   -1,
	}
}

func (it *hashMapIterator) Next() bool {
	if it == nil || len(it.table) == 0 {
		return false
	}

	var bucket []mapEntry
	if it.bucketIdx < len(it.table) {
		bucket = it.table[it.bucketIdx]
		it.pairIdx++
		for bucket == nil || it.pairIdx >= len(bucket) {
			it.bucketIdx++
			if it.bucketIdx >= len(it.table) {
				return false
			} else {
				bucket = it.table[it.bucketIdx]
				it.pairIdx = 0
			}
		}

		if it.pairIdx < len(bucket) {
			return true
		}
	}

	return false
}

func (it *hashMapIterator) Value() atom.Atom {
	if it != nil && len(it.table) > 0 &&
		it.bucketIdx < len(it.table) {
		bucket := it.table[it.bucketIdx]
		pair := bucket[it.pairIdx]

		vec := atom.NewVec()
		vec.Items = append(vec.Items, pair.key, pair.val)
		return vec
	}

	return nil
}
*/
