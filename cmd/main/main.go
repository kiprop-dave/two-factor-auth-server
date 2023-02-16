package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	config "github.com/kiprop-dave/2faAuth/pkg/config"
	routes "github.com/kiprop-dave/2faAuth/pkg/routes"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	config.ConnectToMongo()
	routes.UserRoute(app)
	routes.AuthRoute(app)

	app.Listen(":3000")
}
