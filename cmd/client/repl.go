package main

import "strings"

type cliCommand struct {
	name        string
	description string
	usage       string
	callback    func(*config, []string) error
}

type config struct {
	username      string
	email         string
	jwt           string
	serverAddress string
	commands      map[string]cliCommand
}

func cleanInput(text string) []string {
	// This takes user's input, divides it by word into a slice and
	// lowercases the first item, which is the command
	fields := strings.Fields(text)
	fields[0] = strings.ToLower(fields[0])

	return fields
}

func listCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"register": {
			name:        "register",
			description: "Create a new user in the system.",
			usage:       "register <username> <email> <password>",
			callback:    commandCreateUser,
		},
		"login": {
			name:        "login",
			description: "Log in to the system.",
			usage:       "login <email> <password>",
			callback:    commandLogin,
		},
		"help": {
			name:        "help",
			description: "Lists all available commands or usage information for a given command.",
			usage:       "help <command>",
			callback:    commandHelp,
		},
	}
}
