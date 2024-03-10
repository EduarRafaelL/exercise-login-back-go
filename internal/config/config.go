package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver string
	DBSource string
}

func LoadConfig() (config Config, err error) {
	// Carga el archivo .env desde el directorio actual
	err = godotenv.Load()
	if err != nil {
		log.Println("No config file found")
	}

	config = Config{
		DBDriver: os.Getenv("DB_DRIVER"),
		DBSource: os.Getenv("DB_SOURCE"),
	}

	return config, nil
}
