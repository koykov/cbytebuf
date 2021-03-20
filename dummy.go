package cbytebuf

type DummyMetrics struct{}

func (m *DummyMetrics) PoolAcquire(cap uint64) {}

func (m *DummyMetrics) PoolRelease(cap uint64) {}
