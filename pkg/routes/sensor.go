package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/kiprop-dave/2faAuth/pkg/controllers"
	middle "github.com/kiprop-dave/2faAuth/pkg/middleware"
)

func SensorRoute(app *fiber.App) {
	app.Post("/sensor", middle.VerifyToken, controllers.CreateSensor)
}
