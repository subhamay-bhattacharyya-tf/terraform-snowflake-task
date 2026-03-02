// File: test/single_task_test.go
package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// TestSingleTask tests creating a single task via the module
func TestSingleTask(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())
	dbName := fmt.Sprintf("TT_DB_%s", unique)
	schemaName := fmt.Sprintf("TT_SCHEMA_%s", unique)
	whName := fmt.Sprintf("TT_WH_%s", unique)
	taskName := fmt.Sprintf("TT_TASK_%s", unique)

	tfDir := "../examples/single-task"

	// Setup: Create database, schema, and warehouse
	db := openSnowflake(t)
	createTestDatabase(t, db, dbName)
	createTestSchema(t, db, dbName, schemaName)
	createTestWarehouse(t, db, whName)
	defer func() {
		dropTestDatabase(t, db, dbName)
		dropTestWarehouse(t, db, whName)
		_ = db.Close()
	}()

	taskConfigs := map[string]interface{}{
		"test_task": map[string]interface{}{
			"database":         dbName,
			"schema":           schemaName,
			"name":             taskName,
			"warehouse":        whName,
			"sql_statement":    "SELECT 1",
			"schedule_minutes": 60,
			"afters":           []string{},
			"started":          false,
			"comment":          "Terratest single task test",
			"grants":           []interface{}{},
		},
	}

	tfOptions := &terraform.Options{
		TerraformDir: tfDir,
		NoColor:      true,
		Vars: map[string]interface{}{
			"task_configs":                taskConfigs,
			"snowflake_organization_name": os.Getenv("SNOWFLAKE_ORGANIZATION_NAME"),
			"snowflake_account_name":      os.Getenv("SNOWFLAKE_ACCOUNT_NAME"),
			"snowflake_user":              os.Getenv("SNOWFLAKE_USER"),
			"snowflake_role":              os.Getenv("SNOWFLAKE_ROLE"),
			"snowflake_private_key":       os.Getenv("SNOWFLAKE_PRIVATE_KEY"),
		},
	}

	defer terraform.Destroy(t, tfOptions)
	terraform.InitAndApply(t, tfOptions)

	time.Sleep(retrySleep)

	exists := taskExists(t, db, dbName, schemaName, taskName)
	require.True(t, exists, "Expected task %q to exist in Snowflake", taskName)

	props := fetchTaskProps(t, db, dbName, schemaName, taskName)
	require.Equal(t, taskName, props.Name)
	require.Contains(t, props.Comment, "Terratest single task test")
}
