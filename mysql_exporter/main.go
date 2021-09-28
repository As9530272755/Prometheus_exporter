package main

import (
	"mysqlmonitor/collectors"
	"mysqlmonitor/linkmysql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MySQL 存活状态指标采集
func main() {
	webAddr := ":9090"

	// 链接数据库
	db := linkmysql.LinkDB()

	// 调用 MysqlUp 监控指标
	collectors.MysqlUp(db)

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

	// 3.暴露指标
	http.Handle("/metrics/", promhttp.Handler())
	http.ListenAndServe(webAddr, nil)
}
