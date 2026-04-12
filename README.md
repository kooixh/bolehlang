# bolehlang

A programming language inspired by Malaysian English (Manglish). Because why code in boring English when you can code Malaysian style? Malaysia boleh!

```
eh "Hello World! Apa khabar?";
```

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Language Reference](#language-reference)
  - [Variables](#variables)
  - [Data Types](#data-types)
  - [Operators](#operators)
  - [Print](#print)
  - [Control Flow](#control-flow)
  - [Loops](#loops)
  - [Functions](#functions)
  - [Error Handling](#error-handling)
  - [Imports](#imports)
  - [Built-in Functions](#built-in-functions)
  - [Comments](#comments)
- [Examples](#examples)
  - [Hello World](#hello-world)
  - [FizzBuzz](#fizzbuzz)
  - [Closures](#closures)
- [Keyword Reference](#keyword-reference)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Overview

**bolehlang** is a dynamically-typed, interpreted programming language written in Go. It replaces conventional programming keywords with Manglish (Malaysian English) slang and expressions, making it both a fun educational tool and a love letter to Malaysian culture.

The file extension for bolehlang programs is `.boleh` (short for "Malaysian Language").

## Features

- Manglish-flavoured syntax with keywords like `dei`, `jom`, `eh`, and `abuden`
- Variables and constants (`dei` / `memang dei`)
- Control flow with `is it`, `or is it`, and `abuden`
- Three styles of loops: while, range, and foreach
- First-class functions, closures, and recursion
- Arrays with indexing and built-in methods
- Error handling with `try see` / `alamak` / `confirm need`
- Built-in functions for type checking, conversion, and random numbers
- Imports with `tapau` (Manglish for "takeaway") with circular import protection
- Expressive Malaysian-style error messages

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) 1.22 or later

### Install

```bash
git clone https://github.com/kooixh/bolehlang.git
cd bolehlang
./install.sh
```

This builds the project and installs the `boleh` command to `/usr/local/bin` so you can use it from anywhere.

## Usage

```bash
# Run a bolehlang program
boleh hello.boleh

# Show version
boleh --version

# Show help
boleh --help
```

## Language Reference

### Variables

Declare variables with `dei` and constants with `memang dei`. Reassign variables using `tukar`.

```
dei nama = "Ali";
dei umur = 25;

memang dei PI = 3.14159;

tukar umur = 26;
tukar umur += 1;
tukar umur -= 5;
```

Constants cannot be reassigned. Attempting to do so will give you:

```
Aiyo, a memang variable 'PI' cannot be redefined
```

### Data Types

| Type | bolehlang Name | Example |
|------|---------------|---------|
| Number | `nombor` | `42`, `3.14` |
| String | `ayat` | `"Hello"` |
| Boolean | `boleh_tak_boleh` | `can`, `cannot` |
| Array | `senarai` | `[1, 2, 3]` |
| Null | `kosong` | `kosong` |
| Function | `kerja` | `jom add(a, b) { ... }` |

### Operators

#### Arithmetic
| Operator | Description |
|----------|-------------|
| `+` | Addition / string concatenation |
| `-` | Subtraction |
| `*` | Multiplication |
| `/` | Division |
| `%` | Modulo |

#### Comparison
| Operator | Description |
|----------|-------------|
| `sama` | Equal to |
| `tak sama` | Not equal to |
| `<` | Less than |
| `>` | Greater than |
| `<=` | Less than or equal |
| `>=` | Greater than or equal |

#### Logical
| Operator | Description |
|----------|-------------|
| `also` | Logical AND |
| `if not` | Logical OR |
| `!` | Logical NOT |

### Print

Use `eh` to print to the console.

```
eh "Hello!";
eh 42;
eh "Nama: " + nama;
```

### Control Flow

#### If / Else If / Else

Use `is it` for if, `or is it` for else if, and `abuden` for else. Conditions must be in parentheses.

```
is it (umur >= 18) {
    eh "You're an adult already lah";
} or is it (umur >= 13) {
    eh "Teenager only";
} abuden {
    eh "Still small kid";
}
```

### Loops

All loops use the `one by one` keyword. The loop expression must be in parentheses.

#### While-style loop

```
dei count = 0;
one by one (count < 5) {
    eh count;
    tukar count += 1;
}
```

#### Range-style loop

Iterates from a start value to an end value (inclusive), using `sampai` (meaning "until").

```
one by one (i = 1 sampai 10) {
    eh i;
}
```

#### Foreach-style loop

Iterates over each element in an array, using `dalam` (meaning "in").

```
dei fruits = ["durian", "rambutan", "mangosteen"];
one by one (buah dalam fruits) {
    eh buah;
}
```

#### Loop control

- `gaodim` — break out of a loop (literally "enough already")
- `no need` — skip to the next iteration (literally "skip it")

```
one by one (i = 1 sampai 10) {
    is it (i sama 5) {
        gaodim;
    }
    eh i;
}
```

### Functions

Declare functions with `jom` (meaning "let's go") and return values with `nah` (meaning "here, take it").

#### Named functions

```
jom salam(nama) {
    nah "Assalamualaikum, " + nama + "!";
}

eh salam("Ahmad");
```

#### Anonymous functions

```
dei greet = jom (nama) {
    nah "Hello, " + nama;
};

eh greet("Ali");
```

#### Recursion

```
jom factorial(n) {
    is it (n <= 1) {
        nah 1;
    }
    nah n * factorial(n - 1);
}

eh factorial(5);
```

#### Closures

Functions capture their surrounding scope.

```
jom buatCounter(mula) {
    dei count = mula;
    jom tambah() {
        tukar count = count + 1;
        nah count;
    }
    nah tambah;
}

dei counter = buatCounter(0);
eh counter();
eh counter();
eh counter();
```

### Error Handling

Use `try see` to attempt something, `alamak` to catch errors, and `confirm need` for cleanup that always runs.

```
try see {
    eh "Trying something risky...";
    mampus "Something went wrong!";
} alamak (err) {
    eh "Caught error: " + err;
} confirm need {
    eh "This always runs";
}
```

- `try see` — try block (literally "try and see")
- `alamak` — catch block (a Malaysian exclamation of surprise)
- `confirm need` — finally block (always executes)
- `mampus` — throw an error (literally "you're dead")

### Imports

Use `tapau` (meaning "takeaway") to import code from other `.boleh` files. All declarations from the imported file become available in the current scope.

```
tapau "math.boleh";

eh tambah(2, 3);
```

Paths are resolved relative to the importing file. Circular imports are detected and safely skipped — each file is only executed once.

```
// eh: lib/helpers.boleh
jom greet(nama) {
    nah "Hello, " + nama + "!";
}
```

```
// eh: main.boleh
tapau "lib/helpers.boleh";

eh greet("Ali");
```

### Built-in Functions

| Function | Description | Example |
|----------|-------------|---------|
| `len(x)` | Get length of string or array | `len("hello")` returns `5` |
| `apa_ni(x)` | Get the type of a value | `apa_ni(42)` returns `"nombor"` |
| `masuk(arr, val)` | Append value(s) to array | `masuk(arr, "new")` |
| `potong_stim(arr, start, end)` | Slice an array | `potong_stim(arr, 1, 3)` |
| `nombor(x)` | Convert to number | `nombor("42")` returns `42` |
| `string(x)` | Convert to string | `string(42)` returns `"42"` |
| `cincai()` | Random float between 0 and 1 | `cincai()` returns `0.7342...` |
| `cincai(max)` | Random int from 0 to max (exclusive) | `cincai(10)` returns `0`-`9` |
| `cincai(min, max)` | Random int from min to max (inclusive) | `cincai(1, 6)` returns `1`-`6` |

Arrays also support `.panjang` property access:

```
dei arr = [1, 2, 3];
eh arr.panjang;
```

### Comments

```
// eh: this is a single-line comment

alamak mula
    This is a
    multi-line comment
alamak habis
```

## Examples

### Hello World

```
eh "Hello World! Apa khabar?";
```

### FizzBuzz

```
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
```

### Closures

```
jom buatCounter(mula) {
    dei count = mula;
    jom tambah() {
        tukar count = count + 1;
        nah count;
    }
    nah tambah;
}

dei counter = buatCounter(0);
eh counter();
eh counter();
eh counter();
```

Output:
```
1
2
3
```

More examples are available in the [`examples/`](examples/) directory.

## Keyword Reference

A complete mapping of bolehlang keywords to their conventional equivalents:

| bolehlang | Conventional | Meaning |
|-----------|-------------|---------|
| `dei` | `let` / `var` | Declare a variable |
| `memang dei` | `const` | Declare a constant |
| `tukar` | reassign | Reassign a variable |
| `eh` | `print` | Print to console |
| `is it` | `if` | Conditional check |
| `or is it` | `else if` | Alternate condition |
| `abuden` | `else` | Default branch ("obviously lah") |
| `one by one` | `for` / `while` | Loop |
| `sampai` | `to` / `until` | Range end |
| `dalam` | `in` | Iterate over collection |
| `gaodim` | `break` | Break from loop ("enough already") |
| `no need` | `continue` | Skip iteration |
| `jom` | `function` | Declare function ("let's go") |
| `nah` | `return` | Return value ("here, take it") |
| `can` | `true` | Boolean true |
| `cannot` | `false` | Boolean false |
| `kosong` | `null` | Null value ("empty") |
| `sama` | `==` | Equal |
| `tak sama` | `!=` | Not equal |
| `also` | `&&` | Logical AND |
| `if not` | `\|\|` | Logical OR |
| `try see` | `try` | Try block |
| `alamak` | `catch` | Catch block ("oh no!") |
| `confirm need` | `finally` | Finally block |
| `mampus` | `throw` | Throw error ("you're dead") |
| `tapau` | `import` | Import a file ("takeaway") |
| `cincai` | `random` | Random number ("whatever/random") |

## Project Structure

```
bolehlang/
├── main.go                 # Entry point and CLI
├── go.mod                  # Go module definition
├── lexer/
│   └── lexer.go            # Tokenizer — keywords, symbols, literals
├── parser/
│   └── parser.go           # AST builder — precedence climbing parser
├── interpreter/
│   └── interpreter.go      # Tree-walking interpreter and built-ins
└── examples/
    ├── hello.boleh            # Hello world
    ├── variables.boleh        # Variables and constants
    ├── loops.boleh            # While, range, and foreach loops
    ├── functions.boleh        # Functions, recursion, closures
    ├── fizzbuzz.boleh         # Classic FizzBuzz
    ├── trycatch.boleh         # Error handling
    ├── math.boleh             # Importable math utilities
    └── imports.boleh          # Import example using tapau
```

## Architecture

bolehlang follows a classic three-stage interpreter pipeline:

1. **Lexer** (`lexer/lexer.go`) — Scans source code into tokens. Uses longest-first greedy matching for multi-word keywords (e.g., `one by one`, `memang dei`). Handles string escapes, decimal numbers, and both single-line and multi-line comments.

2. **Parser** (`parser/parser.go`) — Transforms the token stream into an Abstract Syntax Tree (AST). Uses recursive descent with precedence climbing for expressions. Supports statements (declarations, assignments, control flow, loops, functions) and expressions (binary, unary, calls, member/index access).

3. **Interpreter** (`interpreter/interpreter.go`) — Tree-walking interpreter that evaluates the AST. Implements lexical scoping with an environment chain, first-class functions with closures, and control flow via Go panic/recover for return, break, continue, and throw signals.

## Contributing

Contributions are welcome! Whether it's:

- Adding new Manglish keywords
- New built-in functions
- Bug fixes
- More example programs
- Documentation improvements

Feel free to open an issue or submit a pull request.

## License

MIT
