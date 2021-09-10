module github.com/open-resource-management/metricsclient

go 1.14

replace github.com/open-resource-management/metricsclient => ../metricsclient

require (
	github.com/prometheus/client_golang v1.11.0
	k8s.io/klog v1.0.0
)
