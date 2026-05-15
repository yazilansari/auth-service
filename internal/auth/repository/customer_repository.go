package repository

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/database"
)

func FindCustomerByPhone(phone string) (*model.Customer, error) {

	var customer model.Customer

	err := database.DB.
		Where("phone = ?", phone).
		First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}
