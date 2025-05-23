package cli

import (
	"flag"
	"strings"
)

func isRoot() bool {
	if len(c.args) == 0 {
		return true
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		return true
	}
	return false
}

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

func printArgs(args []string) {
	for _, arg := range args {
		println(arg)
	}
}
