package lib

import (
	"errors"
	"fmt"
	"gocloj/data"
	"gocloj/data/atom"
	// "gocloj/log"
	"gocloj/runtime"
)

type callablefn struct {
	binding *atom.Vec
	body    *atom.Vec
}

func (c callablefn) Exec(env *runtime.Env, args atom.SeqIterator) (res atom.Atom, err error) {
	res = atom.Nil

	env.PushScope()
	defer env.PopScope()

	for _, binding := range c.binding.Items {
		if !args.Next() {
			err = errors.New("missing required argument")
			return
		}

		var value atom.Atom
		if value, err = env.Eval(args.Value()); err != nil {
			return
		}

		if err = env.Destructure(binding, value); err != nil {
			return
		}
	}

	// pass body on to do
	it := c.body.Iterator()
	res, err = do(env, it)

	return
}

func (c callablefn) String() string {
	// TODO: format fn
	return fmt.Sprintf("(fn ???)")
}

func (c callablefn) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (c callablefn) Hash() uint32 {
	return c.binding.Hash() + c.body.Hash()
}

// Returns whether this Atom is equivalent to a given atom.
func (c callablefn) Equals(atom atom.Atom) bool {
	if val, ok := atom.(*callablefn); ok {
		return c.binding.Equals(val.binding) &&
			c.body.Equals(val.body)
	}

	return false
}

func fn(env *runtime.Env, args atom.SeqIterator) (res atom.Atom, err error) {
	res = atom.Nil

	if !args.Next() {
		err = errors.New("fn requires an args")
		return
	}

	// TODO: metadata
	val := args.Value()
	var binding *atom.Vec
	var ok bool
	if binding, ok = val.(*atom.Vec); !ok {
		err = errors.New("expected fn arg to be a vec of args")
		return
	}

	if err = data.ValidBinding(binding); err != nil {
		return
	}

	body := atom.NewVec()
	for args.Next() {
		body.Items = append(body.Items, args.Value())
	}

	res = &callablefn{
		binding: binding,
		body:    body,
	}

	return
}
