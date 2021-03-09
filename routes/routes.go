package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kpunith8/go-jwt-auth/controllers"
)

// Setup - Setup routes
func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Get("/api/users", controllers.AllUsers)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
}
