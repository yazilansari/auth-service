package kafka

type OTPPayload struct {
	Phone string `json:"phone"`

	OTP string `json:"otp"`
}

type WhatsAppPayload struct {
	Phone string `json:"phone"`

	OTP string `json:"otp"`
}

type EmailPayload struct {
	To string `json:"to"`

	Subject string `json:"subject"`

	Body string `json:"body"`
}
