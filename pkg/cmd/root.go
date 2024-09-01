package cmd

import (
	"fmt"
	"os"
)

// Execute is the entry point for the CLI.
func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migration":
		/*
			if len(os.Args) < 3 {
				fmt.Println("Usage: mycli migration <filename>")
				os.Exit(1)
			}
			//filename := os.Args[2]
			//runMigration(filename)
		*/
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
