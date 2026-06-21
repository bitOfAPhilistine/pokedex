package commands

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"time"
	"math"
	"math/rand"
	"github.com/bitofaphilistine/pokedex/internal/pokecache"
)

var config = struct {
	nextUrl string
	prevUrl string
}{
	nextUrl: "https://pokeapi.co/api/v2/location-area",
	prevUrl: "",
}

var pokedex = make(map[string]Pokemon)

type cliCommand struct {
	name string
	description string
	callback func(args ...string) error
}

var cache = pokecache.NewCache(5 * time.Second)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {"help", "Display this help message", commandHelp},
		"exit": {"exit", "Exit the Pokedex", commandExit},
		"map": {"map", "List 20 locations in the Pokemon world, each subsequent command will display the next 20 locations", commandMap},
		"mapb": {"mapb", "List the previous 20 locations in the Pokemon world", commandMapB},
		"explore": {"explore <location>", "Explore a location in the Pokemon world", commandExplore},
		"catch": {"catch <pokemon>", "Attempt to catch a Pokemon", commandCatch},
		"pokedex": {"pokedex", "Display the Pokedex", commandPokedex},
		"inspect": {"inspect <pokemon>", "Inspect a Pokemon's stats and types", commandInspect},
	}
}

func Command(name string, args ...string) error {
	cmd, ok := getCommands()[name]
	if ok {
		return cmd.callback(args...)
	}
	return fmt.Errorf("Unknown command: %s", name)
}

func commandHelp(_ ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")
	for _, cmd := range getCommands() {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(_ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreaResponse struct {
	Next string `json:"next"`
	Prev string `json:"previous"`
	Results []Location `json:"results"`
}

func commandMap(_ ...string) error {
	cached, isCached := cache.Get(config.nextUrl)

	var locationAreaResponse LocationAreaResponse

	if isCached {
		if err := json.Unmarshal(cached, &locationAreaResponse); err != nil {
			return err
		}
	} else {
		res, err := http.Get(config.nextUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		
		if err := json.NewDecoder(res.Body).Decode(&locationAreaResponse); err != nil {
			return err
		}

		cacheData, err := json.Marshal(locationAreaResponse)
		if err != nil {
			return err
		}
		cache.Add(config.nextUrl, []byte(cacheData))
	}

	for _, location := range locationAreaResponse.Results {
		fmt.Println(location.Name)
	}
	config.nextUrl = locationAreaResponse.Next
	config.prevUrl = locationAreaResponse.Prev
	return nil
}

func commandMapB(_ ...string) error {
	if config.prevUrl == "" {
		return fmt.Errorf("No previous locations available")
	}
	cached, isCached := cache.Get(config.prevUrl)

	var locationAreaResponse LocationAreaResponse

	if isCached {
		if err := json.Unmarshal(cached, &locationAreaResponse); err != nil {
			return err
		}
	} else {
		res, err := http.Get(config.prevUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		
		if err := json.NewDecoder(res.Body).Decode(&locationAreaResponse); err != nil {
			return err
		}

		cacheData, err := json.Marshal(locationAreaResponse)
		if err != nil {
			return err
		}
		cache.Add(config.prevUrl, []byte(cacheData))
	}

	for _, location := range locationAreaResponse.Results {
		fmt.Println(location.Name)
	}
	config.nextUrl = locationAreaResponse.Next
	config.prevUrl = locationAreaResponse.Prev
	return nil
}

type pokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

type LocationResponse struct {
	Encounters []pokemonEncounter `json:"pokemon_encounters"`
}

func commandExplore(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no location specified")
	}
	
	location := args[0]
	fmt.Printf("Exploring %s...\n", location)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)
	cached, isCached := cache.Get(url)
	var pokemonResponse LocationResponse

	if isCached {
		if err := json.Unmarshal(cached, &pokemonResponse); err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch data: %w", err)
		}
		defer res.Body.Close()
		
		if err := json.NewDecoder(res.Body).Decode(&pokemonResponse); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		cacheData, err := json.Marshal(pokemonResponse)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		cache.Add(url, []byte(cacheData))
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemonResponse.Encounters {
		fmt.Println(" -", pokemon.Pokemon.Name)
	}
	return nil
}

type stat struct {
	Value int `json:"base_stat"`
	Stat struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type Type struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type Pokemon struct {
	Name string `json:"name"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []stat `json:"stats"`
	Types []Type `json:"types"`
	Exp int `json:"base_experience"`
}

func commandCatch(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no pokemon specified")
	}

	pokemon := args[0]

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)
	cached, isCached := cache.Get(url)
	var pokemonResponse Pokemon

	if isCached {
		if err := json.Unmarshal(cached, &pokemonResponse); err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch data: %w", err)
		}
		defer res.Body.Close()
		
		if err := json.NewDecoder(res.Body).Decode(&pokemonResponse); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		cacheData, err := json.Marshal(pokemonResponse)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		cache.Add(url, []byte(cacheData))
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonResponse.Name)
	if rand.Float64() * math.Log10(float64(pokemonResponse.Exp)) < 1.0 {
		fmt.Println("Caught!")
		pokedex[pokemon] = pokemonResponse
	} else {
		fmt.Println(pokemonResponse.Name, "escaped!")
	}
	return nil
}

func commandPokedex(_ ...string) error {
	fmt.Println("Pokedex:")
	for name, _ := range pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandInspect(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no pokemon specified")
	}

	pokemon, ok := pokedex[args[0]]
	if !ok {
		return fmt.Errorf("pokemon not found")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.Value)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
	return nil
}