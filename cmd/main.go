package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"japa/internal/config"
	"japa/internal/app/http/handler"
	"japa/internal/infrastructure/db"
	"japa/internal/infrastructure/logging"
	"japa/internal/infrastructure/mail"
	"japa/internal/domain/repository"
	"japa/internal/domain/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/oklog/ulid/v2"
)

// GLOBAL VARIABLES
var (
	BASE_PATH string
	ENV_PATH  string
	Val       *validator.Validate = validator.New() // Initialize validator
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
	ENV_PATH = filepath.Join(BASE_PATH, "internal/config/.env") // ENV file

	// Register custom validator for ulid inputs
	Val.RegisterValidation("ulid", func(fl validator.FieldLevel) bool {
		_, err := ulid.Parse(fl.Field().String())
		return err == nil
	})
}

// Application entry point
func main() {
	// Initialize configurations
	log.Println(ENV_PATH)
	cfg := config.InitConfig(ENV_PATH)

	// Initialize logger
	logger := logging.InitLogger(cfg.Logging)
	defer logger.Sync()

	// Initialize DB
	db := db.NewGormDB(cfg.Database)

	// Initialize Mailer
	mailer := mailer.NewSMTPMailer(cfg.Email)

	// Initialize user functions
	userRepo := repository.NewUserRepository(db)
	UserUsecase := usecase.NewUserUsecase(userRepo, db, mailer)
	userHandler := handlers.NewUserHandler(Val, UserUsecase)

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
		return c.SendString("This is home endpoint")
	})

	// Group: /api/v1/
	v1 := app.Group("/api/v1")
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	// Start the server
	log.Fatal(app.Listen(address))
}
