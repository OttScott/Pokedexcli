package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	config := &Config{
		NextLocationURL:     nil,  // No next URL yet
		PreviousLocationURL: nil,  // No previous URL yet
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
			commandHelp(config)
			continue
		}
		cmd, exists := commands[command[0]]
		if !exists {
			fmt.Println("Unknown command.")
			continue
		}
		if err := cmd.callback(config); err != nil {
			fmt.Printf("Error executing command '%s': %v\n", command, err)
		}

	}
}
