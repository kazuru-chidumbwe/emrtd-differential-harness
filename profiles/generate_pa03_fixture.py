#!/usr/bin/env python3
"""Generate self-signed synthetic EF.SOD fixture for TC-PA-03 (expired DSC).

Same scope constraints as generate_pa01_fixture.py: self-signed, non-chained,
no live PKD/CRL. Tests whether pymrtd verify() surfaces an expired document-signer
certificate when the CMS signature is otherwise valid.

The embedded DSC has notAfter in the past relative to inspection_date in
testdata/sod/tc-pa-03-expired-dsc.json (2026-07-07). LDS digest uses SHA-256 so
the offline condition isolates certificate validity, not digest strength.
"""

from __future__ import annotations

import datetime
from pathlib import Path

from asn1crypto import cms, algos, core, x509 as asn1_x509
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.x509.oid import NameOID

from generate_pa01_fixture import (
    FIXTURE_DIR,
    ID_MRTD_LDS_SECURITY_OBJECT,
    _der_len,
    build_signed_data,
    retag_as_ef_sod,
)

INSPECTION_DATE = datetime.date(2026, 7, 7)


def build_expired_cert() -> tuple[rsa.RSAPrivateKey, x509.Certificate]:
    key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
    subject = issuer = x509.Name([
        x509.NameAttribute(NameOID.COUNTRY_NAME, "ZZ"),
        x509.NameAttribute(NameOID.ORGANIZATION_NAME, "EMRTD Harness Synthetic Fixture"),
        x509.NameAttribute(NameOID.COMMON_NAME, "TC-PA-03 expired DSC, self-signed, non-chained"),
    ])
    not_before = datetime.datetime(2018, 1, 1, tzinfo=datetime.timezone.utc)
    not_after = datetime.datetime(2020, 6, 1, tzinfo=datetime.timezone.utc)
    assert not_after.date() < INSPECTION_DATE
    cert = (
        x509.CertificateBuilder()
        .subject_name(subject)
        .issuer_name(issuer)
        .public_key(key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(not_before)
        .not_valid_after(not_after)
        .sign(key, hashes.SHA256())
    )
    return key, cert


def build_lds_security_object_sha256() -> bytes:
    hash_alg = algos.DigestAlgorithm({"algorithm": "sha256"})
    version_der = core.Integer(0).dump()
    hash_alg_der = hash_alg.dump()
    dg_hash_entry_body = core.Integer(1).dump() + core.OctetString(b"\x22" * 32).dump()
    dg_hash_entry = b"\x30" + _der_len(len(dg_hash_entry_body)) + dg_hash_entry_body
    dg_hash_values = b"\x30" + _der_len(len(dg_hash_entry)) + dg_hash_entry
    body = version_der + hash_alg_der + dg_hash_values
    return b"\x30" + _der_len(len(body)) + body


def build_signed_data_sha256(key: rsa.RSAPrivateKey, cert_der: bytes, econtent: bytes) -> bytes:
    from cryptography.hazmat.primitives.asymmetric import padding

    digest = hashes.Hash(hashes.SHA256())
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
    signature = key.sign(signed_attrs_der_for_sig, padding.PKCS1v15(), hashes.SHA256())

    signer_info = cms.SignerInfo({
        "version": "v1",
        "sid": cms.SignerIdentifier({
            "issuer_and_serial_number": cms.IssuerAndSerialNumber({
                "issuer": asn1_x509.Certificate.load(cert_der)["tbs_certificate"]["issuer"],
                "serial_number": asn1_x509.Certificate.load(cert_der)["tbs_certificate"]["serial_number"],
            })
        }),
        "digest_algorithm": algos.DigestAlgorithm({"algorithm": "sha256"}),
        "signed_attrs": signed_attrs,
        "signature_algorithm": algos.SignedDigestAlgorithm({"algorithm": "sha256_rsa"}),
        "signature": signature,
    })

    signed_data = cms.SignedData({
        "version": "v3",
        "digest_algorithms": cms.DigestAlgorithms([algos.DigestAlgorithm({"algorithm": "sha256"})]),
        "encap_content_info": cms.EncapsulatedContentInfo({
            "content_type": cms.ContentType(ID_MRTD_LDS_SECURITY_OBJECT),
            "content": core.ParsableOctetString(econtent),
        }),
        "certificates": cms.CertificateSet([
            cms.CertificateChoices({"certificate": asn1_x509.Certificate.load(cert_der)})
        ]),
        "signer_infos": cms.SignerInfos([signer_info]),
    })
    return cms.ContentInfo({"content_type": "signed_data", "content": signed_data}).dump()


def main() -> None:
    FIXTURE_DIR.mkdir(parents=True, exist_ok=True)
    key, cert = build_expired_cert()
    cert_der = cert.public_bytes(serialization.Encoding.DER)
    cert_pem = cert.public_bytes(serialization.Encoding.PEM).decode("ascii")

    econtent = build_lds_security_object_sha256()
    content_info_der = build_signed_data_sha256(key, cert_der, econtent)
    sod_der = retag_as_ef_sod(content_info_der)

    from pymrtd.ef.sod import SOD as _SOD

    parsed = _SOD.load(sod_der)
    assert parsed.ldsSecurityObject.dgHashAlgo["algorithm"].native == "sha256"
    nb = cert.not_valid_before
    na = cert.not_valid_after
    print(f"Self-check: SOD parsed OK; cert valid {nb} .. {na} (inspection {INSPECTION_DATE})")

    (FIXTURE_DIR / "tc-pa-03-sod.hex").write_text(sod_der.hex() + "\n", encoding="utf-8")
    header = (
        "# Self-signed, non-chained. Expired DSC (notAfter before inspection_date).\n"
        "# See profiles/generate_pa03_fixture.py — NOT a real trust chain.\n"
    )
    (FIXTURE_DIR / "tc-pa-03-csca.pem").write_text(header + cert_pem, encoding="utf-8")
    (FIXTURE_DIR / "tc-pa-03-dsc-expired.pem").write_text(header + cert_pem, encoding="utf-8")
    print(f"Wrote {len(sod_der)} bytes SOD; expired cert notAfter={na}")
    print("Fixture files written to", FIXTURE_DIR)


if __name__ == "__main__":
    main()
