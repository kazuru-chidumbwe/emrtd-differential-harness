# Frozen schemas (v1) — bump version only with migration notes.

See [SCHEMA.md](../docs/SCHEMA.md) for the Scenario → Artifact pipeline.

| Schema | Version file | Current |
| --- | --- | --- |
| Artifact manifest | [`artifact-manifest-v1.json`](artifact-manifest-v1.json) | 1 |
| Run artifact | [`run-artifact-v1.json`](run-artifact-v1.json) | 1 |
| Provenance block | [`provenance-v1.json`](provenance-v1.json) | 1 |
| Profile catalog | [`profile-catalog-v1.json`](profile-catalog-v1.json) (validates `profiles/catalog.json` → `catalog_version`) | 1 |

## Enforcement

`classifier/verify_manifest.py` validates `artifact-manifest.json` and each referenced per-run JSON against these schemas (requires `jsonschema`; see `classifier/requirements-verify.txt`).

```bash
pip install -r classifier/requirements-verify.txt
python3 classifier/verify_manifest.py logs/suite-... --manifest suites/ac-01-wire.json
python3 -m unittest classifier.test_schemas -v
```
