package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"

	"github.com/ErrorBoi/feedparserbot/db"
)

type Bot struct {
	BotAPI *tgbotapi.BotAPI
	db     *db.DB
	lg     *zap.SugaredLogger
}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonSwitch("2sw", "open 2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

var numericKeyboard2 = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Отправить последние посты"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Подписки"),
		tgbotapi.NewKeyboardButton("Настройки"),
	),
)

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
		if update.Message == nil { // ignore any non-Message Updates
			if update.CallbackQuery != nil {
				fmt.Print(update)

				if update.CallbackQuery.Data == "3" {
					msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "UPD")
					b.BotAPI.Send(msg)
				}

				b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

				b.BotAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
			}
		} else {
			if update.Message.IsCommand() {
				b.ExecuteCommand(update.Message)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				switch update.Message.Text {
				case "test1":
					msg.ReplyMarkup = numericKeyboard
				case "open":
					msg.ReplyMarkup = numericKeyboard2
				}

				b.BotAPI.Send(msg)
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
	default:
		if m.Chat.IsPrivate() {
			msg := tgbotapi.NewMessage(m.Chat.ID, "К сожалению, я не знаю такой команды. "+
				"Напишите /help для получения справки по командам")
			msg.ReplyToMessageID = m.MessageID
			b.BotAPI.Send(msg)
		}
	}
}

func (b *Bot) RunScheduler() {
	fmt.Println("started scheduler")

	// Parse RSS Feeds every 5 minutes
	gocron.Every(5).Minute().Do(b.parseSources)

	// Send posts to users with instant or <=4h sending frequency
	gocron.Every(5).Minute().Do(b.sendPostsQuick)

	// Send posts to users with AM or PM sending frequency
	gocron.Every(1).Day().At("11:00").Do(b.sendPostsAM)
	gocron.Every(1).Day().At("19:00").Do(b.sendPostsPM)

	// Send posts to users with weekly sending frequency
	gocron.Every(1).Day().Do(b.sendPostsDaily)

	// Start all the pending jobs
	<-gocron.Start()
}
