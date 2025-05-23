package cli

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	c = &Cli{
		root:     &Command{},
		commands: make(map[string]*Command),
		args:     os.Args[1:],
	}
}

var (
	c *Cli
)

type (
	Commands map[string]*Command
)

type Cli struct {
	root     *Command
	args     []string
	commands Commands
}

type Command struct {
	name        string
	Usage       string
	Flags       []Flag
	SubCommands Commands
	Action      func() error
	Alias       string
}

type Flag struct {
	Name         string
	Variable     any
	Usage        string
	Required     bool
	DefaultValue string
	Action       func() error
}

func Init(root Command, commands Commands) {
	c.root = &root
	for name, command := range commands {
		addCommand(name, command)
	}
}

func addCommand(name string, command *Command) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic:", r)
		}
	}()

	if isExistingCommand(command.name) {
		panic(fmt.Sprintf("command %s already exists", command.name))
	}
	c.commands[name] = command
	c.commands[name].name = name
	if isExistingCommand(command.Alias) {
		panic(fmt.Sprintf("alias %s already exists", command.Alias))
	}
	c.commands[command.Alias] = command
}

func AddCommand(command Command) {
	addCommand(command.name, &command)
}

func addCommands(commands Commands) {
	for name, command := range commands {
		addCommand(name, command)
	}
}

func AddCommands(cmd Commands) {
	addCommands(cmd)
}

func getCommand() *Command {
	fs := parseFlags(c.root)
	if err := fs.Parse(c.args); err != nil {
		panic(err)
	}
	if len(fs.Args()) == 0 {
		return c.root
	}
	return recurseCommand(fs.Args(), c.commands)
}

func recurseCommand(args []string, commands map[string]*Command) *Command {
	for cmdName, cmd := range commands {
		if cmdName == args[0] || cmd.Alias == args[0] {
			if isFlag(args[0]) {
				fs := parseFlags(cmd)
				if err := fs.Parse(args[1:]); err != nil {
					panic(err)
				}
				if len(fs.Args()) != 0 {
					return recurseCommand(fs.Args(), cmd.SubCommands)
				}
			}
			return cmd
		}
	}
	return nil
}

func parseFlags(cmd *Command) *flag.FlagSet {
	fs := flag.NewFlagSet(cmd.name, flag.ExitOnError)
	for _, flag := range cmd.Flags {
		switch t := flag.Variable.(type) {
		case *string:
			fs.StringVar(t, flag.Name, flag.DefaultValue, flag.Usage)
		case *int:
			fs.IntVar(t, flag.Name, 0, flag.Usage)
		}
	}
	return fs
}

func Run() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic:", r)
		}
	}()

	command := getCommand()
	if command == nil {
		return fmt.Errorf("command not found")
	}

	if command.Action == nil {
		return fmt.Errorf("no action defined for command %s", command.name)
	}

	if err := command.Action(); err != nil {
		return err
	}
	return nil
}

func PrintCommands() {
	for name, cmd := range c.commands {
		fmt.Printf("Command: %s\n", name)
		fmt.Printf("Alias: %s\n", cmd.Alias)
		fmt.Printf("Usage: %s\n", cmd.Usage)
		fmt.Println("Flags:")
		for _, flag := range cmd.Flags {
			fmt.Printf("  -%s: %s (default: %s)\n", flag.Name, flag.Usage, flag.DefaultValue)
		}
	}
}
