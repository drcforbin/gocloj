package lib

import (
	"errors"
	"fmt"
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

		if err = env.Destructure(binding, value); err != nil {
			return
		}
	}

	// pass body on to do
	res, err = do(env, args)

	return
}

func assert(env *runtime.Env, args data.SeqIterator) (data.Atom, error) {
	// Evaluates expr and throws an exception if it does not evaluate to
	// logical true.
	/*
	  ([x]
	     (when *assert*
	       `(when-not ~x
	          (throw (new AssertionError (str "Assert failed: " (pr-str '~x)))))))
	  ([x message]
	     (when *assert*
	       `(when-not ~x
	          (throw (new AssertionError (str "Assert failed: " ~message "\n" (pr-str '~x))))))))
	*/
	if !args.Next() {
		return data.Nil, errors.New("assert requires an arg")
	}

	x := args.Value()
	if value, err := env.Eval(x); err == nil {
		if !data.Truthy(value) {
			// do we have a message?
			if args.Next() {
				panic(fmt.Sprintf("assert failed: %s; %s",
					args.Value().String(), x.String()))
			} else {
				// no message
				panic(fmt.Sprintf("assert failed: %s", x.String()))
			}
		}
	} else {
		return nil, err
	}

	return data.Nil, nil
}
