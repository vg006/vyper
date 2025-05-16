package cli

import (
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Usage       string
	Flags       []Flag
	SubCommands []Command
	Action      func(args []string) error
	Alias       string
}

type Flag struct {
	Name         string
	Usage        string
	Required     bool
	DefaultValue string
	Action       func(value string) error
}

type Cli struct {
	commands []Command
	flags    []Flag
}

func NewCli() *Cli {
	return &Cli{}
}

func (c *Cli) AddCommand(command Command) *Command {
	c.commands = append(c.commands, command)
	return &c.commands[len(c.commands)-1]
}

func (cmd *Command) AddSubCommand(subCommand Command) *Command {
	cmd.SubCommands = append(cmd.SubCommands, subCommand)
	return &cmd.SubCommands[len(cmd.SubCommands)-1]
}

func isNestedCommand(args []string, commands []Command) (*Command, []string) {
	if len(args) == 0 {
		return nil, nil
	}

	for _, command := range commands {
		if command.Name == args[0] || command.Alias == args[0] {
			if len(command.SubCommands) > 0 && len(args) > 1 {
				subCommand, remainingArgs := isNestedCommand(args[1:], command.SubCommands)
				if subCommand != nil {
					return subCommand, remainingArgs
				}
			}
			return &command, args[1:]
		}
	}

	return nil, nil
}

func (c *Cli) Run() error {
	args := os.Args[1:]

	if len(args) == 0 {
		return nil
	}

	command, remainingArgs := isNestedCommand(args, c.commands)
	if command == nil {
		return nil
	}

	if command.Action != nil {
		if err := command.Action(remainingArgs); err != nil {
			return err
		}
	}

	return nil
}

func (cli *Cli) PrintCommands() {
	fmt.Println("Commands:")
	for _, command := range cli.commands {
		fmt.Printf("  %s: %s\n", command.Name, command.Usage)
		if len(command.SubCommands) > 0 {
			command.PrintSubCommands()
		}
	}
}

func (cmd *Command) PrintSubCommands() {
	fmt.Printf("Subcommands of %s:\n", cmd.Name)
	for _, subCommand := range cmd.SubCommands {
		fmt.Printf("  %s: %s\n", subCommand.Name, subCommand.Usage)
		if len(subCommand.SubCommands) > 0 {
			subCommand.PrintSubCommands()
		}
	}
}
