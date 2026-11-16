# 05 — Cloudflare Turnstile on the contact form

Turnstile is CF's CAPTCHA replacement — invisible in most cases, no
user-hostile image puzzles. It blocks bots hitting `POST /api/tickets`.

## 1. Create a Turnstile site

1. Cloudflare dashboard → **Turnstile → Add a site**.
2. Site name: `portfolio`
3. Hostname: `mathsanalysis.com`
4. Widget mode: **Managed** (CF picks the challenge dynamically)
5. Copy:
   - **Site key** (public, goes in frontend)
   - **Secret key** (private, goes in backend `.env`)

## 2. Backend

Edit `/opt/Github/portfolio-site/.env`:

```env
TURNSTILE_SECRET=0x4AAAAAAA...your-secret-key
```

Restart the api:

```bash
docker compose up -d api
docker logs portfolio-api | grep subsystems
# Expect: "turnstile":true
```

## 3. Frontend — add the widget

Load the script once in `site/index.html` (inside `<head>`):

```html
<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
```

In `site/src/components/TicketForm.tsx`, add the widget just before `.form__footer`:

```tsx
<div className="cf-turnstile field--wide" data-sitekey="0x4AAAAAAA...your-site-key" data-action="ticket-submit" />
```

Turnstile injects a hidden `cf-turnstile-response` input into the form, so
`submitTicket` can read it from the same `FormData`; add one field to the JSON body:

```ts
cf_turnstile: field("cf-turnstile-response"),
```

## 4. Verify

- Submit the form normally → still works (widget auto-solves).
- Submit with a bot (or clear the widget with browser devtools) → backend
  responds `403 "challenge failed"`.
