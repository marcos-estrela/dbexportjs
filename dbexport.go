package dbexport

import  (
	"io/ioutil"
	"log"
)

func WriteSqlToFile(filePath, sql string) {
	content := []byte(sql)
	err := ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func GetSqlForTable(tableName string) string {
	return "CREATE TABLE " + tableName
}
