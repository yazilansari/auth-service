package dto

type SendOTPRequest struct {
	Mobile string `json:"mobile" validate:"required"`
	Flag   string `json:"flag"`
}
