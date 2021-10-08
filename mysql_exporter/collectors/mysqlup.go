package collectors

import (
	"database/sql"
	"net"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func MysqlUp(db *sql.DB, host, port string) {

	// 1.检查 mysql 是否存活，是通过 gauage 指标类型定义,
	// 2.prometheus.MustRegister 注册 mysql_up 指标
	prometheus.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "mysql",
		Name:        "mysql_up",
		Help:        "Mysql 存活检查",
		ConstLabels: prometheus.Labels{"addr": net.JoinHostPort(host, port)},
	}, func() float64 {
		// 采集数据通过 db.ping 来对数据库进行存活检查,如果 err 为 nil 表示连接成功返回 1 否则返回 0
		if err := db.Ping(); err == nil {
			return 1
		} else {
			// 添加 mysql 指标 err 信息
			logrus.WithFields(logrus.Fields{
				"metric": "mysql_up_info",
			}).Error(err)
			return 0
		}
	}))
}
