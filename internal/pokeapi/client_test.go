package pokeapi

import (
	"testing"
	"time"
	"github.com/OttScott/pokedexcli/internal/pokecache"
)

func TestGetLocationAreas(t *testing.T) {
	// This test would ideally mock the HTTP response and the cache behavior.
	// For simplicity, we will just check if the function can be called without error.
	cache := pokecache.NewCache(time.Second * 5)
	url := "https://pokeapi.co/api/v2/location-area?limit=20"

	nextURL, prevURL, locs, err := GetLocationAreas(cache, &url)
	if err != nil {
		t.Fatalf("GetLocationAreas returned an error: %v", err)
	}

	if locs == nil || len(locs.Results) == 0 {
		t.Fatalf("Expected location areas, got none")
	}

	if nextURL == nil && prevURL == nil {
		t.Logf("No next or previous URL available, which is acceptable for this test.")
	}
}

func TestGetLocationAreas_Cache(t *testing.T) {
	cache := pokecache.NewCache(time.Second * 5)
	url := "https://pokeapi.co/api/v2/location-area?limit=20"

	// First call to populate the cache
	_, _, _, err := GetLocationAreas(cache, &url)
	if err != nil {
		t.Fatalf("First GetLocationAreas call returned an error: %v", err)
	}

	// Second call should hit the cache
	_, _, locs, err := GetLocationAreas(cache, &url)
	if err != nil {
		t.Fatalf("Second GetLocationAreas call returned an error: %v", err)
	}

	if locs == nil || len(locs.Results) == 0 {
		t.Fatalf("Expected location areas from cache, got none")
	}
}

func TestGetLocationAreas_InvalidURL(t *testing.T) {
	cache := pokecache.NewCache(time.Second * 5)
	invalidURL := "https://pokeapi.co/api/v2/invalid-endpoint"

	_, _, _, err := GetLocationAreas(cache, &invalidURL)
	if err == nil {
		t.Fatalf("Expected error for invalid URL, got none")
	}
}

func TestGetLocationAreas_NilURL(t *testing.T) {
	cache := pokecache.NewCache(time.Second * 5)

	nextURL, prevURL, locs, err := GetLocationAreas(cache, nil)
	if err != nil {
		t.Fatalf("GetLocationAreas with nil URL returned an error: %v", err)
	}

	if locs == nil || len(locs.Results) == 0 {
		t.Fatalf("Expected location areas with nil URL, got none")
	}

	if nextURL == nil && prevURL == nil {
		t.Logf("No next or previous URL available with nil URL, which is acceptable for this test.")
	}
}