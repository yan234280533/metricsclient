package prom

import (
	"testing"
	"time"
)

func TestNewPrometheusClient(t *testing.T) {
	var address = "127.0.0.1"
	var timeout = 10 * time.Second
	var keepAlive = 10 * time.Second
	var queryConcurrency = 10

	_, err := NewPrometheusClient(address, timeout, keepAlive, queryConcurrency, false, false, &ClientAuth{})
	if err != nil {
		t.Fatalf("NewPrometheusClient failed %s", err.Error())
	}

	_, err = NewPrometheusClient(address, timeout, keepAlive, queryConcurrency, true, false, &ClientAuth{})
	if err != nil {
		t.Fatalf("NewPrometheusClient failed %s", err.Error())
	}

	t.Logf("TestNewPrometheusClient succeed")
}
