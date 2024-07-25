package database

import (
	"github.com/go-pg/pg"
	"github.com/notamedia-org/backend/internal/config"
)

func StartPolling(config *config.Config) (*pg.DB, error) {
	opt, err := pg.ParseURL(config.PgConnection)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opt)

	return db, nil
}
