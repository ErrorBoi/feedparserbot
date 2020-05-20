package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) start(m *tgbotapi.Message) {
	u, err := b.cli.User.Create().SetTgID(m.From.ID).Save(context.Background())
	if err != nil {
		b.lg.Errorf("failed creating user: %v", err)
	}
	b.lg.Info("user was created: ", u)

	msg := tgbotapi.NewMessage(m.Chat.ID, startText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}

func (b *Bot) help(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, helpText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}
