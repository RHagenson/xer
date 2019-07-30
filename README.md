[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.3355641.svg)](https://doi.org/10.5281/zenodo.3355641)

# `xer`

`xer` simple utility that Xs out every visual character to force looking at code layout over other elements of code structure. The idea spawned from a talk by Kevlin Henney and reading far too much code with triple-digit line lengths.

Users can also selectively unmask by select Unicode categories or by regex. This functionality is to allow the meta elements to be unmasked (parenthesis, brackets, curly braces, constants, and so on).

## Philosophy on Use

Code should suggest logical structure by how it is written. By masking all but
the whitespace we force ourselves to look at the structure differently.
Some code is inherently messy, however most code can be cleaned up.

### Use case
A good use case for running this utility is checking nested statements.
Statements that are nested under one another should appear nested.

Example:

```
if in, err = os.Open(*read); err != nil {
  log.Fatal(err)
}
```

versus

```
if in, err = os.Open(*read); err != nil {log.Fatal(err)}
```

becomes

```
XX XXX XXX X XXXXXXXXXXXXXXX XXX XX XXX X
  XXXXXXXXXXXXXX
X
```

versus

```
XX XXX XXX X XXXXXXXXXXXXXXX XXX XX XXX XXXXXXXXXXXXXXXX
```

The former suggests the logical structure more so than the latter, which
simply appears as a long line.
