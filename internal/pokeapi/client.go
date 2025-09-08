package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/OttScott/pokedexcli/internal/pokecache"
)

type LocationAreaList struct {
    Count    int                 `json:"count"`
    Next     *string             `json:"next"`
    Previous *string             `json:"previous"`
    Results  []NamedAPIResource  `json:"results"`
}

type NamedAPIResource struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}

func GetLocationAreas(cache *pokecache.Cache, url *string) (*string, *string, *LocationAreaList, error) {
	var locationAreaList LocationAreaList
	var nextURL, previousURL *string
	defaultLocationAreaURL := "https://pokeapi.co/api/v2/location-area?limit=20"

	if url == nil {
		url = &defaultLocationAreaURL
	}

	// Check cache for existing data before making an API call
	cachedData, found := cache.Get(*url)
	if found {
		// Parse cached JSON data
		if err := json.Unmarshal(cachedData, &locationAreaList); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to parse cached data: %v", err)
		}
		
		if locationAreaList.Next != nil {
			nextURL = locationAreaList.Next
		}
		if locationAreaList.Previous != nil {
			previousURL = locationAreaList.Previous
		}
		
		return nextURL, previousURL, &locationAreaList, nil
	}

	resp, err := http.Get(*url)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("failed to fetch location areas: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Store in cache
	cache.Add(*url, body)

	// Parse the JSON data
	if err := json.Unmarshal(body, &locationAreaList); err != nil {
		return nil, nil, nil, err
	}
	fmt.Printf("Fetched location areas: %+v\n", locationAreaList)
	if locationAreaList.Next != nil {
		nextURL = locationAreaList.Next
	}
	if locationAreaList.Previous != nil {
		previousURL = locationAreaList.Previous
	}

	return nextURL, previousURL, &locationAreaList, nil
}

type LocationAreaDetail struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             NamedAPIResource      `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedAPIResource `json:"encounter_method"`
	VersionDetails  []VersionDetail  `json:"version_details"`
}

type VersionDetail struct {
	Rate    int              `json:"rate"`
	Version NamedAPIResource `json:"version"`
}

type Name struct {
	Language NamedAPIResource `json:"language"`
	Name     string           `json:"name"`
}

type PokemonEncounter struct {
	Pokemon        NamedAPIResource `json:"pokemon"`
	VersionDetails []VersionDetail  `json:"version_details"`
}

func GetLocationPokemons(cache *pokecache.Cache, locationAreaName string) ([]string, error) {
	if locationAreaName == "" {
		return nil, fmt.Errorf("location area name cannot be empty when fetching location pokemons")
	}
	
	// Fix the URL format - remove the extra slash and query parameter
	LocationAreaURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationAreaName)
	url := &LocationAreaURL
	
	// Check cache for existing data before making an API call
	cachedData, found := cache.Get(*url)
	if found {
		var locationDetail LocationAreaDetail
		// Parse cached JSON data
		if err := json.Unmarshal(cachedData, &locationDetail); err != nil {
			return nil, fmt.Errorf("failed to parse cached data: %v", err)
		}

		// Extract Pokemon names from the encounters
		var pokemonNames []string
		for _, encounter := range locationDetail.PokemonEncounters {
			pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
		}

		return pokemonNames, nil
	}

	resp, err := http.Get(*url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location area: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch location area: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Store in cache
	cache.Add(*url, body)

	// Parse the JSON data
	var locationDetail LocationAreaDetail
	if err := json.Unmarshal(body, &locationDetail); err != nil {
		return nil, fmt.Errorf("failed to parse location area JSON: %v", err)
	}

	// Extract Pokemon names from the encounters
	var pokemonNames []string
	for _, encounter := range locationDetail.PokemonEncounters {
		pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
	}

	return pokemonNames, nil
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}

type PokemonType struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"type"`
}

type PokemonInfo struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	BaseExperience int           `json:"base_experience"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

func GetPokemonInfo(cache *pokecache.Cache, pokemonName string) (*PokemonInfo, error) {
	if pokemonName == "" {
		return nil, fmt.Errorf("pokemon name cannot be empty when fetching pokemon info")
	}

	// Fix the URL format - remove the extra slash and query parameter
	PokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	url := &PokemonURL

	// Check cache for existing data before making an API call
	cachedData, found := cache.Get(*url)
	if found {
		var pokemon PokemonInfo
		// Parse cached JSON data
		if err := json.Unmarshal(cachedData, &pokemon); err != nil {
			return nil, fmt.Errorf("failed to parse cached data: %v", err)
		}
		return &pokemon, nil
	}

	resp, err := http.Get(*url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pokemon: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Store in cache
	cache.Add(*url, body)

	// Parse the JSON data
	var pokemon PokemonInfo
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return nil, fmt.Errorf("failed to parse pokemon JSON: %v", err)
	}

	return &pokemon, nil
}