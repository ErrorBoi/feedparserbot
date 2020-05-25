package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/source"
	"github.com/ErrorBoi/feedparserbot/ent/user"
)

func (b *Bot) start(m *tgbotapi.Message) {
	ctx := context.Background()

	_, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).Only(ctx)
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			// Subscribe user to all parent sources
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
		default:
			b.lg.Errorf("failed querying user: %v", err)
		}
	}

	msg := tgbotapi.NewMessage(m.Chat.ID, startText)
	msg.ReplyMarkup = mainKeyboard
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}

func (b *Bot) help(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, helpText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}
