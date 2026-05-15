package model

import "time"

type Tenant struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	TenantCode  string     `gorm:"column:tenant_code"`
	CountryCode string     `gorm:"column:country_code"`
	Name        string     `gorm:"column:name"`
	Domain      string     `gorm:"column:domain"`
	Currency    *string    `gorm:"column:currency"`
	Timezone    *string    `gorm:"column:timezone"`
	Active      bool       `gorm:"column:active"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (Tenant) TableName() string {
	return "tenants"
}
