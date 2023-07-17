package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/my_project")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
