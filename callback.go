package cbytebuf

// Types of callback functions for pool acquire/release methods.
type PoolAcquireCallbackFn func(cap uint64)
type PoolReleaseCallbackFn func(cap uint64)

var (
	// Default instances of callback functions.
	poolAckCb *PoolAcquireCallbackFn
	poolRelCb *PoolReleaseCallbackFn

	// Suppress go vet warnings.
	_, _ = RegisterPoolAckCbFn, RegisterPoolRelCbFn
)

// Register pool acquire callback.
func RegisterPoolAckCbFn(fn PoolAcquireCallbackFn) {
	poolAckCb = &fn
}

// Register pool release callback.
func RegisterPoolRelCbFn(fn PoolReleaseCallbackFn) {
	poolRelCb = &fn
}
