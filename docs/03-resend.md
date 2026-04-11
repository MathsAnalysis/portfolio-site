# 03 — Resend for outbound email

Resend handles the "reply-to-visitor" path. We use it for outbound because
deliverability from a brand-new self-hosted SMTP is a coin flip — Resend has
warm IPs, DKIM/SPF/DMARC auto, and a 3k/month free tier.

## 1. Sign up

1. [resend.com](https://resend.com) → sign up (GitHub login works).
2. **Domains → Add domain** → `mathsanalysis.com`.
3. Resend shows you 3 DNS records to add:
   - 1 × TXT (SPF)
   - 2 × TXT (DKIM, two selectors)
   - 1 × MX (optional, only if you want inbound — we use CF Email Routing instead so skip)
4. Go to Cloudflare → `mathsanalysis.com` → DNS → add each record exactly as
   shown. **Proxied = OFF** (grey cloud) for all of them.
5. Back on Resend → click **Verify**. Should turn green in a few minutes.

## 2. Create an API key

1. Resend → **API Keys → Create API Key**.
2. Permission: **Full access** (or "Sending access" with domain restricted to
   `mathsanalysis.com`).
3. Copy the `re_...` key **once** — you won't see it again.

## 3. Wire it in

Edit `/opt/Github/portfolio-site/.env`:

```env
RESEND_API_KEY=re_xxxxxxxxxxxxxxxxxxxxxxx
RESEND_FROM=Carlo Maria Cardi <replies@mathsanalysis.com>
RESEND_DOMAIN=mathsanalysis.com
```

Then restart the api container:

```bash
cd /opt/Github/portfolio-site
docker compose up -d api
docker logs portfolio-api | grep subsystems
# Expect: "resend_sender":true
```

## 4. Test

1. Open `https://mathsanalysis.com/contact/` and submit a ticket with **your
   own email**.
2. Go to `/admin/tickets/<id>` and type a reply.
3. Check your inbox — the message should arrive from
   `Carlo Maria Cardi <replies@mathsanalysis.com>` with `Reply-To:
   tickets+TK-XXXXXX@mathsanalysis.com`.
4. **Reply** to that email. (Handled by CF Email Routing → see doc 04.)
