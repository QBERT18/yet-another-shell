package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/QBERT18/yet-another-shell/gui"
	"github.com/QBERT18/yet-another-shell/shell"
)

func main() {
	cliMode := flag.Bool("cli", false, "launch CLI mode")
	flag.Parse()

	if *cliMode {
		runCLI()
	} else {
		gui.Run()
	}
}

func runCLI() {
	engine := shell.NewEngine(os.Stdout, os.Stderr)
	reader := bufio.NewReader(os.Stdin)

	for {
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

		err = engine.Execute(input)
		if err != nil {
			if errors.Is(err, shell.ErrExit) {
				os.Exit(0)
			}
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
