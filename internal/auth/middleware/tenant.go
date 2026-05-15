package middleware

import (
	tenantModel "auth-service/internal/tenant/model"
	tenantRepo "auth-service/internal/tenant/repository"
	tenantService "auth-service/internal/tenant/service"

	"github.com/gofiber/fiber/v2"
)

func TenantMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {

		host := c.Hostname()

		// Local Development Fallback

		if host == "localhost" ||
			host == "127.0.0.1" {

			c.Locals("tenant_code", "AE")
			c.Locals("country_code", "AE")

			return c.Next()
		}

		cacheKey := "tenant:" + host

		var tenant tenantModel.Tenant

		// Try Redis Cache

		err := tenantService.GetTenantCache(
			cacheKey,
			&tenant,
		)

		if err == nil {

			c.Locals("tenant_code", tenant.TenantCode)
			c.Locals("country_code", tenant.CountryCode)

			return c.Next()
		}

		// DB Fallback

		tenantData, err := tenantRepo.FindTenantByDomain(host)

		if err != nil {

			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid tenant domain",
			})
		}

		// Cache Tenant

		tenantService.SetTenantCache(
			cacheKey,
			tenantData,
		)

		c.Locals("tenant_code", tenantData.TenantCode)
		c.Locals("country_code", tenantData.CountryCode)

		return c.Next()
	}
}
