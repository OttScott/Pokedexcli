package main

import (
	"strings"
	"fmt"
	"os"
	"math/rand"
	"github.com/OttScott/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

func commandExit(cfg *Config, commands []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config, commands []string) error {
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

func commandMapb(cfg *Config, commands []string) error {
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

func commandExplore(cfg *Config, commands []string) error {
	if len(commands) < 1 {
		return fmt.Errorf("explore command requires a location area name as an argument")
	}
	
	locationAreaName := commands[0]
	fmt.Printf("Exploring location area: %s\n", locationAreaName)
	pokemons, err := pokeapi.GetLocationPokemons(cfg.cache, locationAreaName)
	if err != nil {
		return fmt.Errorf("failed to fetch Pokémon for location area '%s': %v", locationAreaName, err)
	}
	fmt.Printf("Found Pokémon in %s:\n", locationAreaName)
	for _, pokemon := range pokemons {
		fmt.Printf(" - %s\n", pokemon)
	}
	return nil
}

func commandCatch(cfg *Config, commands []string) error {
	if len(commands) < 1 {
		return fmt.Errorf("catch command requires a Pokémon name as an argument")
	}

	pokemonName := commands[0]
	
	// Check if Pokemon is already caught
	if _, exists := cfg.PokemonCaught[pokemonName]; exists {
		fmt.Printf("You already caught %s!\n", pokemonName)
		return nil
	}
	
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	pokeInfo, err := pokeapi.GetPokemonInfo(cfg.cache, pokemonName)
	if err != nil {
		return fmt.Errorf("failed to fetch Pokémon info for '%s': %v", pokemonName, err)
	}

	baseExperience := pokeInfo.BaseExperience
	if baseExperience > 100 {
		fmt.Printf("%s appears to be strong (base experience: %d)!\n", pokemonName, baseExperience)
	} else {
		fmt.Printf("%s appears to be manageable (base experience: %d).\n", pokemonName, baseExperience)
	}
	
	// Attempt to catch a Pokémon with a 50% success rate
	if rand.Intn(2) == 0 {
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
		cfg.PokemonCaught[pokemonName] = *pokeInfo
	} else {
		fmt.Printf("%s escaped the Pokeball!\n", pokemonName)
	}
	
	return nil
}

func commandPokedex(cfg *Config, commands []string) error {
	if len(cfg.PokemonCaught) == 0 {
		fmt.Println("You haven't caught any Pokémon yet!")
		return nil
	}
	
	fmt.Println("Your Pokedex:")
	for name := range cfg.PokemonCaught {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandInspect(cfg *Config, commands []string) error {
	if len(commands) < 1 {
		return fmt.Errorf("inspect command requires a Pokémon name as an argument")
	}

	pokemonName := commands[0]
	
	// Check if the Pokemon has been caught
	pokemon, exists := cfg.PokemonCaught[pokemonName]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	
	// Display Pokemon information
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokemonType.Type.Name)
	}
	
	return nil
}

var commands_map = map[string]cliCommand{
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
	"explore": {
		name:        "explore",
		description: "Explore a specific location area by name, listing the Pokémon that can be found there. Requires a location area name as an argument.",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Attempt to catch a specific Pokémon by name. Requires a Pokémon name as an argument.",
		callback:    commandCatch,
	},
	"pokedex": {
		name:        "pokedex",
		description: "View all the Pokémon you've caught.",
		callback:    commandPokedex,
	},
	"inspect": {
		name:        "inspect",
		description: "View detailed information about a caught Pokémon. Requires a Pokémon name as an argument.",
		callback:    commandInspect,
	},
}


func commandHelp(cfg *Config, commands []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println(" - help: Display a help message")
	for _, cmd := range commands_map {
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
