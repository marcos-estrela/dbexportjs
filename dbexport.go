package dbexport

import (
	"database/sql"
	"fmt"
	"strings"
)

var TABLES string = "tables"
var VIEWS string = "views"
var PROCEDURES string = "procedures"
var FUNCTIONS string = "functions"
var TRIGGERS string = "triggers"
var EVENTS string = "events"

var queriesForSchema = map[string]string{
	TABLES:     "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = ?",
	VIEWS:      "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ? ",
	PROCEDURES: "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = ?",
	FUNCTIONS:  "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = ?",
	TRIGGERS:   "SELECT TRIGGER_NAME FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_SCHEMA = ?",
	EVENTS:     "SELECT EVENT_NAME FROM INFORMATION_SCHEMA.EVENTS WHERE EVENT_SCHEMA = ?",
}

var queriesForSchemaWhere = map[string]string{
	TABLES:     " AND TABLE_NAME = ?",
	VIEWS:      " AND TABLE_NAME = ?",
	PROCEDURES: " AND ROUTINE_NAME = ?",
	FUNCTIONS:  " AND ROUTINE_NAME = ?",
	TRIGGERS:   " AND TRIGGER_NAME = ?",
	EVENTS:     " AND EVENT_NAME = ?",
}

var sqlForFuncs = map[string]func(string) string{
	TABLES:     GetSqlForTable,
	VIEWS:      GetSqlForView,
	PROCEDURES: GetSqlForProcedure,
	FUNCTIONS:  GetSqlForFunction,
	TRIGGERS:   GetSqlForTrigger,
	EVENTS:     GetSqlForEvent,
}

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

type CreateEvent struct {
	EventName       string
	EventDefinition string
	ExecuteAt       interface{}
	IntervalValue   interface{}
	IntervalField   interface{}
	EventComment    string
	Status          string
	OnCompletion    string
	Starts          interface{}
	Ends            interface{}
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
	return GetDbObjectsFor(VIEWS, viewName)
}

func GetTables(tableName string) []DbObject {
	return GetDbObjectsFor(TABLES, tableName)
}

func GetProcedures(procedureName string) []DbObject {
	return GetDbObjectsFor(PROCEDURES, procedureName)
}

func GetFunctions(functionName string) []DbObject {
	return GetDbObjectsFor(FUNCTIONS, functionName)
}

func GetTriggers(triggerName string) []DbObject {
	return GetDbObjectsFor(TRIGGERS, triggerName)
}

func GetEvents(eventName string) []DbObject {
	return GetDbObjectsFor(EVENTS, eventName)
}

func GetAll() []DbObject {
	var objs []DbObject

	for _, obj := range GetTables("") {
		objs = append(objs, obj)
	}

	for _, obj := range GetViews("") {
		objs = append(objs, obj)
	}

	for _, obj := range GetProcedures("") {
		objs = append(objs, obj)
	}

	for _, obj := range GetFunctions("") {
		objs = append(objs, obj)
	}

	for _, obj := range GetTriggers("") {
		objs = append(objs, obj)
	}

	for _, obj := range GetEvents("") {
		objs = append(objs, obj)
	}

	return objs
}

func GetDbObjectsFor(objType, objName string) []DbObject {
	objects := []DbObject{}
	dbObject := DbObject{}
	objNames := GetObjectsFromSchema(objType, objName)

	for _, objName := range objNames {
		dbObject.Name = objName
		dbObject.Content = GetSqlForObject(objType, objName)
		dbObject.Type = objType
		objects = append(objects, dbObject)
	}

	return objects
}

func GetObjectsFromSchema(objectType string, objectName string) []string {
	var results *sql.Rows
	var searchParameters []interface{}

	query := queriesForSchema[objectType]

	searchParameters = append(searchParameters, "george")

	if objectType == TABLES {
		searchParameters = append(searchParameters, "BASE TABLE")
	}

	if objectType == PROCEDURES {
		searchParameters = append(searchParameters, "PROCEDURE")
	}

	if objectType == FUNCTIONS {
		searchParameters = append(searchParameters, "FUNCTION")
	}

	if objectName != "" {
		query += queriesForSchemaWhere[objectType]
		searchParameters = append(searchParameters, objectName)
	}

	results = ExecuteQuery(query, searchParameters...)

	return objectListFromResults(results)
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

func GetSqlForObject(objectType string, objectName string) string {
	funcSqlFor := sqlForFuncs[objectType]
	return funcSqlFor(objectName)
}

func GetSqlForTable(tableName string) string {
	query := fmt.Sprintf("SHOW CREATE TABLE %s", tableName)
	result := ExecuteQueryRow(query)
	return formatResultForTable(result)
}

func formatResultForTable(result *sql.Row) string {
	var createStatement CreateTable
	err := result.Scan(&createStatement.Table, &createStatement.CreateTable)
	if err != nil {
		panic(err.Error())
	}
	return createStatement.CreateTable
}

func GetSqlForView(viewName string) string {
	query := fmt.Sprintf("SELECT VIEW_DEFINITION FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?")
	result := ExecuteQueryRow(query, "george", viewName)
	return formatResultForView(result, viewName)
}

func formatResultForView(result *sql.Row, viewName string) string {
	var definition string
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

func GetSqlForProcedure(procedureName string) string {
	query := fmt.Sprintf("SHOW CREATE PROCEDURE %s", procedureName)
	result := ExecuteQueryRow(query)
	return formatResultForProcedure(result)
}

func formatResultForProcedure(result *sql.Row) string {
	var createDefinition CreateProcedure
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

func GetSqlForFunction(functionName string) string {
	query := fmt.Sprintf("SHOW CREATE FUNCTION %s", functionName)
	result := ExecuteQueryRow(query)
	return formatResultForFunction(result)
}

func formatResultForFunction(result *sql.Row) string {
	var createDefinition CreateFunction
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

func GetSqlForTrigger(triggerName string) string {
	query := "SELECT TRIGGER_NAME, ACTION_STATEMENT, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORIENTATION FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_SCHEMA = ? AND TRIGGER_NAME = ?"
	result := ExecuteQueryRow(query, "george", triggerName)
	return formatResultForTrigger(result)
}

func formatResultForTrigger(result *sql.Row) string {
	var createTrigger CreateTrigger
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

func GetSqlForEvent(eventName string) string {
	query := "SELECT EVENT_NAME, EVENT_DEFINITION, EXECUTE_AT, INTERVAL_VALUE, INTERVAL_FIELD, EVENT_COMMENT, STATUS, ON_COMPLETION, STARTS, ENDS FROM INFORMATION_SCHEMA.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?"
	result := ExecuteQueryRow(query, "george", eventName)
	return formatResultForEvent(result)
}

func formatResultForEvent(result *sql.Row) string {
	var event CreateEvent
	err := result.Scan(
		&event.EventName,
		&event.EventDefinition,
		&event.ExecuteAt,
		&event.IntervalValue,
		&event.IntervalField,
		&event.EventComment,
		&event.Status,
		&event.OnCompletion,
		&event.Starts,
		&event.Ends,
	)

	if err != nil {
		panic(err.Error())
	}

	definition := formatEvent(event)

	return definition
}

func formatEvent(event CreateEvent) string {
	eventSqlStruct := `CREATE EVENT %s
	ON SCHEDULE %s
	%s
	%s
	ON COMPLETION %s
	COMMENT '%s'
	DO
	%s;
	ALTER EVENT %s
	%s`
	var starts string = ""
	var ends string = ""

	schedule := fmt.Sprintf(`EVERY %s %s`, event.IntervalValue, event.IntervalField)
	if event.ExecuteAt != "" {
		schedule = fmt.Sprintf(`AT %s`, event.ExecuteAt)
	}

	if event.Starts != "" {
		starts = fmt.Sprintf("STARTS '%s'", event.Starts)
	}

	if event.Ends != "" {
		ends = fmt.Sprintf("ENDS '%s'", event.Ends)
	}

	return fmt.Sprintf(eventSqlStruct, event.EventName, schedule, event.OnCompletion, starts, ends, event.EventComment, event.EventDefinition, event.EventName, event.Status)
}
