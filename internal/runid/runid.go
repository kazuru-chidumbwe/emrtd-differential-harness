package runid

import (
	"fmt"
	"sync/atomic"
	"time"
)

var seq uint64

// New returns a unique run_id: {prefix}-{utc_timestamp}-{monotonic_counter}.
func New(prefix string) string {
	n := atomic.AddUint64(&seq, 1)
	ts := time.Now().UTC().Format("20060102T150405.000000Z")
	return fmt.Sprintf("%s-%s-%06d", prefix, ts, n)
}
