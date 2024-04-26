package cli

import "fmt"

func PrintHelp() {
	fmt.Printf(`
THDS migration Tool usage:

--- Commands --- 
create - Creates a new migration file
help - Shows this help
revert - Reverts the last run migration
up - Runs all pending migrations
`)
}

func PrintCreateHelp() {
	fmt.Printf(`
Command: CREATE

--- Description --- 
Creates a new migration file

--- Usage --- 
create <migration-description>
`)
}
