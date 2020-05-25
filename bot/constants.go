package bot

const (
	startText = `Привет! Я бот-агрегатор новостей. Умею собирать новости из разных источников на основе ваших предпочтений.
Введите /help для получения справки.`
	helpText = `Список доступных команд:
* /help - показать справку по командам
* /start - начать/возобновить работу бота
* /add - добавить источник новостей
* /urgent - установить список "срочных" слов
* /banned - установить "Чёрный список"`
)

var (
	FrequencyTextDict = map[string]string{
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
)
