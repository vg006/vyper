package cli

import (
	"flag"
	"fmt"
	"os"
)

var c Cli

type Command struct {
	Name        string
	Usage       string
	Flags       []Flag
	SubCommands []Command
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

type Cli struct {
	commands []Command
}

func AddCommand(command Command) {
	c.commands = append(c.commands, command)
	return
}

func (cmd *Command) AddSubCommand(subCommand Command) {
	cmd.SubCommands = append(cmd.SubCommands, subCommand)
	return
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

func Run() error {
	args := os.Args[1:]
	var err error

	command, remainingArgs := isNestedCommand(args, c.commands)
	if command == nil {
		return nil
	}

	err = validateArgs(args, command)
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet(command.Name, flag.ContinueOnError)
	for _, f := range command.Flags {
		switch v := f.Variable.(type) {
		case *string:
			fs.StringVar(v, f.Name, f.DefaultValue, f.Usage)
		}
	}

	if err := fs.Parse(remainingArgs); err != nil {
		return err
	}

	if command.Action != nil {
		if err := command.Action(); err != nil {
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
