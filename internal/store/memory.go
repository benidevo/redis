package store

var store = make(map[string]string)

// Set stores a key-value pair in the in-memory store.
//
// This function adds or updates a key-value pair in the Redis-compatible
// in-memory data store. If the key already exists, its value will be overwritten.
func Set(key, value string) {
	store[key] = value
}

// Get retrieves the value associated with a given key from the in-memory store.
//
// This function searches for a value in the Redis-compatible in-memory data store
// based on the provided key. If the key exists, the corresponding value is returned.
// If the key does not exist, the function returns an empty string.
func Get(key string) string {
	return store[key]
}
