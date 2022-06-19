package main

import (
	"fmt"
	"os"
	"tunme/pkg/cmd/tunme_cat"
	"tunme/pkg/cmd/tunme_relay"
)

var subCommands = makeSubCommands()

func makeSubCommands() map[string]func(string, []string) {

	subCommands := make(map[string]func(string, []string))

	subCommands["cat"] = tunme_cat.Main
	subCommands["relay"] = tunme_relay.Main

	return subCommands
}

func exitBadUsage(msg string) {

	_, _ = fmt.Fprintf(os.Stderr, "%s Available subcommands are:\n", msg)
	for subCommand := range subCommands {
		_, _ = fmt.Fprintf(os.Stderr, " *  %s\n", subCommand)
	}

	os.Exit(1)
}

func main() {

	if len(os.Args) < 2 {
		exitBadUsage("Missing subcommand.")
	}
	subCommandStr := os.Args[1]

	subCommand, ok := subCommands[subCommandStr]
	if !ok {
		exitBadUsage("Unknown subcommand.")
	}

	subCommand(fmt.Sprintf("%s %s", os.Args[0], subCommandStr), os.Args[2:])
}
