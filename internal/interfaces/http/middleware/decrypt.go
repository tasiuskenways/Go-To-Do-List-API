package middleware

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"tasius.my.id/todolistapi/internal/utils/crypto"
)

// RequestBody represents the expected request body format with encrypted data
type RequestBody struct {
	Data string `json:"data"` // Base64 encoded encrypted data
}

// decryptData handles the decryption of the request data
func decryptData(privateKeyPath, encryptedInput string) ([]byte, error) {
	// Load private key
	privateKey, err := crypto.LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Split the input into encrypted key and data
	parts := strings.Split(encryptedInput, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid input format. Expected format: <encrypted_key_hex>:<encrypted_data_hex>")
	}

	encryptedKey, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted key: %w", err)
	}

	encryptedData, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	// Decrypt the AES key with RSA
	aesKey, err := crypto.DecryptWithPrivateKey(privateKey, encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt AES key: %w", err)
	}

	// Decrypt the message with AES
	decryptedData, err := crypto.DecryptAES(aesKey, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message: %w", err)
	}

	return decryptedData, nil
}

// validateAndProcessRequest handles the request validation and processing
func validateAndProcessRequest(c *fiber.Ctx, body []byte) (*RequestBody, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("request body is empty")
	}

	var reqBody RequestBody
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return nil, fmt.Errorf("invalid request format: %w", err)
	}

	// Validate that the data field exists and has the expected format
	if reqBody.Data == "" {
		return nil, fmt.Errorf("encrypted data is required in 'data' field")
	}

	// Check if the data appears to be in the expected format: <encrypted_key_hex>:<encrypted_data_hex>
	parts := strings.Split(reqBody.Data, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid encrypted data format. Expected format: <encrypted_key_hex>:<encrypted_data_hex>")
	}

	// Validate that both parts are valid hex strings
	if _, err := hex.DecodeString(parts[0]); err != nil {
		return nil, fmt.Errorf("invalid encrypted key format: not a valid hex string")
	}

	if _, err := hex.DecodeString(parts[1]); err != nil {
		return nil, fmt.Errorf("invalid encrypted data format: not a valid hex string")
	}

	return &reqBody, nil
}

// DecryptMiddleware decrypts the request body before passing it to the handler
func DecryptMiddleware(privateKeyPath string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only process JSON requests
		if c.Get("Content-Type") != "application/json" {
			return c.Next()
		}

		// Read and validate request
		reqBody, err := validateAndProcessRequest(c, c.Request().Body())
		if err != nil {
			log.Printf("Request validation failed: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Decrypt the data
		decryptedData, err := decryptData(privateKeyPath, reqBody.Data)
		if err != nil {
			log.Printf("Decryption failed: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to decrypt data. Invalid or corrupted data.",
			})
		}

		// Validate JSON
		if !json.Valid(decryptedData) {
			log.Print("Decrypted data is not valid JSON")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Decrypted data is not valid JSON",
			})
		}

		// Update request with decrypted data
		c.Request().SetBody(decryptedData)
		c.Request().Header.SetContentType("application/json")
		c.Request().SetRequestURI(c.OriginalURL())
		c.Locals("decryptedBody", decryptedData)

		return c.Next()
	}
}
