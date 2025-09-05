package main

import (
	"testing"
	"time"
	"github.com/OttScott/pokedexcli/internal/pokecache"
)

func TestConfigInitialization(t *testing.T) {
	cache := pokecache.NewCache(time.Second * 5)
	
	config := &Config{
		NextLocationURL:     nil,
		PreviousLocationURL: nil,
		cache:               cache,
	}

	// Test that config is properly initialized
	if config.cache == nil {
		t.Error("Cache should not be nil")
	}
	
	if config.NextLocationURL != nil {
		t.Error("NextLocationURL should be nil initially")
	}
	
	if config.PreviousLocationURL != nil {
		t.Error("PreviousLocationURL should be nil initially")
	}
}

func TestCacheIntegration(t *testing.T) {
	cache := pokecache.NewCache(time.Millisecond * 100) // Short TTL for testing
	
	// Test cache functionality
	testKey := "test-url"
	testValue := []byte(`{"test": "data"}`)
	
	// Add to cache
	cache.Add(testKey, testValue)
	
	// Retrieve from cache
	retrievedValue, found := cache.Get(testKey)
	if !found {
		t.Error("Expected to find value in cache")
	}
	
	if string(retrievedValue) != string(testValue) {
		t.Errorf("Expected %s, got %s", string(testValue), string(retrievedValue))
	}
	
	// Wait for TTL to expire
	time.Sleep(time.Millisecond * 150)
	
	// Should be expired now
	_, found = cache.Get(testKey)
	if found {
		t.Error("Expected cache entry to be expired")
	}
}

func TestCleanInputIntegration(t *testing.T) {
	// Test cases for the cleanInput function used in main
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "help",
			expected: []string{"help"},
		},
		{
			input:    "  MAP  ",
			expected: []string{"map"},
		},
		{
			input:    "exit",
			expected: []string{"exit"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("For input '%s': expected length %d, got %d", c.input, len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("For input '%s': expected %v, got %v", c.input, c.expected, actual)
				break
			}
		}
	}
}

// Benchmark test for cache performance
func BenchmarkCacheOperations(b *testing.B) {
	cache := pokecache.NewCache(time.Minute * 5)
	testData := []byte(`{"count": 20, "results": [{"name": "test-location", "url": "https://example.com"}]}`)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		key := "test-key"
		cache.Add(key, testData)
		_, _ = cache.Get(key)
	}
}

// Test that commands map is properly initialized
func TestCommandsInitialization(t *testing.T) {
	expectedCommands := []string{"exit", "map", "mapb"}
	
	for _, cmdName := range expectedCommands {
		cmd, exists := commands[cmdName]
		if !exists {
			t.Errorf("Command '%s' should exist in commands map", cmdName)
		}
		
		if cmd.name != cmdName {
			t.Errorf("Command name mismatch: expected '%s', got '%s'", cmdName, cmd.name)
		}
		
		if cmd.description == "" {
			t.Errorf("Command '%s' should have a description", cmdName)
		}
		
		if cmd.callback == nil {
			t.Errorf("Command '%s' should have a callback function", cmdName)
		}
	}
}

// Example of testing command functionality (you'd need to mock HTTP calls for full testing)
func TestCommandHelpFunction(t *testing.T) {
	cache := pokecache.NewCache(time.Second * 5)
	config := &Config{
		NextLocationURL:     nil,
		PreviousLocationURL: nil,
		cache:               cache,
	}
	
	// Test that commandHelp doesn't return an error
	err := commandHelp(config)
	if err != nil {
		t.Errorf("commandHelp should not return an error, got: %v", err)
	}
}
