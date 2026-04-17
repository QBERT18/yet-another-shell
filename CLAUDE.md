# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Go-based shell (CodeCrafters "Build Your Own Shell" challenge). Supports Unix-style and PowerShell-style commands. Ships as a Fyne-based GUI by default with an optional CLI mode. Cross-platform (Windows + Unix), but execution of non-builtin commands is delegated differently per OS (see Architecture).

## Common Commands

```bash
# Run in GUI mode (default)
go run main.go

# Run in CLI mode (REPL on stdin/stdout)
go run main.go -cli

# Build standalone Windows GUI executable (no console window)
go build -ldflags "-H windowsgui" -o yet-another-shell.exe .

# Plain build
go build -o yet-another-shell .

# Dependencies
go mod tidy
```

No test suite exists yet.

## Architecture

Four packages, one-way dependency graph: `main` → `gui`/`shell` → `command`/`lexer`.

**Execution flow** — all input (GUI or CLI) funnels through a single `shell.Engine.Execute(input)`. The engine:

1. Tokenizes input via `lexer.Scanner` (keeps only non-whitespace tokens).
2. If `tokens[0]` matches a registered builtin in `Engine.Commands`, runs it.
3. Otherwise delegates to the OS:
   - **Windows**: wraps the **original raw input string** (not tokens) into `powershell.exe -Command "[Console]::OutputEncoding = UTF8; <input>"`. This is intentional — tokenizing would destroy PowerShell syntax (see commit `484dff3`). Do not "fix" this by passing tokens.
   - **Unix**: `exec.Command(tokens[0], tokens[1:]...)` — a tokenized exec.
4. Window-hiding for child processes is OS-specific: `exec_windows.go` sets `CREATE_NO_WINDOW`, `exec_other.go` is a no-op stub. Use the build-tag pattern when adding OS-specific Cmd handling.

**Builtin commands** (`echo`, `cd`, `pwd`, `type`, `exit`) live as `BuiltinCommand` structs inside `Engine.registerCommands()`. `exit` signals termination by **panicking with `ErrExit`**, which `runBuiltin` recovers into a returned error. Preserve this contract if adding new builtins that must terminate the shell.

**GUI (`gui/gui.go`)** — Fyne app. An `outputWriter` (mutex-protected `strings.Builder` backing a `widget.Label`) is passed as both stdout and stderr to `NewEngine`. Each submit runs `engine.Execute` in a goroutine so the UI stays responsive; the input is disabled during execution and re-enabled + refocused after. `ErrExit` closes the window.

**Lexer (`lexer/`)** — Hand-written scanner in the style of Gopher Academy's parser/lexer guide. Defines tokens for shell operators (`|`, `&`, `;`, `<`, `>`, `$`, backtick, `-`, quotes, `*`, `,`) and keyword tokens for both Unix (`ECHO`, `CD`, `LS`, …) and PowerShell equivalents (`WRITE-OUTPUT`, `SET-LOCATION`, `GET-CHILDITEM`, …). Keyword matching is case-insensitive via `strings.ToUpper`. Note: the engine currently only uses the `lit` string (not the `Token` type) when tokenizing — builtin dispatch is by literal name, not token kind.

**Command interface (`command/command.go`)** — `Command` interface with `BuiltinCommand` and `CustomCommand` implementations (same shape, different `GetType()`). `command.SearchProgramInPath` walks `$PATH` and checks the executable bit (`mode & 0111`) — the exec-bit check is a no-op on Windows but harmless.

## Gotchas

- On Windows the non-builtin path runs the **raw input** through PowerShell, so builtin shadowing matters: if you add a builtin named the same as a common PowerShell cmdlet, the builtin wins.
- `cd` uses `os.Getenv("HOME")` for both no-arg and `~`. This is broken on Windows (Windows uses `USERPROFILE`). If touching `cd`, fix both.
- `cd` joins args with empty separator (`strings.Join(input[1:], "")`) — paths containing spaces already lost their spaces by tokenization time, so quoted paths won't work through the builtin.
- Builtins run on the calling goroutine; any `panic` other than `ErrExit` is re-panicked out of `runBuiltin`.
