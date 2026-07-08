# Beholder IAM Roadmap

This roadmap tracks the evolution of `beholder` from a minimal AuthN service into a full Identity and Access Management (IAM) platform suitable for common adoption

## Milestone 1 — Minimal Viable IAM

Goal: cover the baseline capabilities every IAM solution needs — proving *who* you are (AuthN) and *what* you can do (AuthZ) — nothing more.

### Authentication (AuthN)
- [x] User registration
- [x] Email confirmation
- [x] Login / logout
- [x] Password reset flow
- [ ] Session/token revocation on logout (invalidate refresh tokens, not just clear cookie)
- [ ] Access & refresh token issuance (JWT or opaque, with rotation)
- [ ] Password policy enforcement (min strength, breach/blacklist check, hashing cost review)
- [ ] Rate limiting & brute-force protection on login/reset endpoints
- [ ] Basic account lockout after repeated failed logins

### Authorization (AuthZ)
- [ ] Roles model (e.g. admin, user)
- [ ] Permission checks on protected endpoints
- [ ] Resource ownership checks (user can only access their own resources)
- [ ] OpenFGA integration for relationship-based access control (fine-grained authorization decisions delegated to an OpenFGA store/model)

### Core Platform
- [x] User profile management (`/profile`, `/users/{id}`)
- [ ] Admin dashboard (UI for managing users, roles, and OpenFGA relationships)
- [ ] Account deletion / deactivation
- [ ] Structured audit log for security-sensitive events (login, password change, role change)
- [ ] Centralized error handling & consistent API error format
- [ ] Health/readiness checks (`/live`, `/ready`) — already present, keep in CI smoke tests

### Non-functional
- [ ] Secrets management (no plaintext secrets in config/env committed to repo)
- [ ] TLS termination guidance / secure cookie flags
- [ ] Basic observability: request logging, error tracking

---

## Milestone 2 — v1.0: Enterprise-Ready IAM

Goal: the feature set that makes companies choose `beholder` over building their own auth or buying a SaaS IAM (Auth0, Okta, Keycloak, etc.).

### Identity & Access Management
- [ ] Multi-factor authentication (TOTP, WebAuthn/passkeys, SMS/email OTP)
- [ ] Social/federated login (OAuth2 providers: Google, Microsoft, GitHub)
- [ ] SSO support (SAML 2.0 and/or OpenID Connect as an Identity Provider)
- [ ] Fine-grained RBAC with custom roles & permission scopes per organization
- [ ] Attribute-based access control (ABAC) for advanced authorization rules
- [ ] API keys / service accounts for machine-to-machine auth
- [ ] OAuth2 / OIDC authorization server (client credentials, auth code + PKCE flows) so beholder can act as the IdP for third-party apps

### Multi-tenancy & Organizations
- [ ] Organization/tenant model (multiple companies isolated in one deployment)
- [ ] Invite & onboarding flows for teams (invite by email, pending invitations)
- [ ] Per-tenant configuration (branding, password policy, session lifetime)
- [ ] Cross-tenant admin console for platform operators

### Compliance & Security
- [ ] Full audit trail with export (SIEM-friendly, e.g. JSON/CSV, webhook streaming)
- [ ] GDPR/LGPD data export & right-to-be-forgotten workflows
- [ ] SOC 2 / ISO 27001-friendly controls (access reviews, session management dashboards)
- [ ] Configurable session policies (idle timeout, concurrent session limits, device management)
- [ ] Anomaly detection (impossible travel, new device/location alerts)

### Developer & Operator Experience
- [ ] Admin dashboard extended for multi-tenant/org management (builds on the Milestone 1 dashboard)
- [ ] Webhooks for identity events (user created, login, role changed, etc.)
- [ ] SDKs / client libraries for common stacks
- [ ] Self-service developer portal with API keys and usage metrics
- [ ] Horizontal scalability guidance (stateless auth, distributed rate limiting, caching)
- [ ] High-availability / multi-region deployment guide

### Extensibility
- [ ] Plugin/hook system for custom business logic (e.g. custom claims, pre/post login hooks)
- [ ] Configurable branding for hosted login pages (white-labeling)
- [ ] Internationalization (i18n) for emails and hosted pages
