package lib

import (
	"errors"
	"gocloj/data/atom"
	"gocloj/log"
	"gocloj/runtime"
	"math/big"
)

var mathLogger = log.Get("math")

func add(env *runtime.Env, args atom.SeqIterator) (res atom.Atom, err error) {
	res = atom.Nil
	sum := &big.Int{}

	var (
		valAtom atom.Atom
		num     *atom.Num
		ok      bool
	)

	for args.Next() {
		if valAtom, err = env.Eval(args.Value()); err != nil {
			return
		}

		if num, ok = valAtom.(*atom.Num); ok {
			sum.Add(sum, num.Val)
		} else {
			err = errors.New("+ arg did not evaluate to a num")
			return
		}
	}

	res = &atom.Num{Val: sum}
	return
}

func mul(env *runtime.Env, args atom.SeqIterator) (res atom.Atom, err error) {
	res = atom.Nil
	total := &big.Int{}

	var (
		valAtom atom.Atom
		num     *atom.Num
		ok      bool
	)

	// mathLogger.Infof("args %T %+v", args, args.String())

	first := true
	for args.Next() {
		if valAtom, err = env.Eval(args.Value()); err != nil {
			return
		}

		if num, ok = valAtom.(*atom.Num); ok {
			if first {
				total.Set(num.Val)
				first = false
			} else {
				total.Mul(total, num.Val)
			}
		} else {
			err = errors.New("* arg did not evaluate to a num")
			return
		}
	}

	res = &atom.Num{Val: total}
	return
}

func inc(env *runtime.Env, args atom.SeqIterator) (res atom.Atom, err error) {
	res = atom.Nil

	val := &big.Int{}
	one := big.NewInt(1)

	var (
		valAtom atom.Atom
		num     *atom.Num
		ok      bool
	)

	if args.Next() {
		if valAtom, err = env.Eval(args.Value()); err != nil {
			return
		}

		if num, ok = valAtom.(*atom.Num); ok {
			val.Set(num.Val)
			val = val.Add(val, one)
		} else {
			err = errors.New("inc arg did not evaluate to a num")
			return
		}
	} else {
		err = errors.New("inc requires a num")
	}

	res = &atom.Num{Val: val}
	return
}

func AddMath(env *runtime.Env) {
	env.SetInternalFn("+", add)
	env.SetInternalFn("*", mul)
	env.SetInternalFn("inc", inc)
}
