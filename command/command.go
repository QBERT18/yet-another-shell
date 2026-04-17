package command

import "io"

type Command interface {
	GetName() string
	GetAction(input []string, stdout io.Writer, stderr io.Writer)
	GetType() string
	IsSubCommand() bool
}

type BuiltinCommand struct {
	Name       string
	SubCommand bool
	ActionFunc func(input []string, stdout io.Writer, stderr io.Writer)
}

func (b *BuiltinCommand) GetName() string {
	return b.Name
}

func (b *BuiltinCommand) GetAction(input []string, stdout io.Writer, stderr io.Writer) {
	b.ActionFunc(input, stdout, stderr)
}

func (b *BuiltinCommand) GetType() string {
	return "builtin"
}

func (b *BuiltinCommand) IsSubCommand() bool {
	return b.SubCommand
}

type CustomCommand struct {
	Name       string
	SubCommand bool
	ActionFunc func(input []string, stdout io.Writer, stderr io.Writer)
}

func (c *CustomCommand) GetName() string {
	return c.Name
}

func (c *CustomCommand) GetAction(input []string, stdout io.Writer, stderr io.Writer) {
	c.ActionFunc(input, stdout, stderr)
}

func (c *CustomCommand) GetType() string {
	return "custom"
}

func (c *CustomCommand) IsSubCommand() bool {
	return c.SubCommand
}
