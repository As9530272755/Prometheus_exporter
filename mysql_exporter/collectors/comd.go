package collectors

import (
	"database/sql"
	"mysql_exporter/logs"

	"github.com/prometheus/client_golang/prometheus"
)

type CommandCollector struct {
	mysqlCollector
	desc *prometheus.Desc
}

// 定义采集监控指标
func (c *CommandCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- c.desc
}

// 采集监控指标
func (c *CommandCollector) Collect(metric chan<- prometheus.Metric) {
	// 采集 insert,delete,update,select 语句
	sqlCmd := []string{"com_insert", "com_delete", "com_update", "com_select"}
	for _, sql := range sqlCmd {

		// 添加 debug 信息
		logs.WithFields("sql_CMD", sql)
		metric <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, c.status(sql), sql)
	}
}

// 注册
func CommandRegister(db *sql.DB) {
	prometheus.MustRegister(&CommandCollector{
		mysqlCollector: mysqlCollector{db},
		desc: prometheus.NewDesc(
			"mysql_sql",         // 指标名称
			"sql 语句指标获取",        // 指标帮助
			[]string{"command"}, // 指标可变标签
			nil,                 // 由于是可变标签没有固定标签
		),
	})
}
