package main

import "fmt"

func commandHelp(cfg *config, args []string) error {
	if len(args) == 0 {
		for _, cmd := range cfg.commands {
			fmt.Printf("%s - %s\n", cmd.name, cmd.description)
		}
	} else {
		cmd, ok := cfg.commands[args[0]]
		if !ok {
			return fmt.Errorf("command %s unknown", args[0])
		}
		fmt.Println(cmd.usage)
		fmt.Println(cmd.description)
	}
	return nil
}
