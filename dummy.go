package cbytebuf

// DummyMetrics writer.
// Used by default and does nothing.
type DummyMetrics struct{}

func (m DummyMetrics) PoolAcquire(_ uint64) {}
func (m DummyMetrics) PoolRelease(_ uint64) {}
