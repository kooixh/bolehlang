package tests

import (
	"testing"

	"github.com/kooixh/bolehlang/lexer"
	"github.com/kooixh/bolehlang/parser"
)

func parse(t *testing.T, source string) *parser.Program {
	t.Helper()
	lex := lexer.New(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("lexer error: %s", err)
	}
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("parser error: %s", err)
	}
	return program
}

func parseFail(t *testing.T, source string) {
	t.Helper()
	lex := lexer.New(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		return // lexer error is also acceptable
	}
	p := parser.New(tokens)
	_, err = p.Parse()
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseVarDeclaration(t *testing.T) {
	prog := parse(t, `dei x = 42;`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}
	decl, ok := prog.Body[0].(*parser.VarDeclaration)
	if !ok {
		t.Fatalf("expected VarDeclaration, got %T", prog.Body[0])
	}
	if decl.Name != "x" {
		t.Errorf("expected name 'x', got '%s'", decl.Name)
	}
	num, ok := decl.Value.(*parser.NumberLiteral)
	if !ok {
		t.Fatalf("expected NumberLiteral, got %T", decl.Value)
	}
	if num.Value != 42 {
		t.Errorf("expected 42, got %f", num.Value)
	}
}

func TestParseConstDeclaration(t *testing.T) {
	prog := parse(t, `memang dei PI = 3.14;`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}
	decl, ok := prog.Body[0].(*parser.ConstDeclaration)
	if !ok {
		t.Fatalf("expected ConstDeclaration, got %T", prog.Body[0])
	}
	if decl.Name != "PI" {
		t.Errorf("expected name 'PI', got '%s'", decl.Name)
	}
}



func TestParseStringLiteral(t *testing.T) {
	prog := parse(t, `dei s = "hello";`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	str, ok := decl.Value.(*parser.StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral, got %T", decl.Value)
	}
	if str.Value != "hello" {
		t.Errorf("expected 'hello', got '%s'", str.Value)
	}
}

func TestParseBoolLiterals(t *testing.T) {
	prog := parse(t, `dei a = can; dei b = cannot;`)
	aDecl := prog.Body[0].(*parser.VarDeclaration)
	bDecl := prog.Body[1].(*parser.VarDeclaration)

	aBool := aDecl.Value.(*parser.BoolLiteral)
	if !aBool.Value {
		t.Error("expected 'can' to be true")
	}
	bBool := bDecl.Value.(*parser.BoolLiteral)
	if bBool.Value {
		t.Error("expected 'cannot' to be false")
	}
}

func TestParseNullLiteral(t *testing.T) {
	prog := parse(t, `dei x = kosong;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	_, ok := decl.Value.(*parser.NullLiteral)
	if !ok {
		t.Fatalf("expected NullLiteral, got %T", decl.Value)
	}
}

func TestParseAssignment(t *testing.T) {
	prog := parse(t, `tukar x = 10;`)
	assign, ok := prog.Body[0].(*parser.Assignment)
	if !ok {
		t.Fatalf("expected Assignment, got %T", prog.Body[0])
	}
	if assign.Name != "x" || assign.Operator != "=" {
		t.Errorf("expected x =, got %s %s", assign.Name, assign.Operator)
	}
}

func TestParseAssignmentPlusEquals(t *testing.T) {
	prog := parse(t, `tukar x += 5;`)
	assign := prog.Body[0].(*parser.Assignment)
	if assign.Operator != "+=" {
		t.Errorf("expected +=, got %s", assign.Operator)
	}
}

func TestParseAssignmentMinusEquals(t *testing.T) {
	prog := parse(t, `tukar x -= 3;`)
	assign := prog.Body[0].(*parser.Assignment)
	if assign.Operator != "-=" {
		t.Errorf("expected -=, got %s", assign.Operator)
	}
}

func TestParsePrintStatement(t *testing.T) {
	prog := parse(t, `eh "hello";`)
	stmt, ok := prog.Body[0].(*parser.PrintStatement)
	if !ok {
		t.Fatalf("expected PrintStatement, got %T", prog.Body[0])
	}
	str := stmt.Value.(*parser.StringLiteral)
	if str.Value != "hello" {
		t.Errorf("expected 'hello', got '%s'", str.Value)
	}
}

func TestParseWarnStatement(t *testing.T) {
	prog := parse(t, `walao eh "danger";`)
	_, ok := prog.Body[0].(*parser.WarnStatement)
	if !ok {
		t.Fatalf("expected WarnStatement, got %T", prog.Body[0])
	}
}

func TestParseIfStatement(t *testing.T) {
	prog := parse(t, `is it (can) { eh "yes"; }`)
	stmt, ok := prog.Body[0].(*parser.IfStatement)
	if !ok {
		t.Fatalf("expected IfStatement, got %T", prog.Body[0])
	}
	if len(stmt.Body) != 1 {
		t.Errorf("expected 1 body statement, got %d", len(stmt.Body))
	}
}

func TestParseIfElseIfElse(t *testing.T) {
	prog := parse(t, `
		is it (x sama 1) {
			eh "one";
		} or is it (x sama 2) {
			eh "two";
		} abuden {
			eh "other";
		}
	`)
	stmt := prog.Body[0].(*parser.IfStatement)
	if len(stmt.ElseIfCases) != 1 {
		t.Errorf("expected 1 else-if case, got %d", len(stmt.ElseIfCases))
	}
	if stmt.ElseBody == nil {
		t.Error("expected else body")
	}
}

func TestParseWhileLoop(t *testing.T) {
	prog := parse(t, `one by one (x < 5) { eh x; }`)
	stmt, ok := prog.Body[0].(*parser.WhileStatement)
	if !ok {
		t.Fatalf("expected WhileStatement, got %T", prog.Body[0])
	}
	if len(stmt.Body) != 1 {
		t.Errorf("expected 1 body statement, got %d", len(stmt.Body))
	}
}

func TestParseForRange(t *testing.T) {
	prog := parse(t, `one by one (i = 1 sampai 10) { eh i; }`)
	stmt, ok := prog.Body[0].(*parser.ForStatement)
	if !ok {
		t.Fatalf("expected ForStatement, got %T", prog.Body[0])
	}
	if stmt.VarName != "i" {
		t.Errorf("expected var 'i', got '%s'", stmt.VarName)
	}
}

func TestParseForEach(t *testing.T) {
	prog := parse(t, `one by one (item dalam items) { eh item; }`)
	stmt, ok := prog.Body[0].(*parser.ForEachStatement)
	if !ok {
		t.Fatalf("expected ForEachStatement, got %T", prog.Body[0])
	}
	if stmt.VarName != "item" {
		t.Errorf("expected var 'item', got '%s'", stmt.VarName)
	}
}

func TestParseFunctionDeclaration(t *testing.T) {
	prog := parse(t, `jom add(a, b) { nah a + b; }`)
	stmt, ok := prog.Body[0].(*parser.FunctionDeclaration)
	if !ok {
		t.Fatalf("expected FunctionDeclaration, got %T", prog.Body[0])
	}
	if stmt.Name != "add" {
		t.Errorf("expected name 'add', got '%s'", stmt.Name)
	}
	if len(stmt.Params) != 2 {
		t.Errorf("expected 2 params, got %d", len(stmt.Params))
	}
}

func TestParseFunctionNoParams(t *testing.T) {
	prog := parse(t, `jom greet() { eh "hi"; }`)
	stmt := prog.Body[0].(*parser.FunctionDeclaration)
	if len(stmt.Params) != 0 {
		t.Errorf("expected 0 params, got %d", len(stmt.Params))
	}
}

func TestParseReturnStatement(t *testing.T) {
	prog := parse(t, `jom f() { nah 42; }`)
	fn := prog.Body[0].(*parser.FunctionDeclaration)
	ret, ok := fn.Body[0].(*parser.ReturnStatement)
	if !ok {
		t.Fatalf("expected ReturnStatement, got %T", fn.Body[0])
	}
	num := ret.Value.(*parser.NumberLiteral)
	if num.Value != 42 {
		t.Errorf("expected 42, got %f", num.Value)
	}
}

func TestParseReturnEmpty(t *testing.T) {
	prog := parse(t, `jom f() { nah; }`)
	fn := prog.Body[0].(*parser.FunctionDeclaration)
	ret := fn.Body[0].(*parser.ReturnStatement)
	_, ok := ret.Value.(*parser.NullLiteral)
	if !ok {
		t.Fatalf("expected NullLiteral for empty return, got %T", ret.Value)
	}
}

func TestParseBreakContinue(t *testing.T) {
	prog := parse(t, `one by one (can) { gaodim; }`)
	loop := prog.Body[0].(*parser.WhileStatement)
	_, ok := loop.Body[0].(*parser.BreakStatement)
	if !ok {
		t.Fatalf("expected BreakStatement, got %T", loop.Body[0])
	}

	prog2 := parse(t, `one by one (can) { no need; }`)
	loop2 := prog2.Body[0].(*parser.WhileStatement)
	_, ok = loop2.Body[0].(*parser.ContinueStatement)
	if !ok {
		t.Fatalf("expected ContinueStatement, got %T", loop2.Body[0])
	}
}

func TestParseTryCatch(t *testing.T) {
	prog := parse(t, `try see { eh "try"; } alamak (e) { eh e; } confirm need { eh "finally"; }`)
	stmt, ok := prog.Body[0].(*parser.TryCatchStatement)
	if !ok {
		t.Fatalf("expected TryCatchStatement, got %T", prog.Body[0])
	}
	if len(stmt.TryBody) != 1 {
		t.Errorf("expected 1 try body statement, got %d", len(stmt.TryBody))
	}
	if stmt.CatchVar != "e" {
		t.Errorf("expected catch var 'e', got '%s'", stmt.CatchVar)
	}
	if len(stmt.CatchBody) != 1 {
		t.Errorf("expected 1 catch body statement, got %d", len(stmt.CatchBody))
	}
	if len(stmt.FinallyBody) != 1 {
		t.Errorf("expected 1 finally body statement, got %d", len(stmt.FinallyBody))
	}
}

func TestParseTryCatchOnly(t *testing.T) {
	prog := parse(t, `try see { eh "try"; } alamak (e) { eh e; }`)
	stmt := prog.Body[0].(*parser.TryCatchStatement)
	if stmt.FinallyBody != nil {
		t.Error("expected no finally body")
	}
}

func TestParseThrow(t *testing.T) {
	prog := parse(t, `mampus "error!";`)
	stmt, ok := prog.Body[0].(*parser.ThrowStatement)
	if !ok {
		t.Fatalf("expected ThrowStatement, got %T", prog.Body[0])
	}
	str := stmt.Value.(*parser.StringLiteral)
	if str.Value != "error!" {
		t.Errorf("expected 'error!', got '%s'", str.Value)
	}
}

func TestParseImport(t *testing.T) {
	prog := parse(t, `tapau "math.boleh";`)
	stmt, ok := prog.Body[0].(*parser.ImportStatement)
	if !ok {
		t.Fatalf("expected ImportStatement, got %T", prog.Body[0])
	}
	if stmt.Path != "math.boleh" {
		t.Errorf("expected 'math.boleh', got '%s'", stmt.Path)
	}
}

func TestParseArrayLiteral(t *testing.T) {
	prog := parse(t, `dei arr = [1, 2, 3];`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	arr, ok := decl.Value.(*parser.ArrayLiteral)
	if !ok {
		t.Fatalf("expected ArrayLiteral, got %T", decl.Value)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr.Elements))
	}
}

func TestParseEmptyArray(t *testing.T) {
	prog := parse(t, `dei arr = [];`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	arr := decl.Value.(*parser.ArrayLiteral)
	if len(arr.Elements) != 0 {
		t.Errorf("expected 0 elements, got %d", len(arr.Elements))
	}
}

func TestParseTrailingCommaArray(t *testing.T) {
	prog := parse(t, `dei arr = [1, 2,];`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	arr := decl.Value.(*parser.ArrayLiteral)
	if len(arr.Elements) != 2 {
		t.Errorf("expected 2 elements, got %d", len(arr.Elements))
	}
}

func TestParseBinaryExpression(t *testing.T) {
	prog := parse(t, `dei x = 1 + 2 * 3;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	// Should be 1 + (2 * 3) due to precedence
	add, ok := decl.Value.(*parser.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", decl.Value)
	}
	if add.Operator != "+" {
		t.Errorf("expected +, got %s", add.Operator)
	}
	mul, ok := add.Right.(*parser.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr on right, got %T", add.Right)
	}
	if mul.Operator != "*" {
		t.Errorf("expected *, got %s", mul.Operator)
	}
}

func TestParseComparisonOperators(t *testing.T) {
	tests := []string{
		`dei x = a < b;`,
		`dei x = a > b;`,
		`dei x = a <= b;`,
		`dei x = a >= b;`,
		`dei x = a sama b;`,
		`dei x = a tak sama b;`,
	}
	for _, src := range tests {
		prog := parse(t, src)
		decl := prog.Body[0].(*parser.VarDeclaration)
		_, ok := decl.Value.(*parser.BinaryExpr)
		if !ok {
			t.Errorf("expected BinaryExpr for %q, got %T", src, decl.Value)
		}
	}
}

func TestParseLogicalOperators(t *testing.T) {
	prog := parse(t, `dei x = can also cannot if not can;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	// "if not" has lower precedence than "also"
	or, ok := decl.Value.(*parser.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", decl.Value)
	}
	if or.Operator != "if not" {
		t.Errorf("expected 'if not' at top, got '%s'", or.Operator)
	}
}

func TestParseUnaryExpression(t *testing.T) {
	prog := parse(t, `dei x = -5;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	unary, ok := decl.Value.(*parser.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", decl.Value)
	}
	if unary.Operator != "-" {
		t.Errorf("expected -, got %s", unary.Operator)
	}
}

func TestParseNotExpression(t *testing.T) {
	prog := parse(t, `dei x = !can;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	unary := decl.Value.(*parser.UnaryExpr)
	if unary.Operator != "!" {
		t.Errorf("expected !, got %s", unary.Operator)
	}
}

func TestParseFunctionCall(t *testing.T) {
	prog := parse(t, `foo(1, 2);`)
	exprStmt := prog.Body[0].(*parser.ExpressionStatement)
	call, ok := exprStmt.Expression.(*parser.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", exprStmt.Expression)
	}
	if len(call.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(call.Args))
	}
}

func TestParseFunctionCallNoArgs(t *testing.T) {
	prog := parse(t, `foo();`)
	exprStmt := prog.Body[0].(*parser.ExpressionStatement)
	call := exprStmt.Expression.(*parser.CallExpr)
	if len(call.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(call.Args))
	}
}

func TestParseIndexExpression(t *testing.T) {
	prog := parse(t, `dei x = arr[0];`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	idx, ok := decl.Value.(*parser.IndexExpr)
	if !ok {
		t.Fatalf("expected IndexExpr, got %T", decl.Value)
	}
	num := idx.Index.(*parser.NumberLiteral)
	if num.Value != 0 {
		t.Errorf("expected index 0, got %f", num.Value)
	}
}

func TestParseMemberExpression(t *testing.T) {
	prog := parse(t, `dei x = arr.panjang;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	member, ok := decl.Value.(*parser.MemberExpr)
	if !ok {
		t.Fatalf("expected MemberExpr, got %T", decl.Value)
	}
	if member.Property != "panjang" {
		t.Errorf("expected 'panjang', got '%s'", member.Property)
	}
}

func TestParseGroupedExpression(t *testing.T) {
	prog := parse(t, `dei x = (1 + 2) * 3;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	mul, ok := decl.Value.(*parser.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", decl.Value)
	}
	if mul.Operator != "*" {
		t.Errorf("expected * at top, got %s", mul.Operator)
	}
}

func TestParseAnonFunction(t *testing.T) {
	prog := parse(t, `dei f = jom (x) { nah x + 1; };`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	anon, ok := decl.Value.(*parser.AnonFunction)
	if !ok {
		t.Fatalf("expected AnonFunction, got %T", decl.Value)
	}
	if len(anon.Params) != 1 || anon.Params[0].Name != "x" {
		t.Errorf("expected params [x], got %v", anon.Params)
	}
}

func TestParseSemicolonOptional(t *testing.T) {
	// Semicolons are optional in bolehlang
	prog := parse(t, "dei x = 1\ndei y = 2")
	if len(prog.Body) != 2 {
		t.Errorf("expected 2 statements, got %d", len(prog.Body))
	}
}

func TestParseErrorMissingBrace(t *testing.T) {
	parseFail(t, `is it (can) { eh "yes";`)
}

func TestParseErrorMissingParen(t *testing.T) {
	parseFail(t, `is it (can { eh "yes"; }`)
}

// --- Type annotations ---

func TestParseVarWithType(t *testing.T) {
	prog := parse(t, `dei x: number = 42;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	if decl.TypeName != "number" {
		t.Errorf("expected type 'number', got '%s'", decl.TypeName)
	}
}

func TestParseVarWithoutType(t *testing.T) {
	prog := parse(t, `dei x = 42;`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	if decl.TypeName != "" {
		t.Errorf("expected empty type, got '%s'", decl.TypeName)
	}
}

func TestParseConstWithType(t *testing.T) {
	prog := parse(t, `memang dei PI: number = 3.14;`)
	decl := prog.Body[0].(*parser.ConstDeclaration)
	if decl.TypeName != "number" {
		t.Errorf("expected type 'number', got '%s'", decl.TypeName)
	}
}

func TestParseFunctionWithTypedParams(t *testing.T) {
	prog := parse(t, `jom add(a: number, b: number): number { nah a + b; }`)
	fn := prog.Body[0].(*parser.FunctionDeclaration)
	if len(fn.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Params))
	}
	if fn.Params[0].Name != "a" || fn.Params[0].TypeName != "number" {
		t.Errorf("expected param a:number, got %s:%s", fn.Params[0].Name, fn.Params[0].TypeName)
	}
	if fn.Params[1].Name != "b" || fn.Params[1].TypeName != "number" {
		t.Errorf("expected param b:number, got %s:%s", fn.Params[1].Name, fn.Params[1].TypeName)
	}
	if fn.ReturnType != "number" {
		t.Errorf("expected return type 'number', got '%s'", fn.ReturnType)
	}
}

func TestParseFunctionMixedTypedParams(t *testing.T) {
	prog := parse(t, `jom greet(name: string, times) { eh name; }`)
	fn := prog.Body[0].(*parser.FunctionDeclaration)
	if fn.Params[0].TypeName != "string" {
		t.Errorf("expected param 0 type 'string', got '%s'", fn.Params[0].TypeName)
	}
	if fn.Params[1].TypeName != "" {
		t.Errorf("expected param 1 type '', got '%s'", fn.Params[1].TypeName)
	}
}

func TestParseAnonFunctionWithTypes(t *testing.T) {
	prog := parse(t, `dei f = jom (x: number): number { nah x + 1; };`)
	decl := prog.Body[0].(*parser.VarDeclaration)
	anon := decl.Value.(*parser.AnonFunction)
	if anon.Params[0].TypeName != "number" {
		t.Errorf("expected param type 'number', got '%s'", anon.Params[0].TypeName)
	}
	if anon.ReturnType != "number" {
		t.Errorf("expected return type 'number', got '%s'", anon.ReturnType)
	}
}

func TestParseLaiExpression(t *testing.T) {
	prog := parse(t, `dei name = lai;`)
	decl, ok := prog.Body[0].(*parser.VarDeclaration)
	if !ok {
		t.Fatalf("expected VarDeclaration, got %T", prog.Body[0])
	}
	if _, ok := decl.Value.(*parser.ReadStdinExpr); !ok {
		t.Fatalf("expected ReadStdinExpr, got %T", decl.Value)
	}
}

func TestParseLaiInExpression(t *testing.T) {
	prog := parse(t, `eh "hi " + lai;`)
	stmt, ok := prog.Body[0].(*parser.PrintStatement)
	if !ok {
		t.Fatalf("expected PrintStatement, got %T", prog.Body[0])
	}
	bin, ok := stmt.Value.(*parser.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", stmt.Value)
	}
	if _, ok := bin.Right.(*parser.ReadStdinExpr); !ok {
		t.Fatalf("expected ReadStdinExpr on right, got %T", bin.Right)
	}
}


