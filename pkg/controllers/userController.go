package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	config "github.com/kiprop-dave/2faAuth/pkg/config"
	models "github.com/kiprop-dave/2faAuth/pkg/models"
	responses "github.com/kiprop-dave/2faAuth/pkg/responses"
	utils "github.com/kiprop-dave/2faAuth/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var usersCollections *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
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

	newUser := models.User{
		Name:        user.Name,
		TagId:       user.TagId,
		PhoneNumber: user.PhoneNumber,
		UserId:      utils.GenerateUserId(),
	}

	filter := bson.M{
		"$or": []bson.M{
			{"name": newUser.Name},
			{"tagid": newUser.TagId},
			{"phonenumber": newUser.PhoneNumber},
			{"userid": newUser.UserId},
		},
	}
	var conflict models.User
	confErr := usersCollections.FindOne(ctx, &filter).Decode(&conflict)
	if confErr == nil {
		return c.SendStatus(http.StatusConflict)
	}

	result, err := usersCollections.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(500).JSON(responses.Response{Status: 500, Message: "internal server error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "User created", Data: &fiber.Map{"data": result}})

}

func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Params("id")
	var user models.User

	filter := bson.M{"userid": id}
	projection := bson.M{"_id": 0}
	err := usersCollections.FindOne(ctx, &filter, options.FindOne().SetProjection(projection)).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(http.StatusNoContent)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(200).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Params("id")
	// var user models.User
	filter := bson.M{"userid": id}
	_, err := usersCollections.DeleteOne(ctx, &filter)
	if err != nil {
		log.Print(err.Error())
	}
	return c.SendStatus(http.StatusAccepted)
}
