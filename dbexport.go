package dbexport

import (
	"database/sql"
	"fmt"
)

type DbObject struct {
	Type    string
	Name    string
	Content string
}

type CreateTable struct {
	Table       string
	CreateTable string
}

type ViewDefinition struct {
}

type DatabaseAdapter interface {
	GetTables(tableName string) []DbObject
	GetViews(viewName string) []DbObject
	GetProcedures(procedureName string) []DbObject
	GetTriggers(triggerName string) []DbObject
	GetFunctions(functionName string) []DbObject
	GetEvents(eventName string) []DbObject
}

func GetViews(viewName string) []DbObject {
	return GetDbObjectsFor("view", viewName)
}

func GetTables(tableName string) []DbObject {
	return GetDbObjectsFor("table", tableName)
}

func GetDbObjectsFor(objType, objName string) []DbObject {
	var funcFromSchema func(string) []string
	var funcSqlFor func(string) string
	objects := []DbObject{}
	dbObject := DbObject{}

	if objType == "table" {
		funcFromSchema = GetTablesFromSchema
		funcSqlFor = GetSqlForTable
	} else if objType == "view" {
		funcFromSchema = GetViewsFromSchema
		funcSqlFor = GetSqlForView
	}

	objNames := funcFromSchema(objName)

	for _, objName := range objNames {
		dbObject.Name = objName
		dbObject.Content = funcSqlFor(objName)
		dbObject.Type = objType
		objects = append(objects, dbObject)
	}

	return objects

}

func GetTablesFromSchema(tableName string) []string {
	whereTableName := ""
	var results *sql.Rows
	var tables []string

	if tableName != "" {
		whereTableName = " AND TABLE_NAME = ?"
	}

	sql := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = ? AND TABLE_SCHEMA = ? " + whereTableName

	if tableName != "" {
		results = ExecuteQuery(sql, "BASE TABLE", "george", tableName)
	} else {
		results = ExecuteQuery(sql, "BASE TABLE", "george")
	}

	for results.Next() {
		var tableName string
		// for each row, scan the result into our tag composite object
		err := results.Scan(&tableName)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		tables = append(tables, tableName)
	}

	return tables
}

func GetSqlForTable(tableName string) string {
	var createStatement CreateTable
	query := fmt.Sprintf("SHOW CREATE TABLE %s", tableName)
	result := ExecuteQueryRow(query)
	err := result.Scan(&createStatement.Table, &createStatement.CreateTable)
	if err != nil {
		panic(err.Error())
	}
	return createStatement.CreateTable
}

func GetViewsFromSchema(viewName string) []string {
	whereViewName := ""
	var results *sql.Rows
	var views []string

	if viewName != "" {
		whereViewName = " AND TABLE_NAME = ?"
	}

	sql := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ? " + whereViewName

	if viewName != "" {
		results = ExecuteQuery(sql, "george", viewName)
	} else {
		results = ExecuteQuery(sql, "george")
	}

	for results.Next() {
		var viewName string
		// for each row, scan the result into our tag composite object
		err := results.Scan(&viewName)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		views = append(views, viewName)
	}

	return views
}

func GetSqlForView(viewName string) string {
	var definition string
	query := fmt.Sprintf("SELECT VIEW_DEFINITION FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?")
	result := ExecuteQueryRow(query, "george", viewName)
	err := result.Scan(&definition)
	if err != nil {
		panic(err.Error())
	}

	definition = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS\n%s", viewName, definition)
	return definition
}
