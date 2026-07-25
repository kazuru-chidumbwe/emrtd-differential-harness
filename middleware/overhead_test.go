package middleware

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/simulator"
)

// BenchmarkNegotiatePACEBAC reports reject-path and success-path latency for the
// explicit-reject wrapper (RQ2 overhead).
func BenchmarkNegotiatePACEBAC_rejectPath(b *testing.B) {
	pass := mustPassword(b)
	doc := mustPaceDoc(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
		_ = NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: false})
	}
}

func BenchmarkNegotiatePACEBAC_successPath(b *testing.B) {
	pass := mustPassword(b)
	doc := mustPaceDoc(b)
	old := runPace
	runPace = func(*iso7816.NfcSession, *document.Document, *password.Password) error { return nil }
	defer func() { runPace = old }()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
		_ = NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: false})
	}
}

// TestMiddlewareOverheadSample records a fixed-iteration latency sample for paper tables.
func TestMiddlewareOverheadSample(t *testing.T) {
	const iters = 200
	pass := mustPassword(t)
	doc := mustPaceDoc(t)

	reject := make([]time.Duration, 0, iters)
	for i := 0; i < iters; i++ {
		nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
		start := time.Now()
		_ = NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: false})
		reject = append(reject, time.Since(start))
	}

	old := runPace
	runPace = func(*iso7816.NfcSession, *document.Document, *password.Password) error { return nil }
	defer func() { runPace = old }()
	success := make([]time.Duration, 0, iters)
	for i := 0; i < iters; i++ {
		nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
		start := time.Now()
		_ = NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: false})
		success = append(success, time.Since(start))
	}

	fmt.Fprintf(os.Stderr, "middleware_overhead reject_path %s\n", summarizeDurations(reject))
	fmt.Fprintf(os.Stderr, "middleware_overhead success_path %s\n", summarizeDurations(success))
}

func summarizeDurations(ds []time.Duration) string {
	sort.Slice(ds, func(i, j int) bool { return ds[i] < ds[j] })
	var sum time.Duration
	for _, d := range ds {
		sum += d
	}
	n := len(ds)
	p50 := ds[n/2]
	p95 := ds[(n*95)/100]
	mean := sum / time.Duration(n)
	return fmt.Sprintf("n=%d mean=%s p50=%s p95=%s", n, mean, p50, p95)
}
