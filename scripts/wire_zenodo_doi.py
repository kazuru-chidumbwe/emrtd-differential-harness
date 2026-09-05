#!/usr/bin/env python3
"""Wire minted Zenodo DOIs into harness cite surfaces (live pin remains v1.0.9)."""
from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
VERSION_DOI = "10.5281/zenodo.22097289"
CONCEPT_DOI = "10.5281/zenodo.22095366"
VERSION_URL = f"https://doi.org/{VERSION_DOI}"
CONCEPT_URL = f"https://doi.org/{CONCEPT_DOI}"

CITATION = f"""cff-version: 1.2.0
title: "eMRTD differential harness"
message: >-
  If you use this software or its locked-run artifacts, please cite it using
  the metadata below.
type: software
authors:
  - family-names: Kazuru
    given-names: Seke
    orcid: "https://orcid.org/0009-0002-4099-1059"
    email: kazuruuni@gmail.com
    affiliation: Independent Researcher
repository-code: "https://github.com/kazuru-chidumbwe/emrtd-differential-harness"
url: "https://github.com/kazuru-chidumbwe/emrtd-differential-harness"
abstract: >-
  Differential test harness for open-source eMRTD reader libraries. Identical
  synthetic chip profiles are replayed against multiple stacks; each run captures
  APDU traces and classifies whether negotiation failures or specification-permitted
  downgrades reach a typical application caller (Observability Score).
keywords:
  - eMRTD
  - ICAO Doc 9303
  - differential testing
  - PACE
  - BAC
  - observability
  - JMRTD
  - gmrtd
  - pymrtd
license: MIT
version: "1.0.9"
date-released: "2026-08-25"
identifiers:
  - type: doi
    value: "{VERSION_DOI}"
  - type: doi
    value: "{CONCEPT_DOI}"
preferred-citation:
  type: software
  authors:
    - family-names: Kazuru
      given-names: Seke
      orcid: "https://orcid.org/0009-0002-4099-1059"
  title: "eMRTD differential harness"
  version: "1.0.9"
  doi: "{VERSION_DOI}"
  url: "{VERSION_URL}"
  year: 2026
  month: 8
"""

ZENODO_JSON = {
    "title": "eMRTD differential harness",
    "upload_type": "software",
    "description": (
        "Differential test harness for open-source eMRTD reader libraries. Identical "
        "synthetic chip profiles are replayed against multiple stacks; each run captures "
        "APDU traces and classifies whether negotiation failures or specification-permitted "
        "downgrades reach a typical application caller (Observability Score). Synthetic APDU "
        "profiles only; not RF / silicon. Cite SemVer tag v1.0.9. Version DOI "
        f"{VERSION_DOI}; concept DOI {CONCEPT_DOI}. Defect locked-run trees: Release v1.0.7 "
        "asset emrtd-locked-runs-v1.0.7.zip. Remeasurement (gmrtd v1.1.3): Release v1.0.9 "
        "asset emrtd-gmrtd-v1.1.3-remeasurement-v1.0.8.zip (digest 04cdd3dd…)."
    ),
    "creators": [
        {
            "name": "Kazuru, Seke",
            "affiliation": "Independent Researcher",
            "orcid": "0009-0002-4099-1059",
        }
    ],
    "keywords": [
        "eMRTD",
        "ICAO Doc 9303",
        "differential testing",
        "PACE",
        "BAC",
        "observability",
        "JMRTD",
        "gmrtd",
        "pymrtd",
    ],
    "license": "MIT",
    "access_right": "open",
    "related_identifiers": [
        {
            "identifier": VERSION_DOI,
            "relation": "isIdenticalTo",
            "resource_type": "software",
            "scheme": "doi",
        },
        {
            "identifier": "https://dev.to/kazuru_73322ef9a7d6ed2b18/differential-testing-revealed-what-conformance-testing-missed-a-case-study-with-open-source-emrtd-1nie",
            "relation": "isDocumentedBy",
            "resource_type": "publication-article",
        },
        {
            "identifier": "https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.7",
            "relation": "isSupplementedBy",
            "resource_type": "dataset",
        },
        {
            "identifier": "https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.9",
            "relation": "isSupplementedBy",
            "resource_type": "dataset",
        },
    ],
    "notes": (
        "GitHub→Zenodo auto-mint for evidence tag v1.0.9. Version DOI "
        f"{VERSION_DOI}; concept series {CONCEPT_DOI}. Locked-run digests remain on "
        "GitHub Releases v1.0.7 (defect factorial) and v1.0.9 (gmrtd v1.1.3 remeasurement). "
        "If a peer-reviewed article receives a DOI later, add it as related_identifiers "
        "with relation isDocumentedBy."
    ),
}


def main() -> None:
    import json

    (ROOT / "CITATION.cff").write_text(CITATION, encoding="utf-8", newline="\n")
    (ROOT / ".zenodo.json").write_text(
        json.dumps(ZENODO_JSON, indent=2) + "\n", encoding="utf-8", newline="\n"
    )

    # README replacements
    readme = (ROOT / "README.md").read_text(encoding="utf-8")
    readme = readme.replace(
        "Zenodo DOI **pending**",
        f"Zenodo [{VERSION_DOI}]({VERSION_URL})",
    )
    readme = readme.replace("Zenodo DOI pending.", f"Zenodo {VERSION_DOI}.")
    readme = readme.replace(
        "Zenodo DOI pending;",
        f"Zenodo {VERSION_DOI};",
    )
    readme = readme.replace(
        "GitHub Release locked trees; Zenodo DOI pending).",
        f"GitHub Release locked trees; Zenodo {VERSION_DOI}).",
    )
    (ROOT / "README.md").write_text(readme, encoding="utf-8", newline="\n")

    tags = (ROOT / "docs" / "TAGS.md").read_text(encoding="utf-8")
    tags = tags.replace(
        "- Zenodo DOI **pending**; GitHub Release URLs are the canonical deposits until a DOI is minted.\n",
        f"- Zenodo version DOI **`{VERSION_DOI}`** ({VERSION_URL}); concept series `{CONCEPT_DOI}`.\n"
        f"- GitHub Release locked-run zips remain the byte-identical digests; Zenodo archives the tagged source tree.\n",
    )
    tags = tags.replace(
        "| [`v1.0.9`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.9) | **Live paper cite pin** — cite-surface sync (CITATION/DAS/README); dual-pin docs; remeasurement asset. |\n",
        f"| [`v1.0.9`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.9) | **Live paper cite pin** — Zenodo `{VERSION_DOI}`; dual-pin docs; remeasurement asset. |\n",
    )
    (ROOT / "docs" / "TAGS.md").write_text(tags, encoding="utf-8", newline="\n")

    das = (ROOT / "docs" / "DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md").read_text(
        encoding="utf-8"
    )
    das = das.replace(
        "A Zenodo DOI is **pending**; until minted, the Release URLs are the canonical deposits.",
        f"Zenodo version DOI **`{VERSION_DOI}`** ({VERSION_URL}); concept series `{CONCEPT_DOI}`. "
        "GitHub Release locked-run zips remain the byte-identical digests for `e15f4b57…` / `04cdd3dd…`; "
        "Zenodo archives the tagged source tree (GitHub→Zenodo auto-mint for `v1.0.9`).",
    )
    das = das.replace(
        "Zenodo DOI minting from the Release archives remains pending.",
        f"Zenodo archive: {VERSION_URL} (concept {CONCEPT_URL}).",
    )
    das = das.replace(
        "Verify digests against the GitHub Release zips (or Zenodo DOI once minted).",
        f"Verify digests against the GitHub Release zips; cite the software archive at {VERSION_URL}.",
    )
    if "**Zenodo version DOI:**" not in das:
        das = das.replace(
            "**Evidence pin (live cite):** tag `v1.0.9` — https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.9  \n",
            "**Evidence pin (live cite):** tag `v1.0.9` — https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.9  \n"
            f"**Zenodo version DOI:** `{VERSION_DOI}` — {VERSION_URL}  \n"
            f"**Zenodo concept DOI:** `{CONCEPT_DOI}` — {CONCEPT_URL}  \n",
        )
    # research pin table row if Zenodo pending
    das = das.replace(
        "| Zenodo | pending | DOI TBD |",
        f"| Zenodo | {VERSION_URL} | `{VERSION_DOI}` |",
    )
    # table in DAS research pin section
    if "| Zenodo version DOI |" not in das:
        das = das.replace(
            "| Annotated tag (live cite) | `v1.0.9` |\n",
            "| Annotated tag (live cite) | `v1.0.9` |\n"
            f"| Zenodo version DOI | `{VERSION_DOI}` |\n"
            f"| Zenodo concept DOI | `{CONCEPT_DOI}` |\n",
        )
    (ROOT / "docs" / "DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md").write_text(
        das, encoding="utf-8", newline="\n"
    )

    changelog = (ROOT / "CHANGELOG.md").read_text(encoding="utf-8")
    if "## [1.0.9+doi]" not in changelog and "22097289" not in changelog.split("## [1.0.8]")[0]:
        block = f"""## [Unreleased]

## [1.0.9] — 2026-08-25 (DOI wire)

### Added

- Zenodo GitHub→Zenodo archive for tag `v1.0.9`: version DOI `{VERSION_DOI}`; concept `{CONCEPT_DOI}`.

### Changed

- `CITATION.cff`, `.zenodo.json`, README, DAS, and TAGS now declare the minted DOIs (no longer “pending”).

"""
        changelog = changelog.replace("## [Unreleased]\n\n## [1.0.9] — 2026-08-25\n", block, 1)
        # Keep original 1.0.9 section title for history
        changelog = changelog.replace(
            "## [1.0.9] — 2026-08-25 (DOI wire)\n\n### Added\n\n"
            f"- Zenodo GitHub→Zenodo archive for tag `v1.0.9`: version DOI `{VERSION_DOI}`; concept `{CONCEPT_DOI}`.\n\n"
            "### Changed\n\n"
            "- `CITATION.cff`, `.zenodo.json`, README, DAS, and TAGS now declare the minted DOIs (no longer “pending”).\n\n",
            f"""## [1.0.9] — 2026-08-25

### Added

- Zenodo GitHub→Zenodo archive for tag `v1.0.9`: version DOI `{VERSION_DOI}`; concept `{CONCEPT_DOI}`.
- Cite-surface DOI wire on tip (same SemVer; archival deposit unchanged).

### Changed

- `CITATION.cff`, `.zenodo.json`, README, DAS, and TAGS declare minted DOIs (no longer “pending”).

""",
            1,
        )
    # Fix if we duplicated Unreleased badly
    while "## [Unreleased]\n\n## [Unreleased]\n" in changelog:
        changelog = changelog.replace("## [Unreleased]\n\n## [Unreleased]\n", "## [Unreleased]\n", 1)
    (ROOT / "CHANGELOG.md").write_text(changelog, encoding="utf-8", newline="\n")

    zdoc = ROOT / "docs" / "ZENODO.md"
    zdoc.write_text(
        f"""# Zenodo archive

| Field | Value |
| --- | --- |
| Version DOI | [`{VERSION_DOI}`]({VERSION_URL}) |
| Concept DOI | [`{CONCEPT_DOI}`]({CONCEPT_URL}) |
| GitHub tag | [`v1.0.9`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.9) |
| Mint path | GitHub→Zenodo auto-mint (2026-08-25) |
| Record | https://zenodo.org/records/22097289 |

Locked-run digests (`e15f4b57…`, `b029a9dc…`, `04cdd3dd…`) remain on GitHub Releases `v1.0.7` / `v1.0.9`. Zenodo holds the tagged source tree zip for the software cite.
""",
        encoding="utf-8",
        newline="\n",
    )
    print("wired", VERSION_DOI, CONCEPT_DOI)


if __name__ == "__main__":
    main()
