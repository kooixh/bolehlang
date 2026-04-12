# Examples Walkthrough

Practical bolehlang programs with explanations. All examples are available in the `examples/` directory.

## Hello World

**File:** `examples/hello.boleh`

```
eh "Hello World! Apa khabar?";
```

The simplest bolehlang program. `eh` prints a value to the console.

## Variables

**File:** `examples/variables.boleh`

```
dei nama = "Ali";
dei umur = 25;
dei programmer = can;

eh "Nama: " + nama;
eh "Umur: " + umur;
eh "Programmer: " + programmer;

memang dei PI = 3.14159;
eh "PI = " + PI;

tukar umur = 26;
eh "Happy Birthday! Now umur: " + umur;
```

Output:

```
Nama: Ali
Umur: 25
Programmer: can
PI = 3.14159
Happy Birthday! Now umur: 26
```

This shows variable declaration (`dei`), constants (`memang dei`), reassignment (`tukar`), and string concatenation with the `+` operator.

## FizzBuzz

**File:** `examples/fizzbuzz.boleh`

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

The classic FizzBuzz problem. Demonstrates range loops (`one by one ... sampai`), conditionals (`is it / or is it / abuden`), modulo (`%`), and equality checking (`sama`).

## Loops

**File:** `examples/loops.boleh`

```
// While-style loop
dei count = 0;
one by one (count < 5) {
    eh "Count: " + count;
    tukar count += 1;
}

// Range-style loop
one by one (i = 1 sampai 5) {
    eh "Number: " + i;
}

// Foreach-style loop
dei fruits = ["durian", "rambutan", "mangosteen", "cempedak"];
one by one (buah dalam fruits) {
    eh "Buah: " + buah;
}
```

All three loop styles use the same `one by one` keyword. The syntax inside the parentheses determines which kind of loop it is.

## Functions, Recursion, and Closures

**File:** `examples/functions.boleh`

```
// Basic function
jom salam(nama) {
    nah "Assalamualaikum, " + nama + "!";
}

eh salam("Ahmad");
eh salam("Siti");

// Recursive function
jom factorial(n) {
    is it (n <= 1) {
        nah 1;
    }
    nah n * factorial(n - 1);
}

eh "5! = " + factorial(5);

// Closure
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
Assalamualaikum, Ahmad!
Assalamualaikum, Siti!
5! = 120
1
2
3
```

Demonstrates named functions, return values with `nah`, recursion, and closures that capture surrounding scope.

## Error Handling

**File:** `examples/trycatch.boleh`

```
try see {
    eh "Cuba buat sesuatu...";
    mampus "Alamak something went wrong!";
    eh "This won't print";
} alamak (err) {
    eh "Kantoi! Error: " + err;
} confirm need {
    eh "Mesti buat: cleanup done";
}
```

Output:

```
Cuba buat sesuatu...
Kantoi! Error: Alamak something went wrong!
Mesti buat: cleanup done
```

## Exception Types

**File:** `examples/exceptions.boleh`

```
// Manual exception creation
try see {
    dei err = KantoiException("something went wrong lah!");
    mampus err;
} alamak (e) {
    eh e;
}

// Auto-thrown exceptions
try see {
    dei x = kosong;
    eh x.panjang;
} alamak (e) {
    eh "Caught: " + e;
}

try see {
    dei x = 10 / 0;
} alamak (e) {
    eh "Caught: " + e;
}
```

Shows how to create typed exceptions manually with the constructor functions, and how the interpreter throws them automatically for null access, type errors, and runtime errors.

## Imports

**Files:** `examples/math.boleh` and `examples/imports.boleh`

`math.boleh` defines reusable functions:

```
jom tambah(a, b) {
    nah a + b;
}

jom darab(a, b) {
    nah a * b;
}

jom kuasa(base, exp) {
    dei result = 1;
    one by one (i = 1 sampai exp) {
        tukar result = result * base;
    }
    nah result;
}
```

`imports.boleh` uses them:

```
tapau "math.boleh";

eh "2 + 3 = " + tambah(2, 3);
eh "4 * 5 = " + darab(4, 5);
eh "2 ^ 10 = " + kuasa(2, 10);
```

Output:

```
2 + 3 = 5
4 * 5 = 20
2 ^ 10 = 1024
```

## Type Annotations

Type annotations are optional. When provided, they enforce runtime type checking:

```
// Typed variables
dei x: number = 42;
dei name: string = "Ali";

// Typed function parameters and return type
jom add(a: number, b: number): number {
    nah a + b;
}

eh add(3, 4);    // 7

// Type mismatch throws ApaNiException
try see {
    dei y: number = "hello";
} alamak (e) {
    eh e;    // ApaNiException: line 1: 'y' expected type number, got string
}
```

## Running the examples

```bash
boleh examples/hello.boleh
boleh examples/fizzbuzz.boleh
boleh examples/functions.boleh
boleh examples/loops.boleh
boleh examples/trycatch.boleh
boleh examples/exceptions.boleh
boleh examples/imports.boleh
boleh examples/variables.boleh
```
