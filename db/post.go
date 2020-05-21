package db

import (
	"context"
	"strings"
	"time"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/schema"
	"github.com/ErrorBoi/feedparserbot/ent/source"
)

type Post struct {
	Title               string
	TitleTranslations   schema.Translations
	Subject             *string
	SubjectTranslations *schema.Translations
	Url                 string
	PublishedAt         time.Time
	Description         string
	H1                  string
	Content             string
	SourceId            int
	CreatedAt           time.Time
	UpdatedAt           time.Time
	UpdatedBy           *int
}

type StorePost struct {
	Title       string
	Url         string
	PublishedAt time.Time
	Description string
	H1          string
	Content     string
	SourceURL   string
}

func (db *DB) StorePost(ctx context.Context, post StorePost) (*ent.Post, error) {
	resp, err := db.Tr.Translate("en", post.Title)
	if err != nil {
		return nil, err
	}

	titleTransactions := schema.Translations{
		RU: post.Title,
		EN: resp.Result(),
	}

	src, err := db.Cli.Source.Query().Where(source.URL(post.SourceURL)).Only(ctx)
	if err != nil {
		return nil, err
	}

	var description string
	if len(post.Description) != 0 {
		description = post.Description
	} else {
		description = post.Title
	}

	var h1 string
	if len(post.H1) != 0 {
		h1 = post.H1
	} else {
		h1 = post.Title
	}

	pst, err := db.Cli.Post.
		Create().
		SetTitle(strings.TrimSpace(post.Title)).
		SetTitleTranslations(titleTransactions).
		SetURL(post.Url).
		SetPublishedAt(post.PublishedAt).
		SetDescription(strings.TrimSpace(description)).
		SetH1(strings.TrimSpace(h1)).
		SetContent(strings.TrimSpace(post.Content)).
		SetSource(src).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return pst, nil
}
