package main

import (
	"fmt"
	"internal/cache" // Adjust the import path as necessary
	"os"
)

// Callback function to handle the exit command
func commandExit(args []string, c *config, cache *cache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}
