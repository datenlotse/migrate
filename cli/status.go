package cli

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"text/tabwriter"

	"github.com/jmoiron/sqlx"
	"tech.thds.migrate/migration"
)

func PrintStatus(db *sqlx.DB, migrationFiles []string) {
	appliedMigrations, err := migration.GetAppliedMigrationNames(db)
	if err != nil {
		log.Panicf("Failed to fetch already applied migrations\n%v", err)
	}

	sort.Strings(migrationFiles)
	var parsedMigrations []*migration.Migration
	for _, migrationRaw := range migrationFiles {
		m, err := migration.FromFile(migrationRaw)
		if err != nil {
			log.Panicf("Error parsing migration file %s", migrationRaw)
		}
		parsedMigrations = append(parsedMigrations, m)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)
	fmt.Fprintln(w, "\nMigration\tStatus")
	for _, m := range parsedMigrations {
		i := slices.Index(appliedMigrations, m.Name)
		if i < 0 {
			fmt.Fprintf(w, "%s\t\u001b[31mPending\u001b[0m\t\n", m.Name)
		} else {
			fmt.Fprintf(w, "%s\t\u001b[32mApplied\u001b[0m\t\n", m.Name)
		}
	}
	w.Flush()
}

func PrintStatusHelp() {
	fmt.Println(`
Command: STATUS

--- Description --- 
Prints the status of all migrations

--- Usage --- 
status`)
}
