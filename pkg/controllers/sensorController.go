package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	config "github.com/kiprop-dave/2faAuth/pkg/config"
	models "github.com/kiprop-dave/2faAuth/pkg/models"
	"github.com/kiprop-dave/2faAuth/pkg/responses"
	utils "github.com/kiprop-dave/2faAuth/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var sensorCollection *mongo.Collection = config.GetCollection(config.DB, "sensors")

func CreateSensor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sensor models.Sensor

	if err := c.BodyParser(&sensor); err != nil {
		return c.Status(400).JSON(responses.Response{Status: 400, Message: "error"})
	}

	if len(sensor.Name) < 1 {
		return c.Status(400).JSON(responses.Response{Status: 400, Message: "sensor name required"})
	}

	filter := bson.M{"name": sensor.Name}
	confErr := sensorCollection.FindOne(ctx, filter).Err()
	if confErr == nil {
		return c.SendStatus(http.StatusConflict)
	}

	newSensor := models.Sensor{Name: sensor.Name, ApiKey: utils.GenerateApiKey()}
	_, err := sensorCollection.InsertOne(ctx, newSensor)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Status: 201, Message: "success", Data: &fiber.Map{"apiKey": newSensor.ApiKey}})
}
