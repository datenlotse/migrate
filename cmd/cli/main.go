package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
)

var batchRegex = regexp.MustCompile(`/\\nGO(\\n)?/gmi`)

func main() {
	cString := "sqlserver://SA:Password123!@localhost"
	db, err := sql.Open("sqlserver", cString)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Connected to SQL Server at: " + cString)
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	migrationDir := filepath.Join(cwd, "migrations")
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Panic(err)
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}

	for _, migrationFile := range fileNames {
		log.Printf("%s", migrationFile)
		filePath := filepath.Join(migrationDir, migrationFile)
		c, err := os.ReadFile(filePath)
		if err != nil {
			log.Panicln(err)
		}
		stringContents := string(c)
		fmt.Printf("Running Migration %s\n\n", migrationFile)
		runSql(stringContents, db)
		fmt.Printf("Migration completed")
	}
}

func runSql(c string, db *sql.DB) {
	batches := strings.Split(c, "\nGO")
	fmt.Printf("batches: %v", batches[0])
	for _, s := range batches {
		r, err := db.Exec(s)
		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("Result: %v\n", r)
	}

}
