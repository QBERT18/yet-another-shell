package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/QBERT18/yet-another-shell/command"
	"github.com/QBERT18/yet-another-shell/lexer"
)

var ErrExit = errors.New("exit")

type Engine struct {
	Commands map[string]command.Command
	Stdout   io.Writer
	Stderr   io.Writer
}

func NewEngine(stdout, stderr io.Writer) *Engine {
	e := &Engine{
		Stdout: stdout,
		Stderr: stderr,
	}
	e.Commands = e.registerCommands()
	return e
}

func (e *Engine) Execute(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	var tokens []string
	scanner := lexer.NewScanner(strings.NewReader(input))
	for {
		tok, lit := scanner.Scan()
		if tok == lexer.EOF {
			break
		}
		if tok == lexer.WS {
			continue
		}
		tokens = append(tokens, lit)
	}

	if len(tokens) == 0 || tokens[0] == "" {
		return nil
	}

	name := tokens[0]
	if cmd, found := e.Commands[name]; found {
		err := e.runBuiltin(cmd, tokens)
		return err
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		fullCmd := "[Console]::OutputEncoding = [System.Text.Encoding]::UTF8; " + strings.Join(tokens, " ")
		cmd = exec.Command("powershell.exe", "-Command", fullCmd)
	} else {
		cmd = exec.Command(tokens[0], tokens[1:]...)
	}

	cmd.Stdout = e.Stdout
	cmd.Stderr = e.Stderr
	hideWindow(cmd)
	err := cmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return nil
	}
	return err
}

func (e *Engine) runBuiltin(cmd command.Command, tokens []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if r == ErrExit {
				err = ErrExit
			} else {
				panic(r) // re-panic for unexpected panics
			}
		}
	}()
	cmd.GetAction(tokens, e.Stdout, e.Stderr)
	return nil
}

func (e *Engine) registerCommands() map[string]command.Command {
	commands := make(map[string]command.Command)

	commands["echo"] = &command.BuiltinCommand{
		Name: "echo",
		ActionFunc: func(input []string, stdout io.Writer, stderr io.Writer) {
			fmt.Fprintln(stdout, strings.Join(input[1:], " "))
		},
		SubCommand: false,
	}

	commands["exit"] = &command.BuiltinCommand{
		Name: "exit",
		ActionFunc: func(input []string, stdout io.Writer, stderr io.Writer) {
			panic(ErrExit)
		},
		SubCommand: false,
	}

	commands["pwd"] = &command.BuiltinCommand{
		Name: "pwd",
		ActionFunc: func(input []string, stdout io.Writer, stderr io.Writer) {
			path, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(stderr, "pwd: error retrieving current directory")
				return
			}
			fmt.Fprintln(stdout, path)
		},
		SubCommand: false,
	}

	commands["cd"] = &command.BuiltinCommand{
		Name: "cd",
		ActionFunc: func(input []string, stdout io.Writer, stderr io.Writer) {
			var path string

			if len(input) == 1 {
				path = os.Getenv("HOME")
			} else {
				path = strings.Join(input[1:], "")
			}

			if path == "~" {
				path = os.Getenv("HOME")
			}

			if info, err := os.Stat(path); err != nil || !info.IsDir() {
				fmt.Fprintf(stderr, "cd: %s: No such file or directory\n", path)
				return
			}

			if err := os.Chdir(path); err != nil {
				fmt.Fprintf(stderr, "cd: %s: %v\n", path, err)
			}
		},
		SubCommand: false,
	}

	commands["type"] = &command.BuiltinCommand{
		Name: "type",
		ActionFunc: func(input []string, stdout io.Writer, stderr io.Writer) {
			for _, subCmd := range input[1:] {
				if cmd, found := commands[subCmd]; found {
					fmt.Fprintf(stdout, "%s is a %s command\n", subCmd, cmd.GetType())
				} else {
					foundPath := command.SearchProgramInPath(subCmd)
					if foundPath != "" {
						fmt.Fprintf(stdout, "%s is %s\n", subCmd, foundPath)
					} else {
						fmt.Fprintf(stdout, "%s: not found\n", subCmd)
					}
				}
			}
		},
		SubCommand: true,
	}

	return commands
}
