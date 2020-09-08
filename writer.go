package dbexport

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func WriteSqlToFile(filePath, sql string) bool {
	content := replaceNewLine(sql)
	byteContent := []byte(content)
	path := filepath.Dir(filePath)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	err := ioutil.WriteFile(filePath, byteContent, os.ModePerm)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func replaceNewLine(content string) string {
	strings.ReplaceAll(content, "^M", "")
	lines := strings.Split(content, "\r\n")
	return strings.Join(lines, "\n")
}

func SaveDbObjects(dbObjects []DbObject) []string {
	var savedFiles []string

	for i := range dbObjects {
		dbObject := dbObjects[i]
		filePath := makeDbObjectPath(dbObject)
		if WriteSqlToFile(filePath, dbObject.Content) {
			savedFiles = append(savedFiles, dbObject.Name)
		}
	}

	return savedFiles
}

func makeDbObjectPath(dbObject DbObject) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	return dir + "/export/" + dbObject.Type + "/" + dbObject.Name + ".sql"
}
