package lib

import (
	"errors"
	"gocloj/data"
	"gocloj/runtime"
)

func op_eq(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.False

	var a, b data.Atom

	if !args.Next() {
		return data.Nil, errors.New("= requires two args")
	}

	if a, err = env.Eval(args.Value()); err != nil {
		return
	}

	if !args.Next() {
		return data.Nil, errors.New("= requires two args")
	}

	if b, err = env.Eval(args.Value()); err != nil {
		return
	}

	if a.Equals(b) {
		res = data.True
	}

	return
}

func op_neq(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.False

	var a, b data.Atom

	if !args.Next() {
		return data.Nil, errors.New("not= requires two args")
	}

	if a, err = env.Eval(args.Value()); err != nil {
		return
	}

	if !args.Next() {
		return data.Nil, errors.New("= requires two args")
	}

	if b, err = env.Eval(args.Value()); err != nil {
		return
	}

	if !a.Equals(b) {
		res = data.True
	}

	return
}
func op_not(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.False

	if !args.Next() {
		return data.Nil, errors.New("not requires an arg")
	}

	var val data.Atom
	if val, err = env.Eval(args.Value()); err != nil {
		return
	}

	if !data.Truthy(val) {
		res = data.True
	}

	return
}
