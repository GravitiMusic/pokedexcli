package main

import (
	"fmt"
	"os"
	"time"

	"pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
}

var commands map[string]cliCommand
var pokeClient = pokeapi.NewClient(60 * time.Second)

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name: "explore",
			description: "Displays the names of Pokemon that can be found in a given location area.",
			callback: commandExplore,
		},
	}
}

func commandExit(config *Config, input string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, input string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range commands {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *Config, input string) error {
	if config.Next == "" {
		config.Next = pokeapi.LocationAreaURL
	}

	result, err := pokeClient.GetLocationAreas(config.Next)
	if err != nil {
		return err
	}

	for _, area := range result.Results {
		fmt.Println(area.Name)
	}

	config.Next = result.Next
	config.Previous = result.Previous

	return nil
}

func commandMapb(config *Config, input string) error {
	if config.Previous == "" {
		fmt.Println("No previous page")
		return nil
	}

	result, err := pokeClient.GetLocationAreas(config.Previous)
	if err != nil {
		return err
	}

	for _, area := range result.Results {
		fmt.Println(area.Name)
	}

	config.Next = result.Next
	config.Previous = result.Previous

	return nil
}

func commandExplore(config *Config, location string) error {
	if location == "" {
		return fmt.Errorf("explore command requires a location area name")
	}

	result, err := pokeClient.GetEncounters(pokeapi.LocationAreaURL + location)
	if err != nil {
		return err
	}

	fmt.Println("Exploring " + location + "...\nFound Pokemon:")

	for _, name := range result {
		fmt.Printf("- %s\n", name)
	}

	return nil
}