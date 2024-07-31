package migration

import (
	"crypto/sha1"
	"fmt"
	"os"
	"strings"
	"time"
)

func Generate(migrationName string, columns []string) error {
	timestamp := time.Now().Format("20060102150405")
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(migrationName+timestamp)))
	filename := fmt.Sprintf("migrations/%s_%s.up.sql", hash[:10], migrationName)

	err := os.MkdirAll("migrations", os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating migrations directory: %v", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating migration file: %v", err)
	}
	defer file.Close()

	sqlContent := generateSQL(migrationName, columns, timestamp)
	_, err = file.WriteString(sqlContent)
	if err != nil {
		return fmt.Errorf("error writing to migration file: %v", err)
	}

	fmt.Println("Migration file created:", filename)
	return nil
}

func generateSQL(migrationName string, columns []string, timestamp string) string {
	parts := strings.SplitN(migrationName, "_", 2)
	if len(parts) < 2 {
		return fmt.Sprintf("-- Invalid migration name format: %s\n-- Expected format: <action>_<table_name>\n", migrationName)
	}

	action := strings.ToLower(parts[0])
	tableName := parts[1]

	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n", migrationName, timestamp))

	switch action {
	case "create":
		sqlBuilder.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))
		sqlBuilder.WriteString("    id SERIAL PRIMARY KEY,\n")
		for _, column := range columns {
			columnDetails := strings.Split(column, ":")
			if len(columnDetails) == 2 {
				sqlBuilder.WriteString(fmt.Sprintf("    %s %s,\n", columnDetails[0], columnDetails[1]))
			}
		}
		sqlBuilder.WriteString("    created_at TIMESTAMPTZ DEFAULT NOW(),\n")
		sqlBuilder.WriteString("    updated_at TIMESTAMPTZ DEFAULT NOW()\n")
		sqlBuilder.WriteString(");\n")
	case "alter":
		sqlBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s\n", tableName))
		for i, column := range columns {
			columnDetails := strings.Split(column, ":")
			if len(columnDetails) == 2 {
				if i > 0 {
					sqlBuilder.WriteString(",\n")
				}
				sqlBuilder.WriteString(fmt.Sprintf("    ADD COLUMN %s %s", columnDetails[0], columnDetails[1]))
			}
		}
		sqlBuilder.WriteString(";\n")
	case "drop":
		sqlBuilder.WriteString(fmt.Sprintf("DROP TABLE %s;\n", tableName))
	case "rename":
		if len(columns) != 1 {
			return fmt.Sprintf("-- Invalid migration name format for rename: %s\n-- Expected format: rename_<old_table_name> <new_table_name>\n", migrationName)
		}
		newTableName := columns[0]
		sqlBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;\n", tableName, newTableName))
	case "drop_column":
		sqlBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s\n", tableName))
		for i, column := range columns {
			if i > 0 {
				sqlBuilder.WriteString(",\n")
			}
			sqlBuilder.WriteString(fmt.Sprintf("    DROP COLUMN %s", column))
		}
		sqlBuilder.WriteString(";\n")
	case "rename_column":
		if len(columns) != 2 {
			return fmt.Sprintf("-- Invalid migration name format for rename_column: %s\n-- Expected format: rename_column_<table_name> <old_column_name>:<new_column_name>\n", migrationName)
		}
		columnDetails := strings.Split(columns[0], ":")
		if len(columnDetails) != 2 {
			return fmt.Sprintf("-- Invalid column rename format: %s\n-- Expected format: <old_column_name>:<new_column_name>\n", columns[0])
		}
		oldColumnName := columnDetails[0]
		newColumnName := columnDetails[1]
		sqlBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;\n", tableName, oldColumnName, newColumnName))
	case "modify":
		sqlBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s\n", tableName))
		for i, column := range columns {
			columnDetails := strings.Split(column, ":")
			if len(columnDetails) == 2 {
				if i > 0 {
					sqlBuilder.WriteString(",\n")
				}
				sqlBuilder.WriteString(fmt.Sprintf("    ALTER COLUMN %s TYPE %s", columnDetails[0], columnDetails[1]))
			}
		}
		sqlBuilder.WriteString(";\n")
	case "add_index":
		if len(columns) != 2 {
			return fmt.Sprintf("-- Invalid migration name format for add_index: %s\n-- Expected format: add_index_<table_name> <index_name>:<column_name>\n", migrationName)
		}
		indexName := columns[0]
		columnName := columns[1]
		sqlBuilder.WriteString(fmt.Sprintf("CREATE INDEX %s ON %s (%s);\n", indexName, tableName, columnName))
	case "drop_index":
		if len(columns) != 1 {
			return fmt.Sprintf("-- Invalid migration name format for drop_index: %s\n-- Expected format: drop_index_<table_name> <index_name>\n", migrationName)
		}
		indexName := columns[0]
		sqlBuilder.WriteString(fmt.Sprintf("DROP INDEX %s;\n", indexName))
	default:
		sqlBuilder.WriteString(fmt.Sprintf("-- Unsupported action: %s\n", action))
	}

	return sqlBuilder.String()
}
