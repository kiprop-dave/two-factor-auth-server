package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type envVariable struct {
	dbUri       string
	accessToken string
}

func loadEnvVariables() envVariable {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can't find .env file")
	}
	db := os.Getenv("MONGO_URI")
	aToken := os.Getenv("ACCESS_TOKEN")

	if db == "" || aToken == "" {
		log.Fatal("Cant find environment variables")
	}

	env := envVariable{dbUri: db, accessToken: aToken}
	return env
}

var EnvVariable envVariable = loadEnvVariables()
