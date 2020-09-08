package dbexport

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conn() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:33006)/george")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	return db
}

func ExecuteQuery(query string, args ...interface{}) *sql.Rows {
	db := Conn()
	results, err := db.Query(query, args...)

	defer db.Close()

	if err != nil {
		panic(err.Error())
	}

	return results
}

func ExecuteQueryRow(query string, args ...interface{}) *sql.Row {
	db := Conn()
	result := db.QueryRow(query, args...)
	return result
}
