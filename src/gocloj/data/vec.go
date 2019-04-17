package data

import (
	"strings"
)

type Vec struct {
	Items []Atom
}

func NewVec() *Vec {
	return &Vec{Items: []Atom{}}
}

func (v Vec) String() string {
	var builder strings.Builder

	builder.WriteString("[")
	for i, val := range v.Items {
		if i != 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(val.String())
	}
	builder.WriteString("]")

	return builder.String()
}

func (v *Vec) IsNil() bool {
	return false
}

func (v Vec) Hash() uint32 {
	hash := uint32(1)

	for _, item := range v.Items {
		hash += (31 * hash) + item.Hash()
	}

	return mixCollHash(hash, uint32(len(v.Items)))
}

func (v Vec) Length() int {
	if !v.IsNil() {
		return len(v.Items)
	}

	return 0
}

func (v Vec) Item(idx int) Atom {
	if !v.IsNil() {
		return v.Items[idx]
	}

	return nil
}

type vecIterator struct {
	items []Atom
	idx   int
}

func (v Vec) Iterator() SeqIterator {
	return &vecIterator{items: v.Items, idx: -1}
}

func (it *vecIterator) Next() bool {
	if it != nil && it.items != nil && len(it.items) > 0 &&
		it.idx < len(it.items)-1 {
		it.idx++
		return true
	}

	return false
}

func (it *vecIterator) Value() Atom {
	if it != nil && it.items != nil && len(it.items) > 0 &&
		it.idx < len(it.items) {
		return it.items[it.idx]
	}

	return nil
}
