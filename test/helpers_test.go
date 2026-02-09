// File: test/helpers_test.go
package test

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

type TaskProps struct {
	Name      string
	Database  string
	Schema    string
	State     string
	Schedule  string
	Comment   string
}

func openSnowflake(t *testing.T) *sql.DB {
	t.Helper()

	orgName := mustEnv(t, "SNOWFLAKE_ORGANIZATION_NAME")
	accountName := mustEnv(t, "SNOWFLAKE_ACCOUNT_NAME")
	user := mustEnv(t, "SNOWFLAKE_USER")
	privateKeyPEM := mustEnv(t, "SNOWFLAKE_PRIVATE_KEY")
	role := os.Getenv("SNOWFLAKE_ROLE")

	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	require.NotNil(t, block, "Failed to decode PEM block from private key")

	var privateKey *rsa.PrivateKey
	var err error

	// Try PKCS8 first, then PKCS1
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		require.NoError(t, err, "Failed to parse private key")
	} else {
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		require.True(t, ok, "Private key is not RSA")
	}

	// Build account identifier: orgname-accountname
	account := fmt.Sprintf("%s-%s", orgName, accountName)

	config := gosnowflake.Config{
		Account:       account,
		User:          user,
		Authenticator: gosnowflake.AuthTypeJwt,
		PrivateKey:    privateKey,
	}

	if role != "" {
		config.Role = role
	}

	dsn, err := gosnowflake.DSN(&config)
	require.NoError(t, err, "Failed to build DSN")

	db, err := sql.Open("snowflake", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

func taskExists(t *testing.T, db *sql.DB, database, schema, taskName string) bool {
	t.Helper()

	q := fmt.Sprintf("SHOW TASKS LIKE '%s' IN SCHEMA %s.%s;", escapeLike(taskName), database, schema)
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	return rows.Next()
}

func fetchTaskProps(t *testing.T, db *sql.DB, database, schema, taskName string) TaskProps {
	t.Helper()

	q := fmt.Sprintf("SHOW TASKS LIKE '%s' IN SCHEMA %s.%s;", escapeLike(taskName), database, schema)
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	cols, err := rows.Columns()
	require.NoError(t, err)

	// Find column indices
	nameIdx, dbIdx, schemaIdx, stateIdx, scheduleIdx, commentIdx := -1, -1, -1, -1, -1, -1
	for i, col := range cols {
		switch col {
		case "name":
			nameIdx = i
		case "database_name":
			dbIdx = i
		case "schema_name":
			schemaIdx = i
		case "state":
			stateIdx = i
		case "schedule":
			scheduleIdx = i
		case "comment":
			commentIdx = i
		}
	}
	require.NotEqual(t, -1, nameIdx, "name column not found in SHOW TASKS output")

	require.True(t, rows.Next(), "No task found matching %s", taskName)

	// Create slice to hold all column values
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	err = rows.Scan(valuePtrs...)
	require.NoError(t, err)

	// Extract the values we need
	getValue := func(idx int) string {
		if idx == -1 {
			return ""
		}
		v := values[idx]
		if v == nil {
			return ""
		}
		if s, ok := v.(string); ok {
			return s
		}
		if b, ok := v.([]byte); ok {
			return string(b)
		}
		return fmt.Sprintf("%v", v)
	}

	return TaskProps{
		Name:     getValue(nameIdx),
		Database: getValue(dbIdx),
		Schema:   getValue(schemaIdx),
		State:    getValue(stateIdx),
		Schedule: getValue(scheduleIdx),
		Comment:  getValue(commentIdx),
	}
}

func createTestDatabase(t *testing.T, db *sql.DB, dbName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	require.NoError(t, err)
}

func createTestSchema(t *testing.T, db *sql.DB, dbName, schemaName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s.%s", dbName, schemaName))
	require.NoError(t, err)
}

func createTestWarehouse(t *testing.T, db *sql.DB, whName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("CREATE WAREHOUSE IF NOT EXISTS %s WAREHOUSE_SIZE = 'X-SMALL' AUTO_SUSPEND = 60 AUTO_RESUME = TRUE INITIALLY_SUSPENDED = TRUE", whName))
	require.NoError(t, err)
}

func dropTestDatabase(t *testing.T, db *sql.DB, dbName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	require.NoError(t, err)
}

func dropTestWarehouse(t *testing.T, db *sql.DB, whName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("DROP WAREHOUSE IF EXISTS %s", whName))
	require.NoError(t, err)
}

func mustEnv(t *testing.T, key string) string {
	t.Helper()
	v := strings.TrimSpace(os.Getenv(key))
	require.NotEmpty(t, v, "Missing required environment variable %s", key)
	return v
}

func escapeLike(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
