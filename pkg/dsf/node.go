package dsf

import (
	"fmt"
	"time"

	"github.com/open-resource-management/metricsclient/pkg/stats"
	"github.com/open-resource-management/metricsclient/pkg/types"
)

// nodeCpuUsage return cpu usage for node
func GetNodeCpuUsage(nodeStats stats.NodeStats) types.MetricValues {
	result, err := nodeStats.GetResourceStats(types.CpuUsageMetrics)
	if err != nil {
		return types.MetricValues{}
	}
	res := result.(*stats.NodeCpu)
	values := types.MetricValues{
		{
			Value:     res.CpuTotal,
			Labels:    []string{"total"},
			Timestamp: res.Timestamp,
		},
	}
	for i, v := range res.CpuPerCore {
		values = append(values, types.MetricValue{
			Value:     v,
			Labels:    []string{fmt.Sprintf("cpu%02d", i)},
			Timestamp: res.Timestamp,
		})
	}
	return values
}

// nodeMemoryUsage return memory usage for node
func GetNodeMemoryUsage(nodeStats stats.NodeStats) types.MetricValues {
	result, err := nodeStats.GetResourceStats(types.MemoryUsageMetrics)
	if err != nil {
		return types.MetricValues{}
	}
	res := result.(*stats.NodeMemory)
	return types.MetricValues{
		{
			Value:     bytesToGi(res.UsageTotal),
			Timestamp: res.Timestamp,
			Labels:    []string{"total"},
		},
		{
			Value:     bytesToGi(res.UsageRss),
			Timestamp: res.Timestamp,
			Labels:    []string{"rss"},
		},
		{
			Value:     bytesToGi(res.UsageCache),
			Timestamp: res.Timestamp,
			Labels:    []string{"cache"},
		},
	}
}

// bytesToGi translates bytes to Gi, and also keep three digital.
func bytesToGi(mBytes float64) float64 {
	return float64(int64(float64(int64(mBytes/1024/1024+0.5))/1024*1000)) / 1000
}
