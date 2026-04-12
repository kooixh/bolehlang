# Installation

Get bolehlang up and running on your machine.

## Prerequisites

- [Go](https://go.dev/dl/) 1.22 or later

## Install from source

```bash
git clone https://github.com/kooixh/bolehlang.git
cd bolehlang
./install.sh
```

This builds the project and installs the `boleh` command to `/usr/local/bin`.

## Manual build

If you prefer not to use the install script:

```bash
git clone https://github.com/kooixh/bolehlang.git
cd bolehlang
go build -o boleh .
```

Then move the `boleh` binary somewhere on your `$PATH`.

## Verify installation

```bash
boleh --version
```

You should see:

```
bolehlang v0.1.0
```

## Next steps

- [Getting Started](02-Getting-Started.md) — write your first bolehlang program
- [VS Code Extension](03-VSCode-Extension.md) — set up syntax highlighting
