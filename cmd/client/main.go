package main

import (
	"bufio"
	"fmt"
	"os"
)

// A simple REPL client for the purpose of testing the server

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// Config struct containing some persistent state information
	cfg := config{
		username:      "",
		email:         "",
		jwt:           "",
		serverAddress: "http://localhost:8080",
		commands:      listCommands(),
	}

	for {
		// Command prompt
		fmt.Print("> ")
		scanner.Scan()
		cleanedText := cleanInput(scanner.Text())
		cmdWord := cleanedText[0]
		cmdArgs := cleanedText[1:]
		cmd, ok := cfg.commands[cmdWord]
		if !ok {
			fmt.Println("Command unknown")
			continue
		}
		err := cmd.callback(&cfg, cmdArgs)
		if err != nil {
			fmt.Println("Unable to process command:", err)
			fmt.Println("Usage:", cmd.usage)
		}
	}
}
