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
			continue
		}

		if update.Message.IsCommand() {
			b.ExecuteCommand(update.Message)
		}

		b.lg.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
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
	gocron.Every(1).Minute().Do(b.parseSources)

	// Start all the pending jobs
	<-gocron.Start()
}
