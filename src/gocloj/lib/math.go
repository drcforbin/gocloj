package lib

import (
	"errors"
	"gocloj/data"
	"gocloj/log"
	"gocloj/runtime"
	"math/big"
)

var mathLogger = log.Get("math")

func add(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil
	sum := &big.Int{}

	var (
		num *data.Num
		ok  bool
	)

	for args.Next() {
		val := args.Value()
		if num, ok = val.(*data.Num); ok {
			sum.Add(sum, num.Val)
		} else {
			err = errors.New("+ arg did not evaluate to a num")
			return
		}
	}

	res = &data.Num{Val: sum}
	return
}

func mul(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil
	total := &big.Int{}

	var (
		num *data.Num
		ok  bool
	)

	// mathLogger.Infof("args %T %+v", args, args.String())

	first := true
	for args.Next() {
		val := args.Value()
		if num, ok = val.(*data.Num); ok {
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

	res = &data.Num{Val: total}
	return
}

func inc(env *runtime.Env, args data.SeqIterator) (res data.Atom, err error) {
	res = data.Nil

	val := &big.Int{}
	one := big.NewInt(1)

	var (
		num *data.Num
		ok  bool
	)

	if args.Next() {
		if num, ok = args.Value().(*data.Num); ok {
			val.Set(num.Val)
			val = val.Add(val, one)
		} else {
			err = errors.New("inc arg did not evaluate to a num")
			return
		}
	} else {
		err = errors.New("inc requires a num")
	}

	res = &data.Num{Val: val}
	return
}

func AddMath(env *runtime.Env) {
	env.SetInternalFn("+", add)
	env.SetInternalFn("*", mul)
	env.SetInternalFn("inc", inc)
}
