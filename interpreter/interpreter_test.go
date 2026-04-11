package interpreter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/kooixh/bolehlang/lexer"
	"github.com/kooixh/bolehlang/parser"
)

// captureOutput runs bolehlang source and captures stdout.
func captureOutput(t *testing.T, source string) string {
	t.Helper()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSource(source)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("runtime error: %s\nsource: %s", err, source)
	}
	return buf.String()
}

// captureStderr runs bolehlang source and captures stderr.
func captureStderr(t *testing.T, source string) string {
	t.Helper()

	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := runSource(source)

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("runtime error: %s", err)
	}
	return buf.String()
}

func runSource(source string) error {
	lex := lexer.New(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		return fmt.Errorf("lexer: %w", err)
	}
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		return fmt.Errorf("parser: %w", err)
	}
	interp := New("test.boleh")
	return interp.Run(program)
}

func runSourceExpectError(t *testing.T, source string) error {
	t.Helper()
	err := runSource(source)
	if err == nil {
		t.Fatal("expected runtime error")
	}
	return err
}

// --- Variable declaration and assignment ---

func TestVarDeclaration(t *testing.T) {
	out := captureOutput(t, `dei x = 42; eh x;`)
	if strings.TrimSpace(out) != "42" {
		t.Errorf("expected '42', got %q", out)
	}
}

func TestStringVar(t *testing.T) {
	out := captureOutput(t, `dei s = "hello"; eh s;`)
	if strings.TrimSpace(out) != "hello" {
		t.Errorf("expected 'hello', got %q", out)
	}
}

func TestBoolVars(t *testing.T) {
	out := captureOutput(t, `eh can; eh cannot;`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 || lines[0] != "can" || lines[1] != "cannot" {
		t.Errorf("expected 'can\\ncannot', got %q", out)
	}
}

func TestNullVar(t *testing.T) {
	out := captureOutput(t, `dei x = kosong; eh x;`)
	if strings.TrimSpace(out) != "kosong" {
		t.Errorf("expected 'kosong', got %q", out)
	}
}

func TestConstDeclaration(t *testing.T) {
	out := captureOutput(t, `memang dei PI = 3.14; eh PI;`)
	if strings.TrimSpace(out) != "3.14" {
		t.Errorf("expected '3.14', got %q", out)
	}
}

func TestConstReassignFails(t *testing.T) {
	err := runSourceExpectError(t, `memang dei PI = 3.14; tukar PI = 3;`)
	if !strings.Contains(err.Error(), "memang") {
		t.Errorf("expected memang error, got: %s", err)
	}
}

func TestAssignment(t *testing.T) {
	out := captureOutput(t, `dei x = 1; tukar x = 2; eh x;`)
	if strings.TrimSpace(out) != "2" {
		t.Errorf("expected '2', got %q", out)
	}
}

func TestPlusAssignment(t *testing.T) {
	out := captureOutput(t, `dei x = 10; tukar x += 5; eh x;`)
	if strings.TrimSpace(out) != "15" {
		t.Errorf("expected '15', got %q", out)
	}
}

func TestMinusAssignment(t *testing.T) {
	out := captureOutput(t, `dei x = 10; tukar x -= 3; eh x;`)
	if strings.TrimSpace(out) != "7" {
		t.Errorf("expected '7', got %q", out)
	}
}

func TestAssignUndeclaredFails(t *testing.T) {
	runSourceExpectError(t, `tukar x = 1;`)
}

func TestUndeclaredVarFails(t *testing.T) {
	runSourceExpectError(t, `eh x;`)
}

// --- Arithmetic ---

func TestArithmetic(t *testing.T) {
	tests := []struct {
		expr     string
		expected string
	}{
		{`eh 2 + 3;`, "5"},
		{`eh 10 - 4;`, "6"},
		{`eh 3 * 4;`, "12"},
		{`eh 15 / 3;`, "5"},
		{`eh 10 % 3;`, "1"},
		{`eh (2 + 3) * 4;`, "20"},
		{`eh -5;`, "-5"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			out := captureOutput(t, tt.expr)
			if strings.TrimSpace(out) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strings.TrimSpace(out))
			}
		})
	}
}

func TestDivideByZero(t *testing.T) {
	runSourceExpectError(t, `eh 1 / 0;`)
}

func TestModuloByZero(t *testing.T) {
	runSourceExpectError(t, `eh 1 % 0;`)
}

// --- String operations ---

func TestStringConcatenation(t *testing.T) {
	out := captureOutput(t, `eh "hello" + " " + "world";`)
	if strings.TrimSpace(out) != "hello world" {
		t.Errorf("expected 'hello world', got %q", out)
	}
}

func TestStringNumberConcat(t *testing.T) {
	out := captureOutput(t, `eh "age: " + 25;`)
	if strings.TrimSpace(out) != "age: 25" {
		t.Errorf("expected 'age: 25', got %q", out)
	}
}

// --- Comparison and equality ---

func TestComparisons(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{`eh 1 < 2;`, "can"},
		{`eh 2 > 1;`, "can"},
		{`eh 1 <= 1;`, "can"},
		{`eh 1 >= 1;`, "can"},
		{`eh 3 < 2;`, "cannot"},
		{`eh 1 sama 1;`, "can"},
		{`eh 1 tak sama 2;`, "can"},
		{`eh 1 sama 2;`, "cannot"},
		{`eh "a" sama "a";`, "can"},
		{`eh "a" sama "b";`, "cannot"},
		{`eh can sama can;`, "can"},
		{`eh kosong sama kosong;`, "can"},
	}
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			out := captureOutput(t, tt.source)
			if strings.TrimSpace(out) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strings.TrimSpace(out))
			}
		})
	}
}

// --- Logical operators ---

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{`eh can dan can;`, "can"},
		{`eh can dan cannot;`, "cannot"},
		{`eh cannot atau can;`, "can"},
		{`eh cannot atau cannot;`, "cannot"},
		{`eh !can;`, "cannot"},
		{`eh !cannot;`, "can"},
	}
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			out := captureOutput(t, tt.source)
			if strings.TrimSpace(out) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strings.TrimSpace(out))
			}
		})
	}
}

// --- If/else ---

func TestIfTrue(t *testing.T) {
	out := captureOutput(t, `is it (can) { eh "yes"; }`)
	if strings.TrimSpace(out) != "yes" {
		t.Errorf("expected 'yes', got %q", out)
	}
}

func TestIfFalse(t *testing.T) {
	out := captureOutput(t, `is it (cannot) { eh "yes"; }`)
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected empty, got %q", out)
	}
}

func TestIfElse(t *testing.T) {
	out := captureOutput(t, `is it (cannot) { eh "yes"; } abuden { eh "no"; }`)
	if strings.TrimSpace(out) != "no" {
		t.Errorf("expected 'no', got %q", out)
	}
}

func TestIfElseIf(t *testing.T) {
	out := captureOutput(t, `
		dei x = 2;
		is it (x sama 1) {
			eh "one";
		} or is it (x sama 2) {
			eh "two";
		} abuden {
			eh "other";
		}
	`)
	if strings.TrimSpace(out) != "two" {
		t.Errorf("expected 'two', got %q", out)
	}
}

// --- Loops ---

func TestWhileLoop(t *testing.T) {
	out := captureOutput(t, `
		dei i = 0;
		one by one (i < 3) {
			eh i;
			tukar i += 1;
		}
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	expected := []string{"0", "1", "2"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %q", len(expected), len(lines), out)
	}
	for i, e := range expected {
		if lines[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, lines[i])
		}
	}
}

func TestForRange(t *testing.T) {
	out := captureOutput(t, `
		one by one (i = 1 sampai 3) {
			eh i;
		}
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	expected := []string{"1", "2", "3"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}
	for i, e := range expected {
		if lines[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, lines[i])
		}
	}
}

func TestForEach(t *testing.T) {
	out := captureOutput(t, `
		dei arr = ["a", "b", "c"];
		one by one (item dalam arr) {
			eh item;
		}
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	expected := []string{"a", "b", "c"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}
	for i, e := range expected {
		if lines[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, lines[i])
		}
	}
}

func TestBreak(t *testing.T) {
	out := captureOutput(t, `
		dei i = 0;
		one by one (i < 10) {
			is it (i sama 3) { gaodim; }
			eh i;
			tukar i += 1;
		}
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	expected := []string{"0", "1", "2"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %v", len(expected), len(lines), lines)
	}
}

func TestContinue(t *testing.T) {
	out := captureOutput(t, `
		one by one (i = 1 sampai 5) {
			is it (i sama 3) { no need; }
			eh i;
		}
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	expected := []string{"1", "2", "4", "5"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %v", len(expected), len(lines), lines)
	}
}

func TestForEachOnNonArrayFails(t *testing.T) {
	runSourceExpectError(t, `
		dei x = 5;
		one by one (item dalam x) {
			eh item;
		}
	`)
}

// --- Functions ---

func TestFunctionDeclarationAndCall(t *testing.T) {
	out := captureOutput(t, `
		jom greet(name) {
			nah "Hello, " + name + "!";
		}
		eh greet("Ali");
	`)
	if strings.TrimSpace(out) != "Hello, Ali!" {
		t.Errorf("expected 'Hello, Ali!', got %q", out)
	}
}

func TestFunctionReturnValue(t *testing.T) {
	out := captureOutput(t, `
		jom add(a, b) { nah a + b; }
		eh add(3, 4);
	`)
	if strings.TrimSpace(out) != "7" {
		t.Errorf("expected '7', got %q", out)
	}
}

func TestFunctionNoReturn(t *testing.T) {
	out := captureOutput(t, `
		jom doNothing() { dei x = 1; }
		eh doNothing();
	`)
	if strings.TrimSpace(out) != "kosong" {
		t.Errorf("expected 'kosong', got %q", out)
	}
}

func TestRecursion(t *testing.T) {
	out := captureOutput(t, `
		jom factorial(n) {
			is it (n <= 1) { nah 1; }
			nah n * factorial(n - 1);
		}
		eh factorial(5);
	`)
	if strings.TrimSpace(out) != "120" {
		t.Errorf("expected '120', got %q", out)
	}
}

func TestClosure(t *testing.T) {
	out := captureOutput(t, `
		jom makeCounter(start) {
			dei count = start;
			jom inc() {
				tukar count = count + 1;
				nah count;
			}
			nah inc;
		}
		dei c = makeCounter(0);
		eh c();
		eh c();
		eh c();
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	expected := []string{"1", "2", "3"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}
	for i, e := range expected {
		if lines[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, lines[i])
		}
	}
}

func TestAnonFunction(t *testing.T) {
	out := captureOutput(t, `
		dei double = jom (x) { nah x * 2; };
		eh double(5);
	`)
	if strings.TrimSpace(out) != "10" {
		t.Errorf("expected '10', got %q", out)
	}
}

func TestWrongArgCountFails(t *testing.T) {
	runSourceExpectError(t, `
		jom f(a, b) { nah a + b; }
		f(1);
	`)
}

func TestCallNonFunctionFails(t *testing.T) {
	runSourceExpectError(t, `
		dei x = 42;
		x();
	`)
}

// --- Arrays ---

func TestArrayLiteral(t *testing.T) {
	out := captureOutput(t, `dei arr = [1, 2, 3]; eh arr;`)
	if strings.TrimSpace(out) != "[1, 2, 3]" {
		t.Errorf("expected '[1, 2, 3]', got %q", out)
	}
}

func TestEmptyArray(t *testing.T) {
	out := captureOutput(t, `dei arr = []; eh arr;`)
	if strings.TrimSpace(out) != "[]" {
		t.Errorf("expected '[]', got %q", out)
	}
}

func TestArrayIndex(t *testing.T) {
	out := captureOutput(t, `dei arr = ["a", "b", "c"]; eh arr[1];`)
	if strings.TrimSpace(out) != "b" {
		t.Errorf("expected 'b', got %q", out)
	}
}

func TestArrayIndexOutOfBounds(t *testing.T) {
	runSourceExpectError(t, `dei arr = [1]; eh arr[5];`)
}

func TestStringIndex(t *testing.T) {
	out := captureOutput(t, `dei s = "abc"; eh s[0];`)
	if strings.TrimSpace(out) != "a" {
		t.Errorf("expected 'a', got %q", out)
	}
}

func TestMemberPanjang(t *testing.T) {
	out := captureOutput(t, `dei arr = [1, 2, 3]; eh arr.panjang;`)
	if strings.TrimSpace(out) != "3" {
		t.Errorf("expected '3', got %q", out)
	}
}

func TestStringMemberPanjang(t *testing.T) {
	out := captureOutput(t, `dei s = "hello"; eh s.panjang;`)
	if strings.TrimSpace(out) != "5" {
		t.Errorf("expected '5', got %q", out)
	}
}

// --- Try/catch/finally ---

func TestTryCatch(t *testing.T) {
	out := captureOutput(t, `
		try see {
			mampus "error!";
		} alamak (e) {
			eh "caught: " + e;
		}
	`)
	if strings.TrimSpace(out) != "caught: error!" {
		t.Errorf("expected 'caught: error!', got %q", out)
	}
}

func TestTryCatchFinally(t *testing.T) {
	out := captureOutput(t, `
		try see {
			mampus "fail";
		} alamak (e) {
			eh "caught";
		} confirm need {
			eh "finally";
		}
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 || lines[0] != "caught" || lines[1] != "finally" {
		t.Errorf("expected 'caught\\nfinally', got %q", out)
	}
}

func TestTryNoError(t *testing.T) {
	out := captureOutput(t, `
		try see {
			eh "ok";
		} alamak (e) {
			eh "caught";
		}
	`)
	if strings.TrimSpace(out) != "ok" {
		t.Errorf("expected 'ok', got %q", out)
	}
}

func TestFinallyAlwaysRuns(t *testing.T) {
	out := captureOutput(t, `
		try see {
			eh "try";
		} alamak (e) {
			eh "catch";
		} confirm need {
			eh "finally";
		}
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 || lines[0] != "try" || lines[1] != "finally" {
		t.Errorf("expected 'try\\nfinally', got %q", out)
	}
}

func TestUncaughtThrow(t *testing.T) {
	// An uncaught throw should cause a panic that propagates
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic from uncaught throw")
		}
	}()
	runSource(`mampus "uncaught!";`)
}

// --- Warn statement ---

func TestWarnStatement(t *testing.T) {
	out := captureStderr(t, `walao eh "danger!";`)
	if !strings.Contains(out, "WALAO EH") || !strings.Contains(out, "danger!") {
		t.Errorf("expected WALAO EH warning, got %q", out)
	}
}

// --- Built-in functions ---

func TestBuiltinPanjang(t *testing.T) {
	out := captureOutput(t, `eh panjang("hello");`)
	if strings.TrimSpace(out) != "5" {
		t.Errorf("expected '5', got %q", out)
	}

	out = captureOutput(t, `eh panjang([1, 2, 3]);`)
	if strings.TrimSpace(out) != "3" {
		t.Errorf("expected '3', got %q", out)
	}
}

func TestBuiltinPanjangWrongArgs(t *testing.T) {
	runSourceExpectError(t, `panjang();`)
	runSourceExpectError(t, `panjang(1);`)
}

func TestBuiltinJenis(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{`eh jenis(42);`, "nombor"},
		{`eh jenis("hello");`, "ayat"},
		{`eh jenis(can);`, "boleh_tak_boleh"},
		{`eh jenis([1,2]);`, "senarai"},
		{`eh jenis(kosong);`, "kosong"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			out := captureOutput(t, tt.source)
			if strings.TrimSpace(out) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strings.TrimSpace(out))
			}
		})
	}
}

func TestBuiltinJenisFunction(t *testing.T) {
	out := captureOutput(t, `jom f() {} eh jenis(f);`)
	if strings.TrimSpace(out) != "kerja" {
		t.Errorf("expected 'kerja', got %q", out)
	}

	out = captureOutput(t, `eh jenis(panjang);`)
	if strings.TrimSpace(out) != "kerja" {
		t.Errorf("expected 'kerja', got %q", out)
	}
}

func TestBuiltinMasuk(t *testing.T) {
	out := captureOutput(t, `
		dei arr = [1, 2];
		tukar arr = masuk(arr, 3);
		eh arr;
	`)
	if strings.TrimSpace(out) != "[1, 2, 3]" {
		t.Errorf("expected '[1, 2, 3]', got %q", out)
	}
}

func TestBuiltinMasukMultiple(t *testing.T) {
	out := captureOutput(t, `
		dei arr = [];
		tukar arr = masuk(arr, 1, 2, 3);
		eh arr;
	`)
	if strings.TrimSpace(out) != "[1, 2, 3]" {
		t.Errorf("expected '[1, 2, 3]', got %q", out)
	}
}

func TestBuiltinPotong(t *testing.T) {
	out := captureOutput(t, `
		dei arr = [10, 20, 30, 40, 50];
		eh potong(arr, 1, 3);
	`)
	if strings.TrimSpace(out) != "[20, 30]" {
		t.Errorf("expected '[20, 30]', got %q", out)
	}
}

func TestBuiltinPotongNoEnd(t *testing.T) {
	out := captureOutput(t, `
		dei arr = [10, 20, 30];
		eh potong(arr, 1);
	`)
	if strings.TrimSpace(out) != "[20, 30]" {
		t.Errorf("expected '[20, 30]', got %q", out)
	}
}

func TestBuiltinNombor(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{`eh nombor("42");`, "42"},
		{`eh nombor(42);`, "42"},
		{`eh nombor(can);`, "1"},
		{`eh nombor(cannot);`, "0"},
	}
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			out := captureOutput(t, tt.source)
			if strings.TrimSpace(out) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strings.TrimSpace(out))
			}
		})
	}
}

func TestBuiltinNomborInvalidString(t *testing.T) {
	runSourceExpectError(t, `nombor("abc");`)
}

func TestBuiltinAyat(t *testing.T) {
	out := captureOutput(t, `eh ayat(42);`)
	if strings.TrimSpace(out) != "42" {
		t.Errorf("expected '42', got %q", out)
	}
}

func TestBuiltinCincai(t *testing.T) {
	// Just test it runs without error — output is random
	captureOutput(t, `dei x = cincai();`)
	captureOutput(t, `dei x = cincai(10);`)
	captureOutput(t, `dei x = cincai(1, 10);`)
}

func TestBuiltinCincaiInvalidArgs(t *testing.T) {
	runSourceExpectError(t, `cincai(0);`)
	runSourceExpectError(t, `cincai(5, 2);`) // min >= max
	runSourceExpectError(t, `cincai(1, 2, 3, 4);`)
}

// --- Scoping ---

func TestBlockScoping(t *testing.T) {
	out := captureOutput(t, `
		dei x = "outer";
		is it (can) {
			dei x = "inner";
			eh x;
		}
		eh x;
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 || lines[0] != "inner" || lines[1] != "outer" {
		t.Errorf("expected 'inner\\nouter', got %q", out)
	}
}

// --- Truthiness ---

func TestTruthiness(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{`is it (1) { eh "yes"; } abuden { eh "no"; }`, "yes"},
		{`is it (0) { eh "yes"; } abuden { eh "no"; }`, "no"},
		{`is it ("hello") { eh "yes"; } abuden { eh "no"; }`, "yes"},
		{`is it ("") { eh "yes"; } abuden { eh "no"; }`, "no"},
		{`is it (kosong) { eh "yes"; } abuden { eh "no"; }`, "no"},
		{`is it ([1]) { eh "yes"; } abuden { eh "no"; }`, "yes"},
		{`is it ([]) { eh "yes"; } abuden { eh "no"; }`, "no"},
	}
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			out := captureOutput(t, tt.source)
			if strings.TrimSpace(out) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strings.TrimSpace(out))
			}
		})
	}
}

// --- Stringify ---

func TestStringify(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{`eh 42;`, "42"},
		{`eh 3.14;`, "3.14"},
		{`eh can;`, "can"},
		{`eh cannot;`, "cannot"},
		{`eh kosong;`, "kosong"},
		{`eh [1, "two", can];`, `[1, two, can]`},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			out := captureOutput(t, tt.source)
			if strings.TrimSpace(out) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strings.TrimSpace(out))
			}
		})
	}
}

// --- FizzBuzz integration test ---

func TestFizzBuzz(t *testing.T) {
	out := captureOutput(t, `
		one by one (i = 1 sampai 15) {
			is it (i % 15 sama 0) {
				eh "FizzBuzz";
			} or is it (i % 3 sama 0) {
				eh "Fizz";
			} or is it (i % 5 sama 0) {
				eh "Buzz";
			} abuden {
				eh i;
			}
		}
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	expected := []string{
		"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8",
		"Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz",
	}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}
	for i, e := range expected {
		if lines[i] != e {
			t.Errorf("line %d: expected %q, got %q", i+1, e, lines[i])
		}
	}
}

// --- Comments ---

func TestSingleLineComment(t *testing.T) {
	out := captureOutput(t, `
		// eh: this is a comment
		eh "hello";
	`)
	if strings.TrimSpace(out) != "hello" {
		t.Errorf("expected 'hello', got %q", out)
	}
}

func TestBlockComment(t *testing.T) {
	out := captureOutput(t, `
		alamak mula
		this is a block comment
		alamak habis
		eh "hello";
	`)
	if strings.TrimSpace(out) != "hello" {
		t.Errorf("expected 'hello', got %q", out)
	}
}

// --- Short-circuit evaluation ---

func TestShortCircuitDan(t *testing.T) {
	// "dan" should not evaluate right side if left is falsy
	out := captureOutput(t, `
		dei x = cannot dan can;
		eh x;
	`)
	if strings.TrimSpace(out) != "cannot" {
		t.Errorf("expected 'cannot', got %q", out)
	}
}

func TestShortCircuitAtau(t *testing.T) {
	// "atau" should not evaluate right side if left is truthy
	out := captureOutput(t, `
		dei x = can atau cannot;
		eh x;
	`)
	if strings.TrimSpace(out) != "can" {
		t.Errorf("expected 'can', got %q", out)
	}
}

// --- Edge cases ---

func TestPlusAssignWithString(t *testing.T) {
	out := captureOutput(t, `
		dei s = "hello";
		tukar s += " world";
		eh s;
	`)
	if strings.TrimSpace(out) != "hello world" {
		t.Errorf("expected 'hello world', got %q", out)
	}
}

func TestNestedFunctionCalls(t *testing.T) {
	out := captureOutput(t, `
		jom double(x) { nah x * 2; }
		jom inc(x) { nah x + 1; }
		eh double(inc(4));
	`)
	if strings.TrimSpace(out) != "10" {
		t.Errorf("expected '10', got %q", out)
	}
}

func TestArrayOfFunctions(t *testing.T) {
	out := captureOutput(t, `
		jom a() { nah "a"; }
		jom b() { nah "b"; }
		dei fns = [a, b];
		eh fns[0]();
		eh fns[1]();
	`)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 || lines[0] != "a" || lines[1] != "b" {
		t.Errorf("expected 'a\\nb', got %q", out)
	}
}
