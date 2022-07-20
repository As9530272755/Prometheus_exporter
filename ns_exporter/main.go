package main

import (
	"fmt"
	"kube_expoter/controller"
	"kube_expoter/ex_config"
	"kube_expoter/msgErr"
	"net/http"

	_ "github.com/CodyGuo/godaemon"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//mkconfigDir.MakeDir()
	mux := http.NewServeMux()

	options, err := ex_config.ParseConfig("./etc/exporter.yaml")
	msgErr.ErrInfo(err)

	addr := fmt.Sprintf(":" + options.Web.ListenPort)

	controller.PrometheusCotroller(options)

	mux.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(addr, mux)
}
