# Yet Another Shell

A cross-platform shell written in Go, featuring a Fyne-based desktop GUI and a classic CLI REPL. It speaks both Unix-style and PowerShell-style command names, and delegates unknown commands to the host OS (PowerShell on Windows, direct `exec` on Unix).

## Purpose

- Explore how a shell works end-to-end: tokenizer, command dispatch, builtins, and OS delegation.
- Provide a single binary that feels at home on both Windows and Linux/macOS.
- Ship a small, readable Go codebase that is easy to extend with new builtins.

## Features

- **Two modes**: desktop GUI (Fyne) by default, or a terminal REPL via `-cli`.
- **Builtins**: `echo`, `cd`, `pwd`, `type`, `exit`.
- **PowerShell-style aliases**: lexer recognizes `Write-Output`, `Set-Location`, `Get-ChildItem`, `Get-Content`, `New-Item`, `Remove-Item`, `Copy-Item`, `Move-Item`, `Get-Location`, `Get-Command`, and more.
- **Host fallback**: anything that isn't a builtin is executed by PowerShell on Windows and by `exec.Command` on Unix.
- **UTF-8 safe on Windows**: forces `[Console]::OutputEncoding = UTF8` for PowerShell child processes.
- **No console flash**: Windows child processes are spawned with `CREATE_NO_WINDOW`.

## Requirements

- **Go 1.23+**
- **CGO toolchain** (required by Fyne for the GUI):
  - **Windows**: MSYS2 + MinGW-w64 (`gcc`), or TDM-GCC. Make sure `gcc` is on `PATH`.
  - **Linux**: `gcc`, plus X11/OpenGL dev packages. On Debian/Ubuntu:
    ```bash
    sudo apt install gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev
    ```
  - **macOS**: Xcode Command Line Tools (`xcode-select --install`).

If you only want the CLI mode you still need CGO for now, because the `gui` package is imported by `main`.

## Getting Started

```bash
# 1. Clone
git clone git@github.com:QBERT18/yet-another-shell.git
cd yet-another-shell

# 2. Fetch dependencies
go mod tidy

# 3. Run the GUI
go run main.go

# 4. Or run the CLI REPL
go run main.go -cli
```

## Building a Standalone Binary

```bash
# Generic build
go build -o yet-another-shell .

# Windows GUI executable (no attached console window)
go build -ldflags "-H windowsgui" -o yet-another-shell.exe .
```

## Usage Examples

### Builtins (both modes)

```
> pwd
/home/user/projects/yet-another-shell

> echo Hello, world!
Hello, world!

> cd ..
> pwd
/home/user/projects

> type echo
echo is a builtin command

> type git
git is /usr/bin/git

> exit
```

### Unix-style commands (Linux/macOS)

Non-builtins are executed directly:

```
> ls -la
> cat go.mod
> git status
```

### PowerShell-style commands (Windows)

On Windows, anything that isn't a builtin is passed to `powershell.exe`, so full PowerShell syntax works:

```
> Get-ChildItem
> Get-Process | Where-Object { $_.CPU -gt 10 }
> (Get-Content README.md).Length
```

The lexer also recognizes PowerShell-style names for its builtins, so these are equivalent:

| Unix style | PowerShell style |
|------------|------------------|
| `echo`     | `Write-Output`   |
| `cd`       | `Set-Location`   |
| `pwd`      | `Get-Location`   |

## Project Layout

```
.
├── main.go              # entry point; picks GUI or CLI (-cli flag)
├── gui/                 # Fyne desktop UI
├── shell/               # Engine: tokenize, dispatch, delegate to OS
│   ├── engine.go
│   ├── exec_windows.go  # CREATE_NO_WINDOW for child processes
│   └── exec_other.go
├── command/             # Command interface + PATH lookup helpers
└── lexer/               # Hand-written scanner and token definitions
```

## Architecture at a Glance

1. `main` parses `-cli` and launches either `gui.Run()` or a bufio-based REPL.
2. Both modes funnel input into a single `shell.Engine.Execute(input)`.
3. The engine tokenizes via `lexer.Scanner`, then:
   - If the first token matches a registered builtin, it runs in-process.
   - Otherwise on Windows the **raw input** is handed to `powershell.exe -Command` (preserves PowerShell syntax), and on Unix the tokens are handed to `exec.Command`.
4. `exit` is implemented by panicking with a sentinel error that the engine recovers and propagates, which the GUI uses to close the window.

## Known Limitations

- `cd` with no argument uses `$HOME`, which is empty on Windows (use `$env:USERPROFILE` instead or pass a path explicitly).
- The builtin `cd` joins its arguments without spaces, so quoted paths containing spaces won't work through the builtin — they work through the OS fallback though.
- No piping or redirection between builtins yet; the lexer tokenizes `|`, `<`, `>` but the engine doesn't act on them.
- No command history or line editing in CLI mode.

## Roadmap

- Pipes and redirections handled by the engine itself.
- Command history and readline-style editing in CLI mode.
- Configurable prompt and theming.
- More builtins (`ls`, `cat`, `mkdir`, `touch`, ...).

## License

MIT — see source headers.
