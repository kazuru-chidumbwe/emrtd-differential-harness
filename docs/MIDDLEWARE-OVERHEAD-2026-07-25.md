# Middleware overhead — test-server pin (R1)

Answers Reviewer 1 §4: RQ2 is miss-rate + FP controls + **measured cost**, not “every mitigated run scored 2” alone.

**Host:** `test-server` · Go 1.25.0 · `go test ./middleware/ -run TestMiddlewareOverheadSample -count=1 -v`  
**Log:** `logs/middleware-overhead-r1-test-server.txt`

## Latency (n=200 iters each path)

| Path | mean | p50 | p95 |
| --- | ---: | ---: | ---: |
| Reject path (`NegotiatePACEBAC` after failed PACE, no BAC fallback) | 13.56 µs | 11.30 µs | 28.63 µs |
| Success path (PACE succeeds through negotiate wrapper) | 88.65 µs | 83.32 µs | 128.39 µs |

Env: in-process Go simulator stubs (same package tests as FP controls). Not a full APDU round-trip wall clock.

## Binary size delta

| Binary | Size (bytes) |
| --- | ---: |
| `cmd/tc-ac-01` (baseline) | 5 390 398 |
| `cmd/tc-ac-01-mitigated` (negotiate middleware) | 5 429 865 |
| **Delta** | **+39 467 (~38.5 KiB)** |

Source: `logs/middleware-binsize-r1-test-server.txt` (`go build -o` unstripped).

## Framing for RQ2

- Zero missed detections on the PACE-to-BAC failure arm remains the primary safety claim.
- FP controls: unit matrix + prior success-path pin `31aa96db…`.
- Cost: microsecond-scale reject/success overhead and ~39 KiB binary delta on the mitigated driver — negligible vs chip RTT.
- “Every mitigated run scored 2” is the designed outcome of explicit-reject, not the sole contribution.
