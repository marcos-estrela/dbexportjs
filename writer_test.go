package dbexport

import (
	"os"
	"strings"
	"testing"
)

func TestWriteSqlToFile(t *testing.T) {
	var sql = "CREATE OR REPLACE VIEW SELECT id, name, email FROM users;"
	var fileName = "vw_users.sql"
	var filePath = "/tmp/" + fileName

	WriteSqlToFile(filePath, sql)
	assertFileExists(t, filePath)
}

func TestSaveDbObjects(t *testing.T) {
	var dbObjects = []DbObject{
		{Type: "tables", Name: "users", Content: "CREATE TABLE users(id INT(11))"},
		{Type: "views", Name: "vw_users", Content: "CREATE OR REPLACE VIEW vw_users SELECT * FROM users"},
		{Type: "procedures", Name: "sp_users_ins", Content: "CREATE PROCEDURE sp_users_ins(IN p_name VARCHAR(255), IN p_email VARCHAR(255))"},
		{Type: "functions", Name: "fn_users_exists", Content: "CREATE FUNCTION fn_users_exists(p_userename VARCHAR(255))"},
		{Type: "events", Name: "ev_something", Content: "CREATE EVENT ev_count_users DO"},
		{Type: "triggers", Name: "tg_users_ins_after", Content: "CREATE TRIGGER tg_users_ins_after AFTER INSERT"},
	}

	SaveDbObjects(dbObjects)

	for i := range dbObjects {
		dbObject := dbObjects[i]

		filePath := makeDbObjectPath(dbObject)
		assertFileExists(t, filePath)
	}
}

func TestMakeDbObjectPath(t *testing.T) {
	dbObj := DbObject{
		Type:    "triggers",
		Name:    "tg_users_ins_after",
		Content: "CREATE TRIGGER tg_users_ins_after",
	}

	filePath := makeDbObjectPath(dbObj)
	expectedFileName := "export/triggers/tg_users_ins_after.sql"

	if !strings.Contains(filePath, expectedFileName) {
		t.Errorf("file name inocrrect. Found %s expected %s", filePath, expectedFileName)
	}
}

func assertFileExists(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)

	if err != nil {
		t.Errorf("file %s not exists", filePath)
	}

	e := os.Remove(filePath)
	if e != nil {
		t.Error(e)
	}
}
