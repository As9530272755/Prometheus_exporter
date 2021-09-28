package collectors

import "database/sql"

type mysqlCollector struct {
	db *sql.DB
}

// 通过 sql 语句来获取 status 监控项
func (c *mysqlCollector) status(name string) float64 {
	var (
		newname string
		total   float64
	)

	sql := "show global status where variable_name=?"
	c.db.QueryRow(sql, name).Scan(&newname, &total)
	return total
}

// 通过 sql 语句来获取 variable 监控项
func (c *mysqlCollector) variable(name string) float64 {
	var (
		newname string
		total   float64
	)

	sql := "show global variables where variable_name=?"
	c.db.QueryRow(sql, name).Scan(&newname, &total)
	return total

}
