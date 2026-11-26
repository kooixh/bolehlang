package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kooixh/bolehlang/lexer"
	"github.com/kooixh/bolehlang/parser"
)

var stdinReader = bufio.NewReader(os.Stdin)

// SetStdin replaces the reader used by the `lai` keyword. Intended for tests.
func SetStdin(r io.Reader) {
	stdinReader = bufio.NewReader(r)
}

type Interpreter struct {
	global   *Environment
	filePath string
	imported map[string]bool
}

func New(filePath string) *Interpreter {
	global := NewEnvironment(nil)
	registerBuiltins(global)

	absPath, _ := filepath.Abs(filePath)
	imported := map[string]bool{absPath: true}
	return &Interpreter{global: global, filePath: absPath, imported: imported}
}

func (interp *Interpreter) Run(program *parser.Program) (retErr error) {
	defer func() {
		if r := recover(); r != nil {
			if ts, ok := r.(throwSignal); ok {
				if exc, ok := ts.Value.(*Exception); ok {
					retErr = fmt.Errorf("%s: %s", exc.Type, exc.Message)
				} else {
					retErr = fmt.Errorf("uncaught exception: %v", ts.Value)
				}
			} else {
				panic(r)
			}
		}
	}()
	env := interp.global
	for _, stmt := range program.Body {
		if err := interp.execStatement(stmt, env); err != nil {
			return err
		}
	}
	return nil
}

func (interp *Interpreter) execImport(n *parser.ImportStatement, env *Environment) error {
	// Resolve path relative to the current file
	dir := filepath.Dir(interp.filePath)
	importPath := filepath.Join(dir, n.Path)
	absPath, err := filepath.Abs(importPath)
	if err != nil {
		return fmt.Errorf("line %d: cannot resolve import path '%s'", n.Line, n.Path)
	}

	// Check for circular imports
	if interp.imported[absPath] {
		return nil // already imported, skip
	}
	interp.imported[absPath] = true

	// Read the file
	source, err := os.ReadFile(absPath)
	if err != nil {
		panic(throwSignal{Value: &Exception{
			Type:    ExBoJio,
			Message: fmt.Sprintf("cannot tapau '%s' — file not found lah!", n.Path),
		}})
	}

	// Lex
	lex := lexer.New(string(source))
	tokens, err := lex.Tokenize()
	if err != nil {
		return fmt.Errorf("tapau '%s': %s", n.Path, err)
	}

	// Parse
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		return fmt.Errorf("tapau '%s': %s", n.Path, err)
	}

	// Execute in the current environment so declarations are available
	savedPath := interp.filePath
	interp.filePath = absPath
	for _, stmt := range program.Body {
		if err := interp.execStatement(stmt, env); err != nil {
			interp.filePath = savedPath
			return fmt.Errorf("tapau '%s': %s", n.Path, err)
		}
	}
	interp.filePath = savedPath
	return nil
}

func (interp *Interpreter) execStatement(node parser.Node, env *Environment) error {
	switch n := node.(type) {
	case *parser.ImportStatement:
		return interp.execImport(n, env)

	case *parser.VarDeclaration:
		val, err := interp.evalExpr(n.Value, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		if n.TypeName != "" {
			checkType(val, n.TypeName, n.Name, n.Line)
		}
		env.Set(n.Name, val, false)
		env.SetType(n.Name, n.TypeName)
		return nil

	case *parser.ConstDeclaration:
		val, err := interp.evalExpr(n.Value, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		if n.TypeName != "" {
			checkType(val, n.TypeName, n.Name, n.Line)
		}
		env.Set(n.Name, val, true)
		env.SetType(n.Name, n.TypeName)
		return nil

	case *parser.MemberAssignment:
		obj , err := interp.evalExpr(n.Object, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}

		instance , ok := obj.(*StructInstance)

		if !ok {
			panic(throwSignal{Value: &Exception{
				Type: ExApaNi,
				Message: fmt.Sprintf("cannot assign property '%s' on %T — bukan kroni lah!", n.Property, obj),
			}})
		}

		if _, exists := instance.Fields[n.Property]; !exists {
			panic(throwSignal{Value: &Exception{
				Type: ExApaNi,
				Message: fmt.Sprintf("kroni '%s' tak ada property '%s' lah", instance.Def.Name, n.Property),
			}})
		}

		val , err := interp.evalExpr(n.Value , env)

		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}

		switch n.Operator{
		case "=" :
			instance.Fields[n.Property] = val

		case "+=":
			instance.Fields[n.Property] = addValues(instance.Fields[n.Property] , val)

		case "-=":
			instance.Fields[n.Property] = toNumber(instance.Fields[n.Property]) - toNumber(val)
		}
		
		return nil 

	case *parser.Assignment:
		val, err := interp.evalExpr(n.Value, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		switch n.Operator {
		case "=":
			if t := env.GetType(n.Name); t != "" {
				checkType(val, t, n.Name, n.Line)
			}
			if err := env.Assign(n.Name, val); err != nil {
				return fmt.Errorf("line %d: %s", n.Line, err)
			}
		case "+=":
			current, ok := env.Get(n.Name)
			if !ok {
				return fmt.Errorf("line %d: variable '%s' not found", n.Line, n.Name)
			}
			result := addValues(current, val)
			if t := env.GetType(n.Name); t != "" {
				checkType(result, t, n.Name, n.Line)
			}
			if err := env.Assign(n.Name, result); err != nil {
				return fmt.Errorf("line %d: %s", n.Line, err)
			}
		case "-=":
			current, ok := env.Get(n.Name)
			if !ok {
				return fmt.Errorf("line %d: variable '%s' not found", n.Line, n.Name)
			}
			result := toNumber(current) - toNumber(val)
			if t := env.GetType(n.Name); t != "" {
				checkType(result, t, n.Name, n.Line)
			}
			if err := env.Assign(n.Name, result); err != nil {
				return fmt.Errorf("line %d: %s", n.Line, err)
			}
		}
		return nil

	case *parser.PrintStatement:
		val, err := interp.evalExpr(n.Value, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		fmt.Println(stringify(val))
		return nil

	case *parser.WarnStatement:
		val, err := interp.evalExpr(n.Value, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		fmt.Fprintf(os.Stderr, "[WALAO EH] %s\n", stringify(val))
		return nil

	case *parser.IfStatement:
		cond, err := interp.evalExpr(n.Condition, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		if isTruthy(cond) {
			return interp.execBlock(n.Body, NewEnvironment(env))
		}
		for _, elseIf := range n.ElseIfCases {
			cond, err := interp.evalExpr(elseIf.Condition, env)
			if err != nil {
				return err
			}
			if isTruthy(cond) {
				return interp.execBlock(elseIf.Body, NewEnvironment(env))
			}
		}
		if n.ElseBody != nil {
			return interp.execBlock(n.ElseBody, NewEnvironment(env))
		}
		return nil

	case *parser.WhileStatement:
		for {
			cond, err := interp.evalExpr(n.Condition, env)
			if err != nil {
				return fmt.Errorf("line %d: %w", n.Line, err)
			}
			if !isTruthy(cond) {
				break
			}
			err = interp.execBlockWithLoopSignals(n.Body, NewEnvironment(env))
			if err == errBreak {
				break
			}
			if err != nil {
				return err
			}
		}
		return nil

	case *parser.ForStatement:
		startVal, err := interp.evalExpr(n.Start, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		endVal, err := interp.evalExpr(n.End, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		start := toNumber(startVal)
		end := toNumber(endVal)

		for i := start; i <= end; i++ {
			loopEnv := NewEnvironment(env)
			loopEnv.Set(n.VarName, i, false)
			err := interp.execBlockWithLoopSignals(n.Body, loopEnv)
			if err == errBreak {
				break
			}
			if err != nil {
				return err
			}
		}
		return nil

	case *parser.ForEachStatement:
		collection, err := interp.evalExpr(n.Collection, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		arr, ok := collection.([]interface{})
		if !ok {
			panic(throwSignal{Value: &Exception{
				Type:    ExApaNi,
				Message: fmt.Sprintf("setiap loop needs a senarai (array), got %T — check your types lah", collection),
			}})
		}
		for _, item := range arr {
			loopEnv := NewEnvironment(env)
			loopEnv.Set(n.VarName, item, false)
			err := interp.execBlockWithLoopSignals(n.Body, loopEnv)
			if err == errBreak {
				break
			}
			if err != nil {
				return err
			}
		}
		return nil

	case *parser.FunctionDeclaration:
		fn := &Function{
			Name:       n.Name,
			Params:     n.Params,
			ReturnType: n.ReturnType,
			Body:       n.Body,
			Closure:    env,
		}
		env.Set(n.Name, fn, false)
		return nil

	case *parser.ReturnStatement:
		val, err := interp.evalExpr(n.Value, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		panic(returnSignal{Value: val})

	case *parser.BreakStatement:
		panic(breakSignal{})

	case *parser.ContinueStatement:
		panic(continueSignal{})

	case *parser.TryCatchStatement:
		err := interp.execTryCatch(n, env)
		return err

	case *parser.ThrowStatement:
		val, err := interp.evalExpr(n.Value, env)
		if err != nil {
			return fmt.Errorf("line %d: %w", n.Line, err)
		}
		panic(throwSignal{Value: val})

	case *parser.ExpressionStatement:
		_, err := interp.evalExpr(n.Expression, env)
		return err


	case *parser.StructDeclaration:
		 def := &StructDef{
			Name: n.Name,
			Fields: n.Fields,
		 }

		 env.Set(n.Name, def, false)
		 return nil


	default:
		return fmt.Errorf("unknown statement type: %T", node)
	}
}

func (interp *Interpreter) execBlock(stmts []parser.Node, env *Environment) error {
	for _, stmt := range stmts {
		if err := interp.execStatement(stmt, env); err != nil {
			return err
		}
	}
	return nil
}

func (interp *Interpreter) execBlockWithLoopSignals(stmts []parser.Node, env *Environment) (retErr error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case breakSignal:
				retErr = errBreak
			case continueSignal:
				retErr = nil // continue just stops the current iteration
			default:
				panic(r) // re-panic for return signals or actual errors
			}
		}
	}()
	return interp.execBlock(stmts, env)
}

var errBreak = fmt.Errorf("__break__")

func (interp *Interpreter) execTryCatch(n *parser.TryCatchStatement, env *Environment) (retErr error) {
	// Try block
	func() {
		defer func() {
			if r := recover(); r != nil {
				if ts, ok := r.(throwSignal); ok {
					// Execute catch block
					if n.CatchBody != nil {
						catchEnv := NewEnvironment(env)
						catchEnv.Set(n.CatchVar, ts.Value, false)
						retErr = interp.execBlock(n.CatchBody, catchEnv)
					}
				} else {
					panic(r) // re-panic non-throw signals
				}
			}
		}()
		retErr = interp.execBlock(n.TryBody, NewEnvironment(env))
	}()

	// Finally block
	if n.FinallyBody != nil {
		if err := interp.execBlock(n.FinallyBody, NewEnvironment(env)); err != nil {
			retErr = err
		}
	}

	return retErr
}

// Expression evaluation

func (interp *Interpreter) evalExpr(node parser.Node, env *Environment) (interface{}, error) {
	switch n := node.(type) {
	case *parser.NumberLiteral:
		return n.Value, nil

	case *parser.StringLiteral:
		return n.Value, nil

	case *parser.BoolLiteral:
		return n.Value, nil

	case *parser.NullLiteral:
		return nil, nil

	case *parser.Identifier:
		val, ok := env.Get(n.Name)
		if !ok {
			panic(throwSignal{Value: &Exception{
				Type:    ExKantoi,
				Message: fmt.Sprintf("'%s' tu apa? Variable not found lah", n.Name),
			}})
		}
		return val, nil

	case *parser.ReadStdinExpr:
		line, err := stdinReader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(throwSignal{Value: &Exception{
				Type:    ExKantoi,
				Message: fmt.Sprintf("lai tak boleh baca stdin: %s", err),
			}})
		}
		return strings.TrimRight(line, "\r\n"), nil

	case *parser.ArrayLiteral:
		var elements []interface{}
		for _, elem := range n.Elements {
			val, err := interp.evalExpr(elem, env)
			if err != nil {
				return nil, err
			}
			elements = append(elements, val)
		}
		if elements == nil {
			elements = []interface{}{}
		}
		return elements, nil

	case *parser.BinaryExpr:
		return interp.evalBinary(n, env)

	case *parser.UnaryExpr:
		operand, err := interp.evalExpr(n.Operand, env)
		if err != nil {
			return nil, err
		}
		switch n.Operator {
		case "-":
			return -toNumber(operand), nil
		case "!":
			return !isTruthy(operand), nil
		}
		return nil, fmt.Errorf("unknown unary operator: %s", n.Operator)

	case *parser.CallExpr:
		return interp.evalCall(n, env)

	case *parser.IndexExpr:
		obj, err := interp.evalExpr(n.Object, env)
		if err != nil {
			return nil, err
		}
		if obj == nil {
			panic(throwSignal{Value: &Exception{
				Type:    ExMakKauHijau,
				Message: "cannot index into kosong — it's null lah!",
			}})
		}
		index, err := interp.evalExpr(n.Index, env)
		if err != nil {
			return nil, err
		}
		switch o := obj.(type) {
		case []interface{}:
			i := int(toNumber(index))
			if i < 0 || i >= len(o) {
				panic(throwSignal{Value: &Exception{
					Type:    ExKantoi,
					Message: fmt.Sprintf("index %d out of range — senarai cuma ada %d items je", i, len(o)),
				}})
			}
			return o[i], nil
		case string:
			i := int(toNumber(index))
			if i < 0 || i >= len(o) {
				panic(throwSignal{Value: &Exception{
					Type:    ExKantoi,
					Message: fmt.Sprintf("index %d out of range — ayat cuma ada %d characters je", i, len(o)),
				}})
			}
			return string(o[i]), nil

	
		default:
			panic(throwSignal{Value: &Exception{
				Type:    ExApaNi,
				Message: fmt.Sprintf("cannot index into %T — that's not an array or string lah", obj),
			}})
		}

	case *parser.MemberExpr:
		obj, err := interp.evalExpr(n.Object, env)
		if err != nil {
			return nil, err
		}
		if obj == nil {
			panic(throwSignal{Value: &Exception{
				Type:    ExMakKauHijau,
				Message: fmt.Sprintf("cannot access property '%s' on kosong — it's null lah!", n.Property),
			}})
		}
		switch o := obj.(type) {
		case []interface{}:
			if n.Property == "panjang" {
				return float64(len(o)), nil
			}
			panic(throwSignal{Value: &Exception{
				Type:    ExApaNi,
				Message: fmt.Sprintf("senarai tak ada property '%s'", n.Property),
			}})
		case string:
			if n.Property == "panjang" {
				return float64(len(o)), nil
			}
			panic(throwSignal{Value: &Exception{
				Type:    ExApaNi,
				Message: fmt.Sprintf("ayat tak ada property '%s'", n.Property),
			}})

		case *StructInstance:
			val , ok := o.Fields[n.Property]
			if !ok {
				panic(throwSignal{Value: &Exception{
					Type: ExApaNi,
					Message: fmt.Sprintf("kroni '%s' tak ada property '%s' lah", o.Def.Name, n.Property),
				}})
			}

			return val, nil
			
		default:
			panic(throwSignal{Value: &Exception{
				Type:    ExApaNi,
				Message: fmt.Sprintf("cannot access property '%s' on %T", n.Property, obj),
			}})
		}

	case *parser.AnonFunction:
		return &Function{
			Name:       "<anon>",
			Params:     n.Params,
			ReturnType: n.ReturnType,
			Body:       n.Body,
			Closure:    env,
		}, nil

	

	default:
		return nil, fmt.Errorf("unknown expression type: %T", node)
	}
}

func (interp *Interpreter) evalBinary(n *parser.BinaryExpr, env *Environment) (interface{}, error) {
	left, err := interp.evalExpr(n.Left, env)
	if err != nil {
		return nil, err
	}

	// Short-circuit for logical operators
	switch n.Operator {
	case "if not":
		if isTruthy(left) {
			return left, nil
		}
		return interp.evalExpr(n.Right, env)
	case "also":
		if !isTruthy(left) {
			return left, nil
		}
		return interp.evalExpr(n.Right, env)
	}

	right, err := interp.evalExpr(n.Right, env)
	if err != nil {
		return nil, err
	}

	switch n.Operator {
	case "+":
		return addValues(left, right), nil
	case "-":
		return toNumber(left) - toNumber(right), nil
	case "*":
		return toNumber(left) * toNumber(right), nil
	case "/":
		r := toNumber(right)
		if r == 0 {
			panic(throwSignal{Value: &Exception{
				Type:    ExKantoi,
				Message: "cannot divide by zero — maths pun tak reti ke",
			}})
		}
		return toNumber(left) / r, nil
	case "%":
		r := toNumber(right)
		if r == 0 {
			panic(throwSignal{Value: &Exception{
				Type:    ExKantoi,
				Message: "cannot modulo by zero",
			}})
		}
		return float64(int(toNumber(left)) % int(r)), nil
	case "sama":
		return isEqual(left, right), nil
	case "tak sama":
		return !isEqual(left, right), nil
	case "<":
		return toNumber(left) < toNumber(right), nil
	case ">":
		return toNumber(left) > toNumber(right), nil
	case "<=":
		return toNumber(left) <= toNumber(right), nil
	case ">=":
		return toNumber(left) >= toNumber(right), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", n.Operator)
	}
}

func (interp *Interpreter) evalCall(n *parser.CallExpr, env *Environment) (result interface{}, retErr error) {
	callee, err := interp.evalExpr(n.Callee, env)
	if err != nil {
		return nil, err
	}

	var args []interface{}
	for _, arg := range n.Args {
		val, err := interp.evalExpr(arg, env)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	switch fn := callee.(type) {
	case *Function:
		if len(args) != len(fn.Params) {
			panic(throwSignal{Value: &Exception{
				Type:    ExApaNi,
				Message: fmt.Sprintf("function '%s' expects %d args, got %d — check betul betul lah", fn.Name, len(fn.Params), len(args)),
			}})
		}
		fnEnv := NewEnvironment(fn.Closure)
		for i, param := range fn.Params {
			if param.TypeName != "" {
				checkType(args[i], param.TypeName, param.Name, n.Line)
			}
			fnEnv.Set(param.Name, args[i], false)
			fnEnv.SetType(param.Name, param.TypeName)
		}

		// Catch return signal
		defer func() {
			if r := recover(); r != nil {
				if rs, ok := r.(returnSignal); ok {
					result = rs.Value
					if fn.ReturnType != "" {
						checkType(result, fn.ReturnType, fn.Name+" return", n.Line)
					}
					retErr = nil
				} else {
					panic(r)
				}
			}
		}()

		err := interp.execBlock(fn.Body, fnEnv)
		if err != nil {
			return nil, err
		}
		return nil, nil

	case *BuiltinFunc:
		return fn.Fn(args)


	case *StructDef:
		if len(args) != len(fn.Fields){
			panic(throwSignal{
				Value: &Exception{
					Type: ExApaNi,
					Message: fmt.Sprintf("kroni '%s' expects %d args, got %d — check betul betul lah", fn.Name, len(fn.Fields), len(args)),
				}})}

		instanceFields := make(map[string]interface{})
		for i , fieldName := range fn.Fields {
			instanceFields[fieldName] = args[i]
		}

		return &StructInstance{
			Def: fn,
			Fields: instanceFields,
		} , nil

	default:
		panic(throwSignal{Value: &Exception{
			Type:    ExApaNi,
			Message: fmt.Sprintf("'%v' is not a function — tak boleh panggil lah", callee),
		}})
	}
}

// Utility functions

func toNumber(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		var n float64
		fmt.Sscanf(v, "%f", &n)
		return n
	default:
		return 0
	}
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case string:
		return v != ""
	case []interface{}:
		return len(v) > 0
	default:
		return true
	}
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	switch av := a.(type) {
	case float64:
		if bv, ok := b.(float64); ok {
			return av == bv
		}
	case string:
		if bv, ok := b.(string); ok {
			return av == bv
		}
	case bool:
		if bv, ok := b.(bool); ok {
			return av == bv
		}
	}
	return false
}

func addValues(a, b interface{}) interface{} {
	// If either is a string, concatenate
	_, aStr := a.(string)
	_, bStr := b.(string)
	if aStr || bStr {
		return stringify(a) + stringify(b)
	}
	return toNumber(a) + toNumber(b)
}

func typeNameOf(val interface{}) string {
	switch val.(type) {
	case float64:
		return "number"
	case string:
		return "string"
	case bool:
		return "bool"
	case []interface{}:
		return "array"
	case *Function, *BuiltinFunc:
		return "function"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%T", val)
	}
}

func checkType(val interface{}, expected string, name string, line int) {
	actual := typeNameOf(val)
	if actual != expected {
		panic(throwSignal{Value: &Exception{
			Type:    ExApaNi,
			Message: fmt.Sprintf("line %d: '%s' expected type %s, got %s", line, name, expected, actual),
		}})
	}
}

func stringify(val interface{}) string {
	if val == nil {
		return "kosong"
	}
	switch v := val.(type) {
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case string:
		return v
	case bool:
		if v {
			return "can"
		}
		return "cannot"
	case []interface{}:
		parts := make([]string, len(v))
		for i, elem := range v {
			parts[i] = stringify(elem)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case *Function:
		return fmt.Sprintf("<kerja %s>", v.Name)
	case *BuiltinFunc:
		return fmt.Sprintf("<builtin %s>", v.Name)
	case *Exception:
		return fmt.Sprintf("%s: %s", v.Type, v.Message)
	default:
		return fmt.Sprintf("%v", v)
	}
}
