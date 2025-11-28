package main

import (
	"strings"

	"github.com/google/uuid"
)

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
	userID        uuid.UUID
	serverAddress string
	commands      map[string]cliCommand
}

func cleanInput(text string) []string {
	// This takes user's input, divides it by word into a slice and
	// lowercases the first item, which is the command
	fields := strings.Fields(text)
	fields[0] = strings.ToLower(fields[0])

	if len(fields) == 1 {
		return fields
	}

	// We want to allow users to enclose multi-word inputs in parantheses
	new_fields := []string{}
	new_fields = append(new_fields, fields[0])
	multi_field := ""
	var sign byte
	for _, item := range fields[1:] {
		// First, if multi_field isn't empty, it means we're looking for the end of it
		if multi_field != "" {
			if item[len(item)-1] == sign {
				// The end of the item is found, we can empty the multi_field and move on
				item = strings.Trim(item, string(sign))
				multi_field = multi_field + item

				new_fields = append(new_fields, multi_field)
				multi_field = ""
				continue
			}
			multi_field = multi_field + item + " "
			continue
		}

		// If there are no parantheses etc we simply transfer the field into the new slice
		if item[0] != '"' && item[0] != '\'' {
			new_fields = append(new_fields, item)
			continue
		} else {
			sign = item[0]
		}

		// If there's a paranthese in the beginning, we check if there's also
		// one at the end of the current field. If so, we remove both and append the item
		// to the new slice
		if item[len(item)-1] == sign {
			item = strings.Trim(item, string(sign))
			new_fields = append(new_fields, item)
		} else {
			item = strings.Trim(item, string(sign))
			// If there's not a paranthese at the end, we prepare a string to
			// concatenate the next field onto
			multi_field = item + " "
		}
	}
	// If at the end of the loop multi_field isn't empty, we append it now to the slice
	if multi_field != "" {
		new_fields = append(new_fields, strings.TrimSpace(multi_field))
	}

	return new_fields
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
		"update-user": {
			name:        "update-user",
			description: "Updates user's username and e-mail address.",
			usage:       "update-user <username> <email>",
			callback:    commandUpdateUser,
		},
		"update-password": {
			name:        "update-password",
			description: "Updates the user's password.",
			usage:       "update-password <password>",
			callback:    commandUpdatePassword,
		},
		"get-user": {
			name:        "get-user",
			description: "Display all of the user's info.",
			usage:       "get-user <ID>",
			callback:    commandGetUserInfo,
		},
		"create-client": {
			name:        "create-client",
			description: "Create a new client.",
			usage:       "create-client <name> <email> <notes>",
			callback:    commandCreateClient,
		},
		"update-client": {
			name:        "update-client",
			description: "Update a client.",
			usage:       "update-client <old_name> <new_name> <email> <notes>",
			callback:    commandUpdateClient,
		},
		"show-client": {
			name:        "show-client",
			description: "Display basic info about a client",
			usage:       "show-client <client-name>",
			callback:    commandGetClient,
		},
		"list-clients": {
			name:        "list-clients",
			description: "Lists all clients",
			usage:       "list-clients",
			callback:    commandGetAllClients,
		},
		"create-project": {
			name:        "create-project",
			description: "Create a new project",
			usage:       "create-project <title> <client>",
			callback:    commandCreateProject,
		},
		"update-project": {
			name:        "update-project",
			description: "Update a project",
			usage:       "update-project <old_title> <new_title> <client>",
			callback:    commandUpdateProject,
		},
		"show-project": {
			name:        "show-project",
			description: "Display basic info about a project",
			usage:       "show-project <title>",
			callback:    commandGetProjectInfo,
		},
		"list-projects": {
			name:        "list-projects",
			description: "List all projects",
			usage:       "list-projects",
			callback:    commandGetAllProjects,
		},
		"create-episode": {
			name:        "create-episode",
			description: "Creates an episode",
			usage:       "create-episode <project title> <episode number> <episode title>",
			callback:    commandCreateEpisode,
		},
		"get-project-eps": {
			name:        "get-project-eps",
			description: "Returns all episodes for a given project",
			usage:       "get-project-eps <project title>",
			callback:    commandGetEpisodesForProject,
		},
		"create-session": {
			name:        "create-session",
			description: "Creates a new session",
			usage:       "create-session <project title> <episode number> <date> <duration> <part worked on> <activity done> <user1> <user2> etc...",
			callback:    commandCreateSession,
		},
		"get-sessions": {
			name:        "get-sessions",
			description: "Lists some sessions",
			usage:       "get-sessions <how many> <project title> <episode number>",
			callback:    commandGetSessions,
		},
	}
}
