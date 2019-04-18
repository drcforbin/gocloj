package lib

import (
	"gocloj/runtime"
)

func AddCore(env *runtime.Env) {
	env.SetInternalFn("if", if_)
	env.SetInternalFn("quote", quote)
	env.SetInternalFn("def", def)
	env.SetInternalFn("do", do)
	env.SetInternalFn("let", let)
	env.SetInternalFn("fn", fn)
	env.SetInternalFn("assert", assert)
	env.SetInternalFn("=", op_eq)
	env.SetInternalFn("not=", op_neq)
	env.SetInternalFn("not", op_not)
}
