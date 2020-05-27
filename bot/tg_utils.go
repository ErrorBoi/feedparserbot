package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/user"
)

var (
	mainKeyboard = map[string]tgbotapi.ReplyKeyboardMarkup{
		"RU": mainKeyboardRu,
		"EN": mainKeyboardEn,
	}
	mainKeyboardRu = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗞  Источники новостей"),
			tgbotapi.NewKeyboardButton("⚙️ Настройки"),
		),
	)
	mainKeyboardEn = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗞  Sources"),
			tgbotapi.NewKeyboardButton("⚙️ Settings"),
		),
	)

	subsMainKeyboard = map[string]tgbotapi.InlineKeyboardMarkup{
		"RU": subsMainKeyboardRu,
		"EN": subsMainKeyboardEn,
	}
	subsMainKeyboardRu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подписаться", "sub"),
			tgbotapi.NewInlineKeyboardButtonData("Отписаться", "unsub"),
		),
	)
	subsMainKeyboardEn = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Subscribe", "sub"),
			tgbotapi.NewInlineKeyboardButtonData("Unsubscribe", "unsub"),
		),
	)

	subscribeKeyboard = map[string]tgbotapi.InlineKeyboardMarkup{
		"RU": subscribeKeyboardRu,
		"EN": subscribeKeyboardEn,
	}
	subscribeKeyboardRu = tgbotapi.NewInlineKeyboardMarkup(
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
	subscribeKeyboardEn = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "backSubscribeKeyboard"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("VC.ru hubs", "subVCHubs"),
			tgbotapi.NewInlineKeyboardButtonData("RB.ru hubs", "subRBHubs"),
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

	subHubKeyboard = map[string]tgbotapi.InlineKeyboardMarkup{
		"RU": subHubKeyboardRu,
		"EN": subHubKeyboardEn,
	}
	subHubKeyboardRu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "sub"),
		),
	)
	subHubKeyboardEn = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "sub"),
		),
	)

	settingsMainKeyboard = map[string]tgbotapi.InlineKeyboardMarkup{
		"RU": settingsMainKeyboardRu,
		"EN": settingsMainKeyboardEn,
	}
	settingsMainKeyboardRu = tgbotapi.NewInlineKeyboardMarkup(
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
	settingsMainKeyboardEn = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏳ Edit frequency", "frequency"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Edit urgent words", "urgent"),
			tgbotapi.NewInlineKeyboardButtonData("Edit blacklist", "banned_words"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Subscription", "payment"),
			tgbotapi.NewInlineKeyboardButtonData("Edit language", "language"),
		),
	)

	settingsFreqKeyboard = map[string]tgbotapi.InlineKeyboardMarkup{
		"RU": settingsFreqKeyboardRu,
		"EN": settingsFreqKeyboardEn,
	}
	settingsFreqKeyboardRu = tgbotapi.NewInlineKeyboardMarkup(
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
	settingsFreqKeyboardEn = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Instant", "instant"),
			tgbotapi.NewInlineKeyboardButtonData("Once in 1 hour", "1h"),
			tgbotapi.NewInlineKeyboardButtonData("Once in 4 hours", "4h"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Every morning", "am"),
			tgbotapi.NewInlineKeyboardButtonData("Every evening", "pm"),
			tgbotapi.NewInlineKeyboardButtonData("Every Mon", "mon"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Every Tue", "tue"),
			tgbotapi.NewInlineKeyboardButtonData("Every Wed", "wed"),
			tgbotapi.NewInlineKeyboardButtonData("Every Thu", "thu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Every Fri", "fri"),
			tgbotapi.NewInlineKeyboardButtonData("Every Sat", "sat"),
			tgbotapi.NewInlineKeyboardButtonData("Every Sun", "sun"),
		),
	)

	settingsBackKeyboard = map[string]tgbotapi.InlineKeyboardMarkup{
		"RU": settingsBackKeyboardRu,
		"EN": settingsBackKeyboardEn,
	}
	settingsBackKeyboardRu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings"),
		),
	)
	settingsBackKeyboardEn = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "settings"),
		),
	)

	settingsLanguageKeyboard = map[string]tgbotapi.InlineKeyboardMarkup{
		"RU": settingsLanguageKeyboardRu,
		"EN": settingsLanguageKeyboardEn,
	}
	settingsLanguageKeyboardRu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "ru_language"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "en_language"),
		),
	)
	settingsLanguageKeyboardEn = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "ru_language"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "en_language"),
		),
	)

	backSubscribeRow = map[string][]tgbotapi.InlineKeyboardButton{
		"RU": backSubscribeRowRu,
		"EN": backSubscribeRowEn,
	}
	backSubscribeRowRu = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "backSubscribeKeyboard"),
	)
	backSubscribeRowEn = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "backSubscribeKeyboard"),
	)
)

func (b *Bot) getSubsText(chatID int64, language string) string {
	text := fmt.Sprintf("<b>ℹ️ %s</b>\n", SubscriptionsMessage[language])

	ctx := context.Background()
	ss, err := b.db.Cli.User.Query().Where(user.TgID(int(chatID))).QuerySources().All(ctx)
	if err != nil {
		b.lg.Errorf("failed querying sources: %v", err)
	}

	if len(ss) == 0 {
		text += NoSubscriptionsMessage[language]
	} else {
		for _, s := range ss {
			text += fmt.Sprintf("* %s\n", s.Title)
		}
	}

	return text
}

func (b *Bot) subUserToSource(cq *tgbotapi.CallbackQuery, sourceURL string, language string) {
	ctx := context.Background()
	err := b.db.StoreUserSource(ctx, cq.From.ID, sourceURL)
	if err != nil {
		b.lg.Errorf("failed storing user source: %v", err)
		switch {
		case ent.IsConstraintError(err):
			b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, AlreadySubscribedMessage[language]))
		}
	} else {
		b.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(cq.ID, SubscriptionCompletedMessage[language]))
	}
}

func (b *Bot) getUnsubKeyboard(cq *tgbotapi.CallbackQuery, language string) (*tgbotapi.InlineKeyboardMarkup, error) {
	ctx := context.Background()

	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, backSubscribeRow[language])

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
