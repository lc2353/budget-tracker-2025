package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type AppConfig struct {
	FirestoreProjectID      string
	FirestoreCredentialPath string
}

func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	cfg := &AppConfig{
		FirestoreProjectID:      getEnv("FIRESTORE_PROJECT_ID", ""),
		FirestoreCredentialPath: getEnv("FIRESTORE_CREDENTIAL_PATH", ""),
	}

	if cfg.FirestoreProjectID == "" {
		log.Fatal("FIRESTORE_PROJECT_ID is not set in the environment variables or .env file")
	}

	return cfg
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
