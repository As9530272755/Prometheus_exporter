package logs

import (
	"fmt"
	"mysql_exporter/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 定义日志信息
func MysqlLog(options config.Log) func() {

	logger := lumberjack.Logger{
		Filename:   options.FileName,
		MaxAge:     options.Max_age,
		MaxSize:    options.Max_size,
		MaxBackups: options.Max_backups,
		Compress:   options.Compress,
	}

	// 解析 level 传入的日志级别
	logLevel, err := logrus.ParseLevel(options.Level)
	fmt.Println(logLevel, err)

	logrus.SetOutput(&logger)

	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)

	// 返回一个 close
	return func() {
		logger.Close()
	}
}

// 定义日志记录监控项
func WithFields(metrics, sample interface{}) {
	logrus.WithFields(
		logrus.Fields{
			"metric": metrics,
			"sample": sample,
		}).Debug("MySQL Exporter")
}
