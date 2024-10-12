package telegram

import (
	"context"
	"strconv"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

type Notifier struct {
	telService *telegram.Telegram
}

func NewNotifier(botAPIToken, telegramUserID string) *Notifier {
	telegramService, err := telegram.New(botAPIToken)
	if err != nil {
		panic(err)
	}
	userTelegramID, err := strconv.ParseInt(telegramUserID, 10, 64)
	if err != nil {
		panic(err)
	}
	telegramService.AddReceivers(userTelegramID)
	notify.UseServices(telegramService)

	return &Notifier{
		telService: telegramService,
	}
}

func (t *Notifier) Notify(subject, msg string) error {
	err := notify.Send(
		context.Background(),
		subject,
		msg,
	)
	if err != nil {
		return err
	}

	return nil
}
