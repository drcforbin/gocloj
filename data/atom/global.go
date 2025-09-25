package atom

// Global nil constant
// This atom should be used instead of golang's nil on the lisp side.
var Nil = &Const{Name: "nil"}

// Global true constant
var True = &Const{Name: "true"}

// Global false constant
var False = &Const{Name: "false"}
