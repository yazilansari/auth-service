package service

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/auth/repository"
	"auth-service/internal/jwt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Signup(
	tenantCode string,
	countryCode string,
	name string,
	email string,
	phone string,
	password string,
) (string, error) {

	existing, _ := repository.FindCustomerByPhoneOrEmail(
		phone,
		email,
	)

	if existing != nil {
		return "", ErrCustomerExists
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", errors.New("password hashing failed")
	}

	customer := model.Customer{
		Name:        name,
		Email:       email,
		Phone:       phone,
		Password:    string(hash),
		Status:      "activated",
		TenantCode:  tenantCode,
		CountryCode: countryCode,
	}

	err = repository.CreateCustomer(&customer)

	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(
		customer.ID,
		tenantCode,
		countryCode,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}
