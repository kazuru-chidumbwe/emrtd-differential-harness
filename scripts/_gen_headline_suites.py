import json
from pathlib import Path

root = Path(__file__).resolve().parents[1]
full = json.loads((root / "suites/ac-01-sweep-full.json").read_text(encoding="utf-8"))
entries_40 = [e for e in full["entries"] if "-sweep-orig-" in e.get("test_case", "")]
assert len(entries_40) == 40, len(entries_40)
headline = {
    "suite_id": "ac-01-headline-40",
    "description": (
        "Headline N=40: 5 SW × 2 injection × 2 libs × 2 variants; "
        "MRZ fixed to orig (representative triple). Subset of ac-01-sweep-full."
    ),
    "seed": 1,
    "n": 1,
    "mrz_id": "orig",
    "parent_suite": "ac-01-sweep-full",
    "entries": entries_40,
}
(root / "suites/ac-01-headline-40.json").write_text(
    json.dumps(headline, indent=2) + "\n", encoding="utf-8"
)
entries_j = [e for e in entries_40 if e.get("library") == "jmrtd"]
assert len(entries_j) == 20, len(entries_j)
jonly = {
    "suite_id": "ac-01-headline-jmrtd-20",
    "description": (
        "JMRTD half of headline N=40 (5 SW × 2 injection × 2 variants, MRZ=orig). "
        "For Class B pace_exception_class×SW table when gmrtd half runs elsewhere."
    ),
    "seed": 1,
    "n": 1,
    "mrz_id": "orig",
    "parent_suite": "ac-01-headline-40",
    "entries": entries_j,
}
(root / "suites/ac-01-headline-jmrtd-20.json").write_text(
    json.dumps(jonly, indent=2) + "\n", encoding="utf-8"
)
print("wrote suites/ac-01-headline-40.json", len(entries_40))
print("wrote suites/ac-01-headline-jmrtd-20.json", len(entries_j))
