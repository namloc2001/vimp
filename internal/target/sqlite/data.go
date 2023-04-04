package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	// sqlite3 driver
	_ "modernc.org/sqlite"
)

const (
	dataDriver string = "sqlite"
)

var (
	//go:embed sql/*
	f embed.FS

	driverPrefix = fmt.Sprintf("%s://", dataDriver)
)

func getStore(path string) (*sql.DB, error) {
	if path == "" {
		return nil, errors.New("directory not specified")
	}

	// remove driver prefix
	path = strings.Replace(path, driverPrefix, "", 1)

	wasCreated := false
	log.Debug().Msgf("data path: %s", path)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("data file does not exist, creating...")
		wasCreated = true
	}

	db, err := sql.Open(dataDriver, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database: %s", path)
	}

	if wasCreated {
		log.Debug().Msg("creating schema...")

		b, err := f.ReadFile("sql/ddl.sql")
		if err != nil {
			return nil, errors.Wrap(err, "failed to read the schema creation file")
		}
		if _, err := db.Exec(string(b)); err != nil {
			return nil, errors.Wrapf(err, "failed to create database schema in: %s", path)
		}
	}

	log.Debug().Msg("data initialized")
	return db, nil
}