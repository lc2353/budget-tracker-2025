package setup

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type AppConfig struct {
	FirestoreProjectID       string
	FirestoreCredentialsPath string
	CorsAllowedOrigins       []string
}

func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	cfg := &AppConfig{
		FirestoreProjectID:       getEnv("FIRESTORE_PROJECT_ID", ""),
		FirestoreCredentialsPath: getEnv("FIRESTORE_CREDENTIAL_PATH", ""),
		CorsAllowedOrigins:       parseCSVEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
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

func parseCSVEnv(key, defaultValue string) []string {
	value := getEnv(key, defaultValue)
	if value == "" {
		return []string{}
	}
	return splitString(value, ",")
}

func splitString(s, sep string) []string {
	var parts []string
	currentPart := ""
	for _, r := range s {
		if string(r) == sep {
			parts = append(parts, currentPart)
			currentPart = ""
		} else {
			currentPart += string(r)
		}
	}
	parts = append(parts, currentPart)
	return parts
}
