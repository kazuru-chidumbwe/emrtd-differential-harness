"""Load frozen v1 JSON Schemas and validate documents."""

from __future__ import annotations

import json
from functools import lru_cache
from pathlib import Path
from typing import Any

SCHEMA_DIR = Path(__file__).resolve().parent.parent / "schemas"

MANIFEST_SCHEMA_ID = (
    "https://github.com/kazuru-chidumbwe/emrtd-differential-harness/schemas/artifact-manifest-v1.json"
)
RUN_SCHEMA_ID = (
    "https://github.com/kazuru-chidumbwe/emrtd-differential-harness/schemas/run-artifact-v1.json"
)
CATALOG_SCHEMA_ID = (
    "https://github.com/kazuru-chidumbwe/emrtd-differential-harness/schemas/profile-catalog-v1.json"
)


class SchemaUnavailableError(RuntimeError):
    """Raised when jsonschema is not installed or schema files are missing."""


@lru_cache(maxsize=1)
def _registry():
    try:
        from referencing import Registry, Resource
        from referencing.jsonschema import DRAFT202012
    except ImportError as e:
        raise SchemaUnavailableError(
            "jsonschema stack missing — pip install -r classifier/requirements-verify.txt"
        ) from e

    pairs = []
    for name in (
        "provenance-v1.json",
        "run-artifact-v1.json",
        "artifact-manifest-v1.json",
        "profile-catalog-v1.json",
    ):
        path = SCHEMA_DIR / name
        if not path.is_file():
            raise SchemaUnavailableError(f"missing schema file: {path}")
        contents = json.loads(path.read_text(encoding="utf-8"))
        pairs.append(
            (contents["$id"], Resource.from_contents(contents, default_specification=DRAFT202012))
        )
    return Registry().with_resources(pairs)


def validate_document(schema_id: str, document: Any) -> list[str]:
    """Return a list of validation error strings (empty = OK)."""
    try:
        from jsonschema import Draft202012Validator
    except ImportError as e:
        raise SchemaUnavailableError(
            "jsonschema missing — pip install -r classifier/requirements-verify.txt"
        ) from e

    registry = _registry()
    schema = registry.contents(schema_id)
    validator = Draft202012Validator(schema, registry=registry)
    return [f"{e.json_path}: {e.message}" for e in sorted(validator.iter_errors(document), key=str)]


def validate_artifact_manifest(document: Any) -> list[str]:
    return validate_document(MANIFEST_SCHEMA_ID, document)


def validate_run_artifact(document: Any) -> list[str]:
    return validate_document(RUN_SCHEMA_ID, document)


def validate_profile_catalog(document: Any) -> list[str]:
    return validate_document(CATALOG_SCHEMA_ID, document)
