package model

import "time"

type RefreshToken struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	CustomerID uint64 `gorm:"column:customer_id"`

	TokenHash string `gorm:"column:token_hash"`

	TokenFamilyID string `gorm:"column:token_family_id"`

	DeviceID string `gorm:"column:device_id"`

	IPAddress string `gorm:"column:ip_address"`

	UserAgent string `gorm:"column:user_agent"`

	ExpiresAt time.Time `gorm:"column:expires_at"`

	Revoked bool `gorm:"column:revoked"`

	CreatedAt *time.Time `gorm:"column:created_at"`

	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
