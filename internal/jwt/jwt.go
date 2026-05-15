package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(
	customerID uint64,
	tenantCode string,
	countryCode string,
) (string, error) {

	claims := jwt.MapClaims{
		"customer_id":  customerID,
		"tenant_code":  tenantCode,
		"country_code": countryCode,
		"type":         "access",
		"exp":          time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken(
	customerID uint64,
	tenantCode string,
	countryCode string,
) (string, error) {

	claims := jwt.MapClaims{
		"customer_id":  customerID,
		"tenant_code":  tenantCode,
		"country_code": countryCode,
		"type":         "refresh",
		"exp": time.Now().
			Add(time.Hour * 24 * 30).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}
