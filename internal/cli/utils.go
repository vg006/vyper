package cli

import (
	"strings"
)

func isExistingCommand(name string) bool {
	_, exists := c.commands[name]
	return exists
}

func isFlag(arg string) bool {
	if strings.HasPrefix(arg, "-") {
		return false
	}
	return true
}
