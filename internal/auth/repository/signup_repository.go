package repository

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/database"
)

func CreateCustomer(customer *model.Customer) error {

	return database.DB.Create(customer).Error
}

func FindCustomerByPhoneOrEmail(
	phone string,
	email string,
) (*model.Customer, error) {

	var customer model.Customer

	err := database.DB.
		Where("phone = ?", phone).
		Or("email = ?", email).
		First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}
