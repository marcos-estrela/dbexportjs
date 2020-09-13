package dbexport

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func conn() *sql.DB {
	driver := os.Getenv("DB_DRIVER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_DATABASE")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	db, err := sql.Open(driver, connectionString)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	return db
}

func ExecuteQuery(query string, args ...interface{}) *sql.Rows {
	db := conn()
	results, err := db.Query(query, args...)

	defer db.Close()

	if err != nil {
		panic(err.Error())
	}

	return results
}

func ExecuteQueryRow(query string, args ...interface{}) *sql.Row {
	db := conn()
	result := db.QueryRow(query, args...)
	defer db.Close()
	return result
}
