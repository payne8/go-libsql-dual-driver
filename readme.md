# LibSQL Boilerplate Golang

Copy this repo to have a golang application that is already set up to work with libSQL.

## How to use

1. Copy this repo to your new project
2. Modify two files.

   * go.mod needs a new package name.
   * main.go needs to use the package name to import from `{{packageName}}/db`

3. Set up your environment variables

   * LIBSQL_AUTH_TOKEN
   * LIBSQL_DATABASE_URL

## Windows?

This project is designed to use an embedded replica by default but Windows is not supported by libSQL. So this project has special build tags for Windows and sets up a libSQL remote driver. If your build target doesn't include Windows you can remote the file `db/remote-only.go`.