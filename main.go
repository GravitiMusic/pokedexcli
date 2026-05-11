package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var config Config

	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)

		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		command, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		commandInput := ""
		if len(words) > 1 {
			commandInput = strings.Join(words[1:], " ")
		}

		err := command.callback(&config, commandInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
