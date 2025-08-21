package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	// Generate a 32-byte (256-bit) random key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating random key: %v\n", err)
		os.Exit(1)
	}

	// Encode the key to base64 for storage
	encodedKey := base64.RawURLEncoding.EncodeToString(key)
	
	// Output the key in a format that can be directly used in .env
	fmt.Println("JWT_PRIVATE_KEY=" + encodedKey)
}
