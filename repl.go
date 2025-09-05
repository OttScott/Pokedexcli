package main

import (
	"strings"
	"fmt"
	"os"
	"github.com/OttScott/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config) error {
    NextLocURL, prevLocURL, locs, err := pokeapi.GetLocationAreas(cfg.cache, cfg.NextLocationURL)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %v", err)
	}

	for _, loc := range locs.Results {
		fmt.Printf("%s\n", loc.Name)
	}
	cfg.NextLocationURL = NextLocURL
	cfg.PreviousLocationURL = prevLocURL

	return nil
}

func commandMapb(cfg *Config) error {
	if cfg.PreviousLocationURL == nil {
		fmt.Println("you're on the first page")
		return nil
	} else {
		// Fetch the previous page of locations using cfg.PreviousLocationURL
		// Update cfg.NextLocationURL and cfg.PreviousLocationURL accordingly
		NextURL, prevURL, locs, err := pokeapi.GetLocationAreas(cfg.cache, cfg.PreviousLocationURL)
		if err != nil {
			return fmt.Errorf("failed to fetch previous location areas: %v", err)
		}

		for _, loc := range locs.Results {
			fmt.Printf("%s\n", loc.Name)
		}
		cfg.NextLocationURL = NextURL
		cfg.PreviousLocationURL = prevURL
	}
	return nil
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"map": {
		name:        "map",
		description: "Display a list of all Pokémon locations in the Pokedex. (grouped in batches of 20)",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Display the previous page of all Pokémon locations in the Pokedex. (grouped in batches of 20)",
		callback:    commandMapb,
	},
}


func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println(" - help: Display a help message")
	for _, cmd := range commands {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func cleanInput(input string) []string {
	cleaned := strings.TrimSpace(input)
	if cleaned == "" {
		return []string{}
	}
	cleaned = strings.ToLower(cleaned)
	cleaned = strings.ReplaceAll(cleaned, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\t", " ")
	return strings.Fields(cleaned)
}
