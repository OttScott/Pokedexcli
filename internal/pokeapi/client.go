package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetLocationAreas(url *string) (*string, *string, *LocationAreaList, error) {
	var locationAreaList LocationAreaList
	var nextURL, previousURL *string
	defaultLocationAreaURL := "https://pokeapi.co/api/v2/location-area?limit=20"

	if url == nil {
		url = &defaultLocationAreaURL
	}

	resp, err := http.Get(*url)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("failed to fetch location areas: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&locationAreaList); err != nil {
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
