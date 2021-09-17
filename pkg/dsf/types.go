package dsf

import (
	"github.com/open-resource-management/metricsclient/pkg/util"
	"time"
)

type DataSourceType string

const (
	DataSourcePromType       DataSourceType = "prom"
	DataSourceNodeLocaleType DataSourceType = "node-local"
)

// Sample is a single timestamped value of the metric.
type DataSample struct {
	Value     float64
	Timestamp time.Time
}

// TimeSeries represents a metric with given labels, with its values possibly changing in time.
type DataTimeSeries struct {
	Labels  map[string]string
	Samples []DataSample
}

func Vector2Sample(v util.Vector) DataSample {
	return DataSample{Timestamp: time.Unix(int64(v.Timestamp), 0), Value: v.Value}
}
