package pokecache

import (
	"sync"
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

func TestCacheAdd(t *testing.T) {
	cache := NewCache(time.Minute)

	tests := []struct {
		name  string
		key   string
		value []byte
	}{
		{
			name:  "Normal string",
			key:   "test-key",
			value: []byte("test-value"),
		},
		{
			name:  "Empty string",
			key:   "empty",
			value: []byte(""),
		},
		{
			name:  "JSON data",
			key:   "json",
			value: []byte(`{"name": "pikachu", "type": "electric"}`),
		},
		{
			name:  "Binary data",
			key:   "binary",
			value: []byte{0x00, 0x01, 0x02, 0xFF},
		},
		{
			name:  "Large data",
			key:   "large",
			value: make([]byte, 1024), // 1KB of zeros
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.Add(tt.key, tt.value)
			
			retrieved, found := cache.Get(tt.key)
			if !found {
				t.Errorf("Expected to find key %q", tt.key)
			}
			
			if !equal(retrieved, tt.value) {
				t.Errorf("Expected %v, got %v", tt.value, retrieved)
			}
		})
	}
}

func TestCacheGet(t *testing.T) {
	cache := NewCache(time.Minute)

	// Test getting non-existent key
	_, found := cache.Get("non-existent")
	if found {
		t.Error("Expected cache miss for non-existent key")
	}

	// Test getting existing key
	key := "existing-key"
	value := []byte("existing-value")
	cache.Add(key, value)

	retrieved, found := cache.Get(key)
	if !found {
		t.Error("Expected cache hit for existing key")
	}

	if !equal(retrieved, value) {
		t.Errorf("Expected %v, got %v", value, retrieved)
	}
}

func TestCacheOverwrite(t *testing.T) {
	cache := NewCache(time.Minute)

	key := "overwrite-key"
	value1 := []byte("first-value")
	value2 := []byte("second-value")

	// Add first value
	cache.Add(key, value1)
	retrieved1, found1 := cache.Get(key)
	if !found1 || !equal(retrieved1, value1) {
		t.Errorf("Expected first value %v, got %v", value1, retrieved1)
	}

	// Overwrite with second value
	cache.Add(key, value2)
	retrieved2, found2 := cache.Get(key)
	if !found2 || !equal(retrieved2, value2) {
		t.Errorf("Expected second value %v, got %v", value2, retrieved2)
	}

	// Ensure first value is gone
	if equal(retrieved2, value1) {
		t.Error("Cache should have been overwritten")
	}
}

func TestCacheTTLExpiration(t *testing.T) {
	shortTTL := time.Millisecond * 100
	cache := NewCache(shortTTL)

	key := "ttl-key"
	value := []byte("ttl-value")

	// Add value
	cache.Add(key, value)

	// Should be available immediately
	_, found := cache.Get(key)
	if !found {
		t.Error("Value should be available immediately after adding")
	}

	// Wait for expiration
	time.Sleep(shortTTL + time.Millisecond*10)

	// Should be expired now
	_, found = cache.Get(key)
	if found {
		t.Error("Value should be expired after TTL")
	}
}

func TestCacheReapLoop(t *testing.T) {
	shortTTL := time.Millisecond * 50
	cache := NewCache(shortTTL)

	// Add multiple entries
	for i := 0; i < 10; i++ {
		key := string(rune('a' + i))
		value := []byte("value-" + key)
		cache.Add(key, value)
	}

	// Verify all entries exist
	for i := 0; i < 10; i++ {
		key := string(rune('a' + i))
		if _, found := cache.Get(key); !found {
			t.Errorf("Expected to find key %q", key)
		}
	}

	// Wait for TTL expiration and reap loop to run
	time.Sleep(shortTTL + time.Second*6) // reap loop runs every 5 seconds

	// Check that entries were reaped
	expiredCount := 0
	for i := 0; i < 10; i++ {
		key := string(rune('a' + i))
		if _, found := cache.Get(key); !found {
			expiredCount++
		}
	}

	if expiredCount == 0 {
		t.Error("Expected some entries to be reaped by reap loop")
	}
}

func TestCacheConcurrency(t *testing.T) {
	cache := NewCache(time.Minute)
	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup

	// Test concurrent writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := string(rune('a' + (id%26)))
				value := []byte("value-" + string(rune('0'+j%10)))
				cache.Add(key, value)
			}
		}(i)
	}

	// Test concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := string(rune('a' + (id%26)))
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// If we get here without deadlocks or data races, the test passes
	t.Log("Concurrency test completed successfully")
}

func TestCacheConcurrentReadWrite(t *testing.T) {
	cache := NewCache(time.Second)
	const duration = time.Millisecond * 500

	var wg sync.WaitGroup
	done := make(chan bool)

	// Start writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					key := string(rune('w' + id%10))
					value := []byte("writer-" + string(rune('0'+id%10)))
					cache.Add(key, value)
					time.Sleep(time.Millisecond)
				}
			}
		}(i)
	}

	// Start readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					key := string(rune('w' + id%10))
					cache.Get(key)
					time.Sleep(time.Millisecond)
				}
			}
		}(i)
	}

	// Let it run for a short duration
	time.Sleep(duration)
	close(done)
	wg.Wait()

	t.Log("Concurrent read/write test completed successfully")
}

func TestCacheEdgeCases(t *testing.T) {
	cache := NewCache(time.Minute)

	// Test empty key
	cache.Add("", []byte("empty-key-value"))
	if _, found := cache.Get(""); !found {
		t.Error("Should be able to use empty string as key")
	}

	// Test nil value (should work since []byte{} is valid)
	cache.Add("nil-value", nil)
	if retrieved, found := cache.Get("nil-value"); !found || retrieved != nil {
		t.Error("Should be able to store nil value")
	}

	// Test very long key
	longKey := string(make([]byte, 1000))
	cache.Add(longKey, []byte("long-key-value"))
	if _, found := cache.Get(longKey); !found {
		t.Error("Should be able to use very long keys")
	}
}

func TestCacheZeroTTL(t *testing.T) {
	cache := NewCache(0) // Zero TTL means immediate expiration
	
	key := "zero-ttl-key"
	value := []byte("zero-ttl-value")
	
	cache.Add(key, value)
	
	// With zero TTL, the value should expire immediately
	_, found := cache.Get(key)
	if found {
		t.Error("Value with zero TTL should expire immediately")
	}
}

func TestCacheVeryShortTTL(t *testing.T) {
	cache := NewCache(time.Nanosecond) // Very short TTL
	
	key := "short-ttl-key"
	value := []byte("short-ttl-value")
	
	cache.Add(key, value)
	
	// Even with very short TTL, might still be available if checked immediately
	// But should definitely be gone after a small sleep
	time.Sleep(time.Millisecond)
	_, found := cache.Get(key)
	if found {
		t.Error("Value with nanosecond TTL should expire very quickly")
	}
}

func TestCacheMemoryEfficiency(t *testing.T) {
	cache := NewCache(time.Millisecond * 10)
	
	// Add many entries that will expire
	for i := 0; i < 1000; i++ {
		key := string(rune('A' + i%26)) + string(rune('0' + i%10))
		value := make([]byte, 100) // 100 bytes each
		cache.Add(key, value)
	}
	
	// Wait for expiration
	time.Sleep(time.Millisecond * 20)
	
	// Try to access expired entries (should clean them up via Get)
	cleanedCount := 0
	for i := 0; i < 1000; i++ {
		key := string(rune('A' + i%26)) + string(rune('0' + i%10))
		if _, found := cache.Get(key); !found {
			cleanedCount++
		}
	}
	
	if cleanedCount < 900 { // Most should be cleaned up
		t.Errorf("Expected most entries to be cleaned up, only %d were cleaned", cleanedCount)
	}
}

// Benchmark tests
func BenchmarkCacheAdd(b *testing.B) {
	cache := NewCache(time.Minute)
	value := []byte("benchmark-value")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := string(rune('a' + i%26))
		cache.Add(key, value)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache(time.Minute)
	value := []byte("benchmark-value")
	
	// Pre-populate cache
	for i := 0; i < 26; i++ {
		key := string(rune('a' + i))
		cache.Add(key, value)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := string(rune('a' + i%26))
		cache.Get(key)
	}
}

func BenchmarkCacheConcurrentAccess(b *testing.B) {
	cache := NewCache(time.Minute)
	value := []byte("benchmark-value")
	
	// Pre-populate cache
	for i := 0; i < 100; i++ {
		key := string(rune('a' + i%26))
		cache.Add(key, value)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := string(rune('a' + (b.N % 26)))
			if b.N%2 == 0 {
				cache.Get(key)
			} else {
				cache.Add(key, value)
			}
		}
	})
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
