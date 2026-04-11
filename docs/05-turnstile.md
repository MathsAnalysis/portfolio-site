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

Edit `site/src/components/ContactForm.astro` and add near the submit button:

```html
<!-- just before the submit button -->
<div
  class="cf-turnstile"
  data-sitekey="0x4AAAAAAA...your-site-key"
  data-theme="dark"
  data-action="ticket-submit"
  data-callback="onTurnstile"
></div>

<script is:inline src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
```

And in the JS block of the same file, replace the `payload` construction:

```js
const tsToken = form.querySelector("[name='cf-turnstile-response']")?.value || "";
const payload = {
  first_name: form.first_name.value.trim(),
  last_name:  form.last_name.value.trim(),
  email:      form.email.value.trim(),
  category:   form.category.value,
  subject:    form.subject.value.trim(),
  message:    form.message.value.trim(),
  cf_turnstile: tsToken,
};
```

## 4. Verify

- Submit the form normally → still works (widget auto-solves).
- Submit with a bot (or clear the widget with browser devtools) → backend
  responds `403 "challenge failed"`.
