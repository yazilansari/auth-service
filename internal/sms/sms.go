package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func normalizePhone(phone string) string {

	phone = strings.TrimSpace(phone)

	phone = strings.TrimPrefix(phone, "+")

	// UAE Format
	// 0501234567 -> 971501234567

	if strings.HasPrefix(phone, "0") {

		phone = "971" + phone[1:]
	}

	return phone
}

func SendOTP(
	phone string,
	otp string,
) error {

	normalizedPhone := normalizePhone(phone)

	// =========================
	// SMS API
	// =========================

	password := os.Getenv("MIM_SMS_PASSWORD")

	smsURL := fmt.Sprintf(
		"https://myinboxmedia.ae/api/mim/SendSMS?userid=%s&pwd=%s&mobile=%s&sender=%s&msg=%s&msgtype=16",
		url.QueryEscape(os.Getenv("MIM_SMS_USER_ID")),
		url.QueryEscape(password),
		url.QueryEscape(normalizedPhone),
		url.QueryEscape(os.Getenv("MIM_SMS_SENDER")),
		url.QueryEscape(
			otp+" is your OTP for Registration",
		),
	)

	smsClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	smsReq, err := http.NewRequest(
		http.MethodGet,
		smsURL,
		nil,
	)

	if err != nil {
		return err
	}

	smsResp, err := smsClient.Do(
		smsReq,
	)

	if err != nil {

		log.Println(
			"SMS API Error:",
			err.Error(),
		)

		return err
	}

	defer smsResp.Body.Close()

	smsBody, _ := io.ReadAll(
		smsResp.Body,
	)

	log.Println(
		"SMS API Response:",
		string(smsBody),
	)

	// =========================
	// WHATSAPP API
	// =========================

	waPayload := map[string]interface{}{
		"ProfileId": os.Getenv("MIM_WA_PROFILE_ID"),

		"APIKey": os.Getenv("MIM_WA_API_KEY"),

		"MobileNumber": normalizedPhone,

		"templateName": "websiteauthentication",

		"Parameters": []string{
			otp,
		},

		"HeaderType": "Text",

		"Text": "",

		"MediaUrl": "",

		"Latitude": 0,

		"Longitude": 0,

		"isTemplate": "true",

		"ButtonOrListJSON": "",

		"SubClientCode": "",

		"HeaderParameter": "",

		"CTAButtonURLParameter": "",

		"CTAButtonURLParameter2": "",
	}

	jsonPayload, err := json.Marshal(
		waPayload,
	)

	if err != nil {
		return err
	}

	waReq, err := http.NewRequest(
		http.MethodPost,
		"https://waba.myinboxmedia.in/api/sendwaba",
		bytes.NewBuffer(jsonPayload),
	)

	if err != nil {
		return err
	}

	waReq.Header.Set(
		"Content-Type",
		"application/json",
	)

	waClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	waResp, err := waClient.Do(
		waReq,
	)

	if err != nil {

		log.Println(
			"WhatsApp API Error:",
			err.Error(),
		)

		return err
	}

	defer waResp.Body.Close()

	waBody, _ := io.ReadAll(
		waResp.Body,
	)

	log.Println(
		"WhatsApp API Response:",
		string(waBody),
	)

	return nil
}
