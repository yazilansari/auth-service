package main

import (
	authMiddleware "auth-service/internal/auth/middleware"
	authRoutes "auth-service/internal/auth/routes"
	"auth-service/internal/database"
	"auth-service/internal/queue"
	redisClient "auth-service/internal/redis"

	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	fiberCors "github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	fiberRequestID "github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/joho/godotenv"
)

func main() {

	// Load ENV

	err := godotenv.Load()

	if err != nil {
		log.Fatal(".env file not loaded")
	}

	// Connect PostgreSQL

	database.ConnectPostgres()

	log.Println("✅ PostgreSQL Connected")

	// Connect Redis

	redisClient.ConnectRedis()

	redisClient.InitFiberStorage()

	queue.InitQueue()

	log.Println("✅ Redis Connected")

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

	app.Use(fiberRequestID.New())

	// Logger

	app.Use(fiberLogger.New())

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

	log.Println("🚀 Auth Service Running On Port:", port)

	log.Fatal(app.Listen(":" + port))
}
