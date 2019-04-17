package data

import (
	"strings"
)

// TODO: this hashmap is fake. implement it for real

type hashPair struct {
	key Atom
	val Atom
}

type HashMap struct {
	items []hashPair
}

func NewHashMap() *HashMap {
	return &HashMap{
		items: []hashPair{},
	}
}

func (m HashMap) String() string {
	var builder strings.Builder

	builder.WriteString("{")
	for i, pair := range m.items {
		if i != 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(pair.key.String())
		builder.WriteString(" ")
		builder.WriteString(pair.val.String())
	}
	builder.WriteString("}")

	return builder.String()
}

func (m *HashMap) IsNil() bool {
	return false
}

func (m HashMap) Hash() uint32 {
	hash := uint32(0)
	count := uint32(0)

	it := m.Iterator()
	for it.Next() {
		hash += it.Value().Hash()

		count++
	}

	return mixCollHash(hash, count)
}

func (m HashMap) Length() int {
	if !m.IsNil() {
		return len(m.items)
	}

	return 0
}

func (m HashMap) Get(key Atom) Atom {
	if !m.IsNil() {
		for _, pair := range m.items {
			if Equals(key, pair.key) {
				return pair.val
			}
		}
	}

	return Nil
}

func (m HashMap) Set(key Atom, val Atom) {
	if !m.IsNil() {
		pair := hashPair{
			key: key,
			val: val,
		}

		for i, pair := range m.items {
			if Equals(key, pair.key) {
				m.items[i] = pair
				return
			}
		}

		m.items = append(m.items, pair)
	}
}

type mapIterator struct {
	items []hashPair
	idx   int
}

func (m HashMap) Iterator() SeqIterator {
	return &mapIterator{
		items: m.items,
		idx:   -1,
	}
}

func (it *mapIterator) Next() bool {
	if it != nil && it.items != nil && len(it.items) > 0 &&
		it.idx < len(it.items)-1 {
		it.idx++
		return true
	}

	return false
}

func (it *mapIterator) Value() Atom {
	if it != nil && it.items != nil && len(it.items) > 0 &&
		it.idx < len(it.items) {
		vec := NewVec()
		pair := it.items[it.idx]
		vec.Items = append(vec.Items, pair.key, pair.val)
		return vec
	}

	return nil
}
