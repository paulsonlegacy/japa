package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"japa/internal/config"
	"japa/internal/infrastructure/db"
	"japa/internal/repository"
	"japa/internal/services"
	"japa/internal/handlers"
	"japa/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/go-playground/validator/v10"
)

var (
	BASE_PATH string
	ENV_PATH string
)


// init() runs before main function is executed
func init() {
	// init() Usage – Prepares paths and server variables before main()

	// Setting BASE path from this file
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	BASE_PATH = dir // Sets the root directory dynamically


	// Setting env path
	ENV_PATH  = filepath.Join(BASE_PATH, "internal/config/.env") // ENV file 
}


// Application entry point
func main() {
	// Initialize configurations
	log.Println(ENV_PATH)
	cfg := config.InitConfig(ENV_PATH)
	
	// Initialize logger
	logger := logging.InitLogger(cfg.LoggingConfig)
	defer logger.Sync()

	// Initialize DB
	db := db.NewGormDB(cfg.Database)

	// Initialize validator
	val := validator.New()

	// Initialize user functions
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, db)
	userHandler := handlers.NewUserHandler(val, userService)

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
		log.Println("About to test logger");
		logging.Logger.Info("Testing Logger hehe!")
		return c.SendString("This is home endpoint")
	})

	// Group: /api/v1
	v1 := app.Group("/api/v1")
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	// Start the server
	log.Fatal(app.Listen(address))
}