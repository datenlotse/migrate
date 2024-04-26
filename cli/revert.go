package cli

import "github.com/jmoiron/sqlx"

func RevertLastMigration(db *sqlx.DB, fileNames []string, migrationsDir string) {
	panic("Not implemented")
}
