package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"

	"github.com/ErrorBoi/feedparserbot/db"
	"github.com/ErrorBoi/feedparserbot/ent/source"
	"github.com/ErrorBoi/feedparserbot/ent/user"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
)

type Bot struct {
	BotAPI *tgbotapi.BotAPI
	db     *db.DB
	lg     *zap.SugaredLogger
}

// InitBot inits a bot with given Token
func InitBot(BotToken string, db *db.DB, lg *zap.SugaredLogger) (*Bot, error) {
	var err error
	botAPI, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		BotAPI: botAPI,
		db:     db,
		lg:     lg,
	}, nil
}

// InitUpdates inits an Updates Channel
func (b *Bot) InitUpdates(BotToken string) error {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	updates, err := b.BotAPI.GetUpdatesChan(ucfg)
	if err != nil {
		return err
	}

	go b.RunScheduler()

	for update := range updates {
		if update.Message == nil {
			if update.CallbackQuery != nil {
				b.ExecuteCallbackQuery(update.CallbackQuery)
			}
		} else {
			if update.Message.IsCommand() {
				b.ExecuteCommand(update.Message)
			} else {
				b.ExecuteText(update.Message)
			}

			b.lg.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
	}

	return nil
}

// ExecuteCommand distributes commands to go routines
func (b *Bot) ExecuteCommand(m *tgbotapi.Message) {
	command := strings.ToLower(m.Command())

	ctx := context.Background()

	var language string
	us, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).QuerySettings().Only(ctx)
	if err != nil {
		b.lg.Errorf("failed querying user settings: %v", err)
		language = "RU"
	} else {
		language = string(us.Language)
	}

	switch command {
	case "start":
		go b.start(m, language)
	case "help":
		go b.help(m, language)
	case "add":
		go b.add(m, language)
	case "urgent":
		go b.urgent(m, language)
	case "banned":
		go b.bannedWords(m, language)
	case "super":
		go b.super(m, language)
	case "set_editor":
		go b.setEditor(m, language)
	case "remove_editor":
		go b.removeEditor(m, language)
	case "set_clickbait":
		go b.setClickbaitWords(m, language)
	case "rewrite":
		go b.rewritePost(m, language)
	default:
		if m.Chat.IsPrivate() {
			ctx := context.Background()

			us, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).QuerySettings().Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying user settings: %v", err)
			}
			msg := tgbotapi.NewMessage(m.Chat.ID, UnknownCommandMessage[string(us.Language)])
			msg.ReplyToMessageID = m.MessageID
			b.BotAPI.Send(msg)
		}
	}
}

// ExecuteCallbackQuery handles callback queries
func (b *Bot) ExecuteCallbackQuery(cq *tgbotapi.CallbackQuery) {
	ctx := context.Background()

	us, err := b.db.Cli.User.Query().Where(user.TgID(cq.From.ID)).QuerySettings().Only(ctx)
	if err != nil {
		b.lg.Errorf("failed querying user settings: %v", err)
	}
	language := string(us.Language)

	switch cq.Data {
	case "sub":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, SelectSourcesMessage[language])
		keyboard := subscribeKeyboard[language]
		msg.ReplyMarkup = &keyboard
		b.BotAPI.Send(msg)
	case "backSubscribeKeyboard":
		text := b.getSubsText(cq.Message.Chat.ID, language)
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, text)
		keyboard := subsMainKeyboard[language]
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "subVC":
		b.subUserToSource(cq, "https://vc.ru/rss/all", language)
	case "subRB":
		b.subUserToSource(cq, "https://rb.ru/feeds/all", language)
	case "subFontanka":
		b.subUserToSource(cq, "https://www.fontanka.ru/fontanka.rss", language)
	case "subForbes":
		b.subUserToSource(cq, "https://www.forbes.ru/newrss.xml", language)
	case "subRBHubs":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, SubRBHubsMessage[language])
		keyboard := subHubKeyboard[language]
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "subVCHubs":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, SubVCHubsMessage[language])
		keyboard := subHubKeyboard[language]
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "unsub":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
			SelectUnsubSourceMessage[language])
		unsubKeyboard, err := b.getUnsubKeyboard(cq, language)
		if err != nil {
			b.lg.Fatalf("failed generating unsub keyboard: %v", err)
		}
		msg.ReplyMarkup = unsubKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "settings":
		ctx := context.Background()
		us, err := b.db.Cli.User.Query().Where(user.TgID(cq.From.ID)).QuerySettings().Only(ctx)
		if err != nil {
			b.lg.Errorf("error querying user settings: %v", err)
		}

		text := fmt.Sprintf(SettingsMessage[language], FrequencyTextDict[language][string(us.SendingFrequency)], us.UrgentWords, us.BannedWords, us.Language)
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, text)
		keyboard := settingsMainKeyboard[language]
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "frequency":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, SelectFrequencyMessage[language])
		keyboard := settingsFreqKeyboard[language]
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "instant":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyInstant).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "1h":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequency1h).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "4h":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequency4h).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "am":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyAm).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "pm":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyPm).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "mon":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyMon).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "tue":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyTue).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "wed":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyWed).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "thu":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyThu).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "fri":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyFri).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "sat":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyFri).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "sun":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencySun).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, FrequencyUpdatedMessage[language]))
	case "urgent":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, SetUrgentWordsMessage[language])
		keyboard := settingsBackKeyboard[language]
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "banned_words":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
			SetBannedWordsMessage[language])
		keyboard := settingsBackKeyboard[language]
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "language":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, SelectLanguageMessage[language])
		keyboard := settingsLanguageKeyboard[language]
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "ru_language":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetLanguage(usersettings.LanguageRU).
			Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, LanguageChangedMessage["RU"]))

		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, LanguageChangedMessage["RU"])
		msg.ReplyMarkup = mainKeyboard["RU"]
		b.BotAPI.Send(msg)
	case "en_language":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetLanguage(usersettings.LanguageEN).
			Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, LanguageChangedMessage["EN"]))

		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, LanguageChangedMessage["EN"])
		msg.ReplyMarkup = mainKeyboard["EN"]
		b.BotAPI.Send(msg)
	default:
		if strings.HasPrefix(cq.Data, "unsubSource") {
			ctx := context.Background()

			sourceID := strings.TrimPrefix(cq.Data, "unsubSource")
			intSourceID, err := strconv.Atoi(sourceID)
			if err != nil {
				b.lg.Errorf("failed converting string to int: %v", err)
			}

			s, err := b.db.Cli.Source.Query().Where(source.ID(intSourceID)).Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying source: %v", err)
			}

			_, err = b.db.Cli.User.Update().RemoveSources(s).Save(ctx)
			if err != nil {
				b.lg.Errorf("failed updating user: %v", err)
			}

			msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
				SelectUnsubSourceMessage[language])
			unsubKeyboard, err := b.getUnsubKeyboard(cq, language)
			if err != nil {
				b.lg.Fatalf("failed generating unsub keyboard: %v", err)
			}
			msg.ReplyMarkup = unsubKeyboard
			msg.ParseMode = tgbotapi.ModeHTML
			b.BotAPI.Send(msg)
		}
	}
}

// ExecuteText handles text messages
func (b *Bot) ExecuteText(m *tgbotapi.Message) {
	ctx := context.Background()

	us, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).QuerySettings().Only(ctx)
	if err != nil {
		b.lg.Errorf("failed querying user settings: %v", err)
	}
	language := string(us.Language)

	var msg tgbotapi.MessageConfig
	switch m.Text {
	case "🗞  Источники новостей", "🗞  Sources":
		text := b.getSubsText(m.Chat.ID, language)
		msg = tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = subsMainKeyboard[language]
		msg.ParseMode = tgbotapi.ModeHTML
	case "⚙️ Настройки", "⚙️ Settings":
		ctx := context.Background()
		us, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).QuerySettings().Only(ctx)
		if err != nil {
			b.lg.Errorf("error querying user settings: %v", err)
		}

		text := fmt.Sprintf(SettingsMessage[language], FrequencyTextDict[language][string(us.SendingFrequency)], us.UrgentWords, us.BannedWords, us.Language)
		msg = tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = settingsMainKeyboard[language]
		msg.ParseMode = tgbotapi.ModeHTML
	}

	b.BotAPI.Send(msg)
}

func (b *Bot) RunScheduler() {
	fmt.Println("started scheduler")

	// Send posts to users with instant or <=4h sending frequency
	gocron.Every(5).Minute().Do(b.sendPostsQuick)

	// Send posts to users with AM or PM sending frequency
	gocron.Every(1).Day().At("11:00").Do(b.sendPostsAM)
	gocron.Every(1).Day().At("19:00").Do(b.sendPostsPM)

	// Send posts to users with weekly sending frequency
	gocron.Every(1).Day().At("19:00").Do(b.sendPostsDaily)

	// Parse RSS Feeds every 5 minutes
	gocron.Every(5).Minute().Do(b.parseSources)

	// Send posts with urgent words in title

	// Start all the pending jobs
	<-gocron.Start()
}
