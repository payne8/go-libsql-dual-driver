# Go LibSQL DualDriver

Welcome to the Go LibSQL DualDriver! This library provides a seamless interface for connecting to LibSQL databases, whether you’re working on Windows or other platforms. By automatically switching between the remote LibSQL driver for Windows and the embedded driver for other platforms, this wrapper ensures that your development environment is as close to production as possible.

## Key Features

- **Cross-Platform Support:** Automatically chooses the appropriate driver—remote for Windows, embedded for other systems.
- **Built-in Migration Tool:** Simplifies database migrations by allowing you to embed SQL files directly in your Go code.

## Getting Started

### Installation

To start using this library, simply import it and initialize the database connection in your Go application. The library will handle the rest.

### Example Usage

Here’s a quick example of how to set up a database connection with embedded migration files:

```golang
//go:embed migrations/*.sql
var migrationFiles embed.FS

tdb, err := libsqldb.NewLibSqlDB(
    primaryUrl,
    libsqldb.WithMigrationFiles(migrationFiles),
    libsqldb.WithAuthToken(authToken),
    libsqldb.WithLocalDBName("local.db"), // will not be used for remote-only
)
```

## Why Use This Library?
I developed this library to streamline the process of setting up a database connection in Go, complete with built-in migration capabilities. During development, I encountered issues with using the embedded driver on Windows. As I researched, I discovered [others had similar issues](https://github.com/tursodatabase/go-libsql/issues/30). This wrapper solves that problem by automatically selecting the appropriate driver based on your operating system.


### Special Note for Windows Users
This library defaults to using the embedded driver for better performance and closer parity with production environments. However, due to the lack of support for embedded LibSQL on Windows, this library uses a remote driver when running on Windows. Special build tags are included to ensure seamless operation across platforms.

Feel free to explore the repository and experiment with the examples provided. Happy coding!
