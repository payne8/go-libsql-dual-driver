//go:build !windows

package db

import (
	"database/sql"
	"fmt"
	"github.com/tursodatabase/go-libsql"
	"os"
	"path/filepath"
)

type LibSqlDB struct {
	db        *sql.DB
	connector *libsql.Connector // only used for embedded replica
	dir       string            // only used for embedded replica
}

var syncInterval = 200 * time.Millisecond

func NewLibSqlDB(primaryUrl, authToken, localDBName string) (*LibSqlDB, error) {
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return nil, fmt.Errorf("error setting up temporary directory for local database | %w", err)
	}
	//defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, localDBName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating connector | %w", err)
	}

	db := sql.OpenDB(connector)

	err = setupMigrations()
	if err != nil {
		return nil, fmt.Errorf("error setting up migrations | %w", err)
	}

	return &LibSqlDB{
		db:        db,
		connector: connector,
		dir:       dir,
	}, nil
}

func (t *LibSqlDB) Close() error {
	var resultError *multierror.Error

	if err := t.db.Close(); err != nil {
		resultError = multierror.Append(resultError, fmt.Errorf("failed to close database: %w", err))
	}

	if t.connector != nil {
		if err := t.connector.Close(); err != nil {
			resultError = multierror.Append(resultError, fmt.Errorf("failed to close connector: %w", err))
		}
	}

	if t.dir != "" {
		if err := os.RemoveAll(t.dir); err != nil {
			resultError = multierror.Append(resultError, fmt.Errorf("failed to remove directory %s: %w", t.dir, err))
		}
	}

	return resultError.ErrorOrNil()
}
