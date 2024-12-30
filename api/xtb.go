package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kejrak/xtb-client-go/xclient"
)

func XtbClient() (*xclient.Client, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get port from environment variable
	username := os.Getenv("X_USER")
	if username == "" {
		log.Fatal("X_USER is not set")
		return nil, fmt.Errorf("user is not set")
	}

	password := os.Getenv("X_PASSWORD")
	if password == "" {
		log.Fatal("X_PASSWORD is not set")
		return nil, fmt.Errorf("password is not set")
	}

	// Initialize the client for the demo environment
	client, err := xclient.NewClient("demo", username, password)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
		return nil, err
	}

	log.Println("XTB Client initialized")
	return client, nil
}
