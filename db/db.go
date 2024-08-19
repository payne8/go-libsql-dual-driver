package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

type Migrations struct {
	name  string
	query string
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

var migrations []Migrations

// func NewLibSqlDB is defined in embedded.go and remote-only.go files
// these files are used to define the LibSqlDB struct and the NewLibSqlDB function
// they have different initializations based on the environment, embedded or remote-only
// Windows does not currently support the embedded database, so the remote-only file is used

// setupMigrations initializes the filesystem and reads the migration files into the migrations variable
func setupMigrations() error {
	// Walk through the embedded files and read their contents
	err := fs.WalkDir(migrationFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			content, err := migrationFiles.ReadFile(path)
			if err != nil {
				return err
			}

			migration := Migrations{
				name:  filepath.Base(path),
				query: string(content),
			}
			migrations = append(migrations, migration)
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
	// check if migration table exists
	var migrationsCheck string
	//goland:noinspection SqlResolve
	err := t.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'").Scan(&migrationsCheck)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := t.db.Exec("CREATE TABLE migrations (name TEXT NOT NULL)")
			if err != nil {
				return fmt.Errorf("error creating migrations table | %w", err)
			}
		} else {
			return fmt.Errorf("error checking if migrations table exists | %w", err)
		}
	}

	for _, migration := range migrations {
		var migrationInHistory string
		err = t.db.QueryRow("SELECT name FROM migrations WHERE name = ?", migration.name).Scan(&migrationInHistory)
		if err != nil {
			if err == sql.ErrNoRows {
				_, err := t.db.Exec(migration.query)
				if err != nil {
					return fmt.Errorf("error running migration: %s | %w", migration.name, err)
				}
				_, err = t.db.Exec("INSERT INTO migrations (name) VALUES (?)", migration.name)
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
