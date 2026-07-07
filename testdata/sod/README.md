# TC-PA offline fixtures (pymrtd tier — never pooled with wire results)

Synthetic SOD/DSC material for passive-authentication divergence tests. No live PKD/CRL.

| ID | File | Condition |
| --- | --- | --- |
| TC-PA-01 | `tc-pa-01-weak-digest.json` | SOD signed with deprecated digest (policy probe) |
| TC-PA-03 | `tc-pa-03-expired-dsc.json` | Valid chain structure; DSC past `notAfter` |

Run: `scripts/run_offline_pa.sh`
