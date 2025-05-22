package cli

import "fmt"

func (c *Cli) validateCommand() error {

	return nil
}

func validateArgs(args []string, command *Command) error {
	if len(args) == 0 {
		return nil
	}

	for _, flag := range command.Flags {
		if flag.Required && flag.Variable == nil {
			return fmt.Errorf("flag %s is required", flag.Name)
		}
	}

	return nil
}
