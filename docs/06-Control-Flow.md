# Control Flow

## If / Else If / Else

Use `is it` for if, `or is it` for else if, and `abuden` for else. Conditions must be wrapped in parentheses.

### Basic if

```
is it (umur >= 18) {
    eh "Adult already lah";
}
```

### If / else

```
is it (score >= 50) {
    eh "Pass!";
} abuden {
    eh "Fail lah";
}
```

### If / else if / else

```
is it (umur >= 18) {
    eh "Adult";
} or is it (umur >= 13) {
    eh "Teenager";
} abuden {
    eh "Still small kid";
}
```

You can chain as many `or is it` blocks as you need:

```
dei grade = 85;

is it (grade >= 90) {
    eh "A";
} or is it (grade >= 80) {
    eh "B";
} or is it (grade >= 70) {
    eh "C";
} or is it (grade >= 60) {
    eh "D";
} abuden {
    eh "F";
}
```

## Scoping

Each block (`{ }`) creates a new scope. Variables declared inside a block are not visible outside:

```
dei x = "outer";
is it (can) {
    dei x = "inner";
    eh x;    // "inner"
}
eh x;    // "outer"
```

## Next steps

- [Loops](07-Loops.md) — while, range, and foreach loops
