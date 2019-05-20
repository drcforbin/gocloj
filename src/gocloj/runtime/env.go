package runtime

import (
	"errors"
	"fmt"
	//"gocloj/data"
	"gocloj/data/atom"
	"gocloj/data/hashmap"
	"gocloj/log"
	"reflect"
)

var envLogger = log.Get("env")

type Callable interface {
	fmt.Stringer
	Exec(env *Env, args atom.SeqIterator) (res atom.Atom, err error)
}

type CallableFn = func(env *Env, args atom.SeqIterator) (res atom.Atom, err error)

type callableThunk struct {
	fn CallableFn
}

func (thunk callableThunk) Exec(env *Env, args atom.SeqIterator) (atom.Atom, error) {
	return thunk.fn(env, args)
}

func (thunk callableThunk) String() string {
	return fmt.Sprintf("callable thunk %p", thunk.fn)
}

func (thunk callableThunk) IsNil() bool {
	return false
}

// Returns a hash value for this Atom.
func (thunk callableThunk) Hash() uint32 {
	return uint32(reflect.ValueOf(thunk.fn).Pointer())
}

// Returns whether this Atom is equivalent to a given atom.
func (thunk callableThunk) Equals(atom atom.Atom) bool {
	if val, ok := atom.(*callableThunk); ok {
		return reflect.ValueOf(thunk.fn).Pointer() ==
			reflect.ValueOf(val.fn).Pointer()
	}

	return false
}

type Env struct {
	// TODO: namespaces
	// TODO: make internal a slice, to implement dynamic scope
	internal map[string]atom.Atom
	scope    []map[string]atom.Atom
}

func NewEnv() *Env {
	env := &Env{
		internal: map[string]atom.Atom{},
		scope:    []map[string]atom.Atom{},
	}

	// set up globals

	env.internal["true"] = atom.True
	env.internal["false"] = atom.False

	return env
}

func (env *Env) ResolveSym(symName *atom.SymName) (res atom.Atom, err error) {
	var ok bool

	// check scope next
	for i := len(env.scope) - 1; i >= 0; i-- {
		res, ok = env.scope[i][symName.Name]
		if ok {
			break
		}
	}

	if !ok {
		res, ok = env.internal[symName.Name]
	}

	if !ok {
		err = errors.New(fmt.Sprintf(
			"unable to find symbol %s", symName.Name))
	}

	return
}

func (env *Env) evalMap(m hashmap.Map) (res atom.Atom, err error) {
	res = atom.Nil
	resMap := hashmap.NewPersistentHashMap()

	var pairAtom atom.Atom
	it := m.Iterator()
	for it.Next() {
		// assumes map is well formed, where each value is a vec pair

		if pairAtom, err = env.Eval(it.Value()); err != nil {
			return
		}

		pair := pairAtom.(*atom.Vec)
		resMap = resMap.Assoc(pair.Items[0], pair.Items[1])
	}

	res = resMap

	return
}

func (env *Env) evalVec(vec *atom.Vec) (res atom.Atom, err error) {
	res = atom.Nil
	resVec := atom.NewVec()

	var itemAtom atom.Atom
	it := vec.Iterator()
	for it.Next() {
		if itemAtom, err = env.Eval(it.Value()); err != nil {
			return
		}

		resVec.Items = append(resVec.Items, itemAtom)
	}

	res = resVec

	return
}

func (env *Env) evalList(lst *atom.List) (res atom.Atom, err error) {
	res = atom.Nil

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

func (env *Env) Eval(val atom.Atom) (res atom.Atom, err error) {
	res = atom.Nil

	if val == nil || val.IsNil() {
		return
	}

	switch v := val.(type) {
	case *atom.Const, *atom.Str, *atom.Char, *atom.Num, *atom.Keyword:
		res = val

	case *atom.SymName:
		res, err = env.ResolveSym(v)

	case *atom.Vec:
		res, err = env.evalVec(v)

	case *atom.List:
		res, err = env.evalList(v)

	case hashmap.Map:
		res, err = env.evalMap(v)

	default:
		err = errors.New(fmt.Sprintf("unexpected type in eval %T", val))
	}

	return
}

func (env *Env) Destructure(binding atom.Atom, value atom.Atom) (err error) {
	switch binding := binding.(type) {
	case *atom.SymName:
		if len(env.scope) > 0 {
			env.scope[len(env.scope)-1][binding.Name] = value
		} else {
			err = errors.New("attempted to set scope when scope stack was empty")
		}

	case *atom.Vec:
		// we need to make sure that value is a sequence too
		if seq, ok := value.(atom.Seq); ok {
			valit := seq.Iterator()

			for i := 0; i < len(binding.Items); i++ {
				b := binding.Items[i]

				handled := false
				var valval atom.Atom = atom.Nil

				if sym, ok := b.(*atom.SymName); ok {
					switch sym.Name {
					case "&":
						if i > len(binding.Items)-2 {
							err = errors.New("destructuring & requires a symbol")
							return
						}

						// get / consume the symbol
						i++
						b = binding.Items[i]

						if _, ok := b.(*atom.SymName); ok {
							vec := atom.NewVec()
							for valit.Next() {
								vec.Items = append(vec.Items, valit.Value())
							}

							if err = env.Destructure(b, vec); err != nil {
								return
							}

							handled = true
						}
					}
				} else if kw, ok := b.(*atom.Keyword); ok {
					switch kw.Name {
					case ":as":
						if i > len(binding.Items)-2 {
							err = errors.New("destructuring :as requires a symbol")
							return
						}

						// get / consume the symbol
						i++
						b = binding.Items[i]

						// bind symbol to incoming value
						if err = env.Destructure(b, value); err != nil {
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

	case hashmap.Map:
		// TODO: destructure map
		// TODO: :keys, :strs and :syms

	default:
		err = errors.New(fmt.Sprintf(
			"unexpected type for destructure binding %T", binding))
	}
	return
}

func (env *Env) PushScope() {
	env.scope = append(env.scope, map[string]atom.Atom{})
}

func (env *Env) PopScope() {
	env.scope = env.scope[:len(env.scope)-1]
}

func (env *Env) SetInternal(name *atom.SymName, val atom.Atom) {
	env.internal[name.Name] = val
}

func (env *Env) SetInternalFn(name string, fn CallableFn) {
	symName := &atom.SymName{Name: name}
	env.SetInternal(symName, callableThunk{fn: fn})
}
