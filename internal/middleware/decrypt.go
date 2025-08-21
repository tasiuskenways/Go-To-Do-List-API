package middleware

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"tasius.my.id/todolistapi/internal/utils/crypto"
)

// RequestBody represents the expected request body format with encrypted data
type RequestBody struct {
	Data string `json:"data"` // Base64 encoded encrypted data
}

// DecryptMiddleware decrypts the request body before passing it to the handler
func DecryptMiddleware(privateKeyPath string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only process JSON requests
		if c.Get("Content-Type") != "application/json" {
			return c.Next()
		}

		// Read the request body
		body := c.Request().Body()
		if len(body) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Request body is empty",
			})
		}

		// Parse the request body to get the encrypted data
		var reqBody RequestBody
		if err := json.Unmarshal(body, &reqBody); err != nil {
			log.Printf("Failed to parse request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format. Expected {\"data\": \"encrypted_data\"}",
			})
		}

		// Load private key
		privateKey, err := crypto.LoadPrivateKey(privateKeyPath)
		if err != nil {
			log.Fatalf("Failed to load private key: %v", err)
		}

		// Split the input into encrypted key and data
		parts := strings.Split(reqBody.Data, ":")
		if len(parts) != 2 {
			log.Fatal("Invalid input format. Expected format: <encrypted_key_hex>:<encrypted_data_hex>")
		}

		encryptedKeyHex := parts[0]
		encryptedDataHex := parts[1]

		// Decode hex strings to bytes
		encryptedKey, err := hex.DecodeString(encryptedKeyHex)
		if err != nil {
			log.Fatalf("Failed to decode encrypted key: %v", err)
		}

		encryptedData, err := hex.DecodeString(encryptedDataHex)
		if err != nil {
			log.Fatalf("Failed to decode encrypted data: %v", err)
		}

		// Decrypt the AES key with RSA
		aesKey, err := crypto.DecryptWithPrivateKey(privateKey, encryptedKey)
		if err != nil {
			log.Fatalf("Failed to decrypt AES key: %v", err)
		}

		// Decrypt the message with AES
		decryptedData, err := crypto.DecryptAES(aesKey, encryptedData)
		if err != nil {
			log.Fatalf("Failed to decrypt message: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to decrypt data. Invalid or corrupted data.",
			})
		}

		log.Printf("Decrypted data: %s", string(decryptedData))

		// Validate the decrypted data is valid JSON
		if !json.Valid(decryptedData) {
			log.Print("Decrypted data is not valid JSON")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Decrypted data is not valid JSON",
			})
		}

		// Replace the request body with the decrypted data
		c.Request().SetBody(decryptedData)
		c.Request().Header.SetContentType("application/json")
		c.Request().SetRequestURI(c.OriginalURL())

		// Store the decrypted data in locals for debugging/audit purposes
		c.Locals("decryptedBody", decryptedData)

		return c.Next()
	}
}
