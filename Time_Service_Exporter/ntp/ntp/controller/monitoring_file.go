package controller

import (
	"cicc/ntp/logs"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
)

// 文件监控结构体
type FileWatcher struct {
	watcher         *fsnotify.Watcher
	filesMonitored  prometheus.Gauge
	alertsTriggered prometheus.Counter
}

// New 结构体
func NewFileWatcher(path string, filesMonitored prometheus.Gauge, alertsTriggered prometheus.Counter) (*FileWatcher, error) {
	// new 一个 NewWatcher 对象
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// 指定监控文件 path
	err = watcher.Add(path)
	if err != nil {
		return nil, err
	}
	return &FileWatcher{watcher, filesMonitored, alertsTriggered}, nil
}

// 启动协程实现对文件系统的实时监控
func (fw *FileWatcher) Start() {
	go func() {
		for {
			select {
			// 对文件的 events 监控
			case event, ok := <-fw.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					logs.WithFields("写入文件告警\n", event.Name)
					fw.alertsTriggered.Inc()
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					logs.WithFields("删除文件告警\n", event.Name)
					fw.alertsTriggered.Inc()
				}

				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					logs.WithFields("修改文件权限告警\n", event.Name)
					fw.alertsTriggered.Inc()
				}

				if event.Op&fsnotify.Rename == fsnotify.Rename {
					logs.WithFields("重命名文件告警\n", event.Name)
					fw.alertsTriggered.Inc()
				}
			case err, ok := <-fw.watcher.Errors:
				if !ok {
					return
				}
				logs.WithFields("error:", err)
			}
		}
	}()

	// 设置默认值
	fw.filesMonitored.Set(1)
}

func (fw *FileWatcher) Stop() {
	fw.watcher.Close()
	fw.filesMonitored.Set(0)
}
