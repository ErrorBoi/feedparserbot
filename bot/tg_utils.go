package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ErrorBoi/feedparserbot/ent"
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
	subHubKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "sub"),
		),
	)
	settingsMainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏳ Изменить периодичность", "frequency"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить срочные слова", "urgent"),
			tgbotapi.NewInlineKeyboardButtonData("Изменить чёрный список", "banned_words"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подписка", "payment"),
			tgbotapi.NewInlineKeyboardButtonData("Изменить язык бота", "language"),
		),
	)
	settingsFreqKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Мгновенно", "instant"),
			tgbotapi.NewInlineKeyboardButtonData("Раз в час", "1h"),
			tgbotapi.NewInlineKeyboardButtonData("Раз в 4 часа", "4h"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Каждое утро", "am"),
			tgbotapi.NewInlineKeyboardButtonData("Каждый вечер", "pm"),
			tgbotapi.NewInlineKeyboardButtonData("Каждый Пн", "mon"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Каждый Вт", "tue"),
			tgbotapi.NewInlineKeyboardButtonData("Каждую Ср", "wed"),
			tgbotapi.NewInlineKeyboardButtonData("Каждый Чт", "thu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Каждую Пт", "fri"),
			tgbotapi.NewInlineKeyboardButtonData("Каждую Сб", "sat"),
			tgbotapi.NewInlineKeyboardButtonData("Каждое Вс", "sun"),
		),
	)
	settingsBackKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings"),
		),
	)
	settingsLanguageKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "ru_language"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "en_language"),
		),
	)

	// Messages
	subRBHubsText = "Для того, чтобы подписаться на конкретный раздел сайта RB.ru, найдите ссылку на этот раздел" +
		" в <a href=\"https://rb.ru/list/rss/\">списке</a> и пришлите команду /add + ссылка на раздел. Например: " +
		"/add http://rusbase.com/feeds/tag/bitcoin/\n\n" +
		"<b>Внимание:</b> подписка на любой раздел RB.ru отключит глобальную подписку на ресурс, а глобальная подписка " +
		"отменяет подписки на разделы!"
	subVCHubsText = "Для того, чтобы подписаться на конкретный раздел сайта VC.ru, найдите ссылку на этот раздел" +
		" в <a href=\"https://vc.ru/subs\">списке</a> и пришлите команду /add + ссылка на раздел. Например: " +
		"/add https://vc.ru/marketing"
	settingsText = "<b>👤 Мои настройки</b>\n\nПериодичность отправки: %s\nСрочные слова: %s\nЧёрный список: %s\nЯзык: %s"
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

func (b *Bot) subUserToSource(cq *tgbotapi.CallbackQuery, sourceURL string) {
	ctx := context.Background()
	err := b.db.StoreUserSource(ctx, cq.From.ID, sourceURL)
	if err != nil {
		b.lg.Errorf("failed storing user source: %v", err)
		switch {
		case ent.IsConstraintError(err):
			b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Вы уже подписаны!"))
		}
	} else {
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, "Подписка оформлена"))
	}
}

func (b *Bot) getUnsubKeyboard(cq *tgbotapi.CallbackQuery) (*tgbotapi.InlineKeyboardMarkup, error) {
	ctx := context.Background()

	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "backSubscribeKeyboard"),
	))

	ss, err := b.db.Cli.User.Query().Where(user.TgID(cq.From.ID)).QuerySources().All(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range ss {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.Title, fmt.Sprintf("unsubSource%d", s.ID)),
		))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}
