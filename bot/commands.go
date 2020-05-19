package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (b *Bot) start(m *tgbotapi.Message) {
	b.help(m)
}

func (b *Bot) help(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, helpText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}
