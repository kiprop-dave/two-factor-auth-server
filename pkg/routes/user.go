package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/kiprop-dave/2faAuth/pkg/controllers"
	middle "github.com/kiprop-dave/2faAuth/pkg/middleware"
)

func UserRoute(app *fiber.App) {
	app.Post("/user", middle.VerifyToken, controllers.CreateUser)
	app.Get("/user/:id", middle.VerifyToken, controllers.GetUser)
	app.Delete("/user/:id", middle.VerifyToken, controllers.DeleteUser)
	app.Patch("/user/:id", middle.VerifyToken, controllers.UpdateUser)
	app.Get("/users", middle.VerifyToken, controllers.GetAllUsers)

	app.Post("/admin", controllers.CreateAdmin)
	app.Delete("/admin/:name", middle.VerifyToken, controllers.DeleteAdmin)
}
