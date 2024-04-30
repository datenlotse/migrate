package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func PrintHelp() {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Print("THDS migration Tool usage\n\n")

	fmt.Fprintln(w, "create\tCreates a new migration file\t")
	fmt.Fprintln(w, "help\tShows this help\t")
	fmt.Fprintln(w, "revert\tReverts the last applied migration\t")
	fmt.Fprintln(w, "status\tPrints the status of the migrations\t")
	fmt.Fprintln(w, "up\tRuns all pending migrations\t")
	w.Flush()
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
