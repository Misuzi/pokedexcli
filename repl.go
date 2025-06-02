package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"internal/cache" // Adjust the import path as necessary
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string, *config, *cache.Cache) error
}

type config struct {
	Next_URL     string
	Previous_URL string
}

// Map_response_body represents the structure of the response body for the map command.
type Map_response_body struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func repl() {
	// This function is the main REPL loop for the Pokedex CLI.
	// It reads user input, processes commands, and displays results.
	// The implementation details are omitted for brevity.

	repl_config := config{
		Next_URL:     "",
		Previous_URL: "",
	}

	cliCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display this help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the map of the Pokedex",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous map of the Pokedex",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore the Pokedex",
			callback:    commandExplore,
		},
	}

	cache_ptr := cache.NewCache(5 * time.Second) // 5 seconds in nanoseconds

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		input_slice := cleanInput(input)
		if len(input_slice) == 0 {
			fmt.Println("No input provided. Please enter a command.")
			continue
		}

		command := input_slice[0]
		args := input_slice[1:]
		if cmd, exists := cliCommands[command]; exists {
			if err := cmd.callback(args, &repl_config, cache_ptr); err != nil {
				fmt.Printf("Error executing command '%s': %v\n", command, err)
			}
		} else {
			fmt.Printf("Unknown command\n")
		}

	}
}

func cleanInput(text string) []string {
	output_strings := strings.Fields(text)
	for i, word := range output_strings {
		output_strings[i] = strings.ToLower(word)
	}
	return output_strings
}

func fetch_data(url string, cache *cache.Cache) (Map_response_body, error) {

	if cachedData, found := cache.Get(url); found {
		var data Map_response_body
		err := json.Unmarshal(cachedData, &data)
		if err != nil {
			return Map_response_body{}, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
		return data, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Map_response_body{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return Map_response_body{}, fmt.Errorf("failed to fetch data, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Map_response_body{}, err
	}

	var data Map_response_body
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Map_response_body{}, err
	}

	// Add the fetched data to the cache
	cache.Add(url, body)

	return data, nil
}
