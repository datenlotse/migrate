package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Creates a new migration file in the migrationsDir
func CreateMigrationFile(description string, migrationsDir string) {
	now := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s.sql", now, description)
	fp := filepath.Join(migrationsDir, fileName)

	fh, err := os.Create(fp)
	if err != nil {
		log.Panicf("Error creating file:\n%v", err)
	}
	defer fh.Close()

	fh.Write([]byte("--THDS:Up\n--Implement up part of migration here \n\n--THDS:Down\n--Implement down part of migration here"))
}
