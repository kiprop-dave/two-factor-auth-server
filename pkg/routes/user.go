package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/kiprop-dave/2faAuth/pkg/controllers"
)

func UserRoute(app *fiber.App) {
	app.Post("/user", controllers.CreateUser)
	app.Get("/user/:id", controllers.GetUser)
	app.Delete("/user/:id", controllers.DeleteUser)
	app.Patch("/user/:id", controllers.UpdateUser)
	app.Get("/users", controllers.GetAllUsers)

	app.Post("/admin", controllers.CreateAdmin)
	app.Delete("/admin/:name", controllers.DeleteAdmin)
}
