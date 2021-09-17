package types

import (
	"time"
)

// MetricsNodeConfig is the configuration for node metrics collection
type MetricsNodeConfig struct {
	CollectInterval time.Duration `json:"collect_interval"`
}

// MetricsContainerConfig is the configuration for container metrics collection
type MetricsContainerConfig struct {
	Resources               []string      `json:"resources"`
	Cgroups                 []string      `json:"cgroups"`
	MaxHousekeepingInterval time.Duration `json:"max_housekeeping_interval"`
}

type MetricKind string

// metricValue describes a single metric value for a given set of label values
type MetricValue struct {
	Value     float64
	Labels    []string
	Timestamp time.Time
}

type MetricValues []MetricValue

const (
	CpuUsageMetrics    MetricKind = "cpu"
	MemoryUsageMetrics MetricKind = "memory"
	CpuLoadMetrics     MetricKind = "cpuLoad"
)
