package main

import (
	"fmt"
	"mysql_exporter/collectors"
	"mysql_exporter/handler"
	"mysql_exporter/linkmysql"
	"mysql_exporter/logs"
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
	logClose := logs.MysqlLog(options.Log)
	defer logClose()

	// 通过 options 解析配置文件得到 dsn 信息
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		options.MySql.UserName,
		options.MySql.Password,
		options.MySql.Host,
		options.MySql.Port,
		options.MySql.DB)
	fmt.Println(dsn)
	// 链接数据库
	db := linkmysql.LinkDB(dsn)

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

	// 3.暴露指标,并且通过 auth 先进行验证是否通过验证，通过这访问 Prometheus 的 metrics
	http.Handle("/metrics/", handler.Auth(promhttp.Handler(), user))
	http.ListenAndServe(options.Web.Addr, nil)
}
