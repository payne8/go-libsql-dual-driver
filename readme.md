# LibSQL Boilerplate Golang

Copy this repo to have a golang application that is already set up to work with libSQL.

## How to use

1. Copy this repo to your new project
2. Modify two files.

go.mod needs a new package name.
main.go needs to use the package name to import from `{{packageName}}/db`

3. Set up your environment variables

LIBSQL_AUTH_TOKEN
LIBSQL_DATABASE_URL