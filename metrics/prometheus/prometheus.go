package cbytebuf

import "github.com/prometheus/client_golang/prometheus"

var (
	promCbytebufAcq     *prometheus.CounterVec
	promCbytebufRel     *prometheus.CounterVec
	promCbytebufPool    prometheus.Gauge
	promCbytebufPoolMem prometheus.Gauge

	_ = NewPrometheusMetrics
)

func init() {
	promCbytebufAcq = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cbytebuf_acq",
		Help: "Count of pool acquire.",
	}, []string{})
	promCbytebufRel = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cbytebuf_rel",
		Help: "Count of pool release.",
	}, []string{})
	promCbytebufPool = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cbytebuf_pool",
		Help: "Capacity of cbytebuf pool.",
	})
	promCbytebufPoolMem = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cbytebuf_pool_mem",
		Help: "Capacity of cbytebuf pool in bytes.",
	})

	prometheus.MustRegister(promCbytebufAcq, promCbytebufRel, promCbytebufPool, promCbytebufPoolMem)
}

// PrometheusMetrics implement cbytebuf.MetricsWriter interface.
type PrometheusMetrics struct{}

func NewPrometheusMetrics() *PrometheusMetrics {
	m := &PrometheusMetrics{}
	return m
}

func (m PrometheusMetrics) PoolAcquire(cap uint64) {
	promCbytebufAcq.WithLabelValues().Add(1)
	promCbytebufPool.Sub(1)
	promCbytebufPoolMem.Sub(float64(cap))
}

func (m PrometheusMetrics) PoolRelease(cap uint64) {
	promCbytebufRel.WithLabelValues().Add(1)
	promCbytebufPool.Add(1)
	promCbytebufPoolMem.Add(float64(cap))
}
