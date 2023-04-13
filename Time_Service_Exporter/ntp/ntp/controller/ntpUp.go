package controller

import (
	"cicc/ntp/logs"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
)

// Determine whether the ntp server is alive
func Service_Up(service string) {

	prometheus.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "up",
		Namespace:   "Time_Service",
		Help:        "Proecss UP info 0 表示程序存活 1 表示程序死亡",
		ConstLabels: prometheus.Labels{"Proecss_Name": service},
	}, func() float64 {
		// 通过 pgrep -f 获取到对应的监控进程
		cmd := exec.Command("pgrep", "-f", service)
		// 获取 PID
		output, err := cmd.Output()
		if err != nil {
			logs.WithFields(service, err)
		}
		// 基于获取到的 pid 判断是否为空，如果 pid 为空就表示该进程不存在不为空就是存活
		if string(output) != "" {
			return 0
		} else {
			return 1
		}
	}))
}
