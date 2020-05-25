package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ErrorBoi/feedparserbot/ent/user"
)

var (
	// Keyboards
	numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
	mainKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗞  Источники новостей"),
			tgbotapi.NewKeyboardButton("⚙️ Настройки"),
		),
	)
	subsMainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подписаться", "sub"),
			tgbotapi.NewInlineKeyboardButtonData("Отписаться", "unsub"),
		),
	)
	subscribeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "backSubscribeKeyboard"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Разделы VC.ru", "subVCHubs"),
			tgbotapi.NewInlineKeyboardButtonData("Разделы RB.ru", "subRBHubs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("VC.ru", "subVC"),
			tgbotapi.NewInlineKeyboardButtonData("RB.ru", "subRB"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Fontanka.ru", "subFontanka"),
			tgbotapi.NewInlineKeyboardButtonData("Forbes.ru", "subForbes"),
		),
		)

	// Messages
)

func (b *Bot) getSubsText(chatID int64) string {
	text := fmt.Sprintf("<b>ℹ️ Подписки</b>\n")

	ctx := context.Background()
	ss, err := b.db.Cli.User.Query().Where(user.TgID(int(chatID))).QuerySources().All(ctx)
	if err != nil {
		b.lg.Errorf("failed querying sources: %v", err)
	}

	if len(ss) == 0 {
		text += "У вас нет активных подписок."
	} else {
		for _, s := range ss {
			text += fmt.Sprintf("* %s\n", s.Title)
		}
	}

	return text
}

//func (b *Bot) getSubscribeKeyboard(tgID int) tgbotapi.InlineKeyboardMarkup {
//	ctx := context.Background()
//	ss, err := b.db.Cli.User.Query().Where(user.TgID(tgID)).QuerySources().All(ctx)
//	if err != nil {
//		b.lg.Errorf("failed querying sources: %v", err)
//	}
//}
