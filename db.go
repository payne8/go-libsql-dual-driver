package libsqldb

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"time"
)

type Migrations struct {
	name  string
	query string
}

type Options func(*LibSqlDB) error

// func NewLibSqlDB is defined in embedded.go and remote-only.go files
// these files are used to define the LibSqlDB struct and the NewLibSqlDB function
// they have different initializations based on the environment, embedded or remote-only
// Windows does not currently support the embedded database, so the remote-only file is used

// setupMigrations initializes the filesystem and reads the migration files into the migrations variable
func (t *LibSqlDB) setupMigrations() error {
	// Walk through the embedded files and read their contents
	err := fs.WalkDir(t.migrationFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			content, err := t.migrationFiles.ReadFile(path)
			if err != nil {
				return err
			}

			migration := Migrations{
				name:  filepath.Base(path),
				query: string(content),
			}
			t.migrations = append(t.migrations, migration)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error setting up migrations | %w", err)
	}
	return nil
}

// Migrate updates the connected LibSqlDB to the latest schema based on the given migrations
func (t *LibSqlDB) Migrate() error {
	if !t.useMigrations {
		return fmt.Errorf("migrations not enabled")
	}

	// check if migration table exists
	var migrationsCheck string
	//goland:noinspection SqlResolve
	err := t.DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'").Scan(&migrationsCheck)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := t.DB.Exec("CREATE TABLE migrations (name TEXT NOT NULL)")
			if err != nil {
				return fmt.Errorf("error creating migrations table | %w", err)
			}
		} else {
			return fmt.Errorf("error checking if migrations table exists | %w", err)
		}
	}

	for _, migration := range t.migrations {
		var migrationInHistory string
		err = t.DB.QueryRow("SELECT name FROM migrations WHERE name = ?", migration.name).Scan(&migrationInHistory)
		if err != nil {
			if err == sql.ErrNoRows {
				_, err := t.DB.Exec(migration.query)
				if err != nil {
					return fmt.Errorf("error running migration: %s | %w", migration.name, err)
				}
				_, err = t.DB.Exec("INSERT INTO migrations (name) VALUES (?)", migration.name)
				if err != nil {
					return fmt.Errorf("error inserting migration: %s into migrations table | %w", migration.name, err)
				}
			} else {
				return fmt.Errorf("error checking if migration: %s has been run | %w", migration.name, err)
			}
		}
	}
	return nil
}

// WithLocalDBName sets the local database name for the embedded database
func WithLocalDBName(localDBName string) Options {
	return func(l *LibSqlDB) error {
		l.localDBName = localDBName
		return nil
	}
}

// WithSyncInterval sets the sync interval for the embedded database
func WithSyncInterval(syncInterval time.Duration) Options {
	return func(l *LibSqlDB) error {
		l.syncInterval = syncInterval
		return nil
	}
}

// WithDir sets the directory for the embedded database
func WithDir(dir string) Options {
	return func(l *LibSqlDB) error {
		l.dir = dir
		return nil
	}
}

// WithAuthToken sets the auth token for the database
func WithAuthToken(authToken string) Options {
	return func(l *LibSqlDB) error {
		l.authToken = authToken
		return nil
	}
}

// WithEncryptionKey sets the encryption key for the embedded database
func WithEncryptionKey(key string) Options {
	return func(l *LibSqlDB) error {
		l.encryptionKey = key
		return nil
	}
}

// WithReadYourWrites sets the encryption key for the embedded database
func WithReadYourWrites(readYourWrites bool) Options {
	return func(l *LibSqlDB) error {
		l.readYourWrites = &readYourWrites
		return nil
	}
}

func WithMigrationFiles(migrationFiles embed.FS) Options {
	return func(l *LibSqlDB) error {
		l.useMigrations = true
		l.migrationFiles = migrationFiles
		return nil
	}
}
