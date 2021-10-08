package collectors

import (
	"database/sql"
	"mysql_exporter/logs"

	"github.com/prometheus/client_golang/prometheus"
)

type Flow struct {
	mysqlCollector
	desc *prometheus.Desc
}

func (c *Flow) Describe(desc chan<- *prometheus.Desc) {
	desc <- c.desc
}

func (c *Flow) Collect(metric chan<- prometheus.Metric) {
	// 发送流量
	received := c.status("Bytes_ received")

	// 记录日志
	logs.WithFields("received", received)
	metric <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, received, "发送流量")

	// 接收流量记录日志
	sent := c.status("Bytes_sent")
	logs.WithFields("Bytes_sent", sent)
	metric <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, sent, "接收流量")
}

func FlowRegister(db *sql.DB) {
	prometheus.MustRegister(&Flow{
		mysqlCollector: mysqlCollector{db},
		desc: prometheus.NewDesc(
			"mysql_Folw",
			"mysql 流量监控",
			[]string{"flow_to"},
			nil,
		),
	})
}
