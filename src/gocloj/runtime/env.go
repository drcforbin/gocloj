package runtime

import (
	"errors"
	"fmt"
	"gocloj/data"
	"gocloj/log"
	"reflect"
)

var envLogger = log.Get("env")

// TODO: consider switching args to []*data.Pair
type Callable interface {
	fmt.Stringer
	Exec(env *Env, args data.SeqIterator) (res data.Atom, err error)
}

type CallableFn = func(env *Env, args data.SeqIterator) (res data.Atom, err error)

type callableThunk struct {
	fn CallableFn
}

func (thunk callableThunk) Exec(env *Env, args data.SeqIterator) (data.Atom, error) {
	return thunk.fn(env, args)
}

func (thunk callableThunk) String() string {
	return fmt.Sprintf("callable thunk %s", thunk.fn)
}

func (thunk callableThunk) IsNil() bool {
	return false
}

func (thunk callableThunk) Hash() uint32 {
	return uint32(reflect.ValueOf(thunk.fn).Pointer())
}

type Env struct {
	// TODO: namespaces
	// TODO: make internal a slice, to implement dynamic scope
	internal map[string]data.Atom
	scope    []map[string]data.Atom
}

func NewEnv() *Env {
	env := &Env{
		internal: map[string]data.Atom{},
		scope:    []map[string]data.Atom{},
	}

	// set up globals

	env.internal["true"] = data.True
	env.internal["false"] = data.False

	return env
}

func (env *Env) ResolveSym(symName *data.SymName) (atom data.Atom, err error) {
	var ok bool

	// check scope next
	for i := len(env.scope) - 1; i >= 0; i-- {
		scope := env.scope[i]
		atom, ok = scope[symName.Name]
		if ok {
			break
		}
	}

	if !ok {
		atom, ok = env.internal[symName.Name]
	}

	if !ok {
		err = errors.New(fmt.Sprintf(
			"unable to find symbol %s", symName.Name))
		// atom = symName
	}

	return
}

func (env *Env) evalVec(vec *data.Vec) (res data.Atom, err error) {
	res = data.Nil
	resVec := data.NewVec()

	var atom data.Atom
	it := vec.Iterator()
	for it.Next() {
		if atom, err = env.Eval(it.Value()); err != nil {
			return
		}

		resVec.Items = append(resVec.Items, atom)
	}

	res = resVec

	return
}

func (env *Env) evalSeq(lst *data.List) (res data.Atom, err error) {
	res = data.Nil

	// evaluate the first argument
	it := lst.Iterator()

	// empty list, evals to self
	if !it.Next() {
		res = lst
		return
	}

	first := it.Value()
	if res, err = env.Eval(first); err != nil {
		return
	}

	if call, ok := res.(Callable); ok {
		res, err = call.Exec(env, it)
		if err != nil {
			envLogger.Debugf("no res, err %s", err.Error())
		} else if !res.IsNil() {
			envLogger.Debugf("res %s, no err", res.String())
		} else {
			envLogger.Debug("no res, no err")
		}
	} else {
		err = errors.New("unable to cast first item of list to Callable")
	}

	return
}

func (env *Env) Eval(atom data.Atom) (res data.Atom, err error) {
	res = data.Nil

	if atom == nil || atom.IsNil() {
		return
	}

	switch v := atom.(type) {
	case *data.Const:
		res = atom

	case *data.Str:
		res = atom

	case *data.Num:
		res = atom

	case *data.SymName:
		res, err = env.ResolveSym(v)

	case *data.Vec:
		res, err = env.evalVec(v)

	case *data.List:
		res, err = env.evalSeq(v)

	default:
		err = errors.New(fmt.Sprintf("unexpected type in eval %T", atom))
	}

	return
}

func (env *Env) Destructure(binding data.Atom, value data.Atom) (err error) {
	switch binding := binding.(type) {
	case *data.SymName:
		if len(env.scope) > 0 {
			env.scope[len(env.scope)-1][binding.Name] = value
		} else {
			err = errors.New("attempted to set scope when scope stack was empty")
		}

	case *data.Vec:
		// TODO: handle &, :as, etc.

		// we need to make sure that value is a sequence too
		if seq, ok := value.(data.Seq); ok {
			valit := seq.Iterator()

			for i, b := range binding.Items {
				handled := false
				var valval data.Atom
				valval = data.Nil

				if sym, ok := b.(*data.SymName); ok {
					switch sym.Name {
					case "&":
						if i != len(binding.Items)-2 {
							err = errors.New("destructuring & requires a symbol")
							return
						}

						vec := data.NewVec()
						for valit.Next() {
							vec.Items = append(vec.Items, valit.Value())
						}

						if err = env.Destructure(
							binding.Items[len(binding.Items)-1],
							vec); err != nil {
							return
						}

						handled = true
					}
				}

				if !handled {
					if valit.Next() {
						valval = valit.Value()
					}

					if err = env.Destructure(b, valval); err != nil {
						return
					}
				}
			}
		} else {
			err = errors.New("destructuring vector, unable to iterate over arg")
		}

		// TODO: destructure map

	default:
		err = errors.New(fmt.Sprintf(
			"unexpected type for destructure binding %T", binding))
	}
	return
}

func (env *Env) PushScope() {
	env.scope = append(env.scope, map[string]data.Atom{})
}

func (env *Env) PopScope() {
	env.scope = env.scope[:len(env.scope)-1]
}

func (env *Env) SetInternal(name *data.SymName, val data.Atom) {
	env.internal[name.Name] = val
}

func (env *Env) SetInternalFn(name string, fn CallableFn) {
	symName := &data.SymName{Name: name}
	env.SetInternal(symName, callableThunk{fn: fn})
}
