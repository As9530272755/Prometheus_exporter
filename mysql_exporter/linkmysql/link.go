package linkmysql

import (
	"database/sql"
	"fmt"
	"mysql_exporter/config"

	"github.com/sirupsen/logrus"
)

func LinkDB(options config.MySql) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		options.UserName,
		options.Password,
		options.Host,
		options.Port,
		options.DB)
	// 链接 mysql
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}
	return db, nil
}
