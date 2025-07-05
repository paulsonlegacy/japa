package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"context"
	"path/filepath"
	"os/signal"
	"syscall"

	"japa/internal/config"
	"japa/internal/app/http/handler"
	"japa/internal/infrastructure/db"
	"japa/internal/infrastructure/logging"
	"japa/internal/infrastructure/mail"
	"japa/internal/infrastructure/scraper"
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


func main() { // Application entry point

	// Initialize configurations
	log.Println(ENV_PATH)
	cfg := config.InitConfig(ENV_PATH)

	// Initialize logger
	logger := logging.InitLogger(cfg.LoggingConfig)
	defer logger.Sync()

	// Initialize DB & DB models
	zap.L().Debug("Initializing database connection")
	db := db.NewGormDB(cfg.DBConfig) // Runs migrations internally

	// Global context and graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for interrupt signals (Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Launch a goroutine that waits for shutdown signal
	go func() {
		sig := <-quit
		zap.L().Info("Shutdown signal received", zap.String("signal", sig.String()))
		cancel() // This will cancel ctx and tell all background tasks to exit
	}()

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

	// Initialize scrapers
	japacontentScraper := &scraper.JapaContentScraper{
		Logger: logger,
		DB:     db,
	}
	multiScraper := &scraper.MultiScraper{
		Scrapers: []scraper.Scraper{
			japacontentScraper,
		},
		Logger:   logger,
		Interval: 7 * 24 * time.Hour,
	}

	// Start the scraper with 
	// the same context app uses
	multiScraper.Run(ctx)

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
	v1.Post("/auth/register", userHandler.Register)
	v1.Post("/auth/login", userHandler.Login)
	v1.Get("/auth/logout", userHandler.Logout)
	v1.Post("/auth/refresh", userHandler.RefreshToken)
	v1.Get("/posts",     postHandler.FetchPosts) // api/v1/posts?page=2&limit=20
	v1.Get("/posts/:post_id/:slug", postHandler.FetchPost)  // posts/01JXYZM4T8HR8PQKJS6E4X2C1Z/seo-tips-for-developers


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

	// Initialize Fiber server in background
	zap.S().Debugw("Starting server at port ", cfg.ServerConfig.ServerAddress, "...")
	go func() {
		if err := app.Listen(cfg.ServerConfig.ServerAddress); err != nil {
			zap.L().Fatal("Fiber server error", zap.Error(err))
		}
	}()

	// Wait for shutdown
	<-ctx.Done()
	zap.L().Info("Shutting down services...")
	// Gracefully shut down Fiber server
	if err := app.Shutdown(); err != nil {
		zap.L().Error("Error shutting down Fiber", zap.Error(err))
	}
}
