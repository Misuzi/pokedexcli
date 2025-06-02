package main

import (
	"fmt"
	"internal/cache" // Adjust the import path as necessary
)

// Callback function to handle the help command
func commandHelp(args []string, c *config, cache *cache.Cache) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	//fmt.Println("Available commands:")
	//for _, cmd := range cliCommands {
	//	fmt.Printf("- %s: %s\n", cmd.name, cmd.description)
	//}
	return nil
}
