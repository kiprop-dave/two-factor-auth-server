package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	routes "github.com/kiprop-dave/2faAuth/pkg/routes"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	routes.UserRoute(app)
	routes.SensorRoute(app)
	routes.AuthRoute(app)

	app.Listen(":3000")
}
