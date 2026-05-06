package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the environment variables required for the application.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

// LoadConfig reads the .env file and populates the Config struct with brutal validation.
func LoadConfig() *Config {
	// 1. Context Logging (The Lie Detector)
	// This proves exactly which folder the Go executable thinks it is running inside.
	cwd, _ := os.Getwd()
	log.Printf("[BOOT] Execution Context: %s", cwd)

	// 2. Load .env with Strict Error Exposing
	// We no longer swallow the error. If Windows hid the extension or corrupted
	// the encoding, this line will print the exact reason.
	if err := godotenv.Load(); err != nil {
		log.Printf("[CONFIG WARNING] godotenv could not load .env smoothly. Reason: %v", err)
	}

	// 3. Initialize the Config struct
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET_KEY"),
		Port:        getEnv("PORT", "8081"), // Explicitly setting fallback to 8081
	}

	// 4. Brutal Validation
	// The server is strictly forbidden from booting if these are missing.
	if err := cfg.Validate(); err != nil {
		log.Fatalf("[FATAL] Configuration error: %v", err)
	}

	log.Println("[CONFIG] System variables loaded and validated successfully.")
	return cfg
}

// Validate ensures all critical keys are present before the server is allowed to boot.
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is missing. The Go binary cannot see it in the environment")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET_KEY is missing. The Go binary cannot see it in the environment")
	}
	return nil
}

// getEnv is a robust helper to allow for default values.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}