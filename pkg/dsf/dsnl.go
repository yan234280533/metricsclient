package dsf

import (
	"fmt"
	"github.com/open-resource-management/metricsclient/pkg/stats"
	"github.com/open-resource-management/metricsclient/pkg/types"
	"k8s.io/klog"

	"k8s.io/client-go/tools/cache"
)

type DataNodeLocalSource struct {
	rsi stats.ResourceStatsInterface
}

func NewDataNodeLocalSource(config *DataSourceNodeLocalConfig, podInformer cache.SharedIndexInformer) (*DataNodeLocalSource, error) {
	klog.Infof("NewDataNodeLocalSource")

	rsi, err := stats.NewResourceStats(config.MetricsTTL, config.Node, config.Container, podInformer)
	if err != nil {
		return nil, err
	}

	return &DataNodeLocalSource{rsi: rsi}, nil
}

func (nl *DataNodeLocalSource) GetCpuUsageSample(name DataSourceObjectName) (DataSample, error) {
	if IsNodeDataSourceObject(name) {
		metricValues := GetNodeCpuUsage(nl.rsi.GetNodeStats())
	} else if IsPodDataSourceObject(name) {

	} else if IsContainerDataSourceObject(name) {

	}
	return DataSample{}, fmt.Errorf("the type of metric is only support (node, pod, container)")
}

func (nl *DataNodeLocalSource) GetMemoryUsageSample(name DataSourceObjectName) (DataSample, error) {
	return DataSample{}, fmt.Errorf("the type of metric is only support (node, pod, container)")
}
