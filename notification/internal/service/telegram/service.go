package telegram

import (
	"bytes"
	"context"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/config"
	"text/template"
	"time"

	"go.uber.org/zap"

	"github.com/crafty-ezhik/rocket-factory/notification/internal/client/http"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/model"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/service/telegram/templates"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

type orderPaidTemplateData struct {
	OrderUUID       string
	TransactionUUID string
	PaymentMethod   string
	PaymentDate     string
}

type orderAssembledTemplateData struct {
	OrderUUID    string
	BuildTimeSec int
}

var (
	orderPaidTemplate      = template.Must(template.ParseFS(templates.FS, "order_paid_notification.tmpl"))
	orderAssembledTemplate = template.Must(template.ParseFS(templates.FS, "order_assembled_notification.tmpl"))
)

type service struct {
	tgClient http.TelegramClient
}

func NewService(tgClient http.TelegramClient) *service {
	return &service{
		tgClient: tgClient,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, msg model.OrderPaidEvent) error {
	chatID := config.AppConfig().TgBot.ChatID()
	message, err := s.buildPaidMsg(msg)
	if err != nil {
		return err
	}

	err = s.tgClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}
	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) SendOrderAssembledNotification(ctx context.Context, msg model.OrderAssembledEvent) error {
	chatID := config.AppConfig().TgBot.ChatID()

	message, err := s.buildAssembledMsg(msg)
	if err != nil {
		return err
	}
	err = s.tgClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) buildPaidMsg(msg model.OrderPaidEvent) (string, error) {
	data := orderPaidTemplateData{
		OrderUUID:       msg.OrderUUID.String(),
		TransactionUUID: msg.TransactionUUID.String(),
		PaymentMethod:   msg.PaymentMethod,
		PaymentDate:     time.Now().Format(time.DateTime),
	}

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *service) buildAssembledMsg(msg model.OrderAssembledEvent) (string, error) {
	data := orderAssembledTemplateData{
		OrderUUID:    msg.OrderUUID.String(),
		BuildTimeSec: msg.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := orderAssembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
