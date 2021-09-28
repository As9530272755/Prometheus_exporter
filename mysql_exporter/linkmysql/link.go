package linkmysql

import (
	"database/sql"
	"log"
)

func LinkDB() *sql.DB {
	// 链接 mysql
	dsn := "root:123456@tcp(10.0.0.3:3306)/mysql"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
