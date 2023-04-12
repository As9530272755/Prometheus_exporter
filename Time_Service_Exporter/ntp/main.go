package main

import (
	"cicc/ntp/conf"
	"cicc/ntp/controller"
	"cicc/ntp/logs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func main() {
	config := "./etc/config.yaml"
	options, err := conf.ParseConfig(config)
	if err != nil {
		logrus.Fatal(err)
	}

	logs.Ex_logs(&options.Log)

	// 监控数量
	filesMonitored := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "file_monitored",
		Help: "Number of files being monitored",
	})

	// 触发告警数量
	alertsTriggered := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "file_monitored",
		Namespace:   "Time_Service",
		Help:        "file_monitored:文件监控触发告警正常情况下是 0 如果不等于 0 那就触发告警",
		ConstLabels: prometheus.Labels{"FileName": options.Service.ServiceConfigPath},
	})

	// 注册
	prometheus.MustRegister(filesMonitored)
	prometheus.MustRegister(alertsTriggered)

	// New一个 FileWatcher 结构体
	fw, err := controller.NewFileWatcher(options.Service.ServiceConfigPath, filesMonitored, alertsTriggered)
	if err != nil {
		log.Fatal(err)
		logs.WithFields("NewFileWatcher_err:", err)
	}
	fw.Start()
	defer fw.Stop()

	// 监控 service 是否存活
	controller.Service_Up(options.Service.ServiceName)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(options.Web.Addr, nil)

}
