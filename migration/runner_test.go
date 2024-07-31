package migration

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRunMigrations(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS gomig_migrations`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT version FROM gomig_migrations`).WillReturnRows(sqlmock.NewRows([]string{"version"}))
	mock.ExpectExec(`CREATE TABLE users`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO gomig_migrations`).WillReturnResult(sqlmock.NewResult(1, 1))
	migrationContent := "-- Migration: create_users\nCREATE TABLE users (id SERIAL PRIMARY KEY);"
	os.MkdirAll("migrations", os.ModePerm)
	ioutil.WriteFile("migrations/0000000000_create_users.up.sql", []byte(migrationContent), 0644)
	defer os.RemoveAll("migrations")

	err = RunMigrations(db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
