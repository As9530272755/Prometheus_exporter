package linkmysql

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func LinkDB(dsn string) *sql.DB {
	// 链接 mysql
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatal(err)
	}
	return db
}
