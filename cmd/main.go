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
	"japa/internal/app/http/middleware"

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
	logger := logging.InitLogger(cfg.LoggingConfig)
	defer logger.Sync()

	// Initialize DB
	zap.L().Debug("Initializing database connection")
	db := db.NewGormDB(cfg.DBConfig)

	// Initialize mailing providers
	zap.L().Debug("Initializing mailing providers")
	smtpMailer := mailer.NewSMTPMailer(
		cfg.ServerConfig,
		cfg.SiteConfig,
		cfg.EmailConfig.SMTPConfig,
		logger,
	)
	mailer := &mailer.ResponsiveMailer{
		Providers: []mailer.Mailer{
			//sendgridMailer, // first try
			//mailgunMailer,  // fallback if SendGrid fails
			smtpMailer,     // fallback if Mailgun fails
		},
	}


	// Initialize app functions
	zap.L().Debug("Initializing repositories")
	userRepo := repository.NewUserRepository(db)
	visaRepo := repository.NewVisaRepository(db)
	postRepo := repository.NewPostRepository(db)

	zap.L().Debug("Initializing services")
	userUsecase := usecase.NewUserUsecase(cfg.JWTConfig, userRepo, db, mailer)
	visaUsecase := usecase.NewVisaUsecase(visaRepo, db)
	postUsecase := usecase.NewPostUsecase(postRepo, db)

	zap.L().Debug("Initializing handlers")
	userHandler := handlers.NewUserHandler(Validator, userUsecase)
	visaHandler := handlers.NewVisaHandler(Validator, visaUsecase)
	postHandler := handlers.NewPostHandler(Validator, postUsecase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.ServerConfig, cfg.JWTConfig, db).Handler()

	// Setup server
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

	zap.L().Debug("Linking http routes..")

	// API prefix (base route)
	v1 := app.Group("/api/v1")

	// Public routes
	v1.Get(
		"/protocol", func(c *fiber.Ctx) error {
			return c.SendString(c.Protocol()) // => https
		},
	)
	v1.Post("/register", userHandler.Register)
	v1.Post("/login",    userHandler.Login)
	v1.Get("/posts",     postHandler.FetchPosts) // /api/v1/posts?page=2&limit=20
	//v1.Get("/posts/:id", postHandler.GetPost)

	// Authenticated routes
	accountGroup := v1.Group("/account")
	accountGroup.Use(authMiddleware)

	// Visa routes (authenticated)
	visaGroup :=  accountGroup.Group("/visa")
	visaGroup.Post("/apply", visaHandler.SubmitVisaApplication)

	// Agent routes (authenticated)
	agentGroup := v1.Group("/agent")
	agentGroup.Use(middleware.AgentOnly())

	//agentGroup.Get("/dashboard", agentHandler.GetDashboard)

	// Admin routes (authenticated)
	adminGroup := accountGroup.Group("/admin")
	adminGroup.Use(middleware.AdminOnly())

	adminGroup.Post("/posts/create", postHandler.CreatePost)

	// SuperAdmin routes (authenticated)
	superAdminGroup := accountGroup.Group("/superadmin")
	superAdminGroup.Use(middleware.SuperadminOnly())

	//superAdminGroup.Get("/dashboard", superAdminHandler.GetDashboard)

	// Initialize server
	zap.S().Debugw("Starting server at port ", cfg.ServerConfig.ServerAddress, "...")
	log.Fatal(app.Listen(cfg.ServerConfig.ServerAddress))
}
