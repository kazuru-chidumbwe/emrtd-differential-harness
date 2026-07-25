#!/usr/bin/env python3
"""Generate a combinatorial offline PA fixture grid (≥24 cells).

Axes (deterministic; n=1 per fixture):
  digest:      sha1 | md5 | sha256
  validity:    fresh | expired_1y | expired_5y | not_yet
  chain:       self | csca_dsc

SHA-256 + fresh cells are control fixtures (expect_policy_rejection=false).
All other cells set expect_policy_rejection=true (weak digest and/or bad validity).

Does not overwrite the four smoke anchors under testdata/sod/tc-pa-0*.json.
Outputs: testdata/sod/pa-sweep/*.json + testdata/sod/fixtures/pa-sweep-*
"""

from __future__ import annotations

import datetime
import json
from pathlib import Path

from asn1crypto import cms, algos, core, x509 as asn1_x509
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import padding, rsa
from cryptography.x509.oid import NameOID

ROOT = Path(__file__).resolve().parents[1]
FIXTURE_DIR = ROOT / "testdata" / "sod" / "fixtures"
CASE_DIR = ROOT / "testdata" / "sod" / "pa-sweep"
ID_MRTD_LDS_SECURITY_OBJECT = "2.23.136.1.1.1"
INSPECTION = datetime.datetime(2026, 7, 7, tzinfo=datetime.timezone.utc)

DIGESTS = ("sha1", "md5", "sha256")
VALIDITIES = ("fresh", "expired_1y", "expired_5y", "not_yet")
CHAINS = ("self", "csca_dsc")

DIGEST_OID = {"sha1": "sha1", "md5": "md5", "sha256": "sha256"}
HASH_CLS = {"sha1": hashes.SHA1, "md5": hashes.MD5, "sha256": hashes.SHA256}
DG_LEN = {"sha1": 20, "md5": 16, "sha256": 32}
CMS_SIG = {"sha1": "sha1_rsa", "md5": "md5_rsa", "sha256": "sha256_rsa"}


def _der_len(n: int) -> bytes:
    if n < 0x80:
        return bytes([n])
    b = []
    while n:
        b.insert(0, n & 0xFF)
        n >>= 8
    return bytes([0x80 | len(b)]) + bytes(b)


def retag_as_ef_sod(content_info_der: bytes) -> bytes:
    return b"\x77" + _der_len(len(content_info_der)) + content_info_der


def validity_window(v: str) -> tuple[datetime.datetime, datetime.datetime]:
    if v == "fresh":
        return (
            datetime.datetime(2018, 1, 1, tzinfo=datetime.timezone.utc),
            datetime.datetime(2030, 1, 1, tzinfo=datetime.timezone.utc),
        )
    if v == "expired_1y":
        return (
            datetime.datetime(2018, 1, 1, tzinfo=datetime.timezone.utc),
            datetime.datetime(2025, 7, 1, tzinfo=datetime.timezone.utc),
        )
    if v == "expired_5y":
        return (
            datetime.datetime(2015, 1, 1, tzinfo=datetime.timezone.utc),
            datetime.datetime(2020, 6, 1, tzinfo=datetime.timezone.utc),
        )
    # not_yet
    return (
        datetime.datetime(2027, 1, 1, tzinfo=datetime.timezone.utc),
        datetime.datetime(2035, 1, 1, tzinfo=datetime.timezone.utc),
    )


def build_lds(digest: str) -> bytes:
    hash_alg = algos.DigestAlgorithm({"algorithm": DIGEST_OID[digest]})
    version_der = core.Integer(0).dump()
    hash_alg_der = hash_alg.dump()
    dg_hash_entry_body = core.Integer(1).dump() + core.OctetString(b"\x11" * DG_LEN[digest]).dump()
    dg_hash_entry = b"\x30" + _der_len(len(dg_hash_entry_body)) + dg_hash_entry_body
    dg_hash_values = b"\x30" + _der_len(len(dg_hash_entry)) + dg_hash_entry
    body = version_der + hash_alg_der + dg_hash_values
    return b"\x30" + _der_len(len(body)) + body


def sign_sod(key: rsa.RSAPrivateKey, certs_der: list[bytes], econtent: bytes, digest: str) -> bytes:
    h = hashes.Hash(HASH_CLS[digest]())
    h.update(econtent)
    message_digest = h.finalize()
    signed_attrs = cms.CMSAttributes(
        [
            cms.CMSAttribute(
                {"type": "content_type", "values": [cms.ContentType(ID_MRTD_LDS_SECURITY_OBJECT)]}
            ),
            cms.CMSAttribute({"type": "message_digest", "values": [core.OctetString(message_digest)]}),
        ]
    )
    signature = key.sign(signed_attrs.dump(), padding.PKCS1v15(), HASH_CLS[digest]())
    dsc_der = certs_der[0]
    signer_info = cms.SignerInfo(
        {
            "version": "v1",
            "sid": cms.SignerIdentifier(
                {
                    "issuer_and_serial_number": cms.IssuerAndSerialNumber(
                        {
                            "issuer": asn1_x509.Certificate.load(dsc_der)["tbs_certificate"]["issuer"],
                            "serial_number": asn1_x509.Certificate.load(dsc_der)["tbs_certificate"][
                                "serial_number"
                            ],
                        }
                    )
                }
            ),
            "digest_algorithm": algos.DigestAlgorithm({"algorithm": DIGEST_OID[digest]}),
            "signed_attrs": signed_attrs,
            "signature_algorithm": algos.SignedDigestAlgorithm({"algorithm": CMS_SIG[digest]}),
            "signature": signature,
        }
    )
    cert_set = [
        cms.CertificateChoices({"certificate": asn1_x509.Certificate.load(d)}) for d in certs_der
    ]
    signed_data = cms.SignedData(
        {
            "version": "v3",
            "digest_algorithms": cms.DigestAlgorithms(
                [algos.DigestAlgorithm({"algorithm": DIGEST_OID[digest]})]
            ),
            "encap_content_info": cms.EncapsulatedContentInfo(
                {
                    "content_type": cms.ContentType(ID_MRTD_LDS_SECURITY_OBJECT),
                    "content": core.ParsableOctetString(econtent),
                }
            ),
            "certificates": cms.CertificateSet(cert_set),
            "signer_infos": cms.SignerInfos([signer_info]),
        }
    )
    return cms.ContentInfo({"content_type": "signed_data", "content": signed_data}).dump()


def make_cert(
    *,
    cn: str,
    key: rsa.RSAPrivateKey,
    issuer_key: rsa.RSAPrivateKey,
    issuer_name: x509.Name,
    nb: datetime.datetime,
    na: datetime.datetime,
    ca: bool,
) -> x509.Certificate:
    subject = x509.Name(
        [
            x509.NameAttribute(NameOID.COUNTRY_NAME, "ZZ"),
            x509.NameAttribute(NameOID.ORGANIZATION_NAME, "EMRTD Harness PA Sweep"),
            x509.NameAttribute(NameOID.COMMON_NAME, cn),
        ]
    )
    builder = (
        x509.CertificateBuilder()
        .subject_name(subject)
        .issuer_name(issuer_name)
        .public_key(key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(nb)
        .not_valid_after(na)
        .add_extension(x509.BasicConstraints(ca=ca, path_length=None if not ca else 0), critical=True)
    )
    return builder.sign(issuer_key, hashes.SHA256())


def build_cell(digest: str, validity: str, chain: str) -> dict:
    slug = f"{digest}-{validity}-{chain}".replace("_", "-")
    cid = f"TC-PA-sweep-{slug}"
    nb, na = validity_window(validity)
    dsc_key = rsa.generate_private_key(public_exponent=65537, key_size=2048)

    if chain == "self":
        name = x509.Name(
            [
                x509.NameAttribute(NameOID.COUNTRY_NAME, "ZZ"),
                x509.NameAttribute(NameOID.ORGANIZATION_NAME, "EMRTD Harness PA Sweep"),
                x509.NameAttribute(NameOID.COMMON_NAME, f"PA sweep self {slug}"),
            ]
        )
        dsc = make_cert(
            cn=f"PA sweep self {slug}",
            key=dsc_key,
            issuer_key=dsc_key,
            issuer_name=name,
            nb=nb,
            na=na,
            ca=False,
        )
        # Rebuild so subject == issuer (CN-based subject must match issuer_name bytes).
        dsc = (
            x509.CertificateBuilder()
            .subject_name(name)
            .issuer_name(name)
            .public_key(dsc_key.public_key())
            .serial_number(x509.random_serial_number())
            .not_valid_before(nb)
            .not_valid_after(na)
            .sign(dsc_key, hashes.SHA256())
        )
        dsc_der = dsc.public_bytes(serialization.Encoding.DER)
        certs = [dsc_der]
        csca_pem = dsc.public_bytes(serialization.Encoding.PEM).decode()
        dsc_pem = csca_pem
    else:
        csca_name = x509.Name(
            [
                x509.NameAttribute(NameOID.COUNTRY_NAME, "ZZ"),
                x509.NameAttribute(NameOID.ORGANIZATION_NAME, "EMRTD Harness PA Sweep"),
                x509.NameAttribute(NameOID.COMMON_NAME, f"PA sweep CSCA {slug}"),
            ]
        )
        csca_key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
        csca = (
            x509.CertificateBuilder()
            .subject_name(csca_name)
            .issuer_name(csca_name)
            .public_key(csca_key.public_key())
            .serial_number(x509.random_serial_number())
            .not_valid_before(datetime.datetime(2010, 1, 1, tzinfo=datetime.timezone.utc))
            .not_valid_after(datetime.datetime(2040, 1, 1, tzinfo=datetime.timezone.utc))
            .add_extension(x509.BasicConstraints(ca=True, path_length=0), critical=True)
            .sign(csca_key, hashes.SHA256())
        )
        dsc = make_cert(
            cn=f"PA sweep DSC {slug}",
            key=dsc_key,
            issuer_key=csca_key,
            issuer_name=csca.subject,
            nb=nb,
            na=na,
            ca=False,
        )
        dsc_der = dsc.public_bytes(serialization.Encoding.DER)
        csca_der = csca.public_bytes(serialization.Encoding.DER)
        certs = [dsc_der, csca_der]
        csca_pem = csca.public_bytes(serialization.Encoding.PEM).decode()
        dsc_pem = dsc.public_bytes(serialization.Encoding.PEM).decode()

    econtent = build_lds(digest)
    sod = retag_as_ef_sod(sign_sod(dsc_key, certs, econtent, digest))

    sod_path = FIXTURE_DIR / f"pa-sweep-{slug}-sod.hex"
    csca_path = FIXTURE_DIR / f"pa-sweep-{slug}-csca.pem"
    dsc_path = FIXTURE_DIR / f"pa-sweep-{slug}-dsc.pem"
    sod_path.write_text(sod.hex() + "\n", encoding="utf-8")
    csca_path.write_text(f"# PA sweep {slug}\n" + csca_pem, encoding="utf-8")
    dsc_path.write_text(f"# PA sweep {slug}\n" + dsc_pem, encoding="utf-8")

    # Control: modern digest + fresh validity → no policy rejection expected from digest/expiry.
    expect_reject = not (digest == "sha256" and validity == "fresh")
    condition = f"pa_{digest}_{validity}_{chain}"
    case = {
        "id": cid,
        "name": f"PA sweep digest={digest} validity={validity} chain={chain}",
        "mechanism": "PA",
        "condition": condition,
        "tier": "offline",
        "fixture": str(sod_path.relative_to(ROOT)).replace("\\", "/"),
        "dsc_fixture": str(dsc_path.relative_to(ROOT)).replace("\\", "/"),
        "csca_fixture": str(csca_path.relative_to(ROOT)).replace("\\", "/"),
        "sweep_axes": {"digest": digest, "validity": validity, "chain": chain},
        "inspection_date": INSPECTION.date().isoformat(),
        "expect_policy_rejection": expect_reject,
        "notes": (
            "Combinatorial offline PA grid cell. Live PKD/CRL out of scope. "
            f"expect_policy_rejection={expect_reject}."
        ),
    }
    case_path = CASE_DIR / f"{cid}.json"
    case_path.write_text(json.dumps(case, indent=2) + "\n", encoding="utf-8")
    return {
        "id": cid,
        "path": str(case_path.relative_to(ROOT)).replace("\\", "/"),
        "digest": digest,
        "validity": validity,
        "chain": chain,
        "expect_policy_rejection": expect_reject,
    }


def main() -> None:
    FIXTURE_DIR.mkdir(parents=True, exist_ok=True)
    CASE_DIR.mkdir(parents=True, exist_ok=True)
    entries = []
    for digest in DIGESTS:
        for validity in VALIDITIES:
            for chain in CHAINS:
                entries.append(build_cell(digest, validity, chain))
    index = {
        "generator": "profiles/generate_pa_sweep.py",
        "n_fixtures": len(entries),
        "axes": {"digest": list(DIGESTS), "validity": list(VALIDITIES), "chain": list(CHAINS)},
        "factorial_note": (
            f"{len(DIGESTS)}×{len(VALIDITIES)}×{len(CHAINS)} = {len(entries)} fixtures, n=1 each. "
            "Combinatorial policy coverage, not statistical variance."
        ),
        "cases": entries,
    }
    (CASE_DIR / "index.json").write_text(json.dumps(index, indent=2) + "\n", encoding="utf-8")
    print(f"wrote {len(entries)} PA sweep cases under {CASE_DIR}")


if __name__ == "__main__":
    main()
