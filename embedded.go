//go:build !windows

package libsqldb

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

type LibSqlDB struct {
	DB             *sql.DB
	connector      *libsql.Connector // only used for embedded replica
	dir            string            // only used for embedded replica
	localDBName    string            // only used for embedded replica
	authToken      string
	syncInterval   time.Duration // only used for embedded replica
	encryptionKey  string        // only used for embedded replica
	readYourWrites *bool         // only used for embedded replica
}

var syncInterval = 200 * time.Millisecond

func NewLibSqlDB(primaryUrl string, migrationFiles embed.FS, opts ...Options) (*LibSqlDB, error) {
	l := libSqlDefaults()

	_migrationFiles = migrationFiles

	for _, option := range opts {
		err := option(l)
		if err != nil {
			return nil, fmt.Errorf("error applying options | %w", err)
		}
	}

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return nil, fmt.Errorf("error setting up temporary directory for local database | %w", err)
	}
	//defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, l.localDBName)

	var lsOpts []libsql.Option

	if l.authToken != "" {
		lsOpts = append(lsOpts, libsql.WithAuthToken(l.authToken))
	}

	if l.syncInterval != 0 {
		lsOpts = append(lsOpts, libsql.WithSyncInterval(l.syncInterval))
	}

	if l.encryptionKey != "" {
		lsOpts = append(lsOpts, libsql.WithEncryption(l.encryptionKey))
	}

	if l.readYourWrites != nil {
		lsOpts = append(lsOpts, libsql.WithReadYourWrites(l.readYourWrites))
	}

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		lsOpts,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating connector | %w", err)
	}

	db := sql.OpenDB(connector)

	err = setupMigrations()
	if err != nil {
		return nil, fmt.Errorf("error setting up migrations | %w", err)
	}

	l.db = db
	l.connector = connector
	l.dir = dir

	return l, nil
}

func libSqlDefaults() *LibSqlDB {
	return &LibSqlDB{
		localDBName: "local.db",
	}
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
