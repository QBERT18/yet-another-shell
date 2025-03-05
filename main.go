package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"

	yetcommand "github.com/QBERT18/yet-another-shell/yetCommand"
)

func main() {

	commands := registerCommands()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print()
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		tokens := tokenize(strings.TrimSuffix(input, "\n"))
		if err = execute(tokens, commands); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

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
			path := "~"
			if len(input) > 1 {
				path = strings.Join(input[1:], " ")
			}
			if path == "~" {
				path = os.Getenv("HOME")
			}
			if err := os.Chdir(path); err != nil {
				fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", path)
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
			args := append([]string{"-Command"}, input...)
			cmd = exec.Command("powershell.exe", args...)
		} else {
			cmd = exec.Command(input[0], input[1:]...)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()

	}
	return nil
}

func tokenize(input string) []string {
	var tokens []string
	var token strings.Builder
	inQuotes := false
	quoteChar := rune(0) // Track whether we're in ' or "

	for i, ch := range input {
		switch {
		case ch == '"' || ch == '\'':
			// Handle quoted strings
			if inQuotes && ch == quoteChar {
				inQuotes = false // Closing quote
				tokens = append(tokens, token.String())
				token.Reset()
			} else if !inQuotes {
				inQuotes = true
				quoteChar = ch // Remember quote type
			} else {
				token.WriteRune(ch) // Inside quotes
			}

		case unicode.IsSpace(ch) && !inQuotes:
			// End current token on space (unless inside quotes)
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}

		default:
			token.WriteRune(ch)
		}

		// If it's the last character, flush the token
		if i == len(input)-1 && token.Len() > 0 {
			tokens = append(tokens, token.String())
		}
	}

	return tokens
}
