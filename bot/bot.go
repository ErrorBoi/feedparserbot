package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"

	"github.com/ErrorBoi/feedparserbot/db"
	"github.com/ErrorBoi/feedparserbot/ent"
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
	ctx := context.Background()
	switch cq.Data {
	case "3":
		msg := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "UPD!~")
		b.BotAPI.Send(msg)
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
		err := b.db.StoreUserSource(ctx, cq.From.ID, "https://vc.ru/rss/all")
		if err != nil {
			b.lg.Errorf("failed storing user source: %v", err)
			switch {
			case ent.IsConstraintError(err):
				b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Вы уже подписаны!"))
			}
		} else {
			b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Подписались на VC.ru"))
		}
	default:
		//TODO: realize removing source here, by checking if data has prefix "removeSource"
		// ss = query sources
		// b.db.Cli.User.Update().Where(user.ID(1)).RemoveSources()
	}
}

// ExecuteText handles text messages
func (b *Bot) ExecuteText(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, m.Text)
	switch m.Text {
	case "test1":
		msg.ReplyMarkup = numericKeyboard
	case "open":
		msg.ReplyMarkup = mainKeyboard
	case "🗞  Источники новостей":
		text := b.getSubsText(m.Chat.ID)
		msg = tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = subsMainKeyboard
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

	// Start all the pending jobs
	<-gocron.Start()
}
