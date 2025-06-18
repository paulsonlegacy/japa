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
	"go.uber.org/zap"
)

// GLOBAL VARIABLES
var (
	BASE_PATH string
	ENV_PATH  string
	Validator *validator.Validate = validator.New() // Initialize validator
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
	Validator.RegisterValidation("ulid", func(fl validator.FieldLevel) bool {
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
	zap.L().Debug("Initializing database connection")
	db := db.NewGormDB(cfg.Database)

	// Initialize Mailer
	zap.L().Debug("Initializing mailing service")
	mailer := mailer.NewSMTPMailer(cfg.Email)

	// Initialize app functions
	zap.L().Debug("Initializing repositories")
	userRepo := repository.NewUserRepository(db)
	visaRepo := repository.NewVisaRepository(db)

	zap.L().Debug("Initializing services")
	userUsecase := usecase.NewUserUsecase(userRepo, db, mailer)
	visaUsecase := usecase.NewVisaUsecase(visaRepo, db)

	zap.L().Debug("Initializing handlers")
	userHandler := handlers.NewUserHandler(Validator, userUsecase)
	visaHandler := handlers.NewVisaHandler(Validator, visaUsecase)

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

	zap.L().Debug("Linking http routes")
	app.Get(
		"/protocol", func(c *fiber.Ctx) error {
			return c.SendString(c.Protocol()) // => https
		},
	)

	// Group: /api/v1/
	v1 := app.Group("/api/v1")
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	visaGroup :=  v1.Group("/visa")
	visaGroup.Post("/apply", visaHandler.SubmitVisaApplication)

	// Initialize server
	zap.S().Debugw("Starting server at port ", cfg.Server.ServerAddress, "...")
	log.Fatal(app.Listen(cfg.Server.ServerAddress))
}
