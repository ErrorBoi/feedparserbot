package bot

import (
	"context"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/globalsettings"
	"github.com/ErrorBoi/feedparserbot/ent/source"
	"github.com/ErrorBoi/feedparserbot/ent/user"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
)

func (b *Bot) start(m *tgbotapi.Message, language string) {
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

	msg := tgbotapi.NewMessage(m.Chat.ID, startText[language])
	msg.ReplyMarkup = mainKeyboard[language]
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}

func (b *Bot) help(m *tgbotapi.Message, language string) {
	ctx := context.Background()
	u, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).Only(ctx)
	if err != nil {
		b.lg.Errorf("failed querying user: %v", err)
	}

	var text string
	switch u.Role {
	case user.RoleUser:
		text = helpText[language]
	case user.RoleEditor:
		text = helpText[language] + helpEditorText[language]
	case user.RoleAdmin:
		text = helpText[language] + helpEditorText[language] + helpAdminText[language]
	}

	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}

func (b *Bot) add(m *tgbotapi.Message, language string) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)

	var msg string
	if len(args) == 0 {
		msg = EmptyAddArgsMessage[language]
	} else {
		switch {
		case strings.HasPrefix(args, "https://vc.ru/"):
			ctx := context.Background()
			hubName := strings.TrimPrefix(args, "https://vc.ru/")
			s, err := b.db.Cli.Source.Query().Where(source.URL("https://vc.ru/rss/" + hubName)).Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying source: %v", err)
				if ent.IsNotFound(err) {
					msg = VCHubNotFoundMessage[language]
				}
			} else {
				_, err = b.db.Cli.User.Update().AddSources(s).Save(ctx)
				if err != nil {
					b.lg.Errorf("failed updating user: %v", err)
				}
				msg = SubscriptionCompletedMessage[language]
			}
		case strings.HasPrefix(args, "http://rusbase.com/feeds/tag/"):
			ctx := context.Background()
			s, err := b.db.Cli.Source.Query().Where(source.URL(args)).Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying source: %v", err)
				if ent.IsNotFound(err) {
					msg = RBHubNotFoundMessage[language]
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
				msg = SubscriptionCompletedMessage[language]
			}
		default:
			ctx := context.Background()
			s, err := b.db.Cli.Source.Query().Where(source.URL(args)).Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying source: %v", err)
				if ent.IsNotFound(err) {
					msg = SourceNotFoundMessage[language]
				}
			} else {
				_, err = b.db.Cli.User.Update().AddSources(s).Save(ctx)
				if err != nil {
					b.lg.Errorf("failed updating user: %v", err)
				}
				msg = SubscriptionCompletedMessage[language]
			}
		}
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) urgent(m *tgbotapi.Message, language string) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)

	var msg string
	if len(args) == 0 {
		msg = EmptyUrgentArgsMessage[language]
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
		msg = UrgentWordsSuccessMessage[language]
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) bannedWords(m *tgbotapi.Message, language string) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)

	var msg string
	if len(args) == 0 {
		msg = EmptyBannedArgsMessage[language]
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
		msg = BannedWordsSuccessMessage[language]
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) super(m *tgbotapi.Message, language string) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)
	var msg string
	if len(args) == 0 {
		msg = EmptySuperArgsMessage[language]
	} else {
		if args == Token {
			ctx := context.Background()
			_, err := b.db.Cli.User.Update().Where(user.TgID(m.From.ID)).SetRole(user.RoleAdmin).Save(ctx)
			if err != nil {
				b.lg.Fatalf("failed updating user: %v", err)
			}
			msg = SuperSuccessMessage[language]
		} else {
			msg = SuperValidationErrorMessage[language]
		}
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) setEditor(m *tgbotapi.Message, language string) {
	ctx := context.Background()
	u, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).Only(ctx)
	if err != nil {
		b.lg.Fatalf("failed querying user: %v", err)
	}

	if u.Role == user.RoleAdmin || u.Role == user.RoleEditor {
		args := m.CommandArguments()
		args = strings.TrimSpace(args)
		var msg string
		if len(args) == 0 {
			msg = AddEditorInvalidMessage[language]
		} else {
			tgID, err := strconv.Atoi(args)
			if err != nil {
				b.lg.Fatalf("failed converting str to int: %v", err)
			}
			n, err := b.db.Cli.User.Update().Where(user.TgID(tgID)).SetRole(user.RoleEditor).Save(ctx)
			if err != nil {
				b.lg.Fatalf("failed updating user: %v", err)
			}
			if n == 0 {
				msg = UserNotFoundMessage[language]
			} else {
				msg = AddEditorSuccessMessage[language]
			}
		}
		message := tgbotapi.NewMessage(m.Chat.ID, msg)
		message.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(message)
	}
}

func (b *Bot) removeEditor(m *tgbotapi.Message, language string) {
	ctx := context.Background()
	u, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).Only(ctx)
	if err != nil {
		b.lg.Fatalf("failed querying user: %v", err)
	}

	if u.Role == user.RoleAdmin || u.Role == user.RoleEditor {
		args := m.CommandArguments()
		args = strings.TrimSpace(args)
		var msg string
		if len(args) == 0 {
			msg = RemoveEditorInvalidMessage[language]
		} else {
			tgID, err := strconv.Atoi(args)
			if err != nil {
				b.lg.Fatalf("failed converting str to int: %v", err)
			}
			n, err := b.db.Cli.User.Update().Where(user.TgID(tgID)).SetRole(user.RoleUser).Save(ctx)
			if err != nil {
				b.lg.Fatalf("failed updating user: %v", err)
			}
			if n == 0 {
				msg = UserNotFoundMessage[language]
			} else {
				msg = RemoveEditorSuccessMessage[language]
			}
		}
		message := tgbotapi.NewMessage(m.Chat.ID, msg)
		message.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(message)
	}
}

func (b *Bot) setClickbaitWords(m *tgbotapi.Message, language string) {
	ctx := context.Background()
	u, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).Only(ctx)
	if err != nil {
		b.lg.Fatalf("failed querying user: %v", err)
	}

	if u.Role == user.RoleAdmin || u.Role == user.RoleEditor {
		args := m.CommandArguments()
		args = strings.TrimSpace(args)
		var msg string
		if len(args) == 0 {
			msg = ClickbaitEmptyArgsMessage[language]
		} else {
			arr := strings.Split(args, ",")

			_, err = b.db.Cli.Globalsettings.Update().Where(globalsettings.ID(1)).SetClickbaitWords(arr).Save(ctx)
			if err != nil {
				b.lg.Fatalf("failed updating global settings: %v", err)
			}
			msg = ClickbaitSuccessMessage[language]
		}
		message := tgbotapi.NewMessage(m.Chat.ID, msg)
		message.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(message)
	}
}

func (b *Bot) rewritePost(m *tgbotapi.Message, language string) {
	ctx := context.Background()
	u, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).Only(ctx)
	if err != nil {
		b.lg.Fatalf("failed querying user: %v", err)
	}

	if u.Role == user.RoleAdmin || u.Role == user.RoleEditor {
		args := m.CommandArguments()
		args = strings.TrimSpace(args)
		var msg string
		if len(args) < 2 {
			msg = RewriteEmptyArgsMessage[language]
		} else {
			arr := strings.Split(args, " ")
			postID, err := strconv.Atoi(arr[0])
			if err != nil {
				b.lg.Fatalf("failed converting str to int: %v", err)
			}
			subject := strings.Join(arr[1:], " ")

			err = b.db.RewritePost(ctx, postID, subject, m.From.ID)
			if err != nil {
				b.lg.Fatalf("failed rewriting post: %v", err)
			}

			msg = RewriteSuccessMessage[language]
		}

		message := tgbotapi.NewMessage(m.Chat.ID, msg)
		message.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(message)
	}
}
