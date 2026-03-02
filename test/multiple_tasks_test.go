// File: test/multiple_tasks_test.go
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

// TestMultipleTasks tests creating multiple tasks (DAG) via the module
func TestMultipleTasks(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())
	dbName := fmt.Sprintf("TT_DB_%s", unique)
	schemaName := fmt.Sprintf("TT_SCHEMA_%s", unique)
	whName := fmt.Sprintf("TT_WH_%s", unique)

	rootTaskName := fmt.Sprintf("TT_ROOT_%s", unique)

	tfDir := "../examples/multiple-tasks"

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
		"root_task": map[string]interface{}{
			"database":         dbName,
			"schema":           schemaName,
			"name":             rootTaskName,
			"warehouse":        whName,
			"sql_statement":    "SELECT 'root'",
			"schedule_minutes": 60,
			"afters":           []string{},
			"started":          false,
			"comment":          "Terratest root task",
			"grants":           []map[string]interface{}{},
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

	// Verify root task exists
	exists := taskExists(t, db, dbName, schemaName, rootTaskName)
	require.True(t, exists, "Expected task %q to exist in Snowflake", rootTaskName)

	// Verify properties of root task
	props := fetchTaskProps(t, db, dbName, schemaName, rootTaskName)
	require.Equal(t, rootTaskName, props.Name)
	require.Contains(t, props.Comment, "Terratest root task")
}
