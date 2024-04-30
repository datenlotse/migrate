package migration

import (
	"log"
	"strings"

	"tech.thds.migrate/database"
)

func RunSql(c string, db database.Querier, name string) error {
	log.Printf("Running Migration %s\n\n", name)
	batches := strings.Split(c, "\nGO")
	for _, s := range batches {
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}
	return nil
}
