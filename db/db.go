package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	entsql "github.com/facebookincubator/ent/dialect/sql"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/schema"

	// Go PostgreSQL package
	_ "github.com/lib/pq"
)

// NewDB creates and returns Database
func NewDB(dataSourceName string) (*ent.Client, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB("postgres", db)

	return ent.NewClient(ent.Driver(drv)), nil
}

func CreatePost(ctx context.Context, client *ent.Client) (*ent.Post, error) {
	titleTransactions := schema.TitleTranslations{
		RU: "пример",
		EN: "sample",
	}
	p, err := client.Post.Create().SetTitle("custom title").SetTitleTranslations(titleTransactions).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating post: %v", err)
	}

	posts, err := client.Post.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed reading posts: %v", err)
	}
	log.Println(posts)

	log.Println("post was created: ", p)
	return p, nil
}
