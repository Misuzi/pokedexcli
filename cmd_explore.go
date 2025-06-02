package main

import (
	"encoding/json"
	"fmt"
	"internal/cache" // Adjust the import path as necessary
	"io"
	"net/http"
)

type Explore_response_body struct {
	Id                   int    `json:"count"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			}
		} `json:"version_details"`
	}
	Location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
	}
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int `json:"min_level"`
				MaxLevel        int `json:"max_level"`
				ConditionValues []struct {
				}
				Chance int `json:"chance"`
				Method struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				}
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

// Callback function to handle the explore command
func commandExplore(args []string, c *config, cache *cache.Cache) error {

	if len(args) == 0 {
		return fmt.Errorf("please provide a location area ID or name to explore")
	}

	fetch_URL := "https://pokeapi.co/api/v2/location-area/" + args[0]
	if fetch_URL == "" {
		return fmt.Errorf("please provide a URL to explore")
	}

	// Fetch the explore data from the API
	body_data, err := fetch_location_data(fetch_URL, cache)
	if err != nil {
		return fmt.Errorf("failed to fetch explore data: %v", err)
	}

	// Display the results
	fmt.Printf("Exploring %s...\n", body_data.Name)
	fmt.Printf("Found Pokemon:\n")
	for _, result := range body_data.PokemonEncounters {
		fmt.Printf(" - %s\n", result.Pokemon.Name)
	}

	return nil
}

func fetch_location_data(url string, cache *cache.Cache) (Explore_response_body, error) {

	if cachedData, found := cache.Get(url); found {
		var data Explore_response_body
		err := json.Unmarshal(cachedData, &data)
		if err != nil {
			return Explore_response_body{}, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
		return data, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Explore_response_body{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return Explore_response_body{}, fmt.Errorf("failed to fetch data, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Explore_response_body{}, err
	}

	var data Explore_response_body
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Explore_response_body{}, err
	}

	// Add the fetched data to the cache
	cache.Add(url, body)

	return data, nil
}
