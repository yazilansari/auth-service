package repository

import (
	"auth-service/internal/database"
	"auth-service/internal/tenant/model"
)

func FindTenantByDomain(domain string) (*model.Tenant, error) {

	var tenant model.Tenant

	err := database.DB.
		Where("domain = ?", domain).
		Where("active = ?", true).
		First(&tenant).Error

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}
