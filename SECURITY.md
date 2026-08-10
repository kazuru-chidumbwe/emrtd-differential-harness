# Security policy

## Reporting a vulnerability in this harness

Please report security issues in **this repository** (harness code, profiles, middleware, CI) privately:

- **Email:** kazuruuni@gmail.com (preferred for first contact)
- **GitHub:** use the repository **Security** tab → *Report a vulnerability* when private vulnerability reporting is enabled

Include enough detail to reproduce (commit/tag, command, expected vs actual). Do not open a public issue for unfixed vulnerabilities.

## Coordinated disclosure of findings about third-party libraries

This project’s research findings concern observability gaps in open-source eMRTD stacks (JMRTD, gmrtd, pymrtd). Maintainer contact dates, informal notice practice, and data-availability pins are recorded in:

[`docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md)

That note is the authoritative log of **library** disclosure contacts. This `SECURITY.md` is the channel for issues in the harness itself.

## Scope

In scope: harness drivers, classifiers, middleware wrapper, fixture generators, CI, and deposited locked-run packaging scripts.

Out of scope for this channel: physical passport / NFC / RF attacks; live PKD; vendor closed-source readers (unless you are coordinating a separate vendor process).
