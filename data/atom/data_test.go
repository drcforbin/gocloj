package atom

import (
	"testing"
)

func TestBasicAtoms(t *testing.T) {
	// TODO: Const
	// TODO: Char
	// TODO: Num
	// TODO: SymName
}

func TestKeywordEquals(t *testing.T) {
	// Keyword
	kw1 := &Keyword{Name: ":xyz"}
	kw2 := &Keyword{Name: ":xyz"}
	kw3 := &Keyword{Name: ":abc"}
	sym := &SymName{Name: ":xyz"}
	str := &Str{Val: ":xyz"}

	if !kw1.Equals(kw2) {
		t.Error("expected keywords to be equal")
	}
	if kw1.Equals(kw3) {
		t.Error("expected keywords to differ")
	}
	if kw1.Equals(sym) {
		t.Error("expected keyword to differ from sym")
	}
	if kw1.Equals(str) {
		t.Error("expected keyword to differ from str")
	}
}
