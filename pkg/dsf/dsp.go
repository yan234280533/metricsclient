package dsf

import (
	"fmt"
	"github.com/open-resource-management/metricsclient/pkg/prom"
	"github.com/open-resource-management/metricsclient/pkg/util"
	"k8s.io/klog"
	"time"
)

const (
	defaultMinPerResolution = time.Second * 60
	defaultDuration         = time.Second * 60
)

type DataPromSource struct {
	ctx              *prom.Context
	minPerResolution time.Duration
	duration         time.Duration
}

func NewDataPromSource(config *DataSourcePromConfig) (*DataPromSource, error) {
	klog.Infof("NewDataPromSource")

	client, err := prom.NewPrometheusClient(config.address, config.timeout, config.keepAlive,
		config.queryConcurrency, config.insecureSkipVerify, config.bRateLimit, config.auth)
	if err != nil {
		return nil, err
	}

	ctx := prom.NewNamedContext(client, prom.ClusterContextName)

	return &DataPromSource{ctx: ctx, minPerResolution: defaultMinPerResolution, duration: defaultDuration}, nil
}

func (c *DataPromSource) GetCpuUsageSample(name DataSourceObjectName) (DataSample, error) {
	if IsNodeDataSourceObject(name) {
		durationStr := fmt.Sprintf("%ds", int64(c.duration.Seconds()))
		minPerResolutionStr := fmt.Sprintf("%ds", int64(c.minPerResolution.Seconds()))
		queryNodeCPUUsage := fmt.Sprintf(`1-avg(rate(node_cpu_seconds_total{mode="idle",instance="%s"}[%s:%s])) by (instance)`, name.NodeName, durationStr, minPerResolutionStr)

		results, err := c.ctx.QuerySync(queryNodeCPUUsage)
		if err != nil {
			klog.Errorf("GetCpuUsageSample Query failed, err %s", err.Error())
			return DataSample{}, err
		}

		v, err := GetVectorFromResults(results)
		if err != nil {
			klog.Errorf("GetCpuUsageSample get vector failed, err %s", err.Error())
			return DataSample{}, err
		}

		return Vector2Sample(v), nil

	} else if IsPodDataSourceObject(name) {
		durationStr := fmt.Sprintf("%ds", int64(c.duration.Seconds()))
		minPerResolutionStr := fmt.Sprintf("%ds", int64(c.minPerResolution.Seconds()))

		queryPodCPUUsage := fmt.Sprintf(`rate(container_cpu_usage_seconds_total{pod="%s",container="",namespace="%s"}[%s:%s])`, name.PodName, name.Namespace, durationStr, minPerResolutionStr)

		results, err := c.ctx.QuerySync(queryPodCPUUsage)
		if err != nil {
			klog.Errorf("GetCpuUsageSample Query failed, err %s", err.Error())
			return DataSample{}, err
		}

		v, err := GetVectorFromResults(results)
		if err != nil {
			klog.Errorf("GetCpuUsageSample get vector failed, err %s", err.Error())
			return DataSample{}, err
		}

		return Vector2Sample(v), nil

	} else if IsContainerDataSourceObject(name) {
		durationStr := fmt.Sprintf("%ds", int64(c.duration.Seconds()))
		minPerResolutionStr := fmt.Sprintf("%ds", int64(c.minPerResolution.Seconds()))

		queryContainerCPUUsage := fmt.Sprintf(`rate(container_cpu_usage_seconds_total{pod="%s",container="%s",namespace="%s"}[%s:%s])`, name.PodName, name.ContainerName, name.Namespace, durationStr, minPerResolutionStr)

		results, err := c.ctx.QuerySync(queryContainerCPUUsage)
		if err != nil {
			klog.Errorf("GetCpuUsageSample Query failed, err %s", err.Error())
			return DataSample{}, err
		}

		v, err := GetVectorFromResults(results)
		if err != nil {
			klog.Errorf("GetCpuUsageSample get vector failed, err %s", err.Error())
			return DataSample{}, err
		}

		return Vector2Sample(v), nil

	}
	return DataSample{}, fmt.Errorf("the type of metric is only support (node, pod, container)")
}

func (c *DataPromSource) GetMemoryUsageSample(name DataSourceObjectName) (DataSample, error) {
	return DataSample{}, fmt.Errorf("the type of metric is only support (node, pod, container)")
}

func GetVectorFromResults(results []*prom.QueryResult) (util.Vector, error) {
	if len(results) == 0 {
		return util.Vector{}, fmt.Errorf("QuerySync empty")
	}

	if len(results) > 1 {
		klog.Warningf("QuerySync results length is %d, not single", len(results))
		results = results[0:1]
	}

	if len(results[0].Values) == 0 {
		return util.Vector{}, fmt.Errorf("QuerySync values empty")
	}

	if results[0].Values[0] == nil {
		return util.Vector{}, fmt.Errorf("QuerySync values zero is nil")
	}

	return *(results[0].Values[0]), nil
}
