package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/ErrorBoi/feedparserbot/bot"
	"github.com/ErrorBoi/feedparserbot/db"
	"github.com/ErrorBoi/feedparserbot/scrap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("init logger error: %v", err))
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	token := os.Getenv("TOKEN")

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")
	dbName := os.Getenv("DBNAME")

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbName)
	ytToken := os.Getenv("YTTOKEN")
	Fpdb, err := db.NewDB(dataSourceName, ytToken)
	if err != nil {
		sugar.Infow("Database connection error",
			"host", host,
			"port", port,
			"user", user,
			"pass", pass,
			"dbName", dbName)
		sugar.Fatalf("Database connection error: %v", err)
	}
	defer Fpdb.Cli.Close()

	sugar.Info("Successfully connected to database!")

	// run the auto migration tool.
	if err := Fpdb.Cli.Schema.Create(context.Background()); err != nil {
		sugar.Fatalf("failed creating schema resources: %v", err)
	}

	doScrap, err := strconv.ParseBool(os.Getenv("SCRAP"))
	if err != nil {
		sugar.Errorf("Parse bool error: %v", err)
	}

	if doScrap {
		err = scrap.Sources(Fpdb.Cli)
		if err != nil {
			sugar.Errorf("Scrap sources error: %v", err)
		}

		err = scrap.VCHubs(Fpdb.Cli)
		if err != nil {
			sugar.Errorf("Scrap VC hubs error: %v", err)
		}

		err = scrap.RBHubs(Fpdb.Cli)
		if err != nil {
			sugar.Errorf("Scrap RB hubs error: %v", err)
		}
	}

	Fpbot, err := bot.InitBot(token, Fpdb, sugar)
	if err != nil {
		sugar.Infow("Bot init error",
			"botToken", token)
		sugar.Fatalf("Bot init error: %v", err)
	}

	debugMode, err := strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	if err != nil {
		sugar.Errorf("Parse bool error: %v", err)
	}
	Fpbot.BotAPI.Debug = debugMode

	err = Fpbot.InitUpdates(token)
	if err != nil {
		sugar.Errorf("Init Updates error: %v", err)
	}
}
