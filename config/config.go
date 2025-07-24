package config

import (
	"log"
	"os"

	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	ServerPort    string

	DBName string
	DBHost string
	DBPort string

	DBUser     string
	DBPassword string
}

func InitCFG() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(errs.ErrCFGLoad())
	}

	return &Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		ServerPort:    os.Getenv("SERVER_PORT"),

		DBName: os.Getenv("DB_NAME"),
		DBHost: os.Getenv("DB_HOST"),

		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}
