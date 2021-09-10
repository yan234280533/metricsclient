package main

import (
	"fmt"
	"time"

	"github.com/open-resource-management/metricsclient/pkg/prometheus"
)

func main() {
	var address = "127.0.0.1"
	var timeout = 10 * time.Second
	var keepAlive = 10 * time.Second
	var queryConcurrency = 10

	client, _ := prometheus.NewPrometheusClient(address, timeout, keepAlive, queryConcurrency, false, false, &prometheus.ClientAuth{})

	fmt.Sprintf("client: %+v", client)
	return
}
