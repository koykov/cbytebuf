package cbytebuf

// Dummy metrics writer.
// Used by default and does nothing.
type DummyMetrics struct{}

func (m *DummyMetrics) PoolAcquire(cap uint64) {}

func (m *DummyMetrics) PoolRelease(cap uint64) {}
