package main

import (
	authMiddleware "auth-service/internal/auth/middleware"
	authRoutes "auth-service/internal/auth/routes"
	"auth-service/internal/database"
	"auth-service/internal/logger"
	requestMiddleware "auth-service/internal/middleware"
	"auth-service/internal/queue"
	redisClient "auth-service/internal/redis"

	"os"

	"github.com/gofiber/fiber/v2"

	fiberCors "github.com/gofiber/fiber/v2/middleware/cors"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	// Init Logger

	logger.InitLogger()

	defer logger.Log.Sync()

	// Load ENV

	err := godotenv.Load()

	if err != nil {
		logger.Log.Fatal(".env file not loaded")
	}

	// Connect PostgreSQL

	database.ConnectPostgres()

	logger.Log.Info("✅ PostgreSQL Connected")

	// Connect Redis

	redisClient.ConnectRedis()

	redisClient.InitFiberStorage()

	queue.InitQueue()

	logger.Log.Info("✅ Redis Connected")

	// Fiber App

	app := fiber.New(fiber.Config{
		AppName: "Ahmed Auth Service",
	})

	// =========================
	// Global Middleware
	// =========================

	// Panic Recovery

	app.Use(fiberRecover.New())

	// Request ID

	app.Use(requestMiddleware.RequestContextMiddleware())

	app.Use(requestMiddleware.RequestLoggerMiddleware())

	// CORS

	app.Use(fiberCors.New(fiberCors.Config{
		AllowOrigins: "http://localhost:3000, https://ae.ahmedalmaghribi.com, https://ksa.ahmedalmaghribi.com, https://qa.ahmedalmaghribi.com, https://kw.ahmedalmaghribi.com, https://bh.ahmedalmaghribi.com, https://om.ahmedalmaghribi.com",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Tenant Middleware

	app.Use(authMiddleware.TenantMiddleware())

	// =========================
	// Health Check
	// =========================

	app.Get("/health", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	// =========================
	// Routes
	// =========================

	authRoutes.SetupAuthRoutes(app)

	// =========================
	// Start Server
	// =========================

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	logger.Log.Info("🚀 Auth Service Running",
		zap.String("port", port),
	)

	logger.Log.Fatal("server stopped",
		zap.Error(app.Listen(":"+port)),
	)
}
