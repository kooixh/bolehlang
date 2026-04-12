# Loops

All loops in bolehlang use the `one by one` keyword. There are three styles: while, range, and foreach.

## While-style loop

Loops while the condition is truthy:

```
dei count = 0;
one by one (count < 5) {
    eh count;
    tukar count += 1;
}
```

Output:

```
0
1
2
3
4
```

## Range-style loop

Iterates from a start value to an end value (**inclusive**), using `sampai` ("until"):

```
one by one (i = 1 sampai 5) {
    eh i;
}
```

Output:

```
1
2
3
4
5
```

The loop variable is scoped to each iteration — you don't need to declare it beforehand.

## Foreach-style loop

Iterates over elements in an array, using `dalam` ("in"):

```
dei fruits = ["durian", "rambutan", "mangosteen"];
one by one (buah dalam fruits) {
    eh buah;
}
```

Output:

```
durian
rambutan
mangosteen
```

Using `dalam` on a non-array value throws an `ApaNiException`.

## Loop control

### Break — `gaodim`

Exit a loop early (literally "enough already"):

```
one by one (i = 1 sampai 10) {
    is it (i sama 4) {
        gaodim;
    }
    eh i;
}
```

Output:

```
1
2
3
```

### Continue — `no need`

Skip to the next iteration:

```
one by one (i = 1 sampai 5) {
    is it (i sama 3) {
        no need;
    }
    eh i;
}
```

Output:

```
1
2
4
5
```

## Nested loops

Loops can be nested. `gaodim` and `no need` apply to the innermost loop:

```
one by one (i = 1 sampai 3) {
    one by one (j = 1 sampai 3) {
        eh "" + i + "," + j;
    }
}
```

## Next steps

- [Functions](08-Functions.md) — named functions, anonymous functions, closures, and recursion
