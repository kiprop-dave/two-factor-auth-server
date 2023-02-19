package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kiprop-dave/2faAuth/pkg/config"
	"github.com/kiprop-dave/2faAuth/pkg/models"
	"github.com/kiprop-dave/2faAuth/pkg/responses"
	"go.mongodb.org/mongo-driver/bson"
)

var attemptsCollection = config.GetCollection(config.DB, "attempts")

func LogAttempt(name, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	attempt := models.Attempt{Name: name, UserId: id, Time: time.Now()}

	_, err := attemptsCollection.InsertOne(ctx, attempt)
	return err
}

func GetAllAttempts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var attempts []models.Attempt

	cursor, err := attemptsCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Print(err.Error())
		return c.SendStatus(http.StatusInternalServerError)
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var attempt models.Attempt
		if err := cursor.Decode(&attempt); err != nil {
			log.Print(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}
		attempts = append(attempts, attempt)
	}

	return c.Status(200).JSON(responses.Response{Status: 200, Message: "success", Data: &fiber.Map{"attempts": attempts}})
}
