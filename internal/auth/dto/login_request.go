package dto

type LoginRequest struct {
	Mobile   string `json:"mobile" validate:"required"`
	Password string `json:"password" validate:"required"`
}
