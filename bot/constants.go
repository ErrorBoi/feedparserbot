package bot

import "os"

type Translations struct {
	RU string
	EN string
}

var (
	Token     = os.Getenv("TOKEN")
	Host      = os.Getenv("HOST")
	Port      = os.Getenv("PORT")
	User      = os.Getenv("USER")
	Pass      = os.Getenv("PASS")
	DbName    = os.Getenv("DBNAME")
	YtToken   = os.Getenv("YTTOKEN")
	Scrap     = os.Getenv("SCRAP")
	DebugMode = os.Getenv("DEBUG_MODE")
)

var (
	startText = map[string]string{
		"RU": `Привет! Я бот-агрегатор новостей. Умею собирать новости из разных источников на основе ваших предпочтений.
Введите /help для получения справки.`,
		"EN": `Hi! I'm a news aggregator bot. I can collect news from different sources based on your preferences.
		Enter /help for list of commands.`,
	}
	helpText = map[string]string{
		"RU": `Список доступных команд:
* /help - показать справку по командам
* /start - начать/возобновить работу бота
* /add [URL источника] - добавить источник новостей
* /urgent [список слов] - установить список "срочных" слов
* /banned [список слов] - установить "Чёрный список"`,
		"EN": `List of available commands:
* /help - show help for commands
* /start - start / resume the bot
* /add [Source URL] - add a news source
* /urgent [list of words] - set a list of "urgent" words
* /banned [list of words] - set a "Black list"`,
	}
	helpEditorText = map[string]string{
		"RU": `
* /set_clickbait [список слов] - установить список кликбейтных слов
* /rewrite [post ID] [новый заголовок] - новый заголовок для новости`,
		"EN": `
* /set_clickbait [list of words] - set list of clickbait words
* /rewrite [post ID] [new title] - new title for a post`,
	}
	helpAdminText = map[string]string{
		"RU": `
* /set_editor [telegram ID] - назначить пользователя редактором
* /remove_editor [telegram ID] - забрать у пользователя права редактора`,
		"EN": `
* /set_editor [telegram ID] - give editor rights to a user
* /remove_editor [telegram ID] - remove Editor rights from a user`,
	}
	FrequencyTextDict = map[string]map[string]string{
		"RU": FrequencyTextDictRu,
		"EN": FrequencyTextDictEn,
	}
	FrequencyTextDictRu = map[string]string{
		"instant": "Мгновенно",
		"1h":      "Не чаще 1 раза в час",
		"4h":      "Не чаще 1 раза в 4 часа",
		"am":      "Раз в день (утром)",
		"pm":      "Раз в день (вечером)",
		"mon":     "По понедельникам",
		"tue":     "По вторникам",
		"wed":     "По средам",
		"thu":     "По четвергам",
		"fri":     "По пятницам",
		"sat":     "По субботам",
		"sun":     "По воскресеньям",
	}
	FrequencyTextDictEn = map[string]string{
		"instant": "Instantly",
		"1h":      "No more than 1 time per hour",
		"4h":      "No more than once every 4 hours",
		"am":      "Once a day (in the morning)",
		"pm":      "Once a day (in the evening)",
		"mon":     "On Mondays",
		"tue":     "On Tuesdays",
		"wed":     "On Wednesdays",
		"thu":     "On Thursdays",
		"fri":     "On Fridays",
		"sat":     "On Saturdays",
		"sun":     "On Sundays",
	}
	UnknownCommandMessage = map[string]string{
		"RU": "К сожалению, я не знаю такой команды. Напишите /help для получения справки по командам",
		"EN": "Unfortunately, I don't know this command. Type /help to get list of commands",
	}
	SelectSourcesMessage = map[string]string{
		"RU": "Выберите источники, на которые хотите подписаться:",
		"EN": "Select the sources you want to subscribe to:",
	}
	SubscriptionsMessage = map[string]string{
		"RU": "Подписки",
		"EN": "Subscriptions",
	}
	NoSubscriptionsMessage = map[string]string{
		"RU": "У вас нет активных подписок.",
		"EN": "You don't have any active subscriptions.",
	}
	AlreadySubscribedMessage = map[string]string{
		"RU": "Вы уже подписаны!",
		"EN": "You have already signed up!",
	}
	SubscriptionCompletedMessage = map[string]string{
		"RU": "Подписка оформлена",
		"EN": "Subscription completed",
	}
	SubRBHubsMessage = map[string]string{
		"RU": "Для того, чтобы подписаться на конкретный раздел сайта RB.ru, найдите ссылку на этот раздел" +
			" в <a href=\"https://rb.ru/list/rss/\">списке</a> и пришлите команду /add + ссылка на раздел. Например: " +
			"/add http://rusbase.com/feeds/tag/bitcoin/\n\n" +
			"<b>Внимание:</b> подписка на любой раздел RB.ru отключит глобальную подписку на ресурс, а глобальная подписка " +
			"отменяет подписки на разделы!",
		"EN": "To subscribe to a specific hub of the RB.ru, find a link to this section" +
			"in <a href=\"https://rb.ru/list/rss/\">the list</a> and send command /add + link to the hub. For example:" +
			"/add http://rusbase.com/feeds/tag/bitcoin/\n\n" +
			"<b>Attention:</b> subscribe to any hub of RB.ru disables the global subscription to the resource, and the global subscription" +
			"cancels subscriptions to sections!",
	}
	SubVCHubsMessage = map[string]string{
		"RU": "Для того, чтобы подписаться на конкретный раздел сайта VC.ru, найдите ссылку на этот раздел" +
			" в <a href=\"https://vc.ru/subs\">списке</a> и пришлите команду /add + ссылка на раздел. Например: " +
			"/add https://vc.ru/marketing",
		"EN": "To subscribe to a specific hub of the VC.ru, find a link to this section" +
			"in <a href=\"https://vc.ru/subs\">the list</a> and send command /add + link to the hub. For example:" +
			"/add https://vc.ru/marketing",
	}
	SettingsMessage = map[string]string{
		"RU": "<b>👤 Мои настройки</b>\n\nПериодичность отправки: %s\nСрочные слова: %s\nЧёрный список: %s\nЯзык: %s",
		"EN": "<b>👤 My settings</b>\n\nSending frequency: %s\nUrgent words: %s\nBanned words: %s\nLanguage: %s",
	}
	SelectUnsubSourceMessage = map[string]string{
		"RU": "Кликните на источник, от которого хотите отписаться:",
		"EN": "Click on the source you want to unsubscribe from:",
	}
	SelectFrequencyMessage = map[string]string{
		"RU": "Выберите желаемую периодичность отправки:",
		"EN": "Select the desired frequency of sending:",
	}
	FrequencyUpdatedMessage = map[string]string{
		"RU": "Периодичность изменена",
		"EN": "The frequency changed",
	}
	SetUrgentWordsMessage = map[string]string{
		"RU": "Вы можете установить \"срочные\" слова, при вхождении " +
			"которых в заголовок, новости будут отправлены вам вне зависимости от настроек периодичности. " +
			"Чтобы это сделать - введите /urgent + список слов через запятую. Например:\n /urgent ЛО,Санкт-Петербург,Доллар",
		"EN": "You can set \"urgent \" words, when they appear in the title, news will be sent to you regardless " +
			"of the frequency settings. To do this, enter /urgent + a comma-separated list of words. For example:" +
			"\n /urgent Moscow,Saint Petersburg,Russia",
	}
	SetBannedWordsMessage = map[string]string{
		"RU": "Вы можете установить \"чёрный список\" слов. Если слова из этого списка входят " +
			"в заголовок новости, она не будет вам отправлена. " +
			"Чтобы это сделать - введите /banned + список слов через запятую. Например:\n /banned Коронавирус,Поправки",
		"EN": "You can set a \"blacklist\" of words. If words from this list are included in the news title, " +
			"it will not be sent to you. To do this, enter /banned + a comma-separated list of words. For example:" +
			"/banned Coronavirus,Amendments",
	}
	LanguageChangedMessage = map[string]string{
		"RU": "Язык изменён",
		"EN": "The language is changed",
	}
	EmptyAddArgsMessage = map[string]string{
		"RU": "После команды нужно указать ссылку на источник. Например /add https://vc.ru/marketing",
		"EN": "After the command, you need to specify a link to the source. For example, /add https://vc.ru/marketing",
	}
	VCHubNotFoundMessage = map[string]string{
		"RU": "Раздел VC.ru с таким названием не найден",
		"EN": "VC.ru hub with this name was not found",
	}
	RBHubNotFoundMessage = map[string]string{
		"RU": "Раздел RB.ru с таким названием не найден",
		"EN": "RB.ru hub with this name was not found",
	}
	SourceNotFoundMessage = map[string]string{
		"RU": "Источник с таким названием не найден",
		"EN": "The source with this name was not found",
	}
	EmptyUrgentArgsMessage = map[string]string{
		"RU": "После команды нужно указать список слов через запятую. Например:\n" +
			"/urgent ЛО,Санкт-Петербург,Доллар",
		"EN": "After the command, enter a comma-separated list of words. For example:\n" +
			"/urgent Moscow,Saint Petersburg,Russia",
	}
	UrgentWordsSuccessMessage = map[string]string{
		"RU": "\"Срочные\" слова записаны!",
		"EN": "\"Urgent\" words are saved!",
	}
	EmptyBannedArgsMessage = map[string]string{
		"RU": "После команды нужно указать список слов через запятую. Например:\n/banned Коронавирус,Поправки",
		"EN": "After the command, enter a comma-separated list of words. For example:\n/banned Coronavirus,Amendments",
	}
	BannedWordsSuccessMessage = map[string]string{
		"RU": "\"Чёрный список\" обновлён!",
		"EN": "\"Blacklist\" updated!",
	}
	SelectLanguageMessage = map[string]string{
		"RU": "Выберите желаемый язык:",
		"EN": "Select the desired language:",
	}
	EmptySuperArgsMessage = map[string]string{
		"RU": "После команды нужно указать API Token этого бота для получения прав администратора.",
		"EN": "After the command, you need to specify the API Token of this bot to get admin rights.",
	}
	SuperSuccessMessage = map[string]string{
		"RU": "Права админа получены!",
		"EN": "Admin's rights granted!",
	}
	SuperValidationErrorMessage = map[string]string{
		"RU": "Права админа не получены! Указан неправильный токен",
		"EN": "Admin rights are not granted! Invalid token specified",
	}
	AddEditorInvalidMessage = map[string]string{
		"RU": "После команды нужно указать Telegram ID юзера, которому хотите выдать права редактора.",
		"EN": "After the command, specify the Telegram ID of the user you want to grant editor rights to.",
	}
	UserNotFoundMessage = map[string]string{
		"RU": "Не найден пользователь с таким Telegram ID",
		"EN": "No user was found with this Telegram ID",
	}
	AddEditorSuccessMessage = map[string]string{
		"RU": "Права редактора выданы!",
		"EN": "Editor's rights granted!",
	}
	RemoveEditorInvalidMessage = map[string]string{
		"RU": "После команды нужно указать Telegram ID юзера, у которого вы хотите забрать права редактора.",
		"EN": "After the command, you need to specify the Telegram ID of the user from whom you want to take away the editor's rights.",
	}
	RemoveEditorSuccessMessage = map[string]string{
		"RU": "Права редактора забраны!",
		"EN": "Editor's rights are revoked!",
	}
	ClickbaitEmptyArgsMessage = map[string]string{
		"RU": "После команды нужно указать список слов через запятую. Например:\n/clickbait Коронавирус,Поправки",
		"EN": "After the command, enter a comma-separated list of words. For example:\n/clickbait Coronavirus,Amendments",
	}
	ClickbaitSuccessMessage = map[string]string{
		"RU": "\"Кликбейтные слова\" обновлены!",
		"EN": "\"Clickbait words\" updated!",
	}
	ClickbaitFormatMessage = map[string]string{
		"RU": "Обнаружен пост с кликбейтным заголовком. Пожалуйста, перепишите заголовок на обычный.\n" +
			"Для этого напишите '/rewrite + postID + переписанный заголовок'.\n" +
			"Например: /rewrite 123 Новый Некликбейтный Заголовок\n" +
			"Post ID: %d\n<a href='%s'>%s</a>\n\n%s",
		"EN": "A post with a clickbait header was detected. Please change the title to normal.\n" +
			"To do this, write '/rewrite + postID + rewritten title'.\n" +
			"Example: /rewrite 123 New Not-Clickbait Title\n" +
			"Post ID: %d\n<a href='%s'>%s</a>\n\n%s",
	}
	RewriteEmptyArgsMessage = map[string]string{
		"RU": "После команды нужно указать ID поста и новый заголовок. Например:\n/rewrite 12 Новый Заголовок",
		"EN": "After the command, you need to specify the post ID and a new title. For example:\n/rewrite 12 New Header",
	}
	RewriteSuccessMessage = map[string]string{
		"RU": "Заголовок переписан!",
		"EN": "Title rewritten!",
	}
)
