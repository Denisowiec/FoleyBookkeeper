package main

// A simple REPL client for the purpose of testing the server

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	username      string
	email         string
	jwt           string
	serverAddress string
}

func main() {

}
