package dbexport

import (
	"strings"
	"testing"
)

func TestExportTable(t *testing.T) {
	tableName := "users"
	tableSql := GetSqlForTable(tableName)

	if !strings.Contains(strings.ToUpper(tableSql), "CREATE TABLE") {
		t.Errorf("Expected contains CREATE TABLE actual %s", tableSql)
	}
}
