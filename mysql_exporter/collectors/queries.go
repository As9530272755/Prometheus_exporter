package collectors

import (
	"database/sql"
	"mysql_exporter/logs"

	"github.com/prometheus/client_golang/prometheus"
)

type SlowQueriesCollector struct {
	mysqlCollector // 匿名结构体组合 mysqlCollector
	slowDesc       *prometheus.Desc
}

// 通过调用结构体属性来查询监控指标
func (c *SlowQueriesCollector) Describe(desc chan<- *prometheus.Desc) {
	// 将指标描述信息写入 descs 只写管道中
	desc <- c.slowDesc
}

// 对慢查询指标采集数据
func (c *SlowQueriesCollector) Collect(metric chan<- prometheus.Metric) {
	slow := c.status("slow_queries")
	logs.WithFields("slow_queries", slow)
	metric <- prometheus.MustNewConstMetric(
		c.slowDesc,              // 采集监控指标
		prometheus.CounterValue, // 指标类型为 counter 对递增、递减的指标进行采集
		slow,                    // 采集慢查询指标的值
	)
}

// 注册慢查询监控项
func SlowQueriesRegister(db *sql.DB) {
	prometheus.MustRegister(&SlowQueriesCollector{
		// 赋值 db 到 mysqlCollector 结构体组合中，拿到当前数据库信息
		mysqlCollector: mysqlCollector{db},
		slowDesc: prometheus.NewDesc(
			"Mysql_globel_status_slow_queries", // 指标名称
			"MySQL 慢查询指标",                      // 指标帮助
			nil,                                // 没有可变标签
			nil,                                // 没有固定标签
		),
	})
}

// 定义 Qps 语句的监控指标
type QpsCollector struct {
	mysqlCollector // 匿名结构体组合 mysqlCollector
	desc           *prometheus.Desc
}

// 查询监控指标
func (c *QpsCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.desc
}

// 采集当前执行 sql 语句监控指标
func (c *QpsCollector) Collect(metric chan<- prometheus.Metric) {
	queries := c.status("Queries")
	logs.WithFields("queries", queries)

	metric <- prometheus.MustNewConstMetric(
		c.desc,                  // 采集监控指标
		prometheus.CounterValue, // 采集数据类型
		queries,                 // 采集当前执行 sql 语句指标的值
	)
}

// 注册并给 QpsCollector 结构体赋值
func QpsRegister(db *sql.DB) {
	prometheus.MustRegister(&QpsCollector{
		mysqlCollector: mysqlCollector{db},
		desc: prometheus.NewDesc(
			"MySQL_QPS",        // 采集名称
			"MySQL 当前执行语句总数采集", // 采集帮助信息
			nil,                // 可变标签
			nil,                // 固定标签
		),
	})
}
