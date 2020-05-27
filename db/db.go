package db

import (
	"database/sql"

	translate "github.com/dafanasev/go-yandex-translate"
	entsql "github.com/facebookincubator/ent/dialect/sql"

	"github.com/ErrorBoi/feedparserbot/ent"
	// Go PostgreSQL package
	_ "github.com/lib/pq"
)

type DB struct {
	Cli *ent.Client
	Tr  *translate.Translator
}

// NewDB creates and returns Database
func NewDB(dataSourceName string, ytToken string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB("postgres", db)

	tr := translate.New(ytToken)

	return &DB{Cli: ent.NewClient(ent.Driver(drv)), Tr: tr}, nil
}
