package hashmap

import (
	"fmt"
	"gocloj/data/atom"
)

// mapEntry stores a single entry for a map. Both key and
// val should be data.Nil rather than nil if they are to
// be exposed.
type mapEntry struct {
	key atom.Atom
	val atom.Atom
}

func (m mapEntry) String() string {
	return fmt.Sprintf("(%s, %s)", m.key, m.val)
}
