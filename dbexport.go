package dbexport

func GetSqlForTable(tableName string) string {
	return "CREATE TABLE " + tableName
}
