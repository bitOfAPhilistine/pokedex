package commands

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"time"
	"github.com/bitofaphilistine/pokedex/internal/pokecache"
)

var config = struct {
	nextUrl string
	prevUrl string
}{
	nextUrl: "https://pokeapi.co/api/v2/location-area",
	prevUrl: "",
}

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

type pokemon struct {
	Name string `json:"name"`
}

type pokemonEncounter struct {
	Pokemon pokemon `json:"pokemon"`
}

type LocationResponse struct {
	Encounters []pokemonEncounter `json:"pokemon_encounters"`
}

func commandExplore(args ...string) error {
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