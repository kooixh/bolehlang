# Getting Started

Write and run your first bolehlang program.

## Your first program

Create a file called `hello.boleh`:

```
eh "Hello World! Apa khabar?";
```

Run it:

```bash
boleh hello.boleh
```

Output:

```
Hello World! Apa khabar?
```

## File extension

All bolehlang programs use the `.boleh` file extension.

## Semicolons

Semicolons are **optional** in bolehlang. Both of these are valid:

```
eh "with semicolon";
eh "without semicolon"
```

## Comments

### Single-line comments

```
// eh: this is a comment
eh "Hello";
```

### Multi-line comments

```
alamak mula
    This is a block comment.
    It can span multiple lines.
alamak habis
```

## Printing output

Use `eh` to print to stdout:

```
eh "Hello!";
eh 42;
eh "Nama: " + "Ali";
```

Use `walao eh` to print warnings to stderr:

```
walao eh "Something not right lah!";
```

This outputs:

```
[WALAO EH] Something not right lah!
```

## CLI options

```bash
boleh <file.boleh>    # Run a program
boleh --version       # Show version
boleh --help          # Show help
```

## Next steps

- [Variables and Types](04-Variables-and-Types.md) — learn about data types and variable declaration
