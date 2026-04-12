# Error Handling

## Try / Catch / Finally

Use `try see` to attempt something, `alamak` to catch errors, and `confirm need` for cleanup that always runs.

```
try see {
    eh "Trying something risky...";
    mampus "Something went wrong!";
    eh "This won't print";
} alamak (err) {
    eh "Caught: " + err;
} confirm need {
    eh "This always runs";
}
```

Output:

```
Trying something risky...
Caught: Something went wrong!
This always runs
```

## Throwing errors

Use `mampus` ("you're dead") to throw an error:

```
mampus "Alamak, something broke!";
```

You can throw any value — strings, numbers, or exception objects.

## Exception types

bolehlang has four built-in exception types:

| Exception | Description | When it's thrown |
|-----------|-------------|-----------------|
| `MakKauHijauException` | Null pointer | Accessing a property or index on `kosong` |
| `BoJioException` | Import not found | `tapau` can't find the file |
| `ApaNiException` | Invalid type / wrong usage | Type mismatch, wrong arg count, invalid operations |
| `KantoiException` | General runtime error | Division by zero, variable not found, index out of bounds |

### Creating exceptions manually

Each exception type is available as a constructor function:

```
dei err = KantoiException("something went wrong lah!");
mampus err;
```

### Catching exceptions

The `alamak` block receives the thrown value:

```
try see {
    dei x = kosong;
    eh x.panjang;
} alamak (e) {
    eh e;    // "MakKauHijauException: cannot access property 'panjang' on kosong — it's null lah!"
}
```

## Auto-thrown exceptions

Many operations automatically throw typed exceptions:

```
// MakKauHijauException — null pointer
dei x = kosong;
eh x.panjang;

// ApaNiException — type mismatch
dei y: number = "hello";

// KantoiException — division by zero
dei z = 10 / 0;

// KantoiException — index out of bounds
dei arr = [1, 2];
eh arr[5];

// KantoiException — undefined variable
eh unknownVar;
```

## Try without catch

You can use `try see` with only `confirm need` (no `alamak`):

```
try see {
    mampus "error";
} confirm need {
    eh "cleanup";
}
```

Note: without `alamak`, the exception is not caught — `confirm need` still runs, but the error propagates.

## Try with only catch

`confirm need` is optional:

```
try see {
    mampus "error";
} alamak (e) {
    eh "caught: " + e;
}
```

## Uncaught exceptions

If an exception isn't caught by any `alamak` block, the program exits with an error:

```
mampus "uncaught!";
```

```
Mampus! Runtime error: uncaught exception: uncaught!
```

## Next steps

- [Built-in Functions](10-Built-in-Functions.md) — all available built-in functions
