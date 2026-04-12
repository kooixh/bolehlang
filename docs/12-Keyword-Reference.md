# Keyword Reference

A complete mapping of bolehlang keywords to their conventional equivalents.

## Language keywords

| bolehlang | Conventional | Meaning |
|-----------|-------------|---------|
| `dei` | `let` / `var` | Declare a variable |
| `memang dei` | `const` | Declare a constant ("confirm got") |
| `tukar` | reassign | Reassign a variable ("change") |
| `eh` | `print` | Print to stdout |
| `walao eh` | `warn` | Print to stderr ("oh my god") |
| `lai` | `readline` | Read a line from stdin as a string ("come") |
| `is it` | `if` | Conditional check |
| `or is it` | `else if` | Alternate condition |
| `abuden` | `else` | Default branch ("obviously lah") |
| `one by one` | `for` / `while` | Loop |
| `sampai` | `to` / `until` | Range end ("until") |
| `dalam` | `in` | Iterate over collection ("inside") |
| `gaodim` | `break` | Break from loop ("done already") |
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
| `try see` | `try` | Try block ("try and see") |
| `alamak` | `catch` | Catch block ("oh no!") |
| `confirm need` | `finally` | Finally block ("must do") |
| `mampus` | `throw` | Throw error ("you're dead") |
| `tapau` | `import` | Import a file ("takeaway") |

## Type annotation keywords

| Keyword | Type |
|---------|------|
| `number` | Numeric values (float64) |
| `string` | Text values |
| `bool` | Boolean (`can` / `cannot`) |
| `array` | Arrays (`[...]`) |
| `function` | Functions |

## Built-in functions

| Function | Description |
|----------|-------------|
| `len` | Get length of string or array |
| `apa_ni` | Get runtime type of a value ("what is this?") |
| `masuk` | Append to array ("enter") |
| `potong_stim` | Slice an array ("cut nicely") |
| `nombor` | Convert to number ("number") |
| `string` | Convert to string |
| `cincai` | Random number ("whatever/random") |

## Exception types

| Exception | Description |
|-----------|-------------|
| `MakKauHijauException` | Null pointer access |
| `BoJioException` | Import file not found |
| `ApaNiException` | Invalid type or wrong usage ("what is this?") |
| `KantoiException` | General runtime error ("busted") |

## Next steps

- [Examples Walkthrough](13-Examples-Walkthrough.md) — practical examples with explanations
