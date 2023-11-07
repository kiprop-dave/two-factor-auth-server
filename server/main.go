package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kiprop-dave/2fa/storage"
	twofa "github.com/kiprop-dave/2fa/twoFa"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
}

func main() {
	db, err := storage.ConnectMongoDb()
	if err != nil {
		log.Fatalln("failed to connect to mongo")
	}

	mongoStore, err := storage.NewMongoStorage(db)
	if err != nil {
		log.Fatalln(err)
	}

	twofa := twofa.Authy{}
	server := NewServer(&twofa, ":8080", mongoStore)
	server.Run()
}
