package runtime

import (
	"gocloj/data/atom"
	"math/big"
	"testing"
)

func TestEvalBasicAtoms(t *testing.T) {
	atoms := []atom.Atom{
		&atom.Const{Name: "alice"},
		&atom.Str{Val: "bob"},
		&atom.Char{Val: 'c'},
		&atom.Num{Val: big.NewInt(1000)},
	}

	for i, a := range atoms {
		env := NewEnv()
		if res, err := env.Eval(a); err != nil {
			t.Errorf("error evaluating atom %d: %s", i, err)
		} else if !a.Equals(res) {
			t.Errorf("eval of atom %d was %s", i, res)
		}
	}
}
