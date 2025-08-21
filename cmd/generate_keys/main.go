package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"tasius.my.id/todolistapi/internal/utils/crypto"
)

func main() {
	outputDir := flag.String("output", "./keys", "Output directory for keys")
	flag.Parse()

	// Create output directory if it doesn't exist
	err := os.MkdirAll(*outputDir, 0700)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate RSA key pair
	privateKey, publicKey, err := crypto.GenerateRSAKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	// Save private key
	privateKeyPath := filepath.Join(*outputDir, "private.pem")
	if err := crypto.SavePEMKey(privateKeyPath, privateKey); err != nil {
		log.Fatalf("Failed to save private key: %v", err)
	}

	// Save public key
	publicKeyPath := filepath.Join(*outputDir, "public.pem")
	if err := crypto.SavePublicPEMKey(publicKeyPath, publicKey); err != nil {
		log.Fatalf("Failed to save public key: %v", err)
	}

	fmt.Printf("Successfully generated key pair:\n  Private key: %s\n  Public key:  %s\n", privateKeyPath, publicKeyPath)
}
