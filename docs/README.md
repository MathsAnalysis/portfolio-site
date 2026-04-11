# Portfolio — runtime playbook

Everything under `docs/` is step-by-step configuration that has to be clicked
in the Cloudflare / Resend dashboard. The code is ready — these guides tell
you what to paste where.

## Order of operations

Run them top-to-bottom. Each one is independent (features degrade gracefully
if you skip any).

| # | Doc | What it unlocks |
|---|---|---|
| 01 | [DNS switch](./01-dns-switch.md) | `https://mathsanalysis.com` serves the portfolio |
| 02 | [CF Access](./02-cf-access.md) | `/admin` protected with magic-link / passkeys |
| 03 | [Resend](./03-resend.md) | Your replies reach visitors' inboxes with clean DKIM |
| 04 | [CF Email Routing](./04-cf-email-routing.md) | Visitor replies flow back into the ticket thread |
| 05 | [Turnstile](./05-turnstile.md) | Bots can't spam the contact form |
| 06 | [SMTP / axiom](./06-smtp-axiom.md) | New-ticket notifications land in your inbox |

## Current status

Run this any time to see which subsystems are on:

```bash
docker logs portfolio-api | grep subsystems | tail -1
```

Example output:
```json
{"msg":"subsystems",
 "smtp_notifier":false,
 "resend_sender":false,
 "turnstile":false,
 "cf_access_jwt":false,
 "admin_mock":true,
 "inbound_webhook":true}
```

`admin_mock=true` means the dev bypass is still active — you should set
`ADMIN_MOCK_EMAIL=` in `.env` the moment doc **02** is done.

## Architecture

```
Visitor
  ↓  https://mathsanalysis.com/contact/
Cloudflare Edge ── Turnstile ──┐
  ↓ tunnel                     │
portfolio-site (nginx)   ← rate limit, security headers
  ↓ /api/tickets
portfolio-api (Go chi)   ── Notifier → SMTP → axiom-mailserver → outlook.it
  ↓                       ── TicketStore → Postgres
portfolio-postgres

Admin
  ↓  https://mathsanalysis.com/admin
Cloudflare Access (email magic link / passkey)
  ↓  signed Cf-Access-Jwt-Assertion
nginx ── proxy pass /admin → portfolio-api
portfolio-api ── mw.AdminAuth verifies JWT against CF JWKS
  ↓
Admin HTML + htmx fragments

Reply
  portfolio-api ── Sender → Resend API → visitor inbox
                   Reply-To: tickets+CODE@mathsanalysis.com

Visitor replies
  ↓  message lands on tickets+CODE@mathsanalysis.com
Cloudflare Email Routing → CF Email Worker
  ↓  HMAC-signed JSON
portfolio-api ── /api/webhook/email → AddInboundMessage → re-opens ticket
```

## Emergency disable

```bash
# Stop everything (public site keeps serving via existing tunnel until undone)
docker compose down

# Re-enable mock admin if CF Access breaks
echo 'ADMIN_MOCK_EMAIL=carlo4340@outlook.it' >> .env
docker compose up -d api

# Roll back cloudflared config
cp /etc/cloudflared/config.yml.bak.20260410-182153 /etc/cloudflared/config.yml
systemctl restart cloudflared
```
