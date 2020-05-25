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

	switch command {
	case "start":
		go b.start(m)
	case "help":
		go b.help(m)
	case "add":
		go b.add(m)
	case "urgent":
		go b.urgent(m)
	case "banned":
		go b.bannedWords(m)
	default:
		if m.Chat.IsPrivate() {
			msg := tgbotapi.NewMessage(m.Chat.ID, "К сожалению, я не знаю такой команды. "+
				"Напишите /help для получения справки по командам")
			msg.ReplyToMessageID = m.MessageID
			b.BotAPI.Send(msg)
		}
	}
}

// ExecuteCallbackQuery handles callback queries
func (b *Bot) ExecuteCallbackQuery(cq *tgbotapi.CallbackQuery) {
	switch cq.Data {
	case "sub":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
			"Выберите источники, на которые хотите подписаться:")
		msg.ReplyMarkup = &subscribeKeyboard
		b.BotAPI.Send(msg)
	case "backSubscribeKeyboard":
		text := b.getSubsText(cq.Message.Chat.ID)
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, text)
		msg.ReplyMarkup = &subsMainKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "subVC":
		b.subUserToSource(cq, "https://vc.ru/rss/all")
	case "subRB":
		b.subUserToSource(cq, "https://rb.ru/feeds/all")
	case "subFontanka":
		b.subUserToSource(cq, "https://www.fontanka.ru/fontanka.rss")
	case "subForbes":
		b.subUserToSource(cq, "https://www.forbes.ru/newrss.xml")
	case "subRBHubs":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, subRBHubsText)
		msg.ReplyMarkup = &subHubKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "subVCHubs":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, subVCHubsText)
		msg.ReplyMarkup = &subHubKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "unsub":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
			"Кликните на источник, от которого хотите отписаться:")
		unsubKeyboard, err := b.getUnsubKeyboard(cq)
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

		text := fmt.Sprintf(settingsText, FrequencyTextDict[string(us.SendingFrequency)], us.UrgentWords, us.BannedWords, us.Language)
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, text)
		msg.ReplyMarkup = &settingsMainKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "frequency":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
			"Выберите желаемую периодичность отправки:")
		msg.ReplyMarkup = &settingsFreqKeyboard
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
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "1h":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequency1h).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "4h":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequency4h).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "am":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyAm).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "pm":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyPm).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "mon":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyMon).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "tue":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyTue).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
	case "wed":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyWed).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "thu":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyThu).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "fri":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyFri).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "sat":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencyFri).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "sun":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetSendingFrequency(usersettings.SendingFrequencySun).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Периодичность изменена"))
	case "urgent":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
			"Вы можете установить \"срочные\" слова, при вхождении "+
				"которых в заголовок, новости будут отправлены вам вне зависимости от настроек периодичности. "+
				"Чтобы это сделать - введите /urgent + список слов через запятую. Например:\n /urgent ЛО,Санкт-Петербург,Доллар")
		msg.ReplyMarkup = &settingsBackKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "banned_words":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
			"Вы можете установить \"чёрный список\" слов. Если слова из этого списка входят "+
				"в заголовок новости, она не будет вам отправлена. "+
				"Чтобы это сделать - введите /banned + список слов через запятую. Например:\n /banned Коронавирус,Поправки")
		msg.ReplyMarkup = &settingsBackKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "language":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID,
			"Выберите желаемый язык:")
		msg.ReplyMarkup = &settingsLanguageKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
		b.BotAPI.Send(msg)
	case "ru_language":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetLanguage(usersettings.LanguageRu).
			Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Язык изменён"))
	case "en_language":
		ctx := context.Background()
		_, err := b.db.Cli.UserSettings.Update().
			Where(usersettings.HasUserWith(user.TgID(cq.From.ID))).
			SetLanguage(usersettings.LanguageEn).
			Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Язык изменён"))
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
				"Кликните на источник, от которого хотите отписаться:")
			unsubKeyboard, err := b.getUnsubKeyboard(cq)
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
	var msg tgbotapi.MessageConfig
	switch m.Text {
	case "🗞  Источники новостей":
		text := b.getSubsText(m.Chat.ID)
		msg = tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = subsMainKeyboard
		msg.ParseMode = tgbotapi.ModeHTML
	case "⚙️ Настройки":
		ctx := context.Background()
		us, err := b.db.Cli.User.Query().Where(user.TgID(m.From.ID)).QuerySettings().Only(ctx)
		if err != nil {
			b.lg.Errorf("error querying user settings: %v", err)
		}

		text := fmt.Sprintf(settingsText, FrequencyTextDict[string(us.SendingFrequency)], us.UrgentWords, us.BannedWords, us.Language)
		msg = tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = settingsMainKeyboard
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
