package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AppConfig *Config

type Config struct {
	DBHost        string
	DBPort        string
	DBUsername    string
	DBPassword    string
	DBDatabase    string
	JWTPrivateKey string
	JWTPublicKey  string
}

func InitEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUsername:    os.Getenv("DB_USERNAME"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBDatabase:    os.Getenv("DB_DATABASE"),
		JWTPrivateKey: os.Getenv("JWT_PRIVATE_KEY"),
		JWTPublicKey:  os.Getenv("JWT_PUBLIC_KEY"),
	}
}
