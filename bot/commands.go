package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/source"
	"github.com/ErrorBoi/feedparserbot/ent/user"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
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

func (b *Bot) add(m *tgbotapi.Message) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)

	var msg string
	if len(args) == 0 {
		msg = "После команды нужно указать ссылку на источник. Например /add https://vc.ru/marketing"
	} else {
		switch {
		case strings.HasPrefix(args, "https://vc.ru/"):
			ctx := context.Background()
			hubName := strings.TrimPrefix(args, "https://vc.ru/")
			s, err := b.db.Cli.Source.Query().Where(source.URL("https://vc.ru/rss/" + hubName)).Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying source: %v", err)
				if ent.IsNotFound(err) {
					msg = "Раздел VC.ru с таким названием не найден"
				}
			} else {
				_, err = b.db.Cli.User.Update().AddSources(s).Save(ctx)
				if err != nil {
					b.lg.Errorf("failed updating user: %v", err)
				}
				msg = "Подписка оформлена"
			}
		case strings.HasPrefix(args, "http://rusbase.com/feeds/tag/"):
			ctx := context.Background()
			s, err := b.db.Cli.Source.Query().Where(source.URL(args)).Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying source: %v", err)
				if ent.IsNotFound(err) {
					msg = "Раздел RB.ru с таким названием не найден"
				}
			} else {
				parent, err := s.QueryParent().Only(ctx)
				if parent != nil {
					_, err = b.db.Cli.User.Update().RemoveSources(parent).Save(ctx)
					if err != nil {
						b.lg.Errorf("failed updating user: %v", err)
					}
				}

				_, err = b.db.Cli.User.Update().AddSources(s).Save(ctx)
				if err != nil {
					b.lg.Errorf("failed updating user: %v", err)
				}
				msg = "Подписка оформлена"
			}
		default:
			ctx := context.Background()
			s, err := b.db.Cli.Source.Query().Where(source.URL(args)).Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying source: %v", err)
				if ent.IsNotFound(err) {
					msg = "Источник с таким названием не найден"
				}
			} else {
				_, err = b.db.Cli.User.Update().AddSources(s).Save(ctx)
				if err != nil {
					b.lg.Errorf("failed updating user: %v", err)
				}
				msg = "Подписка оформлена"
			}
		}
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) urgent(m *tgbotapi.Message) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)

	var msg string
	if len(args) == 0 {
		msg = "После команды нужно указать список слов через запятую. Например:\n" +
			"/urgent ЛО,Санкт-Петербург,Доллар"
	} else {
		arr := strings.Split(args, ",")
		ctx := context.Background()

		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(m.From.ID))).
			SetUrgentWords(arr).
			Save(ctx)
		if err != nil {
			b.lg.Fatalf("failed updating user settings: %v", err)
		}
		msg = "\"Срочные\" слова записаны!"
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) bannedWords(m *tgbotapi.Message) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)

	var msg string
	if len(args) == 0 {
		msg = "После команды нужно указать список слов через запятую. Например:\n/banned Коронавирус,Поправки"
	} else {
		arr := strings.Split(args, ",")
		ctx := context.Background()

		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(m.From.ID))).
			SetBannedWords(arr).
			Save(ctx)
		if err != nil {
			b.lg.Fatalf("failed updating user settings: %v", err)
		}
		msg = "\"Чёрный список\" обновлён!"
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}
