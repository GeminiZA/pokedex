package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/GeminiZA/pokedex/internal/pokeAPI"
	pokecache "github.com/GeminiZA/pokedex/internal/pokeCache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, cmdArgs ...string) error
}

type config struct {
	NextMapURL string
	PrevMapURL string
	dex        map[string]*pokeAPI.Pokemon
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
			name:        "explore",
			description: "Display a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List your pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandMap(cfg *config, cmdArgs ...string) error {
	if len(cmdArgs) > 0 {
		return fmt.Errorf("Invalid arguments for map; usage: map")
	}
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

func commandMapB(cfg *config, cmdArgs ...string) error {
	if len(cmdArgs) > 0 {
		return fmt.Errorf("Invalid arguments for mapb; usage: mapb")
	}
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

func commandExplore(cfg *config, cmdArgs ...string) error {
	if len(cmdArgs) != 1 {
		return fmt.Errorf("Invalid arguments for explore; usage: explore <area_name>")
	}
	locationStr := cmdArgs[0]
	fmt.Println("Exploring", locationStr)
	exploreRes, err := api.Explore(locationStr)
	if err != nil {
		return err
	}
	if len(exploreRes.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
		for _, encounter := range exploreRes.PokemonEncounters {
			fmt.Printf(" - %s\n", encounter.Pokemon.Name)
		}
	}
	return nil
}

func commandCatch(cfg *config, cmdArgs ...string) error {
	if len(cmdArgs) != 1 {
		return fmt.Errorf("Invalid arguments for catch; usage: catch <pokemon>")
	}
	nameStr := cmdArgs[0]
	pokemon, err := api.GetPokemon(nameStr)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", nameStr)
	randRes := rand.Int() % 700
	if randRes > pokemon.BaseExperience {
		cfg.dex[nameStr] = pokemon
		fmt.Printf("%s was caught!\n", nameStr)
	} else {
		fmt.Printf("%s escaped!\n", nameStr)
	}
	return nil
}

func commandInspect(cfg *config, cmdArgs ...string) error {
	if len(cmdArgs) != 1 {
		return fmt.Errorf("Invalid arguments for inspect; usage: inspect <pokemon>")
	}
	nameStr := cmdArgs[0]
	if pokemon, ok := cfg.dex[nameStr]; ok {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range pokemon.Stats {
			fmt.Printf("\t-%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, tp := range pokemon.Types {
			fmt.Printf("\t- %s\n", tp.Type.Name)
		}
	} else {
		fmt.Printf("you have not caught that pokemon\n")
	}
	return nil
}

func commandPokedex(cfg *config, cmdArgs ...string) error {
	if len(cmdArgs) > 0 {
		return fmt.Errorf("Invalid arguments for pokedex; usage: pokedex")
	}
	fmt.Println("Your Pokedex:")
	for k := range cfg.dex {
		fmt.Printf("\t- %s\n", k)
	}
	return nil
}

func commandExit(cfg *config, cmdArgs ...string) error {
	if len(cmdArgs) > 0 {
		return fmt.Errorf("Invalid arguments for exit; usage: exit")
	}
	defer os.Exit(0)
	return nil
}

func commandHelp(cfg *config, cmdArgs ...string) error {
	if len(cmdArgs) > 1 {
		return fmt.Errorf("Invalid arguments for help; usage: help")
	}
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
		dex:        make(map[string]*pokeAPI.Pokemon),
	}
	for true {
		fmt.Print("\nPokedex > ")
		if scanner.Scan() {
			line := scanner.Text()
			commandArgs := strings.Split(line, " ")
			if v, ok := commands[commandArgs[0]]; ok {
				err := v.callback(cfg, commandArgs[1:]...)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Invalid command:", line)
				fmt.Println("Use 'help' for more information")
			}
		}
	}
}
