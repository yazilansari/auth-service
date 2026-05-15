package dto

type VerifyOTPRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`

	OTP  string `json:"otp"`
	Flag string `json:"flag"`
}
