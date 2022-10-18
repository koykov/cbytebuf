package cbytebuf

type MetricsWriter interface {
	// PoolAcquire registers acquire of cbyte object from pool.
	PoolAcquire(cap uint64)
	// PoolRelease registers release of cbyte object back to pool.
	PoolRelease(cap uint64)
}

var (
	// Builtin instance of metrics writer.
	// By default is a DummyMetrics object that does nothing on call.
	metricsHandler MetricsWriter = &DummyMetrics{}

	_ = RegisterMetricsHandler
)

// RegisterMetricsHandler registers new metrics handler.
func RegisterMetricsHandler(handler MetricsWriter) {
	metricsHandler = handler
}
