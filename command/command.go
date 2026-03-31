package command

type Command interface {
	GetName() string
	GetAction(input []string)
	GetType() string
	IsSubCommand() bool
}

type BuiltinCommand struct {
	Name       string
	SubCommand bool
	ActionFunc func(input []string)
}

func (b *BuiltinCommand) GetName() string {
	return b.Name
}

func (b *BuiltinCommand) GetAction(input []string) {
	b.ActionFunc(input)
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
	ActionFunc func(input []string)
}

func (c *CustomCommand) GetName() string {
	return c.Name
}

func (c *CustomCommand) GetAction(input []string) {
	c.ActionFunc(input)
}

func (c *CustomCommand) GetType() string {
	return "custom"
}

func (c *CustomCommand) IsSubCommand() bool {
	return c.SubCommand
}
