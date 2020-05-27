package main

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/ErrorBoi/feedparserbot/bot"
	"github.com/ErrorBoi/feedparserbot/db"
	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/globalsettings"
	"github.com/ErrorBoi/feedparserbot/scrap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("init logger error: %v", err))
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", bot.User, bot.Pass, bot.Host, bot.Port, bot.DbName)
	Fpdb, err := db.NewDB(dataSourceName, bot.YtToken)
	if err != nil {
		sugar.Infow("Database connection error",
			"host", bot.Host,
			"port", bot.Port,
			"user", bot.User,
			"pass", bot.Pass,
			"dbName", bot.DbName)
		sugar.Fatalf("Database connection error: %v", err)
	}
	defer Fpdb.Cli.Close()

	// Entity for global settings (id = 1)
	ctx := context.Background()
	_, err = Fpdb.Cli.Globalsettings.Query().Where(globalsettings.ID(1)).Only(ctx)
	if ent.IsNotFound(err) {
		_, err = Fpdb.Cli.Globalsettings.Create().Save(ctx)
		if err != nil {
			sugar.Fatalf("Failed creating global settings: %v", err)
		}
	}

	sugar.Info("Successfully connected to database!")

	// run the auto migration tool.
	if err := Fpdb.Cli.Schema.Create(context.Background()); err != nil {
		sugar.Fatalf("failed creating schema resources: %v", err)
	}

	doScrap, err := strconv.ParseBool(bot.Scrap)
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

	Fpbot, err := bot.InitBot(bot.Token, Fpdb, sugar)
	if err != nil {
		sugar.Infow("Bot init error",
			"botToken", bot.Token)
		sugar.Fatalf("Bot init error: %v", err)
	}

	debugMode, err := strconv.ParseBool(bot.DebugMode)
	if err != nil {
		sugar.Errorf("Parse bool error: %v", err)
	}
	Fpbot.BotAPI.Debug = debugMode

	err = Fpbot.InitUpdates(bot.Token)
	if err != nil {
		sugar.Errorf("Init Updates error: %v", err)
	}
}
