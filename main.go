package main

import (
	"fmt"
	"os"
)

func main() {
	primaryUrl := os.Getenv("LIBSQL_DATABASE_URL")
	authToken := os.Getenv("LIBSQL_AUTH_TOKEN")

	tdb, err := libsqlDB.NewLibSqlDB(
		primaryUrl,
		libsqlDB.WithAuthToken(authToken),
		libsqlDB.WithLocalDBName("local.db"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", primaryUrl, err)
		os.Exit(1)
	}

	err = tdb.Migrate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to migrate db %s: %s", primaryUrl, err)
		os.Exit(1)
	}

	defer tdb.Close()
}
