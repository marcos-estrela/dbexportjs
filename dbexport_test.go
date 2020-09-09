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

func TestGetAllProceduresFromSchema(t *testing.T) {
	objs := GetProceduresFromSchema("")
	qtd := len(objs)
	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneProcedureFromSchema(t *testing.T) {
	objs := GetProceduresFromSchema("sp_users_ins")
	qtd := len(objs)

	if qtd != 1 {
		t.Errorf("Expected 1 found %d", qtd)
	}
}

func TestGetProcedures(t *testing.T) {
	objs := GetProcedures("")
	qtd := len(objs)

	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneProcedure(t *testing.T) {
	objs := GetProcedures("sp_users_ins")
	qtd := len(objs)

	if qtd != 1 {
		t.Errorf("Expected 1 found %d", qtd)
	}
}

func TestExportProcedure(t *testing.T) {
	objName := "sp_users_ins"
	sql := GetSqlForProcedure(objName)

	if !strings.Contains(strings.ToUpper(sql), "CREATE PROCEDURE") {
		t.Errorf("Expected contains CREATE PROCEDURE actual %s", sql)
	}
}

func TestGetAllFunctionsFromSchema(t *testing.T) {
	objs := GetFunctionsFromSchema("")
	qtd := len(objs)
	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneFunctionFromSchema(t *testing.T) {
	objs := GetFunctionsFromSchema("fn_users_exists")
	qtd := len(objs)

	if qtd != 1 {
		t.Errorf("Expected 1 found %d", qtd)
	}
}

func TestGetFunctions(t *testing.T) {
	objs := GetFunctions("")
	qtd := len(objs)

	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneFunction(t *testing.T) {
	objs := GetFunctions("fn_users_exists")
	qtd := len(objs)

	if qtd != 1 {
		t.Errorf("Expected 1 found %d", qtd)
	}
}

func TestExportFunction(t *testing.T) {
	objName := "fn_users_exists"
	sql := GetSqlForFunction(objName)

	if !strings.Contains(strings.ToUpper(sql), "CREATE FUNCTION") {
		t.Errorf("Expected contains CREATE FUNCTION actual %s", sql)
	}
}
