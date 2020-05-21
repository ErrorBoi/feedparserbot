package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ErrorBoi/feedparserbot/ent/source"
)

func (b *Bot) start(m *tgbotapi.Message) {
	ctx := context.Background()

	ss, err := b.db.Cli.Source.Query().Where(source.Not(source.HasParent())).All(ctx)
	if err != nil {
		b.lg.Errorf("failed querying parent sources: %v", err)
	}

	u, err := b.db.Cli.User.
		Create().
		SetTgID(m.From.ID).
		AddSources(ss...).
		Save(ctx)
	if err != nil {
		b.lg.Errorf("failed creating user: %v", err)
	} else {
		b.lg.Info("user was created: ", u)
	}

	us, err := b.db.Cli.UserSettings.
		Create().
		SetUser(u).
		Save(ctx)
	if err != nil {
		b.lg.Errorf("failed creating user settings: %v", err)
	} else {
		b.lg.Info("user settings were created: ", us)
	}

	msg := tgbotapi.NewMessage(m.Chat.ID, startText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}

func (b *Bot) help(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, helpText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}
