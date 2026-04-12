# Imports

## Importing files

Use `tapau` ("takeaway") to import code from other `.boleh` files. All declarations from the imported file become available in the current scope.

```
tapau "math.boleh";

eh tambah(2, 3);
```

## Path resolution

Import paths are resolved **relative to the importing file**, not the working directory.

Given this structure:

```
project/
  main.boleh
  lib/
    helpers.boleh
    utils.boleh
```

From `main.boleh`:

```
tapau "lib/helpers.boleh";
```

From `lib/utils.boleh`, to import a sibling:

```
tapau "helpers.boleh";
```

## Circular imports

bolehlang automatically detects circular imports and safely skips them. Each file is only executed once, even if imported from multiple places.

```
// a.boleh
tapau "b.boleh";

// b.boleh
tapau "a.boleh";    // safely skipped — already imported
```

## Example

### Library file — `math.boleh`

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

### Main file — `imports.boleh`

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

## Missing files

If the imported file doesn't exist, a `BoJioException` is thrown:

```
tapau "doesnt_exist.boleh";
// BoJioException: cannot tapau 'doesnt_exist.boleh' — file not found lah!
```

## Next steps

- [Keyword Reference](12-Keyword-Reference.md) — complete mapping of all bolehlang keywords
