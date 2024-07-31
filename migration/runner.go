package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

// RunMigrations executes all pending migrations.
func RunMigrations(db *sql.DB) error {
	err := createMigrationsTable(db)
	if err != nil {
		return err
	}

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	for _, file := range files {
		if !appliedMigrations[file.Name()] {
			err := applyMigration(db, file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS gomig_migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}
	return nil
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM gomig_migrations")
	if err != nil {
		return nil, fmt.Errorf("error querying migrations table: %v", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("error scanning migrations table: %v", err)
		}
		appliedMigrations[version] = true
	}

	return appliedMigrations, nil
}

func applyMigration(db *sql.DB, file os.FileInfo) error {
	content, err := ioutil.ReadFile(filepath.Join("migrations", file.Name()))
	if err != nil {
		return fmt.Errorf("error reading migration file %s: %v", file.Name(), err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("error executing migration %s: %v", file.Name(), err)
	}

	_, err = db.Exec("INSERT INTO gomig_migrations (version) VALUES ($1)", file.Name())
	if err != nil {
		return fmt.Errorf("error recording migration %s: %v", file.Name(), err)
	}

	fmt.Println("Migration applied:", file.Name())
	return nil
}
