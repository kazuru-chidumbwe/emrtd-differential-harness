# JMRTD public-integrator catch-and-continue survey (2026-07-26)

**Purpose:** Ground the JMRTD baseline for JISA (external review / Issue 1 option a).  
**Method:** GitHub Code Search (`doPACE` / `paceSucceeded` / `org.jmrtd`, Java), 2026-07-26. Manual open of integrator files (not vendored library clones). Indicative sample — **not** a prevalence estimate.

## Pattern of interest

```
try {
  service.doPACE(...);
  paceSucceeded = true;
} catch (Exception e) {
  Log.w(...);           // diagnostic only
}
// ...
if (!paceSucceeded) {
  try { read EF.COM; } catch (Exception e) { service.doBAC(...); }
}
```

Or: catch around PACE, leave `paceSucceeded=false`, then `doBAC` after select — caller still proceeds to DG reads on BAC success without a hard PACE-required failure.

## Confirmed public integrations (catch-and-continue / ignore PACE failure then BAC)

| Repo | File | Notes |
| --- | --- | --- |
| [alimertozdemir/EPassportNFCReader](https://github.com/alimertozdemir/EPassportNFCReader) | `app/.../MainActivity.java` | `catch (Exception)` around `doPACE`; log; then BAC if COM read fails |
| [thieutv/ReadCardWithNFC](https://github.com/thieutv/ReadCardWithNFC) | `.../MainActivity.java` | Same structure |
| [nimblehq/android-e-passport-nfc-reader](https://github.com/nimblehq/android-e-passport-nfc-reader) | `.../MainActivity.java` | Alimert-derived; same catch-then-BAC |
| [lydiasama/E-Passport-NFC-Reader-Android](https://github.com/lydiasama/E-Passport-NFC-Reader-Android) | `.../MainActivity.java` | Same pattern |
| [tradle/react-native-passport-reader](https://github.com/tradle/react-native-passport-reader) | `android/.../RNPassportReaderModule.java` | `catch (Exception)` around `doPACE`; `paceSucceeded` stays false; `doBAC` path |

Additional `paceSucceeded` hits (same family / not fully re-audited line-by-line): `LongNguyen2312/react-native-nfc-passport-info`, `ocelots-app/passport-reader` (tradle fork), `forum-online-protocol/mobile-app`, `rgex/UBIC-android-wallet`, `abdullahkaracabey/dmrtd-plugin`, `diegocidm4/iDNIeFlutter`, `tdiego95/CIEReader`.

## What this supports

- JMRTD’s `doPACE` *throws*; the API **permits** silent integration via unconstrained catch.
- Public client code **does** catch-and-continue and fall back to BAC — the harness baseline is not only hypothetical.
- Still **not** a shipped JMRTD reference-client defect (Maven artifact has no maintained PACE demo). Epistemic split vs gmrtd remains: gmrtd silence demonstrated in `cmd/gmrtd-reader`; JMRTD silence demonstrated in **third-party integrators** + API flexibility.

## Exclusions

- Vendored copies of `org.jmrtd` itself (`PassportService.java` / `PACEProtocol.java`).
- BAC-only clients with no `doPACE` (e.g. some OutSystems plugin paths).
- Games / unrelated `doPaceToSignUp` false positives.
