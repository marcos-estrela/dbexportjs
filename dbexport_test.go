package dbexport

import (
	"strings"
	"testing"
)

// Tests for tables
func TestGetAllTablesFromSchema(t *testing.T) {
	tables := GetObjectsFromSchema(TABLES, "")
	qtdTables := len(tables)
	if qtdTables < 2 {
		t.Errorf("Expected grater than 2 found %d", qtdTables)
	}
}

func TestGetOneTableFromSchema(t *testing.T) {
	tables := GetObjectsFromSchema(TABLES, "users")
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

// Tests for views
func TestGetAllViewsFromSchema(t *testing.T) {
	tables := GetObjectsFromSchema(VIEWS, "")
	qtdTables := len(tables)
	if qtdTables < 2 {
		t.Errorf("Expected grater than 2 found %d", qtdTables)
	}
}

func TestGetOneViewFromSchema(t *testing.T) {
	views := GetObjectsFromSchema(VIEWS, "vw_users")
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

// Tests for procedures
func TestGetAllProceduresFromSchema(t *testing.T) {
	objs := GetObjectsFromSchema(PROCEDURES, "")
	qtd := len(objs)
	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneProcedureFromSchema(t *testing.T) {
	objs := GetObjectsFromSchema(PROCEDURES, "sp_users_ins")
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

// Tests for functions
func TestGetAllFunctionsFromSchema(t *testing.T) {
	objs := GetObjectsFromSchema(FUNCTIONS, "")
	qtd := len(objs)
	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneFunctionFromSchema(t *testing.T) {
	objs := GetObjectsFromSchema(FUNCTIONS, "fn_users_exists")
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

// Tests for triggers
func TestGetAllTriggersFromSchema(t *testing.T) {
	objs := GetObjectsFromSchema(TRIGGERS, "")
	qtd := len(objs)
	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneTriggerFromSchema(t *testing.T) {
	objs := GetObjectsFromSchema(TRIGGERS, "tg_users_ins_after")
	qtd := len(objs)

	if qtd != 1 {
		t.Errorf("Expected 1 found %d", qtd)
	}
}

func TestGetTriggers(t *testing.T) {
	objs := GetTriggers("")
	qtd := len(objs)

	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneTrigger(t *testing.T) {
	objs := GetTriggers("tg_users_ins_after")
	qtd := len(objs)

	if qtd != 1 {
		t.Errorf("Expected 1 found %d", qtd)
	}
}

func TestExportTrigger(t *testing.T) {
	objName := "tg_users_ins_after"
	sql := GetSqlForTrigger(objName)

	if !strings.Contains(strings.ToUpper(sql), "CREATE TRIGGER") {
		t.Errorf("Expected contains TRIGGER actual %s", sql)
	}
}

// Tests for events
func TestGetAllEventsFromSchema(t *testing.T) {
	objs := GetObjectsFromSchema(EVENTS, "")
	qtd := len(objs)
	if qtd < 2 {
		t.Errorf("Expected grater than 2 found %d", qtd)
	}
}

func TestGetOneEventFromSchema(t *testing.T) {
	objs := GetObjectsFromSchema(EVENTS, "ev_delete_test_data")
	qtd := len(objs)

	if qtd != 1 {
		t.Errorf("Expected 1 found %d", qtd)
	}
}

func TestGetEvents(t *testing.T) {
	objs := GetEvents("")
	qtd := len(objs)

	if qtd < 1 {
		t.Errorf("Expected grater than 1 found %d", qtd)
	}
}

func TestGetOneEvent(t *testing.T) {
	objs := GetEvents("ev_delete_test_data")
	qtd := len(objs)

	if qtd != 1 {
		t.Errorf("Expected 1 found %d", qtd)
	}
}

func TestExportEvent(t *testing.T) {
	objName := "ev_delete_test_data"
	sql := GetSqlForEvent(objName)

	if !strings.Contains(strings.ToUpper(sql), "CREATE EVENT") {
		t.Errorf("Expected contains CREATE EVENT actual %s", sql)
	}
}
