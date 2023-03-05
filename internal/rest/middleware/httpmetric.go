package middleware

import (
	"fmt"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/metric"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func HTTPMetric(metrics *metric.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(rw http.ResponseWriter, r *http.Request) {
			now := time.Now()
			lrw := newLoggingResponseWriter(rw)

			defer func() {
				metrics.RequestDuration.With(prometheus.Labels{
					"method": r.Method,
					"status": fmt.Sprintf("%d", lrw.statusCode)},
				).Observe(time.Since(now).Seconds())
				metrics.RequestCount.With(prometheus.Labels{"type": r.URL.RequestURI()}).Inc()
			}()
			next.ServeHTTP(lrw, r)
		}

		return http.HandlerFunc(fn)
	}
}
