package model

import "time"

type Customer struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	Name        string     `gorm:"column:name"`
	Email       string     `gorm:"column:email"`
	Password    string     `gorm:"column:password"`
	Phone       string     `gorm:"column:phone"`
	Status      string     `gorm:"column:status"`
	CountryCode string     `gorm:"column:country_code"`
	TenantCode  string     `gorm:"column:tenant_code"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (Customer) TableName() string {
	return "customers"
}
