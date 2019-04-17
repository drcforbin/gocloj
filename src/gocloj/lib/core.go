package lib

import (
	"errors"
	// "fmt"
	"gocloj/data"
	"gocloj/log"
	"gocloj/runtime"
)

var coreLogger = log.Get("core")

func if_(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil

	if !args.Next() {
		err = errors.New("if missing predicate")
		return
	}

	// evaluate the test
	var testRes data.Atom
	if testRes, err = env.Eval(args.Value()); err != nil {
		return
	}

	// find true clause
	var truePart data.Atom
	if args.Next() {
		val := args.Value()
		truePart = val
	}

	if !truePart.IsNil() {
		// if the test is true, return args.Val.Key
		if data.Truthy(testRes) {
			res, err = env.Eval(truePart)
		} else {
			// find false clause
			var falsePart data.Atom
			if args.Next() {
				val := args.Value()
				falsePart = val
			}

			res, err = env.Eval(falsePart)
		}
	} else {
		err = errors.New("if missing true clause")
	}

	return
}

func quote(env *runtime.Env, args data.SeqIterator) (data.Atom, error) {
	if args.Next() {
		val := args.Value()
		if !val.IsNil() {
			return val, nil
		}
	}

	return data.Nil, nil
}

/*
func cmpAtom(expected, actual data.Atom) bool {
	if expected.IsNil() || actual.IsNil() {
		if expected.IsNil() != actual.IsNil() {
			return false
		} else {
			// both nil
		}
	} else {
		var ok bool

		switch v1 := expected.(type) {
		case *data.Num:
			var v2 *data.Num
			v2, ok = actual.(*data.Num)
			if !ok || v1.Val.Cmp(v2.Val) != 0 {
				return false
			}

		case *data.SymName:
			var v2 *data.SymName
			v2, ok = actual.(*data.SymName)
			if !ok || v1.Name != v2.Name {
				return false
			}

		case *data.Pair:
			var v2 *data.Pair
			v2, ok = actual.(*data.Pair)
			if !ok ||
				!cmpAtom(v1.Key, v2.Key) ||
				!cmpAtom(v1.Val, v2.Val) {
				return false
			}

		default:
			return false
		}

	}

	return true
}
*/

func def(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil

	if !args.Next() {
		res = &data.Str{Val: "unexpected nil args for def"}
		return
	}

	val := args.Value()
	if sym, ok := val.(*data.SymName); ok {
		// advance to next arg
		if args.Next() {
			var val data.Atom
			if val, err = env.Eval(args.Value()); err != nil {
				return
			}

			env.SetInternal(sym, val)

			res = sym
		} else {
			err = errors.New("def must have at least two args")
		}
	} else {
		err = errors.New("expected def arg to be a symbol")
	}

	return
}

func do(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil

	for args.Next() {
		if res, err = env.Eval(args.Value()); err != nil {
			return
		}
	}

	return
}

func let(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil

	if !args.Next() {
		err = errors.New("let requires an args")
		return
	}

	// TODO: metadata
	val := args.Value()
	var vec *data.Vec
	var ok bool
	if vec, ok = val.(*data.Vec); !ok {
		err = errors.New("expected let arg to be a vec of args")
		return
	}

	if len(vec.Items)%2 != 0 {
		err = errors.New("let requires even number of args")
		return
	}

	env.PushScope()
	defer env.PopScope()

	for i := 0; i < len(vec.Items); i += 2 {
		binding := vec.Items[i]
		if err = data.ValidBinding(binding); err != nil {
			return
		}

		value := vec.Items[i+1]
		if value, err = env.Eval(value); err != nil {
			return
		}
		coreLogger.Info("let ", binding, value)

		if err = env.Destructure(binding, value); err != nil {
			return
		}
	}

	// pass body on to do
	res, err = do(env, args)

	return
}
