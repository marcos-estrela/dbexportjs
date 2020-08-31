package dbexport

import (
	"strings"
	"testing"
)

func TestGetTables(t *testing.T) {
	tables := GetTables()
	qtdTables := len(tables)

	if (qtdTables != 1) {
		t.Errorf("Expected 1 found %d", qtdTables)
	}
}

func TestExportTable(t *testing.T) {
	tableName := "users"
	tableSql := GetSqlForTable(tableName)

	if !strings.Contains(strings.ToUpper(tableSql), "CREATE TABLE") {
		t.Errorf("Expected contains CREATE TABLE actual %s", tableSql)
	}
}
