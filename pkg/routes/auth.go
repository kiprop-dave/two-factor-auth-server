package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/kiprop-dave/2faAuth/pkg/controllers"
	"github.com/kiprop-dave/2faAuth/pkg/middleware"
)

func AuthRoute(app *fiber.App) {
	app.Post("/auth/admin", controllers.AuthAdmin)
	app.Post("/auth/user", controllers.AuthUser)
	app.Get("/auth/attempts", middleware.VerifyToken, controllers.GetAllAttempts)
}
