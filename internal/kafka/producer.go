package kafka

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"auth-service/internal/logger"

	kafkaGo "github.com/segmentio/kafka-go"

	"go.uber.org/zap"
)

var OTPWriter *kafkaGo.Writer

var WhatsAppWriter *kafkaGo.Writer

var EmailWriter *kafkaGo.Writer

func InitKafkaProducer() {

	logger.Log.Info(
		"initializing kafka producers",
	)

	OTPWriter = &kafkaGo.Writer{
		Addr: kafkaGo.TCP(
			os.Getenv("KAFKA_BROKER"),
		),

		Topic: OTPTopic,

		Balancer: &kafkaGo.LeastBytes{},
	}

	WhatsAppWriter = &kafkaGo.Writer{
		Addr: kafkaGo.TCP(
			os.Getenv("KAFKA_BROKER"),
		),

		Topic: WhatsAppTopic,

		Balancer: &kafkaGo.LeastBytes{},
	}

	EmailWriter = &kafkaGo.Writer{
		Addr: kafkaGo.TCP(
			os.Getenv("KAFKA_BROKER"),
		),

		Topic: EmailTopic,

		Balancer: &kafkaGo.LeastBytes{},
	}

	logger.Log.Info(
		"kafka producers initialized",
	)
}

// =========================
// OTP PRODUCER
// =========================

func PublishOTPEvent(
	phone string,
	otp string,
) error {

	payload := OTPPayload{
		Phone: phone,
		OTP:   otp,
	}

	return publish(
		OTPWriter,
		phone,
		payload,
	)
}

// =========================
// WHATSAPP PRODUCER
// =========================

func PublishWhatsAppEvent(
	phone string,
	otp string,
) error {

	payload := WhatsAppPayload{
		Phone: phone,
		OTP:   otp,
	}

	return publish(
		WhatsAppWriter,
		phone,
		payload,
	)
}

// =========================
// EMAIL PRODUCER
// =========================

func PublishEmailEvent(
	to string,
	subject string,
	body string,
) error {

	payload := EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	return publish(
		EmailWriter,
		to,
		payload,
	)
}

// =========================
// COMMON PUBLISHER
// =========================

func publish(
	writer *kafkaGo.Writer,
	key string,
	payload interface{},
) error {

	data, err := json.Marshal(
		payload,
	)

	if err != nil {

		logger.Log.Error(
			"failed to marshal kafka payload",

			zap.Error(err),
		)

		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	err = writer.WriteMessages(
		ctx,

		kafkaGo.Message{
			Key:   []byte(key),
			Value: data,
			Time:  time.Now(),
		},
	)

	if err != nil {

		logger.Log.Error(
			"failed to publish kafka event",

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"kafka event published",

		zap.String(
			"key",
			key,
		),
	)

	return nil
}
