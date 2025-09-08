package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"github.com/OttScott/pokedexcli/internal/pokecache"
	"github.com/OttScott/pokedexcli/internal/pokeapi"
)

type Config struct {
	cache			    *pokecache.Cache
	NextLocationURL     *string
	PreviousLocationURL *string
	PokemonCaught       map[string]pokeapi.PokemonInfo
}

func main() {
	cache := pokecache.NewCache(time.Second * 30)

	config := &Config{
		NextLocationURL:     nil,  // No next URL yet
		PreviousLocationURL: nil,  // No previous URL yet
		cache:               cache,
		PokemonCaught:       make(map[string]pokeapi.PokemonInfo),
	}

	// Basic REPL loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("POKEDEX > ")
		scanner.Scan()
		input := scanner.Text()
		command := cleanInput(input)
		if len(command) == 0 {
			continue
		}
		if command[0] == "help" {
			commandHelp(config, command[1:])
			continue
		}
		cmd, exists := commands_map[command[0]]
		if !exists {
			fmt.Println("Unknown command.")
			continue
		}
		if err := cmd.callback(config, command[1:]); err != nil {
			fmt.Printf("Error executing command '%s': %v\n", command, err)
		}

	}
}
