package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/datenlotse/migrate/database"
)

var ErrFailedToSplit = fmt.Errorf("failed to split migration.\nDoes it contain a --THDS:Up and --THDS:Down comment in that order?")
var NoTransactionRegexp = regexp.MustCompile(`/--\s?THDS: No transaction/img`)

type Migration struct {
	Up                     *string
	Down                   *string
	Name                   string
	IsTransactionMigration bool
}

// Creates a migration struct from the provided path
func FromFile(path string) (*Migration, error) {
	fileName := filepath.Base(path)
	migrationName := strings.Split(fileName, ".sql")[0]

	contents, err := os.ReadFile(path)
	contentsString := string(contents)
	if err != nil {
		return nil, err
	}

	// Split the migration into Up and Down
	splits := strings.Split(contentsString, "\n--THDS:Down\n")
	if len(splits) != 2 {
		return nil, ErrFailedToSplit
	}

	upR := strings.Replace(splits[0], "--THDS:Up\n", "", -1)
	m := &Migration{
		Up:                     &upR,
		Down:                   &splits[1],
		Name:                   migrationName,
		IsTransactionMigration: isTransactionMigration(contentsString),
	}

	return m, nil
}

// Checks if the migration should run as transaction or not
// Returns true if the transaction has the "No transaction" annotation
func isTransactionMigration(c string) bool {
	return !NoTransactionRegexp.MatchString(c)
}

func InsertMigrationIntoHistory(db database.Querier, migrationName string) error {
	_, err := db.Exec(`INSERT INTO dbo.MigrationsHistory (MigrationId) VALUES (@p1)`, migrationName)
	return err
}

func RemoveMigrationFromHistory(db database.Querier, migrationName string) error {
	_, err := db.Exec(`DELETE FROM dbo.MigrationsHistory WHERE MigrationId = @p1`, migrationName)
	return err
}

func GetAppliedMigrationNames(db database.Querier) ([]string, error) {
	appliedMigrationNames := make([]string, 0)
	err := db.Select(&appliedMigrationNames, "SELECT MigrationId FROM dbo.MigrationsHistory")
	if err != nil {
		if err == sql.ErrNoRows {
			return appliedMigrationNames, nil
		}
		return appliedMigrationNames, err
	}

	return appliedMigrationNames, err
}
