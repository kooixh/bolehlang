package interpreter

import "github.com/kooixh/bolehlang/parser"

// Control flow signals (using panics caught at boundaries)

type returnSignal struct {
	Value interface{}
}

type breakSignal struct{}

type continueSignal struct{}

type throwSignal struct {
	Value interface{}
}

// Function value (user-defined)

type Function struct {
	Name    string
	Params  []string
	Body    []parser.Node
	Closure *Environment
}

// BuiltinFunc is a built-in callable.

type BuiltinFunc struct {
	Name string
	Fn   func(args []interface{}) (interface{}, error)
}
