package prometheus

import (
	"fmt"
	"github.com/open-resource-management/metricsclient/pkg/util"
	"testing"
	"time"
)

func TestNewNamedContext(t *testing.T) {

	var address = "http://159.75.100.72:12345"
	var timeout = 10 * time.Second
	var keepAlive = 10 * time.Second
	var queryConcurrency = 10

	client, err := NewPrometheusClient(address, timeout, keepAlive, queryConcurrency, false, false, &ClientAuth{})
	if err != nil {
		t.Fatalf("NewPrometheusClient failed %s", err.Error())
	}

	_ = NewNamedContext(client, ClusterContextName)

	t.Logf("TestNewNamedContext succeed")
}

func TestContext_Query(t *testing.T) {
	var address = "http://159.75.100.72:12345"
	var timeout = 10 * time.Second
	var keepAlive = 10 * time.Second
	var queryConcurrency = 10

	client, err := NewPrometheusClient(address, timeout, keepAlive, queryConcurrency, false, false, &ClientAuth{})
	if err != nil {
		t.Fatalf("NewPrometheusClient failed %s", err.Error())
	}

	ctx := NewNamedContext(client, ClusterContextName)

	start := time.Minute * 2
	offset := time.Minute * 2
	minsPerResolution := 1

	durationStr := fmt.Sprintf("%dm", int64(start.Minutes()))
	offsetStr := fmt.Sprintf(" offset %dm", int64(offset.Minutes()))
	queryNodeCPUModeTotal := fmt.Sprintf(`sum(rate(node_cpu_seconds_total[%s:%dm]%s)) by (kubernetes_node, mode)`, durationStr, minsPerResolution, offsetStr)

	t.Logf("TestContext_Query queryNodeCPUModeTotal %s", queryNodeCPUModeTotal)

	resultsChan := ctx.Query(queryNodeCPUModeTotal)

	results, err := resultsChan.Await()
	if err != nil {
		t.Fatalf("TestContext_Query Query failed, err %s", err.Error())
		return
	}

	t.Logf("TestContext_Query request %v", results)

	for key := range results {
		t.Logf("Metrics: %s ", results[key].Metric)
		t.Logf("Values: %s", util.GetStringVerctors(results[key].Values))
	}

	t.Logf("TestContext_Query succeed")
}

func TestContext_QueryRange(t *testing.T) {
	var address = "http://159.75.100.72:12345"
	var timeout = 10 * time.Second
	var keepAlive = 10 * time.Second
	var queryConcurrency = 10

	var end = time.Now()
	var start = end.Add(-time.Minute * 10)

	client, err := NewPrometheusClient(address, timeout, keepAlive, queryConcurrency, false, false, &ClientAuth{})
	if err != nil {
		t.Fatalf("NewPrometheusClient failed %s", err.Error())
	}

	ctx := NewNamedContext(client, ClusterContextName)

	duration := time.Minute * 5
	offset := time.Minute * 5
	minsPerResolution := 1

	durationStr := fmt.Sprintf("%dm", int64(duration.Minutes()))
	offsetStr := fmt.Sprintf(" offset %dm", int64(offset.Minutes()))
	queryNodeCPUModeTotal := fmt.Sprintf(`sum(rate(node_cpu_seconds_total[%s:%dm]%s)) by (kubernetes_node, mode)`, durationStr, minsPerResolution, offsetStr)

	t.Logf("TestContext_QueryRange queryNodeCPUModeTotal %s", queryNodeCPUModeTotal)

	resultsChan := ctx.QueryRange(queryNodeCPUModeTotal, start, end, time.Minute)

	results, err := resultsChan.Await()
	if err != nil {
		t.Fatalf("TestContext_QueryRange Query failed, err %s", err.Error())
		return
	}

	t.Logf("TestContext_QueryRange request %v", results)

	for key := range results {
		t.Logf("Metrics: %s ", results[key].Metric)
		t.Logf("Values: %s", util.GetStringVerctors(results[key].Values))
	}

	t.Logf("TestContext_QueryRange succeed")
}

func TestContext_QueryRangeMemory(t *testing.T) {
	var address = "http://159.75.100.72:12345"
	var timeout = 10 * time.Second
	var keepAlive = 10 * time.Second
	var queryConcurrency = 10

	var end = time.Now()
	var start = end.Add(-time.Minute * 10)

	client, err := NewPrometheusClient(address, timeout, keepAlive, queryConcurrency, false, false, &ClientAuth{})
	if err != nil {
		t.Fatalf("NewPrometheusClient failed %s", err.Error())
	}

	ctx := NewNamedContext(client, ClusterContextName)

	queryNodeMemoryTotal := fmt.Sprintf(`sum(node_memory_MemTotal_bytes) by (instance)`)

	t.Logf("TestContext_QueryRangeMemory queryNodeMemoryTotal %s", queryNodeMemoryTotal)

	resultsChan := ctx.QueryRange(queryNodeMemoryTotal, start, end, time.Minute)

	results, err := resultsChan.Await()
	if err != nil {
		t.Fatalf("TestContext_QueryRangeMemory Query failed, err %s", err.Error())
		return
	}

	t.Logf("TestContext_QueryRangeMemory request %v", results)

	for key := range results {
		t.Logf("Metrics: %s ", results[key].Metric)
		t.Logf("Values: %s", util.GetStringVerctors(results[key].Values))
	}

	t.Logf("TestContext_QueryRangeMemory succeed")
}