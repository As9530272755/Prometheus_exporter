package collectors

import (
	"database/sql"
	"mysql_exporter/logs"

	"github.com/prometheus/client_golang/prometheus"
)

type ConnectionCollector struct {
	mysqlCollector

	// 由于有两个指标所以这里我写到一起
	connectedDesc    *prometheus.Desc
	maxconnectedDesc *prometheus.Desc
}

// 每个收集器都必须实现descripe函数
func (c *ConnectionCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- c.connectedDesc
	desc <- c.maxconnectedDesc
}

// 采集监控指标
func (c *ConnectionCollector) Collect(metric chan<- prometheus.Metric) {
	Threads_connected := c.status("Threads_connected")

	// 调用监控项日志
	logs.WithFields("Threads_connected", Threads_connected)

	// 当前链接数指标
	metric <- prometheus.MustNewConstMetric(
		c.connectedDesc,         // 采集监控指标
		prometheus.CounterValue, // 采集值类型
		Threads_connected,       // 采集 “已连接的线程” 监控指标
	)

	max_connections := c.variable("max_connections")

	// 调用监控项日志
	logs.WithFields("max_connections", max_connections)
	// 最大链接数指标
	metric <- prometheus.MustNewConstMetric(
		c.maxconnectedDesc,
		prometheus.CounterValue,
		max_connections,
	)
}

// 注册
func ConnectionRegister(db *sql.DB) {
	prometheus.MustRegister(&ConnectionCollector{
		mysqlCollector: mysqlCollector{db},
		connectedDesc: prometheus.NewDesc(
			"mysql_connection", // 指标名称
			"mysql 当前链接采集",     // 指标帮助信息
			nil,                // 没有可变标签
			nil,                // 没有固定标签
		),
		maxconnectedDesc: prometheus.NewDesc(
			"mysql_MAX_connection",
			"mysql 最大链接数",
			nil,
			nil,
		),
	})
}
