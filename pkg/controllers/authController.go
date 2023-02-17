package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	config "github.com/kiprop-dave/2faAuth/pkg/config"
	models "github.com/kiprop-dave/2faAuth/pkg/models"
	"github.com/kiprop-dave/2faAuth/pkg/responses"
	"github.com/kiprop-dave/2faAuth/pkg/services"
	"github.com/kiprop-dave/2faAuth/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var env = config.Environment

// Struct to be encoded to jwt
type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func AuthAdmin(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var admin models.Admin
	if err := c.BodyParser(&admin); err != nil {
		c.Status(400).JSON(responses.Response{Status: 400, Message: "error"})
	}

	if validErr := validate.Struct(&admin); validErr != nil {
		c.Status(400).JSON(responses.Response{Status: 400, Message: "error"})
	}

	filter := bson.M{"name": admin.Name}
	var storedAdmin models.Admin
	qerr := adminCollection.FindOne(ctx, filter).Decode(&storedAdmin)
	if qerr != nil {
		if qerr == mongo.ErrNoDocuments {
			return c.Status(http.StatusUnauthorized).JSON(responses.Response{Status: 401})
		}
		log.Print(qerr.Error())
		return c.SendStatus(http.StatusInternalServerError)
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedAdmin.Password), []byte(admin.Password))
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Name: admin.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(env.AccessToken))
	if err != nil {
		log.Print(err.Error())
		return c.SendStatus(500)
	}
	return c.Status(200).JSON(responses.TokenResponse{Token: tokenString})

}

func AuthUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sensor models.AuthUserRequest

	if err := c.BodyParser(&sensor); err != nil {
		return c.Status(400).JSON(responses.Response{Status: 400, Message: "error"})
	}

	if validErr := validate.Struct(&sensor); validErr != nil {
		return c.Status(400).JSON(responses.Response{Status: 400, Message: "error"})
	}

	filter := bson.M{"name": sensor.Name}
	var storedSensor models.Sensor
	err := sensorCollection.FindOne(ctx, &filter).Decode(&storedSensor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(http.StatusUnauthorized)
		}
		log.Print(err.Error())
		return c.SendStatus(http.StatusInternalServerError)
	}
	if storedSensor.ApiKey != sensor.ApiKey {
		return c.SendStatus(http.StatusUnauthorized)
	}

	filter = bson.M{"tagId": sensor.TagId}
	var user models.User
	err = usersCollections.FindOne(ctx, &filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(http.StatusForbidden)
		}
		log.Print(err.Error())
		return c.SendStatus(http.StatusInternalServerError)
	}

	err = LogAttempt(user.Name, user.UserId)
	if err != nil {
		log.Print(err.Error())
		return c.SendStatus(http.StatusInternalServerError)
	}

	otp := utils.GenerateOtp()
	err = services.SendSms(user.PhoneNumber, otp)
	if err != nil {
		log.Print(err.Error())
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.Status(200).JSON(responses.TokenResponse{Token: otp})
}
