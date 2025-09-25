package lib

import (
	"errors"
	"gocloj/data"
	"gocloj/data/atom"
	"gocloj/runtime"
)

func op_eq(env *runtime.Env, args atom.SeqIterator) (res atom.Atom, err error) {
	res = atom.False

	var a, b atom.Atom

	if !args.Next() {
		return atom.Nil, errors.New("= requires two args")
	}

	if a, err = env.Eval(args.Value()); err != nil {
		return
	}

	if !args.Next() {
		return atom.Nil, errors.New("= requires two args")
	}

	if b, err = env.Eval(args.Value()); err != nil {
		return
	}

	if a.Equals(b) {
		res = atom.True
	}

	return
}

func op_neq(env *runtime.Env, args atom.SeqIterator) (res atom.Atom, err error) {
	res = atom.False

	var a, b atom.Atom

	if !args.Next() {
		return atom.Nil, errors.New("not= requires two args")
	}

	if a, err = env.Eval(args.Value()); err != nil {
		return
	}

	if !args.Next() {
		return atom.Nil, errors.New("= requires two args")
	}

	if b, err = env.Eval(args.Value()); err != nil {
		return
	}

	if !a.Equals(b) {
		res = atom.True
	}

	return
}
func op_not(env *runtime.Env, args atom.SeqIterator) (res atom.Atom, err error) {
	res = atom.False

	if !args.Next() {
		return atom.Nil, errors.New("not requires an arg")
	}

	var val atom.Atom
	if val, err = env.Eval(args.Value()); err != nil {
		return
	}

	if !data.Truthy(val) {
		res = atom.True
	}

	return
}
