package db

import (
	"context"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/source"
	"github.com/ErrorBoi/feedparserbot/ent/user"
)

func (db *DB) StoreUserSource(ctx context.Context, tgID int, sourceURL string) error {
	u, err := db.Cli.User.Query().Where(user.TgID(tgID)).Only(ctx)
	if err != nil {
		return err
	}

	s, err := db.Cli.Source.Query().Where(source.URL(sourceURL)).Only(ctx)
	if err != nil {
		return err
	}

	var children []*ent.Source
	if sourceURL == "https://rb.ru/feeds/all" {
		children, err = s.QueryChildren().All(ctx)
		if err != nil {
			return err
		}
	}

	_, err = u.Update().AddSources(s).RemoveSources(children...).Save(ctx)
	if err != nil {
		return err
	}

	return nil
}
