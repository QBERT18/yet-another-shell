package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	yetcommand "github.com/QBERT18/yet-another-shell/yetCommand"
	yetshellwords "github.com/QBERT18/yet-another-shell/yetShellWords"
)

func main() {

	p := yetshellwords.NewParser()
	p.ParseBacktick = true

	reader := bufio.NewReader(os.Stdin)

	commands := registerCommands()

	for {
		fmt.Println()
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			continue
		}

		inputParts, err := p.Parse(strings.TrimSpace(input))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing input: ", err)
		}

		handleCommand(inputParts, commands)
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

func handleCommand(input []string, commands map[string]yetcommand.Command) {
	if len(input) == 0 || input[0] == "" {
		return
	}

	command := input[0]
	if cmd, found := commands[command]; found {
		cmd.GetAction(input)
	} else {
		cmd := exec.Command(command, input[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()

		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				fmt.Fprintf(os.Stderr, "%s: command failed with exit code %d\n", command, exitError.ExitCode())
			} else if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
			} else {
				fmt.Fprintf(os.Stderr, "%s: execution error: %v\n", command, err)
			}
		}
	}
}
