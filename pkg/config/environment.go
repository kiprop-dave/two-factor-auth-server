package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariable struct {
	DbUri       string
	AccessToken string
	TwilioSid   string
	TwilioToken string
	PhoneNumber string
}

func loadEnvVariables() EnvVariable {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can't find .env file")
	}
	db := os.Getenv("MONGO_URI")
	aToken := os.Getenv("ACCESS_TOKEN")
	twilioSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhone := os.Getenv("TWILIO_NUMBER")

	if db == "" || aToken == "" || twilioSid == "" || twilioToken == "" || twilioPhone == "" {
		log.Fatal("Cant find environment variables")
	}

	env := EnvVariable{DbUri: db, AccessToken: aToken, TwilioSid: twilioSid, TwilioToken: twilioToken, PhoneNumber: twilioPhone}
	return env
}

var Environment EnvVariable = loadEnvVariables()
