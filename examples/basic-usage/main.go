package main

import (
	"embed"
	"fmt"
	"git.sa.vin/payne8/libsqldb"
	"os"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	primaryUrl := os.Getenv("LIBSQL_DATABASE_URL")
	authToken := os.Getenv("LIBSQL_AUTH_TOKEN")

	// Open the database
	tdb, err := libsqldb.NewLibSqlDB(
		primaryUrl,
		migrationFiles,
		libsqldb.WithAuthToken(authToken),
		libsqldb.WithLocalDBName("local.db"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", primaryUrl, err)
		os.Exit(1)
	}

	// Migrate the database
	err = tdb.Migrate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to migrate db %s: %s", primaryUrl, err)
		os.Exit(1)
	}

	defer tdb.Close()
}
