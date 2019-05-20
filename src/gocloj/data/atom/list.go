package atom

import (
	"strings"
)

type ListNode struct {
	Value Atom
	Next  *ListNode
}

type List struct {
	Head *ListNode
}

func NewList() *List {
	return &List{}
}

func (l List) String() string {
	var builder strings.Builder
	builder.WriteString("(")

	first := true
	it := l.Iterator()
	for it.Next() {
		if !first {
			builder.WriteString(" ")
		} else {
			first = false
		}

		builder.WriteString(it.Value().String())
	}
	/*
		if l.Head.Next != nil && l.Head.Next.Value !={
			// both are non-nil
			if item, ok := p.Val.(*Pair); ok {
				// if it's a pair, display simply
				builder.WriteString(p.Key.String())

				for item != nil && item.Key != nil {
					builder.WriteString(" ")
					builder.WriteString(item.Key.String())

					if item, ok = item.Val.(*Pair); !ok {
						break
					}
				}

				builder.WriteString(")")
			} else {
				// not a pair, write it out
				builder.WriteString("(")
				builder.WriteString(p.Key.String())
				builder.WriteString(" ")
				builder.WriteString(p.Val.String())
			}
		} else {
			// only key is non-nil
			builder.WriteString(p.Key.String())
		}*/

	builder.WriteString(")")
	return builder.String()
}

func (l *List) IsNil() bool {
	// TODO: this right?
	return false
}

// Returns a hash value for this Atom.
func (l *List) Hash() uint32 {
	hash := uint32(1)
	count := uint32(0)

	node := l.Head
	for node != nil {
		hash += (31 * hash) + node.Value.Hash()
		count++

		node = node.Next
	}

	return mixCollHash(hash, count)
}

// Returns whether this Atom is equivalent to a given atom.
func (l *List) Equals(atom Atom) bool {
	if val, ok := atom.(Seq); ok {
		return SeqEquals(l, val)
	}

	return false
}

func (l *List) Length() int {
	count := 0

	node := l.Head
	for node != nil {
		count++

		node = node.Next
	}

	return count
}

func (l *List) Item(idx int) Atom {
	it := l.Iterator()

	for ; idx >= 0 && it.Next(); idx-- {
		if idx == 0 {
			return it.Value()
		}
	}

	return nil
}

type listIterator struct {
	list *List
	item *ListNode
}

func (l *List) Iterator() SeqIterator {
	return &listIterator{list: l}
}

func (it *listIterator) Next() bool {
	if it != nil {
		if it.list != nil && it.item == nil {
			it.item = it.list.Head
			it.list = nil
		} else if it.item != nil {
			it.item = it.item.Next
		}

		return it.item != nil
	}

	return false
}

func (it *listIterator) Value() Atom {
	if it != nil && it.item != nil {
		return it.item.Value
	}
	return nil
}
