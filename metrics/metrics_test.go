// Copyright Red Hat

package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
)

func TestPrometheusMiddleware(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	handlerToTest := PrometheusMiddleware(nextHandler)

	// mock the request
	req := httptest.NewRequest("GET", "http://test", nil)

	// call handler with mock response recorder
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)

	// confirm http_requests_total was incremented
	g.Expect(testutil.CollectAndCount(httpReqs)).To(gomega.Equal(1))

	// confirm http_response_time_seconds captured the duration
	metric := &dto.Metric{}
	h := httpDuration.With(prometheus.Labels{"path": ""})
	h.(prometheus.Histogram).Write(metric)
	g.Expect(int(*metric.Histogram.SampleCount)).To(gomega.Equal(1))
}
