package main

import (
	"mysql_exporter/handler"
	"mysql_exporter/linkmysql"
	"mysql_exporter/logs"
	"mysql_exporter/monitor"
	"net/http"

	"mysql_exporter/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// MySQL 存活状态指标采集
func main() {

	// 调用配置文件
	options, err := config.ParseConfig("./etc/mysql_exporter.yaml")
	if err != nil {
		logrus.Fatal(err)
	}

	// 通过 options 解析配置文件得到 user 授权信息
	user := map[string]string{options.Web.Auth.UserName: options.Web.Auth.Password}

	// 通过 options 解析配置文件得到日志信息,并接受 close 返回值因为我们需要在 main 程序中延迟关闭
	logClose := logs.InitMysqlLog(options.Log)
	defer logClose()

	// 通过 options 解析配置文件得到 dsn 信息
	// 链接数据库
	db, err := linkmysql.InitLinkDB(options.MySql)
	if err != nil {
		logrus.Fatal(err)
		return
	}

	monitor.MonitorController(db, *options)

	// 3.暴露指标,并且通过 auth 先进行验证是否通过验证，通过这访问 Prometheus 的 metrics
	http.Handle("/metrics/", handler.Auth(promhttp.Handler(), user))
	http.ListenAndServe(options.Web.Addr, nil)
}
