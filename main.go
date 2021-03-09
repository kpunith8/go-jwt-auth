package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/kpunith8/go-jwt-auth/database"
	"github.com/kpunith8/go-jwt-auth/routes"
)

func main() {
	client := database.Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close connection at the end of usage, otherwise connection would be leaked
	defer client.Disconnect(ctx)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":3000")
}
