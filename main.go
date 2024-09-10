package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/GeminiZA/pokedex/internal/pokeAPI"
	pokecache "github.com/GeminiZA/pokedex/internal/pokeCache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
}

type config struct {
	NextMapURL string
	PrevMapURL string
}

var (
	commands map[string]cliCommand
	api      pokeAPI.API
)

func initVars() {
	api = pokeAPI.API{
		ApiCache: pokecache.NewCache(5 * time.Second),
	}
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
			description: "Display the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    commandMapB,
		},
    "explore": {
      name: "explote",
      description: "Display the ",
    }
	}
}

func commandMap(cfg *config) error {
	ret, err := api.Map(cfg.NextMapURL)
	if err != nil {
		return err
	}
	if ret.Next != nil {
		cfg.NextMapURL = *ret.Next
	}
	if ret.Prev != nil {
		cfg.PrevMapURL = *ret.Prev
	}
	for _, res := range ret.Results {
		fmt.Printf("%s\n", res.Name)
	}
	return nil
}

func commandMapB(cfg *config) error {
	ret, err := api.Map(cfg.PrevMapURL)
	if err != nil {
		return err
	}
	if ret.Next != nil {
		cfg.NextMapURL = *ret.Next
	}
	if ret.Prev != nil {
		cfg.PrevMapURL = *ret.Prev
	}
	for _, res := range ret.Results {
		fmt.Printf("%s\n", res.Name)
	}
	return nil
}

func commandExit(cfg *config) error {
	defer os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Print("Welcome to the Pokedex\nUsage:\n\n")
	for k, v := range commands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func main() {
	initVars()
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		NextMapURL: "",
		PrevMapURL: "",
	}
	for true {
		fmt.Print("\nPokedex > ")
		if scanner.Scan() {
			line := scanner.Text()
			if v, ok := commands[line]; ok {
				v.callback(cfg)
			} else {
				fmt.Println("Invalid command:", line)
				fmt.Println("Use 'help' for more information")
			}
		}
	}
}
