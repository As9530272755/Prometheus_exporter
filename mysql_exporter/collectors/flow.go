package collectors

import (
	"database/sql"

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
	metric <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, c.status("Bytes_ received"), "发送流量")
	metric <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, c.status("Bytes_sent"), "接收流量")
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
