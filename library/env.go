package library

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("env")
	filename := ".env"
	if env == "prd" || env == "production" {
		filename = ".env.production"
	}
	if env == "bak" || env == "backup" {
		filename = ".env.bak"
	}
	if env == "dev" || env == "development" {
		filename = ".env.development"
	}
	if env == "local" {
		filename = ".env.local"
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading embedded .env file: %v", err)
	}
	// 將讀取的內容轉換為環境變量
	envMap, err := godotenv.Unmarshal(string(data))
	if err != nil {
		log.Fatalf("failed to parse .env file: %v", err)
	}
	for k, v := range envMap {
		log.Printf("setting env var %s=%s", k, v)
		os.Setenv(k, v)
	}
}
