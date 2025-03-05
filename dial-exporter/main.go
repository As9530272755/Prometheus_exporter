package main

import (
	"context"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	nfsPortUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "nfs_tcp_port_up",
			Help: "NFS port  (1=up, 0=down)",
		},
		[]string{"target"},
	)

	healthGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_port_up",
			Help: "Health status of the host (1=up, 0=down)",
		},
		[]string{"target"},
	)
	icmpPingUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "icmp_ping_up",
			Help: "ICMP ping status (1=up, 0=down)",
		},
		[]string{"target"},
	)
)

func init() {
	prometheus.MustRegister(nfsPortUp)
	prometheus.MustRegister(healthGauge)
	prometheus.MustRegister(icmpPingUp)

}

// exporter config
func Exporter_Config() string {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file, %s", err)
	}

	return viper.GetString("exporter.port")

}

func Web_Config() []string {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file, %s", err)
	}

	return viper.GetStringSlice("http.url")

}

func Tcp_Config() []string {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file, %s", err)
	}

	return viper.GetStringSlice("tcp.server")

}

func Icmp_Config() []string {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
	}
	return viper.GetStringSlice("icmp.targets")
}

// checkHealth 对目标地址进行探活检查
func checkHealth(target string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	for {
		resp, err := client.Get(target)
		isHealthy := false
		if err == nil && resp.StatusCode == 200 {
			isHealthy = true
		}
		healthGauge.WithLabelValues(target).Set(float64(btoi(isHealthy)))
		time.Sleep(10 * time.Second) // 每10秒检查一次
	}
}

func probeNFSPort(ctx context.Context, target string, wg *sync.WaitGroup) {
	defer wg.Done()
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := dialer.DialContext(ctx, "tcp", target)
			if err == nil {
				nfsPortUp.WithLabelValues(target).Set(1)
				conn.Close()
			} else {
				nfsPortUp.WithLabelValues(target).Set(0)
			}
			time.Sleep(15 * time.Second) // 探测间隔
		}
	}
}

func probeICMP(ctx context.Context, target string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			pinger, err := ping.NewPinger(target)
			if err != nil {
				icmpPingUp.WithLabelValues(target).Set(0)
				continue
			}
			pinger.Count = 3
			pinger.Timeout = 5 * time.Second
			pinger.SetPrivileged(true) // Linux需要特权模式
			err = pinger.Run()
			if err != nil {
				icmpPingUp.WithLabelValues(target).Set(0)
			} else {
				if pinger.Statistics().PacketsRecv > 0 {
					icmpPingUp.WithLabelValues(target).Set(1)
				} else {
					icmpPingUp.WithLabelValues(target).Set(0)
				}
			}
			time.Sleep(15 * time.Second)
		}
	}
}

// btoi converts a boolean to an integer (1 for true, 0 for false)
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	fmt.Println(Web_Config())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	httpTargets := Web_Config()
	for _, target := range httpTargets {
		wg.Add(1)
		go checkHealth(target, &wg)
	}

	icmpTargets := Icmp_Config()
	for _, target := range icmpTargets {
		wg.Add(1)
		go probeICMP(ctx, target, &wg)
	}

	TcpTargets := Tcp_Config()
	for _, target := range TcpTargets {
		wg.Add(1)
		go probeNFSPort(ctx, target, &wg)
	}

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Beginning to serve on port : http://127.0.0.1:9090/metrics")
	http.ListenAndServe(":"+Exporter_Config(), nil)
}
