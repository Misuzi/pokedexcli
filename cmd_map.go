package main

import (
	"fmt"
	"internal/cache" // Adjust the import path as necessary
)

func commandMap(args []string, c *config, cache *cache.Cache) error {
	// This function is a placeholder for the map command.
	// It would typically display a map of the Pokedex or related information.
	//fmt.Println("Displaying the map...")

	// If map has been fetched before, use the stored URLs
	fetch_URL := "https://pokeapi.co/api/v2/location-area?limit=20"
	if c.Next_URL != "" {
		fetch_URL = c.Next_URL
	}

	// Fetch the map data from the API// Fetch the map data from the API
	body_data, err := fetch_data(fetch_URL, cache)
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

func commandMapb(args []string, c *config, cache *cache.Cache) error {
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
	body_data, err := fetch_data(fetch_URL, cache)
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
