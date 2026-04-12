# bolehlang

A dynamically typed, interpreted programming language inspired by Malaysian English (Manglish). Build by  🇲🇾 for 🇲🇾. 

```
eh "Hello World! Apa khabar?";
```

## Documentation

Full language reference and tutorials live in **[`docs/`](docs/)**:

| Topic | Doc |
|--------|-----|
| Install & prerequisites | [Installation](docs/01-Installation.md) |
| Run your first program, CLI | [Getting started](docs/02-Getting-Started.md) |
| VS Code syntax highlighting | [VS Code extension](docs/03-VSCode-Extension.md) |
| Variables, types, literals | [Variables and types](docs/04-Variables-and-Types.md) |
| Arithmetic, comparison, logic | [Operators](docs/05-Operators.md) |
| `is it` / `or is it` / `abuden` | [Control flow](docs/06-Control-Flow.md) |
| `one by one`, `sampai`, `dalam`, break/continue | [Loops](docs/07-Loops.md) |
| `jom`, `nah`, closures | [Functions](docs/08-Functions.md) |
| `try see`, `alamak`, `mampus` | [Error handling](docs/09-Error-Handling.md) |
| Built-ins (`len`, `apa_ni`, …) | [Built-in functions](docs/10-Built-in-Functions.md) |
| `tapau` imports | [Imports](docs/11-Imports.md) |
| Keyword ↔ meaning cheat sheet | [Keyword reference](docs/12-Keyword-Reference.md) |
| Walkthrough examples | [Examples walkthrough](docs/13-Examples-Walkthrough.md) |

## Quick start

This Project requries Go. 

```bash
git clone https://github.com/kooixh/bolehlang.git
cd bolehlang
./install.sh          # installs `boleh` to /usr/local/bin
boleh examples/hello.boleh
```

See [Getting started](docs/02-Getting-Started.md) for `boleh` CLI flags and more samples. Extra programs are under [`examples/`](examples/).

## Features (overview)

Manglish-style keywords (`dei`, `jom`, `eh`, `abuden`, …), lexical scoping, first-class functions and closures, arrays, try/catch/finally-style errors, `tapau` imports, and built-ins for types, conversion, and randomness. Details are in the docs linked above.

## Project layout

```
bolehlang/
├── main.go              # CLI entry
├── lexer/               # Grammar and Syntax
├── parser/              # AST
├── interpreter/         # Evaluator + builtins functions
├── vscode/              # VS Code extension package
├── examples/            # Sample .boleh files
└── docs/                # Language documentation
└── tests/                # Tests
```

## Architecture

Classic pipeline: **lexer** → **parser** (recursive descent + precedence climbing) → **tree-walking interpreter** with environment chains for scope; control flow for loops/returns uses bounded panic/recover patterns inside the interpreter.

## Contributing

Issues and pull requests are welcome—language ideas, builtins, examples, or doc fixes.

## License

[MIT](LICENSE)
