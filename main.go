package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	yetcommand "github.com/QBERT18/yet-another-shell/yetCommand"
)

var paths []string

func init() {
	paths = strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	commands := registerCommands()

	for {
		fmt.Print("$ ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			continue
		}

		inputParts := yetcommand.CustomSplit(strings.TrimSpace(input))

		fmt.Println("inputParts: ", inputParts)
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

	commands["type"] = &yetcommand.BuiltinCommand{
		Name: "type",
		ActionFunc: func(input []string) {
			for _, SubCommand := range input[1:] {
				if cmd, found := commands[SubCommand]; found {
					fmt.Printf("%s is a %s command\n", SubCommand, cmd.GetType())
				} else {
					foundPath := yetcommand.SearchProgramInPath(SubCommand, paths)
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
			fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
		}
	}
}
