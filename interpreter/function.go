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
	Name       string
	Params     []parser.Param
	ReturnType string
	Body       []parser.Node
	Closure    *Environment
}

// BuiltinFunc is a built-in callable.

type BuiltinFunc struct {
	Name string
	Fn   func(args []interface{}) (interface{}, error)
}

// Exception represents a typed error in bolehlang.

type Exception struct {
	Type    string
	Message string
}

// Default exception types
const (
	ExMakKauHijau = "MakKauHijauException" // NullPointer
	ExBoJio       = "BoJioException"       // ImportNotFound
	ExApaNi       = "ApaNiException"       // Invalid types
	ExKantoi      = "KantoiException"      // General runtime error
)


// structdef represents the blueprint/constructor of a kroni
type StructDef struct {
	Name string
	Fields []string
}


// represents an instantiated kroni
type StructInstance struct {
	Def *StructDef
	Fields map[string]interface{}
}
