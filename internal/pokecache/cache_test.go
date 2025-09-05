package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(time.Second * 5)

	// Test adding and retrieving a value
	key := "testKey"
	value := []byte("testValue")
	cache.Add(key, value)

	if got, _ := cache.Get(key); !equal(got, value) {
		t.Errorf("Expected %q, got %q", value, got)
	}

	// Test TTL expiration
	time.Sleep(time.Second * 6)
	if _, ok := cache.Get(key); ok {
		t.Errorf("Expected cache miss for expired key %q", key)
	}
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
