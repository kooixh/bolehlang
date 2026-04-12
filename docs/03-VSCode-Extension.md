# VS Code Extension

Set up syntax highlighting for bolehlang in Visual Studio Code.

## Installation

### From source

1. Symlink the `vscode` directory into your VS Code extensions folder:

   **macOS / Linux:**

   ```bash
   ln -s "$(pwd)/vscode" ~/.vscode/extensions/bolehlang
   ```

   **Windows (PowerShell):**

   ```powershell
   New-Item -ItemType SymbolicLink -Path "$env:USERPROFILE\.vscode\extensions\bolehlang" -Target "$(Get-Location)\vscode"
   ```

2. Reload VS Code (`Cmd+Shift+P` > "Developer: Reload Window")

3. Open any `.boleh` file — syntax highlighting should be active.

## Highlighted elements

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
| Type annotations | `: number`, `: string`, `: bool`, `: array`, `: function` |
| Comments | `// eh: ...`, `alamak mula ... alamak habis` |

## Next steps

- [Variables and Types](04-Variables-and-Types.md) — learn about data types and variable declaration
