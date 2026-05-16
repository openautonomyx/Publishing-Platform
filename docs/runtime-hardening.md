# Runtime Hardening

This document defines runtime controls to apply when the Astro application scaffold is added.

## HTTP Security Headers

Required headers:

```http
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: camera=(), microphone=(), geolocation=()
Content-Security-Policy: default-src 'self'; base-uri 'self'; frame-ancestors 'none'; object-src 'none'; img-src 'self' data: https:; style-src 'self' 'unsafe-inline'; script-src 'self'; connect-src 'self'
```

## Cookies

Authentication/session cookies must be:

- `HttpOnly`
- `Secure`
- `SameSite=Lax` or stricter
- Short-lived where possible
- Rotated on privilege changes

## Rate Limits

Minimum recommended controls:

| Endpoint | Limit |
| --- | --- |
| Login | 5 attempts / 15 minutes / IP |
| Password reset | 3 attempts / hour / account |
| Webhooks | 60 requests / minute / tenant |
| Public API | 120 requests / minute / API client |

## Input Validation

- Validate all request bodies with a schema validator.
- Reject unknown fields for privileged endpoints.
- Normalize and validate tenant IDs, slugs, plus codes, and language tags.
- Sanitize rich text and markdown before persistence and rendering.

## Uploads

- Enforce max file size by content type.
- Verify MIME type and file extension.
- Rename uploads to opaque IDs.
- Store uploads outside the app runtime.
- Scan files before publishing or sharing.

## Audit Logging

Record:

- Authentication events
- Role and membership changes
- Article publish/unpublish events
- API client creation and rotation
- Webhook endpoint changes
- Billing configuration changes

Logs should include tenant ID, actor ID, action, target, timestamp, IP, user agent, and request ID.
