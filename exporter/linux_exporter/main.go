package main

import (
	"linux_exporter/linuxconllector"
	"linux_exporter/sshlink"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	addr := ":9999"

	client := sshlink.Link()
	linuxconllector.LinuxRegister(client)

	defer client.Close()

	http.Handle("/metrics/", promhttp.Handler())
	http.ListenAndServe(addr, nil)
}
