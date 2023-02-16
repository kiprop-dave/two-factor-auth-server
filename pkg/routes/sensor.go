package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/kiprop-dave/2faAuth/pkg/controllers"
)

func SensorRoute(app *fiber.App) {
	app.Post("/sensor", controllers.CreateSensor)
}
