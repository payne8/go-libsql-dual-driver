//go:build windows

package libsqldb

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type LibSqlDB struct {
	DB *sql.DB

	dir            string // only used for embedded replica
	localDBName    string // only used for embedded replica but needs to be here for consistency
	authToken      string
	syncInterval   time.Duration // only used for embedded replica
	encryptionKey  string        // only used for embedded replica
	readYourWrites *bool         // only used for embedded replica
}

func NewLibSqlDB(primaryUrl string, migrationFiles embed.FS, opts ...Options) (*LibSqlDB, error) {
	l := &LibSqlDB{}

	_migrationFiles = migrationFiles

	for _, option := range opts {
		err := option(l)
		if err != nil {
			return nil, fmt.Errorf("error applying options | %w", err)
		}
	}

	url := primaryUrl
	if l.authToken != "" {
		url = primaryUrl + "?authToken=" + l.authToken
	}

	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, fmt.Errorf("error setting up LibSQL db | %w", err)
	}

	err = setupMigrations()
	if err != nil {
		return nil, fmt.Errorf("error setting up migrations | %w", err)
	}

	return &LibSqlDB{
		DB: db,
	}, nil
}

func (t *LibSqlDB) Close() error {
	var resultError *multierror.Error

	if err := t.DB.Close(); err != nil {
		resultError = multierror.Append(resultError, fmt.Errorf("failed to close database: %w", err))
	}

	return resultError.ErrorOrNil()
}
