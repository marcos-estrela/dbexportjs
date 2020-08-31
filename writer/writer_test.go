package writer

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
		{Type: "views", Name: "vw_users", Content: "CREATE OR REPLACE VIEW vw_users SELECT * FROM users"},
		{Type: "procedures", Name: "sp_users_ins", Content: "CREATE PROCEDURE sp_users_ins(IN p_name VARCHAR(255), IN p_email VARCHAR(255))"},
	}

	SaveDbObjects(dbObjects)

	for i := range dbObjects {
		dbObject := dbObjects[i]

		filePath := makeObjectPath(dbObject)
		assertFileExists(t, filePath)
	}
}

func TestMakeObjectPath(t *testing.T) {
	dbObj := DbObject{
		Type:    "triggers",
		Name:    "tg_users_ins_after",
		Content: "CREATE TRIGGER tg_users_ins_after",
	}

	filePath := makeObjectPath(dbObj)
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
