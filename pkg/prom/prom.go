package prom

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/open-resource-management/metricsclient/pkg/util/atomic"
	"github.com/open-resource-management/metricsclient/pkg/util/httputil"
	"github.com/open-resource-management/metricsclient/pkg/util/queue"

	prometheusapi "github.com/prometheus/client_golang/api"
	"k8s.io/klog"
)

const (
	// PrometheusClientID is the identifier used when creating the client that
	// targets prometheus. This can be used to check a specific client instance
	// by calling prom.IsClientID(client, prom.PrometheusClientID)
	PrometheusClientID string = "Prometheus"
)

//--------------------------------------------------------------------------
//  QueryParamsDecorator
//--------------------------------------------------------------------------

// QueryParamsDecorator is used to decorate and return query parameters for
// outgoing requests
type QueryParamsDecorator = func(path string, values url.Values) url.Values

//--------------------------------------------------------------------------
//  PrometheusClient
//--------------------------------------------------------------------------

func NewPrometheusClient(address string, timeout, keepAlive time.Duration, queryConcurrency int, insecureSkipVerify bool,
	needRateLimit bool, auth *ClientAuth) (prometheusapi.Client, error) {

	tlsConfig := &tls.Config{InsecureSkipVerify: insecureSkipVerify}

	pc := prometheusapi.Config{
		Address: address,
		RoundTripper: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: keepAlive,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     tlsConfig,
		},
	}

	if needRateLimit {
		return newRateLimitedClient(PrometheusClientID, pc, queryConcurrency, auth, nil)
	} else {
		return newPrometheusClientImp(PrometheusClientID, pc, auth, nil)
	}
}

type PrometheusClient struct {
	id        string
	client    prometheusapi.Client
	auth      *ClientAuth
	decorator QueryParamsDecorator
}

func newPrometheusClientImp(id string, config prometheusapi.Config, auth *ClientAuth, decorator QueryParamsDecorator) (prometheusapi.Client, error) {
	c, err := prometheusapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	nlpc := &PrometheusClient{
		id:        id,
		client:    c,
		decorator: decorator,
		auth:      auth,
	}

	return nlpc, nil
}

// ID is used to identify the type of client
func (nlpc *PrometheusClient) ID() string {
	return nlpc.id
}

// Passthrough to the prometheus client API
func (nlpc *PrometheusClient) URL(ep string, args map[string]string) *url.URL {
	return nlpc.client.URL(ep, args)
}

//passthrough to prometheus client API
func (nlpc *PrometheusClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	nlpc.auth.Apply(req)
	return nlpc.client.Do(ctx, req)
}

//--------------------------------------------------------------------------
//  RateLimitedPrometheusClient
//--------------------------------------------------------------------------

// RateLimitedPrometheusClient is a prometheus client which limits the total number of
// concurrent outbound requests allowed at a given moment.
type RateLimitedPrometheusClient struct {
	id        string
	client    prometheusapi.Client
	auth      *ClientAuth
	queue     queue.BlockingQueue
	decorator QueryParamsDecorator
	outbound  *atomic.AtomicInt32
}

// requestCounter is used to determine if the prometheus client keeps track of
// the concurrent outbound requests
type requestCounter interface {
	TotalQueuedRequests() int
	TotalOutboundRequests() int
}

// NewRateLimitedClient creates a prometheus client which limits the number of concurrent outbound
// prometheus requests.
func newRateLimitedClient(id string, config prometheusapi.Config, maxConcurrency int, auth *ClientAuth, decorator QueryParamsDecorator) (prometheusapi.Client, error) {
	c, err := prometheusapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	queue := queue.NewBlockingQueue()
	outbound := atomic.NewAtomicInt32(0)

	rlpc := &RateLimitedPrometheusClient{
		id:        id,
		client:    c,
		queue:     queue,
		decorator: decorator,
		outbound:  outbound,
		auth:      auth,
	}

	// Start concurrent request processing
	for i := 0; i < maxConcurrency; i++ {
		go rlpc.worker()
	}

	return rlpc, nil
}

// ID is used to identify the type of client
func (rlpc *RateLimitedPrometheusClient) ID() string {
	return rlpc.id
}

// TotalRequests returns the total number of requests that are either waiting to be sent and/or
// are currently outbound.
func (rlpc *RateLimitedPrometheusClient) TotalQueuedRequests() int {
	return rlpc.queue.Length()
}

// TotalOutboundRequests returns the total number of concurrent outbound requests, which have been
// sent to the server and are awaiting response.
func (rlpc *RateLimitedPrometheusClient) TotalOutboundRequests() int {
	return int(rlpc.outbound.Get())
}

// Passthrough to the prometheus client API
func (rlpc *RateLimitedPrometheusClient) URL(ep string, args map[string]string) *url.URL {
	return rlpc.client.URL(ep, args)
}

// workRequest is used to queue requests
type workRequest struct {
	ctx      context.Context
	req      *http.Request
	start    time.Time
	respChan chan *workResponse
	// used as a sentinel value to close the worker goroutine
	closer bool
	// request metadata for diagnostics
	contextName string
	query       string
}

// workResponse is the response payload returned to the Do method
type workResponse struct {
	res  *http.Response
	body []byte
	err  error
}

// worker is used as a consumer goroutine to pull workRequest from the blocking queue and execute them
func (rlpc *RateLimitedPrometheusClient) worker() {
	for {
		// blocks until there is an item available
		item := rlpc.queue.Dequeue()

		// Ensure the dequeued item was a workRequest
		if we, ok := item.(*workRequest); ok {
			// if we need to shut down all workers, we'll need to submit sentinel values
			// that will force the worker to return
			if we.closer {
				return
			}

			ctx := we.ctx
			req := we.req

			// decorate the raw query parameters
			if rlpc.decorator != nil {
				req.URL.RawQuery = rlpc.decorator(req.URL.Path, req.URL.Query()).Encode()
			}

			// measure time in queue
			timeInQueue := time.Since(we.start)

			// Increment outbound counter
			rlpc.outbound.Increment()

			// Execute Request
			roundTripStart := time.Now()
			res, body, err := rlpc.client.Do(ctx, req)

			// Decrement outbound counter
			rlpc.outbound.Decrement()
			LogQueryRequest(req, timeInQueue, time.Since(roundTripStart))

			// Pass back response data over channel to caller
			we.respChan <- &workResponse{
				res:  res,
				body: body,
				err:  err,
			}
		}
	}
}

// Rate limit and passthrough to prometheus client API
func (rlpc *RateLimitedPrometheusClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	rlpc.auth.Apply(req)

	respChan := make(chan *workResponse)
	defer close(respChan)

	// request names are used as a debug utility to identify requests in queue
	contextName := "<none>"
	if n, ok := httputil.GetName(req); ok {
		contextName = n
	}
	query, _ := httputil.GetQuery(req)

	rlpc.queue.Enqueue(&workRequest{
		ctx:         ctx,
		req:         req,
		start:       time.Now(),
		respChan:    respChan,
		closer:      false,
		contextName: contextName,
		query:       query,
	})

	workRes := <-respChan
	return workRes.res, workRes.body, workRes.err
}

// LogQueryRequest logs the query that was send to prom/thanos with the time in queue and total time after being sent
func LogQueryRequest(req *http.Request, queueTime time.Duration, sendTime time.Duration) {
	qp := httputil.NewQueryParams(req.URL.Query())
	query := qp.Get("query", "<Unknown>")

	klog.Infof("[Queue: %fs, Outbound: %fs][Query: %s]", queueTime.Seconds(), sendTime.Seconds(), query)
}
