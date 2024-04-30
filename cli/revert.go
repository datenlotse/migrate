package cli

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/datenlotse/migrate/database"
	"github.com/datenlotse/migrate/migration"
	"github.com/jmoiron/sqlx"
)

func RevertLastMigration(db *sqlx.DB, migrationsDir string) {
	lastMigrationName, err := GetLastAppliedMigrationName(db)
	if err != nil {
		log.Panicf("Error getting last applied migration\n%v", err)
	}

	if lastMigrationName == nil {
		log.Println("No migration applied yet")
		return
	}

	fileName := fmt.Sprintf("%s.sql", *lastMigrationName)
	filePath := filepath.Join(migrationsDir, fileName)
	m, err := migration.FromFile(filePath)
	if err != nil {
		log.Panicf("Error parsing migration file: %s\n%v", *lastMigrationName, err)
	}

	// TODO: Check if it is a transaction migration or not
	var q database.Querier
	var tx *sqlx.Tx
	if m.IsTransactionMigration {
		tx, err = db.Beginx()
		q = tx
	} else {
		q = db
	}

	downPart := *m.Down
	if err != nil {
		// TODO: Handle Transaction creation failure
		log.Panicf("Error creating transaction:\n%v", err)
	}

	err = migration.RunSql(downPart, q, *lastMigrationName)
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		log.Panicf("Error running migration %s\n%v", *lastMigrationName, err)
	}

	err = migration.RemoveMigrationFromHistory(q, *lastMigrationName)
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		log.Panicf("Migration removal from history was not successfull\n%v", *lastMigrationName)
	}

	// Apply the migration, if in transaction mode
	if tx != nil {
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			log.Panicf("Error running migration %s\n%v", *lastMigrationName, err)
		}
	}

	fmt.Printf("\u001b[32mSuccessfully reverted %s\u001b[0m", *lastMigrationName)
}

func GetLastAppliedMigrationName(db database.Querier) (*string, error) {
	var name *string
	err := db.Get(&name, "SELECT TOP 1 MigrationId FROM dbo.MigrationsHistory ORDER BY MigrationId DESC")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return name, nil
}
