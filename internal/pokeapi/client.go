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

	if locationAreaList.Next != nil {
		nextURL = locationAreaList.Next
	}
	if locationAreaList.Previous != nil {
		previousURL = locationAreaList.Previous
	}

	return nextURL, previousURL, &locationAreaList, nil
}
