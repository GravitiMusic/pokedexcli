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
var dex = pokeapi.Pokedex{Caught: make(map[string]pokeapi.Pokemon)}
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
		"catch": {
			name: "catch",
			description: "Attempts to catch a Pokemon by name.",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "Displays detailed information about a caught Pokemon.",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "Lists all caught Pokemon.",
			callback: commandPokedex,
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

func commandCatch(config *Config, target string) error {
	if target == "" {
		return fmt.Errorf("catch command requires a Pokemon name")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", target)

	success, err := pokeClient.CatchPokemon(pokeapi.PokemonURL+target, target, &dex)
	if err != nil {
		return err
	}
	
	if success > 0 {
		fmt.Printf("%s was caught!\n", target)
	} else {
		fmt.Printf("%s escaped!\n", target)
	}
	return nil
}

func commandInspect(config *Config, target string) error {
	if target == "" {
		return fmt.Errorf("inspect command requires a Pokemon name")
	}

	pokemon, ok := dex.Caught[target]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(config *Config, input string) error {
	for _, pokemon := range dex.Caught {
		fmt.Printf("- %s\n", pokemon.Name)
	}
	return nil
}