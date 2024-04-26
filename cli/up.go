package cli

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"tech.thds.migrate/database"
	"tech.thds.migrate/migration"
)

func UpAllMigrations(db *sqlx.DB, migrationFiles []string, migrationDir string) {
	for _, migrationFile := range migrationFiles {
		filePath := filepath.Join(migrationDir, migrationFile)
		c, err := os.ReadFile(filePath)
		if err != nil {
			log.Panicln(err)
		}
		m, err := migration.FromFile(filePath)
		fmt.Printf("%p", m)
		continue

		if utf8.RuneCountInString(migrationFile) > 255 {
			// TODO: Handle file to long error
			log.Panicf("File name %v to long", migrationFile)
		}
		migrationName := strings.Split(migrationFile, ".sql")[0]
		alreadyApplied, err := migrationIsAlreadyApplied(db, migrationName)
		if err != nil {
			log.Panic(err)
			// TODO: Handle cannot get Migration status
		}
		if alreadyApplied {
			continue
		}

		s := string(c)

		var q database.Querier
		var tx *sqlx.Tx
		if !isNoTransactionMigration(s) {
			tx, err = db.Beginx()
			q = tx
		} else {
			q = db
		}

		if err != nil {
			// TODO: Handle Transaction creation failure
			log.Panicf("Error creating transaction:\n%v", err)
		}
		runSql(s, q, migrationName)
		err = insertMigrationIntoHistory(q, migrationName)
		if err != nil {
			tx.Rollback()
			log.Panicf("Migration insertion was not successfull\n%v", migrationName)
		}

		if tx != nil {
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				log.Panicf("Error running migration %s\n%v", migrationName, err)
			}
		}

		fmt.Printf("Migration %s completed", migrationName)
	}
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

	log.Printf("Migration %s already applied -> Skipping", name)
	return true, nil
}

// Runs the SQL of the provided string. It splits at "GO" into batches
func runSql(c string, db database.Querier, name string) {
	fmt.Printf("Running Migration %s\n\n", name)
	batches := strings.Split(c, "\nGO")
	for _, s := range batches {
		r, err := db.Exec(s)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("Result: %v\n", r)
	}
}

// Checks if the migration should run as transaction or not
func isNoTransactionMigration(c string) bool {
	return strings.HasPrefix(c, "--THDS: No transaction")
}

func insertMigrationIntoHistory(db database.Querier, migrationName string) error {
	_, err := db.Exec(`INSERT INTO dbo.MigrationsHistory (MigrationId) VALUES (@p1)`, migrationName)
	return err
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
