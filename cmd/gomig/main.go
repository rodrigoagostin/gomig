package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rodrigoagostin/gomig/migration"
	"github.com/spf13/cobra"
)

var dbURL string

func main() {
	rootCmd := &cobra.Command{Use: "gomig"}

	generateCmd := &cobra.Command{
		Use:   "generate [migration_name] [columns...]",
		Short: "Generate a new migration file",
		Args:  cobra.MinimumNArgs(1),
		Run:   generateMigration,
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run all pending migrations",
		Run:   runMigrations,
	}

	rootCmd.PersistentFlags().StringVar(&dbURL, "db", "", "Database connection URL")
	rootCmd.MarkPersistentFlagRequired("db")

	rootCmd.AddCommand(generateCmd, migrateCmd)
	rootCmd.Execute()
}

func generateMigration(cmd *cobra.Command, args []string) {
	migrationName := args[0]
	columns := args[1:]
	err := migration.Generate(migrationName, columns)
	if err != nil {
		fmt.Println("Error generating migration:", err)
		os.Exit(1)
	}
}

func runMigrations(cmd *cobra.Command, args []string) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	err = migration.RunMigrations(db)
	if err != nil {
		fmt.Println("Error running migrations:", err)
		os.Exit(1)
	}
}
