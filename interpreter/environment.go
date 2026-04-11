package interpreter

import "fmt"

// Environment (lexical scoping)

type Environment struct {
	vars      map[string]interface{}
	constants map[string]bool
	parent    *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		vars:      make(map[string]interface{}),
		constants: make(map[string]bool),
		parent:    parent,
	}
}

func (e *Environment) Set(name string, value interface{}, isConst bool) {
	e.vars[name] = value
	if isConst {
		e.constants[name] = true
	}
}

func (e *Environment) Get(name string) (interface{}, bool) {
	if val, ok := e.vars[name]; ok {
		return val, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

func (e *Environment) Assign(name string, value interface{}) error {
	if _, ok := e.vars[name]; ok {
		if e.constants[name] {
			return fmt.Errorf("Aiyo, a memang variable '%s' cannot be redefined", name)
		}
		e.vars[name] = value
		return nil
	}
	if e.parent != nil {
		return e.parent.Assign(name, value)
	}
	return fmt.Errorf("variable '%s' not found — you got declare or not?", name)
}
