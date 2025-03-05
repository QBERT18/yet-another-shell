package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/QBERT18/yet-another-shell/lexer"
	yetcommand "github.com/QBERT18/yet-another-shell/yetCommand"
)

var tokens []string

func main() {
	commands := registerCommands()

	reader := bufio.NewReader(os.Stdin)
	for {
		tokens = []string{}

		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		scanner := lexer.NewScanner(strings.NewReader(input))
		for {
			tok, lit := scanner.Scan()
			if tok == lexer.EOF {
				break
			}

			// fmt.Printf("Token %s: %s, len:%d\n", tok.String(), lit, len(lit)) // Debugging

			tokens = append(tokens, lit)
		}

		// fmt.Println("Tokens array:", tokens) // Debugging

		err = execute(tokens, commands)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execute(input []string, commands map[string]yetcommand.Command) error {
	if len(input) == 0 || input[0] == "" {
		return nil
	}

	command := input[0]
	if cmd, found := commands[command]; found {
		cmd.GetAction(input)
	} else {
		var cmd *exec.Cmd

		if runtime.GOOS == "windows" {
			// Run commands inside PowerShell
			args := append([]string{"-Command"}, strings.Join(input, ""))
			cmd = exec.Command("powershell.exe", args...)
			fmt.Println("cmd:", cmd)
		} else {
			// Unix-based systems (Linux, macOS)
			cmd = exec.Command(input[0], strings.Join(input[1:], ""))
			fmt.Println("cmd:", cmd)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()

	}
	return nil
}

func registerCommands() map[string]yetcommand.Command {
	commands := make(map[string]yetcommand.Command)

	commands["echo"] = &yetcommand.BuiltinCommand{
		Name: "echo",
		ActionFunc: func(input []string) {
			fmt.Println(strings.Join(input[1:], " "))
		},
		SubCommand: false,
	}

	commands["exit"] = &yetcommand.BuiltinCommand{
		Name: "exit",
		ActionFunc: func(input []string) {
			os.Exit(0)
		},
		SubCommand: false,
	}

	commands["pwd"] = &yetcommand.BuiltinCommand{
		Name: "pwd",
		ActionFunc: func(input []string) {
			path, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "pwd: error retrieving current directory")
				return
			}
			fmt.Println(path)
		},
		SubCommand: false,
	}

	commands["cd"] = &yetcommand.BuiltinCommand{
		Name: "cd",
		ActionFunc: func(input []string) {
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
				fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", path)
				return
			}

			if err := os.Chdir(path); err != nil {
				fmt.Fprintf(os.Stderr, "cd: %s: %v\n", path, err)
			}
		},
		SubCommand: false,
	}

	commands["type"] = &yetcommand.BuiltinCommand{
		Name: "type",
		ActionFunc: func(input []string) {
			for _, SubCommand := range input[1:] {
				if cmd, found := commands[SubCommand]; found {
					fmt.Printf("%s is a %s command\n", SubCommand, cmd.GetType())
				} else {
					foundPath := yetcommand.SearchProgramInPath(SubCommand)
					if foundPath != "" {
						fmt.Printf("%s is %s\n", SubCommand, foundPath)
					} else {
						fmt.Printf("%s: not found\n", SubCommand)
					}
				}
			}
		},
		SubCommand: true,
	}

	return commands
}
