package monitor

import (
	"database/sql"
	"mysql_exporter/collectors"
	"mysql_exporter/config"
)

// 监控控制器
func MonitorController(db *sql.DB, options config.Options) {
	// 调用 MysqlUp 监控指标
	collectors.MysqlUp(db, options.MySql.Host, options.MySql.Port)

	// 调用慢查询监控指标
	collectors.SlowQueriesRegister(db)

	// 调用 qps 语句监控指标
	collectors.QpsRegister(db)

	// 调用 sql 执行语句监控指标
	collectors.CommandRegister(db)

	// 调用连接数监控项
	collectors.ConnectionRegister(db)

	// 调用流量监测监控项
	collectors.FlowRegister(db)
}
