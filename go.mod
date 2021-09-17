module github.com/open-resource-management/metricsclient

go 1.14

replace github.com/open-resource-management/metricsclient => ../metricsclient

require (
	github.com/google/cadvisor v0.40.0
	github.com/prometheus/client_golang v1.11.0
	github.com/shirou/gopsutil v3.21.8+incompatible
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.20.6
	k8s.io/klog v1.0.0
)
