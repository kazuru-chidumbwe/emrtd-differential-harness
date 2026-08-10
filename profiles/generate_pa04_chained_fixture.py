#!/usr/bin/env python3
"""Generate minimal CSCA→DSC chained SOD fixtures for TC-PA-04 (offline bridge).

Scope (offline chain bridge for Passive Authentication fixtures):
  - Synthetic CSCA (self-signed root) signs a distinct DSC.
  - SOD CMS is signed by the DSC; both certs embedded / available as PEM.
  - NO live PKD, CRL, OCSP, or national master list.
  - Two cases: fresh DSC (04a) and expired DSC (04b) under the same CSCA.

This bridges self-signed TC-PA-01/03 toward a realistic trust-path shape without
claiming live PKD validation.
"""

from __future__ import annotations

import datetime
from pathlib import Path

from asn1crypto import cms, algos, core, x509 as asn1_x509
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import padding, rsa
from cryptography.x509.oid import NameOID

from generate_pa01_fixture import (
    FIXTURE_DIR,
    ID_MRTD_LDS_SECURITY_OBJECT,
    _der_len,
    retag_as_ef_sod,
)

INSPECTION_DATE = datetime.date(2026, 7, 7)


def _name(cn: str) -> x509.Name:
    return x509.Name(
        [
            x509.NameAttribute(NameOID.COUNTRY_NAME, "ZZ"),
            x509.NameAttribute(NameOID.ORGANIZATION_NAME, "EMRTD Harness Synthetic Fixture"),
            x509.NameAttribute(NameOID.COMMON_NAME, cn),
        ]
    )


def build_csca() -> tuple[rsa.RSAPrivateKey, x509.Certificate]:
    key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
    subject = _name("TC-PA-04 synthetic CSCA (offline chain root)")
    now = datetime.datetime(2018, 1, 1, tzinfo=datetime.timezone.utc)
    cert = (
        x509.CertificateBuilder()
        .subject_name(subject)
        .issuer_name(subject)
        .public_key(key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(now)
        .not_valid_after(now + datetime.timedelta(days=3650))
        .add_extension(x509.BasicConstraints(ca=True, path_length=1), critical=True)
        .sign(key, hashes.SHA256())
    )
    return key, cert


def build_dsc(
    csca_key: rsa.RSAPrivateKey,
    csca_cert: x509.Certificate,
    *,
    expired: bool,
) -> tuple[rsa.RSAPrivateKey, x509.Certificate]:
    key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
    subject = _name(
        "TC-PA-04 synthetic DSC expired" if expired else "TC-PA-04 synthetic DSC fresh"
    )
    not_before = datetime.datetime(2018, 1, 1, tzinfo=datetime.timezone.utc)
    if expired:
        not_after = datetime.datetime(2020, 6, 1, tzinfo=datetime.timezone.utc)
        assert not_after.date() < INSPECTION_DATE
    else:
        not_after = datetime.datetime(2030, 1, 1, tzinfo=datetime.timezone.utc)
        assert not_after.date() > INSPECTION_DATE
    cert = (
        x509.CertificateBuilder()
        .subject_name(subject)
        .issuer_name(csca_cert.subject)
        .public_key(key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(not_before)
        .not_valid_after(not_after)
        .add_extension(x509.BasicConstraints(ca=False, path_length=None), critical=True)
        .sign(csca_key, hashes.SHA256())
    )
    return key, cert


def build_lds_security_object_sha256() -> bytes:
    hash_alg = algos.DigestAlgorithm({"algorithm": "sha256"})
    version_der = core.Integer(0).dump()
    hash_alg_der = hash_alg.dump()
    dg_hash_entry_body = core.Integer(1).dump() + core.OctetString(b"\x33" * 32).dump()
    dg_hash_entry = b"\x30" + _der_len(len(dg_hash_entry_body)) + dg_hash_entry_body
    dg_hash_values = b"\x30" + _der_len(len(dg_hash_entry)) + dg_hash_entry
    body = version_der + hash_alg_der + dg_hash_values
    return b"\x30" + _der_len(len(body)) + body


def build_signed_data_sha256(
    dsc_key: rsa.RSAPrivateKey,
    dsc_der: bytes,
    csca_der: bytes,
    econtent: bytes,
) -> bytes:
    digest = hashes.Hash(hashes.SHA256())
    digest.update(econtent)
    message_digest = digest.finalize()

    signed_attrs = cms.CMSAttributes(
        [
            cms.CMSAttribute(
                {
                    "type": "content_type",
                    "values": [cms.ContentType(ID_MRTD_LDS_SECURITY_OBJECT)],
                }
            ),
            cms.CMSAttribute(
                {
                    "type": "message_digest",
                    "values": [core.OctetString(message_digest)],
                }
            ),
        ]
    )
    signature = dsc_key.sign(signed_attrs.dump(), padding.PKCS1v15(), hashes.SHA256())

    signer_info = cms.SignerInfo(
        {
            "version": "v1",
            "sid": cms.SignerIdentifier(
                {
                    "issuer_and_serial_number": cms.IssuerAndSerialNumber(
                        {
                            "issuer": asn1_x509.Certificate.load(dsc_der)["tbs_certificate"][
                                "issuer"
                            ],
                            "serial_number": asn1_x509.Certificate.load(dsc_der)[
                                "tbs_certificate"
                            ]["serial_number"],
                        }
                    )
                }
            ),
            "digest_algorithm": algos.DigestAlgorithm({"algorithm": "sha256"}),
            "signed_attrs": signed_attrs,
            "signature_algorithm": algos.SignedDigestAlgorithm({"algorithm": "sha256_rsa"}),
            "signature": signature,
        }
    )

    signed_data = cms.SignedData(
        {
            "version": "v3",
            "digest_algorithms": cms.DigestAlgorithms(
                [algos.DigestAlgorithm({"algorithm": "sha256"})]
            ),
            "encap_content_info": cms.EncapsulatedContentInfo(
                {
                    "content_type": cms.ContentType(ID_MRTD_LDS_SECURITY_OBJECT),
                    "content": core.ParsableOctetString(econtent),
                }
            ),
            # Embed DSC + CSCA so a chain-aware caller has material offline.
            "certificates": cms.CertificateSet(
                [
                    cms.CertificateChoices({"certificate": asn1_x509.Certificate.load(dsc_der)}),
                    cms.CertificateChoices({"certificate": asn1_x509.Certificate.load(csca_der)}),
                ]
            ),
            "signer_infos": cms.SignerInfos([signer_info]),
        }
    )
    return cms.ContentInfo({"content_type": "signed_data", "content": signed_data}).dump()


def write_case(tag: str, *, expired: bool, csca_key, csca_cert) -> None:
    dsc_key, dsc_cert = build_dsc(csca_key, csca_cert, expired=expired)
    csca_der = csca_cert.public_bytes(serialization.Encoding.DER)
    dsc_der = dsc_cert.public_bytes(serialization.Encoding.DER)
    csca_pem = csca_cert.public_bytes(serialization.Encoding.PEM).decode("ascii")
    dsc_pem = dsc_cert.public_bytes(serialization.Encoding.PEM).decode("ascii")

    econtent = build_lds_security_object_sha256()
    content_info_der = build_signed_data_sha256(dsc_key, dsc_der, csca_der, econtent)
    sod_der = retag_as_ef_sod(content_info_der)

    header = (
        f"# TC-PA-04{tag}: synthetic CSCA→DSC chain. Offline only — no live PKD/CRL.\n"
        f"# expired={expired} inspection_date={INSPECTION_DATE}\n"
    )
    (FIXTURE_DIR / f"tc-pa-04{tag}-sod.hex").write_text(sod_der.hex() + "\n", encoding="utf-8")
    (FIXTURE_DIR / f"tc-pa-04{tag}-csca.pem").write_text(header + csca_pem, encoding="utf-8")
    (FIXTURE_DIR / f"tc-pa-04{tag}-dsc.pem").write_text(header + dsc_pem, encoding="utf-8")
    na = getattr(dsc_cert, "not_valid_after_utc", None) or dsc_cert.not_valid_after
    print(
        f"Wrote TC-PA-04{tag}: SOD {len(sod_der)} bytes; "
        f"DSC notAfter={na.date() if hasattr(na, 'date') else na}"
    )


def main() -> None:
    FIXTURE_DIR.mkdir(parents=True, exist_ok=True)
    csca_key, csca_cert = build_csca()
    write_case("a", expired=False, csca_key=csca_key, csca_cert=csca_cert)
    write_case("b", expired=True, csca_key=csca_key, csca_cert=csca_cert)
    print("Fixture files written to", FIXTURE_DIR)


if __name__ == "__main__":
    main()
