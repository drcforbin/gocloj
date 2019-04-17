package lib

import (
	"errors"
	"fmt"
	"gocloj/data"
	// "gocloj/log"
	"gocloj/runtime"
)

type callablefn struct {
	binding *data.Vec
	body    *data.Vec
}

func (c callablefn) Exec(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil

	env.PushScope()
	defer env.PopScope()

	for _, binding := range c.binding.Items {
		if !args.Next() {
			err = errors.New("missing required argument")
			return
		}

		var value data.Atom
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

func (c callablefn) Hash() uint32 {
	return c.binding.Hash() + c.body.Hash()
}

func fn(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil

	if !args.Next() {
		err = errors.New("fn requires an args")
		return
	}

	// TODO: metadata
	val := args.Value()
	var binding *data.Vec
	var ok bool
	if binding, ok = val.(*data.Vec); !ok {
		err = errors.New("expected fn arg to be a vec of args")
		return
	}

	if err = data.ValidBinding(binding); err != nil {
		return
	}

	body := data.NewVec()
	for args.Next() {
		body.Items = append(body.Items, args.Value())
	}

	res = &callablefn{
		binding: binding,
		body:    body,
	}

	return
}
