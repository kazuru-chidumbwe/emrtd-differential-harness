#!/usr/bin/env python3
"""Generate a minimal, self-signed, SHA-1-digest synthetic EF.SOD fixture for TC-PA-01.

Scope, stated explicitly per the reviewer's guidance rather than implied:

  - This is a SELF-SIGNED, NON-CHAINED fixture. There is no CSCA/DSC trust chain — the same
    key pair signs and is embedded as the sole certificate in the CMS SignedData. Real eMRTD
    trust validation (CSCA -> DSC chain, PKD lookup, CRL checking) is explicitly out of scope
    for this fixture and for TC-PA-01/03 generally. Do not present this as testing trust-chain
    validation; it tests only whether the digest-algorithm/expiry condition itself is
    surfaced when a library parses and verifies the SOD's signature.
  - The digest algorithm inside the LDS Security Object (the DG hash algorithm ICAO 9303
    requires the SOD to declare) is set to SHA-1, a deprecated algorithm for this purpose.
    The CMS SignedData's own signature is also produced with SHA-1WithRSA for consistency
    (a real weak-digest chip would plausibly use the same deprecated algorithm throughout).
  - No physical travel document, no real country signing key, no live PKD/CRL infrastructure.

Writes:
  testdata/sod/fixtures/tc-pa-01-sod.hex   (hex-encoded DER CMS SignedData)
  testdata/sod/fixtures/tc-pa-01-csca.pem  (self-signed cert, reused as both CSCA and DSC —
                                             see scope note above; this is NOT a real chain)
  testdata/sod/fixtures/tc-pa-01-dsc.pem   (identical cert, present because run_case.py's
                                             fixtures_present() check expects both files)
"""

from __future__ import annotations

import datetime
from pathlib import Path

from asn1crypto import cms, algos, core, x509 as asn1_x509
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.x509.oid import NameOID

ROOT = Path(__file__).resolve().parents[1]
FIXTURE_DIR = ROOT / "testdata" / "sod" / "fixtures"

ID_MRTD_LDS_SECURITY_OBJECT = "2.23.136.1.1.1"
SHA1_OID = "1.3.14.3.2.26"


def build_self_signed_cert() -> tuple[rsa.RSAPrivateKey, x509.Certificate]:
    key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
    subject = issuer = x509.Name([
        x509.NameAttribute(NameOID.COUNTRY_NAME, "ZZ"),
        x509.NameAttribute(NameOID.ORGANIZATION_NAME, "EMRTD Harness Synthetic Fixture"),
        x509.NameAttribute(NameOID.COMMON_NAME, "TC-PA-01 self-signed, non-chained, SHA-1"),
    ])
    now = datetime.datetime.now(datetime.timezone.utc)
    cert = (
        x509.CertificateBuilder()
        .subject_name(subject)
        .issuer_name(issuer)
        .public_key(key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(now - datetime.timedelta(days=1))
        .not_valid_after(now + datetime.timedelta(days=3650))
        .sign(key, hashes.SHA256())  # outer cert signature — SHA-256; the `cryptography`
        # library refuses to produce SHA-1 X.509 signatures at all (hard guardrail, not
        # configurable). The condition under test is the LDS Security Object's declared
        # DG-hash algorithm (see build_lds_security_object) and the CMS SignerInfo signature
        # below, not this outer certificate's own signature — a real weak-digest chip would
        # plausibly still carry a modern DSC certificate signed by a CSCA using SHA-256; the
        # weak-digest condition is specifically about the LDS/CMS layer ICAO 9303 governs,
        # not the X.509 PKI layer around it.
    )
    return key, cert


def build_lds_security_object() -> bytes:
    """DER-encode a minimal ICAO 9303 LDSSecurityObject with a SHA-1 hashAlgorithm and a
    single synthetic DG1 hash entry. Structure mirrors pymrtd.ef.sod.LDSSecurityObject:
    SEQUENCE { version INTEGER, hashAlgorithm AlgorithmIdentifier, dataGroupHashValues
    SEQUENCE OF DataGroupHash }. Built with explicit DER byte construction rather than
    asn1crypto's internal (non-public-API) SOD classes, to keep this generator decoupled
    from pymrtd's implementation details."""
    hash_alg = algos.DigestAlgorithm({"algorithm": "sha1"})
    version_der = core.Integer(0).dump()
    hash_alg_der = hash_alg.dump()

    dg_hash_entry_body = core.Integer(1).dump() + core.OctetString(b"\x11" * 20).dump()
    dg_hash_entry = b"\x30" + _der_len(len(dg_hash_entry_body)) + dg_hash_entry_body
    dg_hash_values = b"\x30" + _der_len(len(dg_hash_entry)) + dg_hash_entry

    body = version_der + hash_alg_der + dg_hash_values
    return b"\x30" + _der_len(len(body)) + body


def _der_len(n: int) -> bytes:
    if n < 0x80:
        return bytes([n])
    b = []
    while n:
        b.insert(0, n & 0xFF)
        n >>= 8
    return bytes([0x80 | len(b)]) + bytes(b)


def retag_as_ef_sod(content_info_der: bytes) -> bytes:
    """Wrap a standard CMS ContentInfo DER encoding in the ICAO 9303 EF.SOD tag
    (APPLICATION 23, constructed, tag byte 0x77). This is an EXPLICIT wrapper — the full
    inner ContentInfo TLV (including its own SEQUENCE tag and length) is preserved
    unchanged inside the new outer tag, not stripped. (An earlier version of this function
    incorrectly treated it as IMPLICIT and stripped the inner tag+length; that produced a
    structure pymrtd's parser rejected with a confusing "extra data" error, since
    `_content_spec.load()` expects to see a complete, self-delimiting ContentInfo — tag and
    all — as the EF's contents.)"""
    if content_info_der[0] != 0x30:
        raise ValueError(f"expected a UNIVERSAL SEQUENCE (0x30) to wrap, got {content_info_der[0]:#x}")
    return b"\x77" + _der_len(len(content_info_der)) + content_info_der


def build_signed_data(key: rsa.RSAPrivateKey, cert_der: bytes, econtent: bytes) -> bytes:
    from cryptography.hazmat.primitives.asymmetric import padding

    digest = hashes.Hash(hashes.SHA1())
    digest.update(econtent)
    message_digest = digest.finalize()

    signed_attrs = cms.CMSAttributes([
        cms.CMSAttribute({
            "type": "content_type",
            "values": [cms.ContentType(ID_MRTD_LDS_SECURITY_OBJECT)],
        }),
        cms.CMSAttribute({
            "type": "message_digest",
            "values": [core.OctetString(message_digest)],
        }),
    ])
    signed_attrs_der_for_sig = signed_attrs.dump()
    # For signing, the SET OF tag (0x31) must be used in the value that gets signed even
    # though it's an implicit [0] IMPLICIT SET in the encoded SignerInfo — asn1crypto's
    # CMSAttributes.dump() already produces the correct SET encoding for signing purposes.

    signature = key.sign(signed_attrs_der_for_sig, padding.PKCS1v15(), hashes.SHA1())

    signer_info = cms.SignerInfo({
        "version": "v1",
        "sid": cms.SignerIdentifier({
            "issuer_and_serial_number": cms.IssuerAndSerialNumber({
                "issuer": asn1_x509.Certificate.load(cert_der)["tbs_certificate"]["issuer"],
                "serial_number": asn1_x509.Certificate.load(cert_der)["tbs_certificate"]["serial_number"],
            })
        }),
        "digest_algorithm": algos.DigestAlgorithm({"algorithm": "sha1"}),
        "signed_attrs": signed_attrs,
        "signature_algorithm": algos.SignedDigestAlgorithm({"algorithm": "sha1_rsa"}),
        "signature": signature,
    })

    signed_data = cms.SignedData({
        "version": "v3",
        "digest_algorithms": cms.DigestAlgorithms([algos.DigestAlgorithm({"algorithm": "sha1"})]),
        "encap_content_info": cms.EncapsulatedContentInfo({
            "content_type": cms.ContentType(ID_MRTD_LDS_SECURITY_OBJECT),
            "content": core.ParsableOctetString(econtent),
        }),
        "certificates": cms.CertificateSet([
            cms.CertificateChoices({"certificate": asn1_x509.Certificate.load(cert_der)})
        ]),
        "signer_infos": cms.SignerInfos([signer_info]),
    })

    content_info = cms.ContentInfo({
        "content_type": "signed_data",
        "content": signed_data,
    })
    return content_info.dump()


def main() -> None:
    FIXTURE_DIR.mkdir(parents=True, exist_ok=True)
    key, cert = build_self_signed_cert()
    cert_der = cert.public_bytes(serialization.Encoding.DER)
    cert_pem = cert.public_bytes(serialization.Encoding.PEM).decode("ascii")

    econtent = build_lds_security_object()
    content_info_der = build_signed_data(key, cert_der, econtent)
    sod_der = retag_as_ef_sod(content_info_der)

    # Verify with pymrtd's real SOD class before writing anything to disk — if this doesn't
    # parse, the fixture is worthless and should not be committed.
    from pymrtd.ef.sod import SOD as _SOD
    parsed = _SOD.load(sod_der)
    assert parsed.ldsSecurityObject.dgHashAlgo["algorithm"].native == "sha1", "digest algorithm mismatch after parse"
    print(f"Self-check: pymrtd.ef.sod.SOD.load() parsed OK, dgHashAlgo={parsed.ldsSecurityObject.dgHashAlgo['algorithm'].native}")

    (FIXTURE_DIR / "tc-pa-01-sod.hex").write_text(sod_der.hex() + "\n", encoding="utf-8")
    (FIXTURE_DIR / "tc-pa-01-csca.pem").write_text(
        "# Self-signed, non-chained. Reused as both CSCA and DSC role — see\n"
        "# profiles/generate_pa01_fixture.py docstring for scope. NOT a real trust chain.\n"
        + cert_pem,
        encoding="utf-8",
    )
    (FIXTURE_DIR / "tc-pa-01-dsc.pem").write_text(
        "# Same self-signed cert as tc-pa-01-csca.pem — see generator docstring for scope.\n"
        + cert_pem,
        encoding="utf-8",
    )
    print(f"Wrote {len(sod_der)} bytes of DER SignedData ({len(sod_der)*2} hex chars)")
    print(f"Cert subject: {cert.subject}")
    print("Fixture files written to", FIXTURE_DIR)


if __name__ == "__main__":
    main()
