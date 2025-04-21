package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(
		fiber.Config{
			EnablePrintRoutes: true,
			JSONEncoder:       json.Marshal,
			JSONDecoder:       json.Unmarshal,
		},
	)

	app.Use(
		cors.New(
			cors.Config{
				AllowOrigins: "*",
				AllowMethods: "GET, POST, DELETE",
			},
		),
	)
	
	address := ":8080"

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("This is default home")
	})

	// Group: /api/v1
	v1 := app.Group("/api/v1")
	v1.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	// Group: /api/v2
	v2 := app.Group("/api/v2")
	v2.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("This is v2 home")
	})

	v2.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		log.Println("User ID:", id)
		return c.SendString("User ID is: " + id)
	})

	// Start the server
	log.Fatal(app.Listen(address))
}