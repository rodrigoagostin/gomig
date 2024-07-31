package migration

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	migrationName := "create_users"
	columns := []string{"name:varchar(200)", "email:varchar(200)", "active:boolean"}

	err := Generate(migrationName, columns)
	assert.NoError(t, err, "Error generating migration")

	expectedFilename := generateExpectedFilename(migrationName)
	_, err = os.Stat(expectedFilename)
	assert.NoError(t, err, "Migration file not created")

	os.Remove(expectedFilename)
}

func generateExpectedFilename(migrationName string) string {
	timestamp := time.Now().Format("20060102150405")
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(migrationName+timestamp)))
	return filepath.Join("migrations", hash[:10]+"_"+migrationName+".up.sql")
}
