package cli

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/datenlotse/migrate/database"
	"github.com/datenlotse/migrate/migration"
	"github.com/jmoiron/sqlx"
)

func UpAllMigrations(db *sqlx.DB, migrationFiles []string, migrationDir string) {
	exists, err := checkIfMigrationsHistoryTableExists(db)
	if err != nil {
		log.Panicf("Error checking if migration history table exists %v", err)
	}

	if !exists {
		err := createMigrationsTable(db)
		if err != nil {
			log.Panicf("Error creating migration history table: %v", err)
		}
	}

	for _, migrationFile := range migrationFiles {
		migrationName := strings.Split(migrationFile, ".sql")[0]

		filePath := filepath.Join(migrationDir, migrationFile)
		m, err := migration.FromFile(filePath)
		if err != nil {
			// TODO: Cannot parse migration
			log.Panicf("Cannot parse migration: %s\n%v", migrationName, err)
		}

		if utf8.RuneCountInString(migrationFile) > 255 {
			// TODO: Handle file to long error
			log.Panicf("File name %s to long", migrationFile)
		}
		alreadyApplied, err := migrationIsAlreadyApplied(db, migrationName)
		if err != nil {
			log.Panicf("Cannot get migration status: %v", err)
		}
		if alreadyApplied {
			continue
		}

		var q database.Querier
		var tx *sqlx.Tx
		if m.IsTransactionMigration {
			tx, err = db.Beginx()
			q = tx
		} else {
			q = db
		}

		if err != nil {
			// TODO: Handle Transaction creation failure
			log.Panicf("Error creating transaction:\n%v", err)
		}

		err = migration.RunSql(*m.Up, q, migrationName)
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			log.Panicf("Error running migration %s\n%v", migrationName, err)
		}

		err = migration.InsertMigrationIntoHistory(q, migrationName)
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			log.Panicf("Migration insertion was not successfull\n%v", migrationName)
		}

		if tx != nil {
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				log.Panicf("Error running migration %s\n%v", migrationName, err)
			}
		}

		fmt.Printf("\u001b[32mMigration %s applied\u001b[0m\n", migrationName)
	}

	fmt.Println("\n\u001b[32mAll pending migrations applied\u001b[0m")
}

// Checks if a migration is already applied
func migrationIsAlreadyApplied(db database.Querier, name string) (bool, error) {
	var result *string
	err := db.Get(&result, `SELECT MigrationId FROM dbo.MigrationsHistory WHERE MigrationId = @p1`, name)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	fmt.Printf("\u001b[33mMigration %s already applied -> Skipping\u001b[0m\n", name)
	return true, nil
}

// Checks if the migrations history table exists
func checkIfMigrationsHistoryTableExists(db database.Querier) (bool, error) {
	var id *int
	err := db.Get(&id, `SELECT OBJECT_ID('dbo.MigrationsHistory', 'U') AS 'ObjectId'`)
	if err != nil {
		return false, err
	}

	return id != nil, nil
}

// Creates the migrations history table
func createMigrationsTable(db database.Querier) error {
	_, err := db.Exec(`
		IF OBJECT_ID('dbo.MigrationsHistory', 'U') IS NULL
		-- Create the table in the specified schema
		CREATE TABLE dbo.MigrationsHistory
		(
			MigrationId NVARCHAR(255)
		);
	`)

	return err
}
