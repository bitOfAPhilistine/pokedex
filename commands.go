package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {"help", "Display this help message", commandHelp},
		"exit": {"exit", "Exit the Pokedex", commandExit},
	}
}

func command(name string) error {
	cmd, ok := getCommands()[name]
	if ok {
		return cmd.callback()
	}
	return fmt.Errorf("Unknown command: %s", name)
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")
	for _, cmd := range getCommands() {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}