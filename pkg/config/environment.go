package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariable struct {
	DbUri       string
	AccessToken string
}

func loadEnvVariables() EnvVariable {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can't find .env file")
	}
	db := os.Getenv("MONGO_URI")
	aToken := os.Getenv("ACCESS_TOKEN")

	if db == "" || aToken == "" {
		log.Fatal("Cant find environment variables")
	}

	env := EnvVariable{DbUri: db, AccessToken: aToken}
	return env
}

var Environment EnvVariable = loadEnvVariables()
