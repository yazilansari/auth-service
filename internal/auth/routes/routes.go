package routes

import (
	"auth-service/internal/auth/handler"
	"auth-service/internal/auth/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {

	api := app.Group("/api/v1/auth")

	api.Post("/signin", handler.Login)

	api.Post("/sendOTP",
		middleware.OTPRateLimiter(),
		handler.SendOTP,
	)

	api.Post("/verifyOTP",
		middleware.LoginRateLimiter(),
		handler.VerifyOTP,
	)

	api.Post("/refresh",
		middleware.RefreshRateLimiter(),
		handler.Refresh,
	)

	api.Post(
		"/signout",
		middleware.Protected(),
		handler.Logout,
	)

	api.Get(
		"/me",
		middleware.Protected(),
		handler.Me,
	)
}
