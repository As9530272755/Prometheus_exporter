package main

import (
	"fmt"
	"kube_expoter/controller"
	"kube_expoter/ex_config"
	"kube_expoter/msgErr"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//mkconfigDir.MakeDir()

	options, err := ex_config.ParseConfig("./etc/exporter.yaml")
	msgErr.ErrInfo(err)

	addr := fmt.Sprintf(":" + options.Web.ListenPort)

	controller.PrometheusCotroller(options)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(addr, nil)
}
