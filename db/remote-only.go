//go:build windows

package db

import (
	"database/sql"
	"fmt"
	"github.com/hashicorp/go-multierror"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type LibSqlDB struct {
	db *sql.DB
}

func NewLibSqlDB(primaryUrl, authToken, localDBName string) (*LibSqlDB, error) {
	url := primaryUrl + "?authToken=" + authToken
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, fmt.Errorf("error setting up LibSQL db | %w", err)
	}

	err = setupMigrations()
	if err != nil {
		return nil, fmt.Errorf("error setting up migrations | %w", err)
	}

	return &LibSqlDB{
		db: db,
	}, nil
}

func (t *LibSqlDB) Close() error {
	var resultError *multierror.Error

	if err := t.db.Close(); err != nil {
		resultError = multierror.Append(resultError, fmt.Errorf("failed to close database: %w", err))
	}

	return resultError.ErrorOrNil()
}
