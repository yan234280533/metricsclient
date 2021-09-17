package dsf

import (
	"fmt"
	"time"

	"github.com/open-resource-management/metricsclient/pkg/prom"
	"github.com/open-resource-management/metricsclient/pkg/types"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

type DataSourceConfig struct {
	*DataSourcePromConfig
	*DataSourceNodeLocalConfig
}

type DataSourcePromConfig struct {
	address            string
	timeout            time.Duration
	keepAlive          time.Duration
	insecureSkipVerify bool
	auth               *prom.ClientAuth

	queryConcurrency int
	bRateLimit       bool
}

type DataSourceNodeLocalConfig struct {
	MetricsTTL time.Duration                `json:"metrics_ttl"`
	Node       types.MetricsNodeConfig      `json:"node"`
	Container  types.MetricsContainerConfig `json:"container"`
}

func DataSourceFactory(t DataSourceType, config DataSourceConfig, podInformer cache.SharedIndexInformer) (DataSource, error) {
	klog.Infof("DataSourceFactory t %s", t)

	switch t {
	case DataSourcePromType:
		{
			if config.DataSourcePromConfig == nil {
				return nil, fmt.Errorf("DataSourcePromConfig is nil")
			}
			return NewDataPromSource(config.DataSourcePromConfig)
		}
	case DataSourceNodeLocaleType:
		{
			if config.DataSourceNodeLocalConfig == nil {
				return nil, fmt.Errorf("DataSourceNodeLocalConfig is nil")
			}

			return NewDataNodeLocalSource(config.DataSourceNodeLocalConfig, podInformer)
		}

	}

	return nil, fmt.Errorf("DataSourceFactory not support type t: %s", t)
}
