package store

import (
	"sync"
	"time"
)

// KeyData represents a value stored in Redis with optional expiration
type KeyData struct {
	Value      string
	Expiration *time.Time
}

var (
	data  = make(map[string]KeyData)
	mutex = &sync.RWMutex{}
)

// Set stores a key-value pair in the in-memory store with an optional expiration time.
//
// This function adds or updates a key-value pair in the Redis-compatible in-memory data store.
// If an expiration time is provided, the key will automatically expire after the specified duration.
func Set(key, value string, expiry time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()

	keyData := KeyData{
		Value: value,
	}

	if expiry > 0 {
		expiryTime := time.Now().Add(expiry)
		keyData.Expiration = &expiryTime
	}

	data[key] = keyData
}

// Get retrieves a value from the in-memory store by its key.
//
// This function looks up a key in the Redis-compatible in-memory data store
// and returns its associated value. If the key does not exist or has expired,
// it returns an empty string and false.
func Get(key string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	keyData, exists := data[key]
	if !exists {
		return "", false
	}

	if keyData.Expiration != nil && time.Now().After(*keyData.Expiration) {
		mutex.RUnlock()
		mutex.Lock()
		delete(data, key)
		mutex.Unlock()
		mutex.RLock()
		return "", false
	}

	return keyData.Value, true
}

// Delete removes a key from the store.
//
// This function removes a key and its associated value from the Redis-compatible
// in-memory data store if it exists.
func Delete(key string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := data[key]; exists {
		delete(data, key)
		return true
	}
	return false
}

// CleanExpired removes all expired keys from the store.
// This would typically be called periodically by a background goroutine.
func CleanExpired() {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now()
	for key, keyData := range data {
		if keyData.Expiration != nil && now.After(*keyData.Expiration) {
			delete(data, key)
		}
	}
}
