package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver  string
	DBSource  string
	SecretKey string
}

// LoadConfig loads the configuration from the environment variables
func LoadConfig() (config Config, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Println("No config file found")
	}

	config = Config{
		DBDriver:  os.Getenv("DB_DRIVER"),
		DBSource:  os.Getenv("DB_SOURCE"),
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

	return config, nil
}
