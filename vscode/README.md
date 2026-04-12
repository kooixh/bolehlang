# bolehlang VS Code Extension

Syntax highlighting for the [bolehlang](https://github.com/kooixh/bolehlang) programming language.

## Features

- Syntax highlighting for all bolehlang keywords, operators, and built-in functions
- Auto-closing pairs for brackets, braces, parentheses, and strings
- Comment toggling support
- Code folding for blocks

## Installation

### Download From VS Code Marketplace (Recommended for general use)

1. Download the bolehlang extension from VS Code Marketplace (https://marketplace.visualstudio.com/items?itemName=bolehlang.bolehlang)

### From source (recommended for development)

1. Copy or symlink the `vscode` directory into your VS Code extensions folder:

   ```bash
   # macOS / Linux
   ln -s "$(pwd)/vscode" ~/.vscode/extensions/bolehlang

   # Windows (PowerShell)
   New-Item -ItemType SymbolicLink -Path "$env:USERPROFILE\.vscode\extensions\bolehlang" -Target "$(Get-Location)\vscode"
   ```

2. Reload VS Code (Cmd+Shift+P > "Developer: Reload Window")

3. Open any `.boleh` file and enjoy syntax highlighting!

## Highlighted Elements

| Element | Examples |
|---------|---------|
| Variables & constants | `dei`, `memang dei`, `tukar` |
| Control flow | `is it`, `or is it`, `abuden`, `one by one` |
| Functions | `jom`, `nah` |
| Booleans & null | `can`, `cannot`, `kosong` |
| Error handling | `try see`, `alamak`, `confirm need`, `mampus` |
| Print | `eh`, `walao eh` |
| Imports | `tapau` |
| Loop keywords | `sampai`, `dalam`, `gaodim`, `no need` |
| Logical & comparison | `sama`, `tak sama`, `also`, `if not` |
| Built-in functions | `len`, `apa_ni`, `masuk`, `potong_stim`, `nombor`, `string`, `cincai` |
| Comments | `// eh: ...`, `alamak mula ... alamak habis` |
