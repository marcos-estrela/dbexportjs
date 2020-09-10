package dbexport

import (
	"database/sql"
	"fmt"
	"strings"
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

type CreateProcedure struct {
	Procedure           string
	SQLMode             string
	CreateProcedure     string
	CharacterSetClient  string
	CollationConnection string
	DatabaseCollation   string
}

type CreateFunction struct {
	Function            string
	SQLMode             string
	CreateFunction      string
	CharacterSetClient  string
	CollationConnection string
	DatabaseCollation   string
}

type CreateTrigger struct {
	TriggerName       string
	ActionStatement   string
	ActionTiming      string
	EventManipulation string
	EventObjectTable  string
	ActionOrientation string
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

func GetProcedures(procedureName string) []DbObject {
	return GetDbObjectsFor("procedure", procedureName)
}

func GetFunctions(functionName string) []DbObject {
	return GetDbObjectsFor("function", functionName)
}

func GetTriggers(triggerName string) []DbObject {
	return GetDbObjectsFor("trigger", triggerName)
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
	} else if objType == "procedure" {
		funcFromSchema = GetProceduresFromSchema
		funcSqlFor = GetSqlForProcedure
	} else if objType == "function" {
		funcFromSchema = GetFunctionsFromSchema
		funcSqlFor = GetSqlForFunction
	} else if objType == "trigger" {
		funcFromSchema = GetTriggersFromSchema
		funcSqlFor = GetSqlForTrigger
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

	if tableName != "" {
		whereTableName = " AND TABLE_NAME = ?"
	}

	sql := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = ? AND TABLE_SCHEMA = ? " + whereTableName

	if tableName != "" {
		results = ExecuteQuery(sql, "BASE TABLE", "george", tableName)
	} else {
		results = ExecuteQuery(sql, "BASE TABLE", "george")
	}

	return objectListFromResults(results)
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

	if viewName != "" {
		whereViewName = " AND TABLE_NAME = ?"
	}

	sql := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ? " + whereViewName

	if viewName != "" {
		results = ExecuteQuery(sql, "george", viewName)
	} else {
		results = ExecuteQuery(sql, "george")
	}

	return objectListFromResults(results)
}

func GetSqlForView(viewName string) string {
	var definition string
	query := fmt.Sprintf("SELECT VIEW_DEFINITION FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?")
	result := ExecuteQueryRow(query, "george", viewName)
	err := result.Scan(&definition)
	if err != nil {
		panic(err.Error())
	}

	definition = formatQuery(definition)

	definition = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS\n%s", viewName, definition)

	return definition
}

func formatQuery(query string) string {
	separedQuery := strings.Split(query, ",")
	query = strings.Join(separedQuery, ",\n ")

	separedQuery = strings.Split(query, "select")
	query = strings.Join(separedQuery, "select\n")

	separedQuery = strings.Split(query, "SELECT")
	query = strings.Join(separedQuery, "SELECT\n")

	separedQuery = strings.Split(query, "from")
	query = strings.Join(separedQuery, "\nfrom\n")

	separedQuery = strings.Split(query, "FROM")
	query = strings.Join(separedQuery, "\nFROM\n")

	return query
}

func GetProceduresFromSchema(procedureName string) []string {
	whereName := ""
	var results *sql.Rows

	if procedureName != "" {
		whereName = " AND ROUTINE_NAME = ?"
	}

	sql := "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = ?" + whereName

	if procedureName != "" {
		results = ExecuteQuery(sql, "george", "PROCEDURE", procedureName)
	} else {
		results = ExecuteQuery(sql, "george", "PROCEDURE")
	}

	return objectListFromResults(results)
}

func GetSqlForProcedure(procedureName string) string {
	var createDefinition CreateProcedure
	query := fmt.Sprintf("SHOW CREATE PROCEDURE %s", procedureName)
	result := ExecuteQueryRow(query)
	err := result.Scan(
		&createDefinition.Procedure,
		&createDefinition.SQLMode,
		&createDefinition.CreateProcedure,
		&createDefinition.CharacterSetClient,
		&createDefinition.CollationConnection,
		&createDefinition.DatabaseCollation,
	)

	if err != nil {
		panic(err.Error())
	}

	definition := formatProcedure(createDefinition.CreateProcedure)

	return definition
}

func formatProcedure(definition string) string {
	separatedDefinition := strings.Split(definition, "PROCEDURE")

	definition = "CREATE PROCEDURE " + separatedDefinition[1]

	return definition
}

func GetFunctionsFromSchema(functionName string) []string {
	whereName := ""
	var results *sql.Rows

	if functionName != "" {
		whereName = " AND ROUTINE_NAME = ?"
	}

	sql := "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = ?" + whereName

	if functionName != "" {
		results = ExecuteQuery(sql, "george", "FUNCTION", functionName)
	} else {
		results = ExecuteQuery(sql, "george", "FUNCTION")
	}

	return objectListFromResults(results)
}

func GetSqlForFunction(functionName string) string {
	var createDefinition CreateFunction
	query := fmt.Sprintf("SHOW CREATE FUNCTION %s", functionName)
	result := ExecuteQueryRow(query)
	err := result.Scan(
		&createDefinition.Function,
		&createDefinition.SQLMode,
		&createDefinition.CreateFunction,
		&createDefinition.CharacterSetClient,
		&createDefinition.CollationConnection,
		&createDefinition.DatabaseCollation,
	)

	if err != nil {
		panic(err.Error())
	}

	definition := formatFunction(createDefinition.CreateFunction)

	return definition
}

func formatFunction(definition string) string {
	separatedDefinition := strings.Split(definition, "FUNCTION")

	definition = "CREATE FUNCTION " + separatedDefinition[1]

	return definition
}

func GetTriggersFromSchema(triggerName string) []string {
	whereName := ""
	var results *sql.Rows

	if triggerName != "" {
		whereName = " AND TRIGGER_NAME = ?"
	}

	sql := "SELECT TRIGGER_NAME FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_SCHEMA = ?" + whereName

	if triggerName != "" {
		results = ExecuteQuery(sql, "george", triggerName)
	} else {
		results = ExecuteQuery(sql, "george")
	}

	return objectListFromResults(results)
}

func GetSqlForTrigger(triggerName string) string {
	var createTrigger CreateTrigger
	query := "SELECT TRIGGER_NAME, ACTION_STATEMENT, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORIENTATION FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_SCHEMA = ? AND TRIGGER_NAME = ?"
	result := ExecuteQueryRow(query, "george", triggerName)
	err := result.Scan(
		&createTrigger.TriggerName,
		&createTrigger.ActionStatement,
		&createTrigger.ActionTiming,
		&createTrigger.EventManipulation,
		&createTrigger.EventObjectTable,
		&createTrigger.ActionOrientation,
	)

	if err != nil {
		panic(err.Error())
	}

	definition := formatTrigger(createTrigger)

	return definition
}

func formatTrigger(trigger CreateTrigger) string {
	definitionStruct := `CREATE TRIGGER %s
	%s %s
	ON %s FOR EACH ROW
	%s`
	definition := fmt.Sprintf(definitionStruct, trigger.TriggerName, trigger.ActionTiming, trigger.EventManipulation, trigger.EventObjectTable, trigger.ActionStatement)
	return definition
}

func objectListFromResults(results *sql.Rows) []string {
	var objctNames []string

	for results.Next() {
		var objectName string
		// for each row, scan the result into our tag composite object
		err := results.Scan(&objectName)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		objctNames = append(objctNames, objectName)
	}

	return objctNames
}
