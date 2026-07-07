# Simulator profiles

One JSON/YAML file per test case ID from [`g1-static-scoping.md`](../../../private/notes/g1-static-scoping.md).

**Dev.to requirement:** profiles must run **without physical hardware** — synthetic chip only. Blog copy: *"no physical passport required."*

**Planned (W2 — blocking for Aug blog §4B):**

- `pace-then-bac-downgrade.json` — **TC-AC-01** (smoke OK 2026-07-07, gmrtd + jmrtd)
- `ca-v1-v2-skew.json` — **TC-CA-01** (EAC-CA MSE fail; gmrtd smoke scaffold)
- `ca-v2-terminal-v1.json` — TC-CA-01
- `fi-cardaccess-truncated.json` — TC-FI-01

**Deliverable:** JMRTD passport applet simulator + profile loader + `scripts/quick_test.sh` (gmrtd path live).
