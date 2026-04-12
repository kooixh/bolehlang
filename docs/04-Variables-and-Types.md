# Variables and Types

## Declaring variables

Use `dei` to declare a variable:

```
dei nama = "Ali";
dei umur = 25;
dei programmer = can;
```

## Constants

Use `memang dei` to declare a constant. Constants cannot be reassigned:

```
memang dei PI = 3.14159;
```

Attempting to reassign a constant gives an error:

```
Aiyo, a memang variable 'PI' cannot be redefined
```

## Reassignment

Use `tukar` to reassign variables. Supports `=`, `+=`, and `-=`:

```
dei x = 10;
tukar x = 20;
tukar x += 5;     // x is now 25
tukar x -= 3;     // x is now 22
```

## Data types

bolehlang is dynamically typed. There are six core types:

| Type | Name | Literal examples |
|------|------|-----------------|
| Number | `number` | `42`, `3.14` |
| String | `string` | `"Hello"`, `"Apa khabar"` |
| Boolean | `bool` | `can`, `cannot` |
| Array | `array` | `[1, 2, 3]`, `[]` |
| Null | — | `kosong` |
| Function | `function` | `jom (x) { nah x; }` |

### Numbers

All numbers are floating point. There is no separate integer type:

```
dei umur = 25;
dei harga = 9.99;
```

### Strings

Strings are enclosed in double quotes and support escape sequences:

```
dei greeting = "Hello\tworld\n";
dei quote = "She said \"boleh\"";
```

Supported escapes: `\n` (newline), `\t` (tab), `\\` (backslash), `\"` (double quote).

### Booleans

`can` is true, `cannot` is false:

```
dei active = can;
dei deleted = cannot;
```

### Null

`kosong` represents null/empty:

```
dei nothing = kosong;
```

### Arrays

Arrays can hold mixed types:

```
dei numbers = [1, 2, 3];
dei mixed = ["hello", 42, can, kosong];
dei empty = [];
```

Access elements by index (zero-based):

```
dei fruits = ["durian", "rambutan", "mangosteen"];
eh fruits[0];    // "durian"
eh fruits[2];    // "mangosteen"
```

Get the length with `.panjang`:

```
eh fruits.panjang;    // 3
```

### Truthiness

bolehlang values have the following truthiness rules:

| Value | Truthy? |
|-------|---------|
| `can` | yes |
| `cannot` | no |
| `0` | no |
| any other number | yes |
| `""` (empty string) | no |
| any other string | yes |
| `[]` (empty array) | no |
| any other array | yes |
| `kosong` | no |

## Type annotations

Type annotations are optional. When provided, they are checked at runtime:

```
dei x: number = 42;
dei name: string = "Ali";
dei active: bool = can;
dei items: array = [1, 2, 3];
dei callback: function = jom (x) { nah x; };
```

If the value doesn't match the declared type, an `ApaNiException` is thrown:

```
dei x: number = "hello";
// ApaNiException: line 1: 'x' expected type number, got string
```

Type annotations are also enforced on reassignment:

```
dei x: number = 42;
tukar x = "hello";    // ApaNiException!
tukar x = 100;        // ok
```

Available type names: `number`, `string`, `bool`, `array`, `function`.

Omitting the type annotation preserves the default untyped behavior — any value can be assigned:

```
dei x = 42;
tukar x = "hello";    // ok, no type constraint
```

## Checking types at runtime

Use the `apa_ni()` built-in to inspect a value's type:

```
eh apa_ni(42);         // "nombor"
eh apa_ni("hello");    // "ayat"
eh apa_ni(can);        // "boleh_tak_boleh"
eh apa_ni([1, 2]);     // "senarai"
eh apa_ni(kosong);     // "kosong"
```

## Next steps

- [Operators](05-Operators.md) — arithmetic, comparison, and logical operators
