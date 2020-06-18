package cbytebuf

// Types of callback functions for pool acquire/release methods.
type PoolAcquireCallbackFn func(cap uint64)
type PoolReleaseCallbackFn func(cap uint64)

var (
	// Default instances of callback functions.
	poolAcqCb *PoolAcquireCallbackFn
	poolRelCb *PoolReleaseCallbackFn

	// Suppress go vet warnings.
	_, _ = RegisterPoolAcqCbFn, RegisterPoolRelCbFn
)

// Register pool acquire callback.
func RegisterPoolAcqCbFn(fn PoolAcquireCallbackFn) {
	poolAcqCb = &fn
}

// Register pool release callback.
func RegisterPoolRelCbFn(fn PoolReleaseCallbackFn) {
	poolRelCb = &fn
}
