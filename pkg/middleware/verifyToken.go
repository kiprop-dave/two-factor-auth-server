package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	config "github.com/kiprop-dave/2faAuth/pkg/config"
	controllers "github.com/kiprop-dave/2faAuth/pkg/controllers"
)

var env = config.Environment

func VerifyToken(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	authHeader, ok := headers["Authorization"]
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}

	clientTkn := strings.Split(authHeader, " ")[1]

	claims := &controllers.Claims{}

	tkn, err := jwt.ParseWithClaims(clientTkn, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(env.AccessToken), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return c.SendStatus(http.StatusUnauthorized)
		}
		if err == jwt.ErrTokenExpired {
			return c.Status(http.StatusBadRequest).JSON("expired token")
		}
		return c.Status(http.StatusBadRequest).JSON("expired token")
	}
	if !tkn.Valid {
		return c.SendStatus(http.StatusUnauthorized)
	}
	return c.Next()
}
