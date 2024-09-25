package monitoring

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ExposeMetrics(port int) {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func ExposeProfiling(port int) {
	logger.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
