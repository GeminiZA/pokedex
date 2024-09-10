package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]*cacheEntry
	mu       sync.RWMutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]*cacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (cache *Cache) Add(key string, val []byte) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if entry, ok := cache.entries[key]; ok {
		entry.createdAt = time.Now()
	} else {
		cache.entries[key] = &cacheEntry{
			val:       val,
			createdAt: time.Now(),
		}
	}
	return nil
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	if entry, ok := cache.entries[key]; ok {
		fmt.Println("Cache hit for:", key)
		return entry.val, true
	} else {
		fmt.Println("Cache miss for:", key)
		return nil, false
	}
}

func (cache *Cache) reapLoop() {
	for {
		time.Sleep(time.Second)
		cache.mu.Lock()
		for k, v := range cache.entries {
			if time.Now().After(v.createdAt.Add(cache.interval)) {
				delete(cache.entries, k)
			}
		}
		cache.mu.Unlock()
	}
}
