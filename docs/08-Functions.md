# Functions

## Declaring functions

Use `jom` ("let's go") to declare a function and `nah` ("here, take it") to return a value:

```
jom salam(nama) {
    nah "Assalamualaikum, " + nama + "!";
}

eh salam("Ahmad");    // "Assalamualaikum, Ahmad!"
```

## Parameters

Functions can take zero or more parameters:

```
jom greet() {
    eh "Hello!";
}

jom add(a, b) {
    nah a + b;
}
```

Calling a function with the wrong number of arguments throws an `ApaNiException`:

```
add(1);    // ApaNiException: function 'add' expects 2 args, got 1
```

## Return values

Use `nah` to return a value. A function without `nah` returns `kosong`:

```
jom doNothing() {
    dei x = 1;
}

eh doNothing();    // "kosong"
```

`nah` without a value returns `kosong`:

```
jom earlyReturn() {
    nah;
}
```

## Typed parameters and return types

Functions support optional type annotations on parameters and return types:

```
jom add(a: number, b: number): number {
    nah a + b;
}

eh add(3, 4);        // 7
eh add("x", 4);     // ApaNiException: 'a' expected type number, got string
```

You can mix typed and untyped parameters:

```
jom greet(name: string, times) {
    one by one (i = 1 sampai times) {
        eh "Hello, " + name + "!";
    }
}
```

Return types are checked when the function returns:

```
jom getName(): number {
    nah "Ali";    // ApaNiException: 'getName return' expected type number, got string
}
```

## Anonymous functions

Use `jom` without a name to create an anonymous function:

```
dei double = jom (x) {
    nah x * 2;
};

eh double(5);    // 10
```

Anonymous functions also support type annotations:

```
dei square = jom (x: number): number {
    nah x * x;
};
```

## Closures

Functions capture their surrounding scope:

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
eh counter();    // 1
eh counter();    // 2
eh counter();    // 3
```

## Recursion

Functions can call themselves:

```
jom factorial(n) {
    is it (n <= 1) {
        nah 1;
    }
    nah n * factorial(n - 1);
}

eh factorial(5);    // 120
```

## Functions as values

Functions are first-class values. They can be stored in variables, passed as arguments, and returned from other functions:

```
jom a() { nah "a"; }
jom b() { nah "b"; }

dei fns = [a, b];
eh fns[0]();    // "a"
eh fns[1]();    // "b"
```

## Next steps

- [Error Handling](09-Error-Handling.md) — try/catch/finally and exceptions
