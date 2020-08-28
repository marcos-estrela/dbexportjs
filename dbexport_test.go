package dbexport

import (
	"testing"
	"os"
	"strings"
)

func TestExportTable(t *testing.T) {
	tableName := "users"
	tableSql := GetSqlForTable(tableName)

	if (!strings.Contains(strings.ToUpper(tableSql), "CREATE TABLE")) {
		t.Errorf("Expected contains CREATE TABLE actual %s", tableSql)
	}
}

func TestWriteSqlToFile(t *testing.T) {
	var sql = "CREATE OR REPLACE VIEW SELECT id, name, email FROM users;"
	var fileName = "vw_users.sql"
	var filePath = "/tmp/"+fileName

	WriteSqlToFile(filePath, sql)
	assertFileExists(t, filePath)
}

func assertFileExists(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)

	if (err != nil) {
		t.Errorf("file %s not exists", filePath)
	}

	e := os.Remove(filePath)
	if e != nil {
		t.Error(e)
	}
}
