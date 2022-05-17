package linuxconllector

import (
	"linux_exporter/cmd"
	"linux_exporter/sshlink"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/ssh"
)

type LinuxConllector struct {
	cmd.LinuxCmd
	memdesc *prometheus.Desc
}

// 每个收集器都必须实现descripe函数
func (c *LinuxConllector) Describe(desc chan<- *prometheus.Desc) {
	desc <- c.memdesc
}

// 监控指标
func (c *LinuxConllector) Collect(metric chan<- prometheus.Metric) {
	command := `cat /proc/meminfo  | grep MemAvailable | awk '{print $2}'`
	metric <- prometheus.MustNewConstMetric(c.memdesc, prometheus.GaugeValue, c.Mem(command), sshlink.Host)
}

// 注册
func LinuxRegister(client *ssh.Client) {
	prometheus.MustRegister(&LinuxConllector{
		LinuxCmd: cmd.LinuxCmd{Client: client},
		memdesc: prometheus.NewDesc(
			"Linux_Memory",
			"Linux 远程主机获取内存,单位为 GB",
			[]string{"Remote_HOST"},
			nil,
		),
	})
}
