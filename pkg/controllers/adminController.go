package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	config "github.com/kiprop-dave/2faAuth/pkg/config"
	models "github.com/kiprop-dave/2faAuth/pkg/models"
	responses "github.com/kiprop-dave/2faAuth/pkg/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var adminCollection *mongo.Collection = config.GetCollection(config.DB, "admins")

func CreateAdmin(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var admin models.Admin
	if err := c.BodyParser(&admin); err != nil {
		return c.Status(400).JSON(responses.Response{Status: 400, Message: "error"})
	}

	if validErr := validate.Struct(&admin); validErr != nil {
		return c.Status(400).JSON(responses.Response{Status: 400, Message: "error"})
	}

	var conflict models.Admin
	confErr := adminCollection.FindOne(ctx, bson.M{"name": admin.Name}).Decode(&conflict)
	if confErr == nil {
		return c.Status(http.StatusConflict).JSON(responses.Response{Status: 409, Message: "error"})
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 8)
	if err != nil {
		return c.Status(500).JSON(responses.Response{Status: 500, Message: "error"})
	}

	newAdmin := models.Admin{
		Name:     admin.Name,
		Password: string(hashedPwd),
	}

	res, err := adminCollection.InsertOne(ctx, newAdmin)
	if err != nil {
		return c.Status(500).JSON(responses.Response{Status: 500, Message: "error"})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Status: 201, Message: "success", Data: &fiber.Map{"data": res}})
}

func DeleteAdmin(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := c.Params("name")
	filter := bson.M{"name": name}
	_, err := adminCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Print(err.Error())
	}
	return c.SendStatus(http.StatusAccepted)
}
