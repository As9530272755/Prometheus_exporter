package logs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 定义日志信息
func MysqlLog(fileName, level string, max_age, max_size, max_backups int, compress bool) {

	logger := lumberjack.Logger{
		Filename:   fileName,
		MaxAge:     max_age,
		MaxSize:    max_size,
		MaxBackups: max_backups,
		Compress:   compress,
	}

	defer logger.Close()

	// 解析 level 传入的日志级别
	logLevel, err := logrus.ParseLevel(level)
	fmt.Println(logLevel, err)

	logrus.SetOutput(&logger)

	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)
}

// 定义日志记录监控项
func WithFields(metrics, sample interface{}) {
	logrus.WithFields(
		logrus.Fields{
			"metric": metrics,
			"sample": sample,
		}).Debug("MySQL Exporter")
}
