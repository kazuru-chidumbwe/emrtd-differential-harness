# TC-PA offline fixtures (pymrtd tier — never pooled with wire results)

Synthetic SOD/DSC material for passive-authentication divergence tests. No live PKD/CRL.

| ID | File | Condition |
| --- | --- | --- |
| TC-PA-01 | `tc-pa-01-weak-digest.json` | SOD signed with deprecated digest (policy probe) |
| TC-PA-03 | `tc-pa-03-expired-dsc.json` | Self-signed; DSC past `notAfter` |
| TC-PA-04a | `tc-pa-04a-chained-fresh.json` | Offline CSCA→DSC chain; fresh DSC |
| TC-PA-04b | `tc-pa-04b-chained-expired.json` | Offline CSCA→DSC chain; expired DSC |

Generate 04: `python profiles/generate_pa04_chained_fixture.py` (no live PKD).

Run: `scripts/run_offline_pa.sh`
