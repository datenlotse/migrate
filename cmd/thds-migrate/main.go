package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/datenlotse/migrate/cli"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	if len(os.Args) < 2 {
		cli.PrintHelp()
		os.Exit(0)
	}

	cString := os.Getenv("THDS_MIGRATION_CONNECTION_STRING")
	if cString == "" {
		panic("Environment variable THDS_MIGRATION_CONNECTION_STRING is not set")
	}

	db, err := sqlx.Open("sqlserver", cString)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

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

	switch os.Args[1] {
	case "up":
		cli.UpAllMigrations(db, fileNames, migrationDir)
	case "revert":
		cli.RevertLastMigration(db, migrationDir)
	case "create":
		if len(os.Args) < 3 {
			cli.PrintCreateHelp()
			break
		}
		description := os.Args[2]
		cli.CreateMigrationFile(description, migrationDir)
	case "status":
		var fileDirs []string
		for _, file := range fileNames {
			fileDirs = append(fileDirs, filepath.Join(migrationDir, file))
		}
		cli.PrintStatus(db, fileDirs)
	default:
		cli.PrintHelp()
	}
}
