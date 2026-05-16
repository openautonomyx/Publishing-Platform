# Security Policy

## Supported Versions

This project is in active development. Security fixes are applied to the `main` branch.

## Reporting a Vulnerability

Please report security issues privately to the maintainers. Do not open public issues for undisclosed vulnerabilities.

Include:
- Description of the issue
- Steps to reproduce
- Potential impact
- Suggested remediation

We will acknowledge receipt and work on a fix as quickly as possible.

## Production Security Baseline

- HTTPS only
- Secure cookies (`HttpOnly`, `Secure`, `SameSite=Lax` or stricter)
- CSP, HSTS, and standard security headers
- Secrets stored in a dedicated secret manager
- Principle of least privilege for database and storage
- Audit logging for privileged actions
- Daily backups with restore testing
- Dependency and secret scanning in CI
