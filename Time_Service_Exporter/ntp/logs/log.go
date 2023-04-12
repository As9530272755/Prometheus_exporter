package logs

import (
	"cicc/ntp/conf"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// exporter 日志功能
func Ex_logs(log *conf.Log) {
	logr := lumberjack.Logger{
		Filename: log.FileName,
		MaxAge:   log.Max_age,
	}

	defer logr.Close()

	logrus.SetOutput(&logr)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

// 自定义日志
func WithFields(metrics, sample interface{}) {
	logrus.WithFields(
		logrus.Fields{
			"metrics": metrics,
			"sample":  sample,
		}).Error("Time_Service err")
}
