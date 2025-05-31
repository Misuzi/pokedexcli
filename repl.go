package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}

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
		if cmd, exists := cliCommands[command]; exists {
			if err := cmd.callback(&repl_config); err != nil {
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

// Callback function to handle the exit command
func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

// Callback function to handle the help command
func commandHelp(c *config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	//fmt.Println("Available commands:")
	//for _, cmd := range cliCommands {
	//	fmt.Printf("- %s: %s\n", cmd.name, cmd.description)
	//}
	return nil
}

func commandMap(c *config) error {
	// This function is a placeholder for the map command.
	// It would typically display a map of the Pokedex or related information.
	//fmt.Println("Displaying the map...")

	// If map has been fetched before, use the stored URLs
	fetch_URL := "https://pokeapi.co/api/v2/location-area?limit=20"
	if c.Next_URL != "" {
		fetch_URL = c.Next_URL
	}

	// Fetch the map data from the API
	body_data, err := fetch_data(fetch_URL)
	if err != nil {
		return fmt.Errorf("failed to fetch map data: %v", err)
	}

	// Set the Next and Previous URLs in the config
	c.Next_URL = body_data.Next
	c.Previous_URL = body_data.Previous

	// Display the results
	for _, result := range body_data.Results {
		fmt.Printf("%s\n", result.Name)
	}

	return nil
}

func commandMapb(c *config) error {
	// This function is a placeholder for the map command.
	// It would typically display a map of the Pokedex or related information.
	//fmt.Println("Displaying the map...")

	// If map has been fetched before, use the stored URLs
	fetch_URL := c.Previous_URL
	if fetch_URL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	// Fetch the map data from the API
	body_data, err := fetch_data(fetch_URL)
	if err != nil {
		return fmt.Errorf("failed to fetch map data: %v", err)
	}

	// Set the Next and Previous URLs in the config
	c.Next_URL = body_data.Next
	c.Previous_URL = body_data.Previous

	// Display the results
	for _, result := range body_data.Results {
		fmt.Printf("%s\n", result.Name)
	}

	return nil
}

func fetch_data(url string) (Map_response_body, error) {
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

	return data, nil
}
