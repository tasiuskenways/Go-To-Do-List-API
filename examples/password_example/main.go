package main

import (
	"fmt"
	"log"

	"tasius.my.id/todolistapi/internal/utils/password"
)

func main() {
	// Example 1: Hashing a password
	plainPassword := "mySecurePassword123!"
	hashedPassword, err := password.HashPassword(plainPassword)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	fmt.Printf("Hashed password: %s\n", hashedPassword)

	// Example 2: Verifying a password
	err = password.CheckPassword(plainPassword, hashedPassword)
	if err != nil {
		log.Fatalf("Password check failed: %v", err)
	}
	fmt.Println("Password verified successfully!")

	// Example 3: Generating a random API key
	apiKey, err := password.GenerateAPIKey("sk", 32)
	if err != nil {
		log.Fatalf("Failed to generate API key: %v", err)
	}
	fmt.Printf("Generated API key: %s\n", apiKey)

	// Example 4: Sanitizing user input
	userInput := "  hello\nworld\t"
	sanitized := password.SanitizeString(userInput)
	fmt.Printf("Sanitized input: '%s'\n", sanitized)
}
