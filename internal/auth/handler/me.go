package handler

import (
	"fmt"

	"auth-service/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Customer struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func Me(c *fiber.Ctx) error {

	// Get JWT token from middleware
	token := c.Locals("user").(*jwt.Token)

	// Extract claims
	claims := token.Claims.(jwt.MapClaims)

	// Get customer_id from token
	customerID := uint64(claims["customer_id"].(float64))

	fmt.Println("Customer ID:", customerID)

	// Customer response object
	var customer Customer

	// Query database
	err := database.DB.Table("customers").
		Select("id, name, email, phone").
		Where("id = ?", customerID).
		First(&customer).Error

	if err != nil {
		fmt.Println("DB Error:", err)

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Customer not found",
		})
	}

	// Return customer data
	return c.JSON(customer)
}
