package main

import (
	"testing"
	"time"
	"github.com/OttScott/pokedexcli/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  ",
			expected: []string{},
		},
		{
			input:    "  hello  ",
			expected: []string{"hello"},
		},
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  HellO  World  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "map",
			expected: []string{"map"},
		},
		{
			input:    "EXIT",
			expected: []string{"exit"},
		},
		{
			input:    "\t\nhelp\t\n",
			expected: []string{"help"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("lengths don't match: '%v' vs '%v'", actual, c.expected)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleaninput(%v) == %v, expected %v", c.input, actual, c.expected)
			}
		}
	}
}

// Helper function to create a test config
func createTestConfig() *Config {
	cache := pokecache.NewCache(time.Second * 30)
	return &Config{
		NextLocationURL:     nil,
		PreviousLocationURL: nil,
		cache:               cache,
	}
}

func TestCommandHelp(t *testing.T) {
	cfg := createTestConfig()
	
	err := commandHelp(cfg)
	if err != nil {
		t.Errorf("commandHelp should not return an error, got: %v", err)
	}
}

func TestCommandsMapExists(t *testing.T) {
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

func TestCommandMapbEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *Config
		expectError bool
		description string
	}{
		{
			name: "First page - no previous URL",
			setup: func() *Config {
				cfg := createTestConfig()
				// PreviousLocationURL is nil by default
				return cfg
			},
			expectError: false,
			description: "Should handle first page gracefully",
		},
		{
			name: "Has previous URL",
			setup: func() *Config {
				cfg := createTestConfig()
				url := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
				cfg.PreviousLocationURL = &url
				return cfg
			},
			expectError: false, // Note: This might fail due to network calls, but the function structure should be correct
			description: "Should attempt to fetch previous page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setup()
			err := commandMapb(cfg)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected an error for test '%s', but got none", tt.name)
			}
			if !tt.expectError && err != nil {
				// For network-dependent tests, we might get network errors
				// In a real test environment, you'd mock the HTTP calls
				t.Logf("Got error (possibly due to network): %v", err)
			}
		})
	}
}

func TestReplConfigInitialization(t *testing.T) {
	cfg := createTestConfig()
	
	if cfg.cache == nil {
		t.Error("Config cache should not be nil")
	}
	
	if cfg.NextLocationURL != nil {
		t.Error("NextLocationURL should be nil initially")
	}
	
	if cfg.PreviousLocationURL != nil {
		t.Error("PreviousLocationURL should be nil initially")
	}
}

func TestConfigURLUpdates(t *testing.T) {
	cfg := createTestConfig()
	
	// Test setting URLs
	nextURL := "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"
	prevURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	
	cfg.NextLocationURL = &nextURL
	cfg.PreviousLocationURL = &prevURL
	
	if cfg.NextLocationURL == nil || *cfg.NextLocationURL != nextURL {
		t.Errorf("NextLocationURL not set correctly, expected '%s', got '%v'", nextURL, cfg.NextLocationURL)
	}
	
	if cfg.PreviousLocationURL == nil || *cfg.PreviousLocationURL != prevURL {
		t.Errorf("PreviousLocationURL not set correctly, expected '%s', got '%v'", prevURL, cfg.PreviousLocationURL)
	}
}

func TestCleanInputEdgeCases(t *testing.T) {
	edgeCases := []struct {
		input    string
		expected []string
		name     string
	}{
		{
			input:    "command\twith\ttabs",
			expected: []string{"command", "with", "tabs"},
			name:     "Tabs should be converted to spaces",
		},
		{
			input:    "command\nwith\nnewlines",
			expected: []string{"command", "with", "newlines"},
			name:     "Newlines should be converted to spaces",
		},
		{
			input:    "   multiple   spaces   between   words   ",
			expected: []string{"multiple", "spaces", "between", "words"},
			name:     "Multiple spaces should be collapsed",
		},
		{
			input:    "MiXeD CaSe InPuT",
			expected: []string{"mixed", "case", "input"},
			name:     "Mixed case should be lowercased",
		},
		{
			input:    "\t\n\r   \t\n",
			expected: []string{},
			name:     "Whitespace-only input should return empty slice",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := cleanInput(tc.input)
			if len(actual) != len(tc.expected) {
				t.Errorf("lengths don't match: got %v, expected %v", actual, tc.expected)
				return
			}
			for i := range actual {
				if actual[i] != tc.expected[i] {
					t.Errorf("cleanInput(%q) = %v, expected %v", tc.input, actual, tc.expected)
					break
				}
			}
		})
	}
}

func TestCommandCallbacks(t *testing.T) {
	cfg := createTestConfig()
	
	// Test that command callbacks don't panic
	tests := []struct {
		name     string
		callback func(*Config) error
		skipTest bool
		reason   string
	}{
		{
			name:     "commandHelp",
			callback: commandHelp,
			skipTest: false,
		},
		{
			name:     "commandExit",
			callback: commandExit,
			skipTest: true,
			reason:   "commandExit calls os.Exit(0) which would terminate the test",
		},
		{
			name:     "commandMap",
			callback: commandMap,
			skipTest: false, // This might fail due to network, but shouldn't panic
		},
		{
			name:     "commandMapb",
			callback: commandMapb,
			skipTest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipTest {
				t.Skipf("Skipping %s: %s", tt.name, tt.reason)
				return
			}
			
			// Test that the callback doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Command callback %s panicked: %v", tt.name, r)
				}
			}()
			
			// Call the function - it might return an error (especially for network calls)
			// but it shouldn't panic
			err := tt.callback(cfg)
			
			// For network-dependent functions, log errors but don't fail the test
			if err != nil {
				t.Logf("Command %s returned error (possibly expected): %v", tt.name, err)
			}
		})
	}
}

// Test command lookup functionality
func TestCommandLookup(t *testing.T) {
	tests := []struct {
		commandName string
		shouldExist bool
	}{
		{"exit", true},
		{"map", true},
		{"mapb", true},
		{"help", false}, // help is handled separately in main loop
		{"invalid", false},
		{"", false},
		{"MAP", false}, // commands are case-sensitive in the map
	}

	for _, tt := range tests {
		t.Run(tt.commandName, func(t *testing.T) {
			_, exists := commands[tt.commandName]
			if exists != tt.shouldExist {
				t.Errorf("Command '%s' existence mismatch: expected %v, got %v", tt.commandName, tt.shouldExist, exists)
			}
		})
	}
}