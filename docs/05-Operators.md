# Operators

## Arithmetic operators

| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Addition / string concatenation | `2 + 3` gives `5` |
| `-` | Subtraction | `10 - 4` gives `6` |
| `*` | Multiplication | `3 * 4` gives `12` |
| `/` | Division | `15 / 3` gives `5` |
| `%` | Modulo | `10 % 3` gives `1` |
| `-` (unary) | Negation | `-5` |

### String concatenation

The `+` operator concatenates when either operand is a string:

```
eh "Hello" + " " + "World";    // "Hello World"
eh "Age: " + 25;               // "Age: 25"
eh "Active: " + can;           // "Active: can"
```

### Division by zero

Division or modulo by zero throws a `KantoiException`:

```
eh 10 / 0;    // KantoiException: cannot divide by zero
```

## Comparison operators

| Operator | Description | Example |
|----------|-------------|---------|
| `sama` | Equal to | `x sama 1` |
| `tak sama` | Not equal to | `x tak sama 2` |
| `<` | Less than | `x < 10` |
| `>` | Greater than | `x > 0` |
| `<=` | Less than or equal | `x <= 100` |
| `>=` | Greater than or equal | `x >= 1` |

Equality checks (`sama`, `tak sama`) compare by value and type — a number is never equal to a string.

## Logical operators

| Operator | Description | Example |
|----------|-------------|---------|
| `also` | Logical AND | `can also can` gives `can` |
| `if not` | Logical OR | `cannot if not can` gives `can` |
| `!` | Logical NOT | `!can` gives `cannot` |

### Short-circuit evaluation

`also` and `if not` use short-circuit evaluation:

- `also` returns the left operand if it's falsy; otherwise evaluates and returns the right.
- `if not` returns the left operand if it's truthy; otherwise evaluates and returns the right.

```
dei x = cannot also someFunction();    // someFunction() is never called
dei y = can if not someFunction();     // someFunction() is never called
```

## Operator precedence

From lowest to highest:

1. `if not` (logical OR)
2. `also` (logical AND)
3. `sama`, `tak sama` (equality)
4. `<`, `>`, `<=`, `>=` (comparison)
5. `+`, `-` (addition/subtraction)
6. `*`, `/`, `%` (multiplication/division/modulo)
7. `!`, `-` (unary NOT, negation)

Use parentheses to override precedence:

```
eh (2 + 3) * 4;    // 20
eh 2 + 3 * 4;      // 14
```

## Next steps

- [Control Flow](06-Control-Flow.md) — conditionals with `is it`, `or is it`, and `abuden`
