package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	config "github.com/kiprop-dave/2faAuth/pkg/config"
	models "github.com/kiprop-dave/2faAuth/pkg/models"
	responses "github.com/kiprop-dave/2faAuth/pkg/responses"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollections *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func CreateUser(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User

	defer cancel()

	// Validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: 400, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if valErr := validate.Struct(&user); valErr != nil {
		return c.Status(400).JSON(responses.Response{Status: 400, Message: "error", Data: &fiber.Map{"data": valErr.Error()}})
	}
}
