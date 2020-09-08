package dbexport

import (
	"strings"
	"testing"
)

func TestGetAllTablesFromSchema(t *testing.T) {
	tables := GetTablesFromSchema("")
	qtdTables := len(tables)
	if qtdTables < 2 {
		t.Errorf("Expected grater than 2 found %d", qtdTables)
	}
}

func TestGetOneTableFromSchema(t *testing.T) {
	tables := GetTablesFromSchema("users")
	qtdTables := len(tables)
	if qtdTables != 1 {
		t.Errorf("Expected 1 found %d", qtdTables)
	}
}

func TestGetTables(t *testing.T) {
	tables := GetTables("")
	qtdTables := len(tables)

	if qtdTables < 2 {
		t.Errorf("Expected grater than 2 found %d", qtdTables)
	}
}

func TestGetOneTables(t *testing.T) {
	tables := GetTables("users")
	qtdTables := len(tables)

	if qtdTables != 1 {
		t.Errorf("Expected 1 found %d", qtdTables)
	}
}

func TestExportTable(t *testing.T) {
	tableName := "users"
	tableSQL := GetSqlForTable(tableName)

	if !strings.Contains(strings.ToUpper(tableSQL), "CREATE TABLE") {
		t.Errorf("Expected contains CREATE TABLE actual %s", tableSQL)
	}
}

func TestGetAllViewsFromSchema(t *testing.T) {
	tables := GetViewsFromSchema("")
	qtdTables := len(tables)
	if qtdTables < 2 {
		t.Errorf("Expected grater than 2 found %d", qtdTables)
	}
}

func TestGetOneViewFromSchema(t *testing.T) {
	views := GetViewsFromSchema("vw_users")
	qtdViews := len(views)

	if qtdViews != 1 {
		t.Errorf("Expected 1 found %d", qtdViews)
	}
}

func TestGetViews(t *testing.T) {
	views := GetViews("")
	qtdviews := len(views)

	if qtdviews < 2 {
		t.Errorf("Expected grater than 2 found %d", qtdviews)
	}
}

func TestGetOneView(t *testing.T) {
	views := GetViews("vw_test_table")
	qtdViews := len(views)

	if qtdViews != 1 {
		t.Errorf("Expected 1 found %d", qtdViews)
	}
}

func TestExportView(t *testing.T) {
	tableName := "vw_test_table"
	sql := GetSqlForView(tableName)

	if !strings.Contains(strings.ToUpper(sql), "CREATE OR REPLACE VIEW") {
		t.Errorf("Expected contains CREATE OR REPLACE VIEW actual %s", sql)
	}
}
