package main

import (
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/ErrorBoi/feedparserbot/bot"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("init logger error: %w", err))
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	token := os.Getenv("TOKEN")

	Fpbot, err := bot.InitBot(token, sugar)
	if err != nil {
		sugar.Infow("Bot init error",
			"botToken", token)
		sugar.Fatalf("Bot init error: %w", err.Error())
	}

	debugMode, err := strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	if err != nil {
		sugar.Errorf("Parse bool error: %w", err)
	}
	Fpbot.BotAPI.Debug = debugMode

	err = Fpbot.InitUpdates(token)
	if err != nil {
		sugar.Errorf("Init Updates error: %w", err)
	}
}
