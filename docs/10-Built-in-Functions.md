# Built-in Functions

## Overview

| Function | Description |
|----------|-------------|
| `len(x)` | Get length of string or array |
| `apa_ni(x)` | Get the runtime type of a value |
| `masuk(arr, val, ...)` | Append value(s) to an array |
| `potong_stim(arr, start, end?)` | Slice an array |
| `nombor(x)` | Convert to number |
| `string(x)` | Convert to string |
| `cincai(...)` | Generate random numbers |

## `len(x)`

Returns the length of a string or array.

```
eh len("hello");     // 5
eh len([1, 2, 3]);   // 3
eh len("");           // 0
eh len([]);           // 0
```

Throws an error if the argument is not a string or array.

Arrays and strings also support `.panjang` property access:

```
dei arr = [1, 2, 3];
eh arr.panjang;      // 3
eh "hello".panjang;  // 5
```

## `apa_ni(x)`

Returns the runtime type of a value as a string.

```
eh apa_ni(42);           // "nombor"
eh apa_ni("hello");      // "ayat"
eh apa_ni(can);          // "boleh_tak_boleh"
eh apa_ni([1, 2]);       // "senarai"
eh apa_ni(kosong);       // "kosong"
```

For functions (both user-defined and built-in):

```
jom f() {}
eh apa_ni(f);            // "kerja"
eh apa_ni(len);          // "kerja"
```

For exceptions, returns the exception type name:

```
dei e = KantoiException("oops");
eh apa_ni(e);            // "KantoiException"
```

## `masuk(arr, val, ...)`

Appends one or more values to an array and returns the new array.

**Important:** `masuk` returns a new array — you must reassign it.

```
dei arr = [1, 2];
tukar arr = masuk(arr, 3);
eh arr;    // [1, 2, 3]
```

Append multiple values at once:

```
dei arr = [];
tukar arr = masuk(arr, 1, 2, 3);
eh arr;    // [1, 2, 3]
```

## `potong_stim(arr, start, end?)`

Slices an array from `start` index to `end` index (exclusive). If `end` is omitted, slices to the end of the array.

```
dei arr = [10, 20, 30, 40, 50];

eh potong_stim(arr, 1, 3);    // [20, 30]
eh potong_stim(arr, 2);       // [30, 40, 50]
eh potong_stim(arr, 0, 2);    // [10, 20]
```

Returns a new array — the original is not modified.

## `nombor(x)`

Converts a value to a number.

```
eh nombor("42");       // 42
eh nombor("3.14");     // 3.14
eh nombor(can);        // 1
eh nombor(cannot);     // 0
eh nombor(42);         // 42 (no-op)
```

Throws an error for strings that can't be parsed as numbers:

```
nombor("abc");    // Error: tak boleh convert 'abc' to nombor lah
```

## `string(x)`

Converts any value to its string representation.

```
eh string(42);         // "42"
eh string(3.14);       // "3.14"
eh string(can);        // "can"
eh string(cannot);     // "cannot"
eh string(kosong);     // "kosong"
eh string([1, 2]);     // "[1, 2]"
```

## `cincai(...)`

Generates random numbers. The behavior depends on the number of arguments:

### No arguments — random float

Returns a random float between 0 (inclusive) and 1 (exclusive):

```
dei x = cincai();    // e.g., 0.7342...
```

### One argument — random int from 0

Returns a random integer from 0 (inclusive) to `max` (exclusive):

```
dei x = cincai(10);    // 0 to 9
```

### Two arguments — random int in range

Returns a random integer from `min` (inclusive) to `max` (inclusive):

```
dei x = cincai(1, 6);    // 1 to 6 (like a dice roll)
```

## Next steps

- [Imports](11-Imports.md) — importing code from other files
