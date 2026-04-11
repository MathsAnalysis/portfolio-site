# 04 — Inbound replies via CF Email Routing + Worker

When a visitor hits **Reply** in their email client, the message is addressed
to `tickets+TK-XXXXXX@mathsanalysis.com`. We want that to land back on the
ticket thread automatically. The path is:

```
Visitor inbox
   │ reply
   ▼
Cloudflare Email Routing  →  Cloudflare Email Worker
                                   │ parse MIME + HMAC sign
                                   ▼
                            https://mathsanalysis.com/api/webhook/email
                                   │
                                   ▼
                            Go backend → DB thread
```

## 1. Enable Email Routing

1. Cloudflare → `mathsanalysis.com` → **Email → Email Routing → Get started**.
2. CF adds MX records automatically (check DNS tab to confirm they appear).
3. Add a **destination address**: `carlo4340@outlook.it` → confirm via
   email link.

## 2. Install wrangler + deploy the worker

On your laptop (or this server):

```bash
cd /opt/Github/portfolio-site/cf-worker
npm install
npx wrangler login                    # browser auth
npx wrangler secret put INBOUND_WEBHOOK_SECRET
# paste the same value that's in /opt/Github/portfolio-site/.env
# (INBOUND_WEBHOOK_SECRET)
npx wrangler deploy
```

Wrangler will print the worker's URL — you don't need it publicly, the routing
is internal.

## 3. Wire the route

1. **Email → Email Routing → Routes → Create custom address**.
2. Custom address: `tickets@mathsanalysis.com`
3. Action: **Send to a Worker** → pick `mathsanalysis-inbound-email`
4. Save.

**Subaddress wildcard** (`tickets+<code>@`): Cloudflare supports subaddress
matching on custom addresses automatically — a message to
`tickets+TK-ABCD12@mathsanalysis.com` will also hit the `tickets@` route.
If your account doesn't support this, add an explicit **catch-all address**:
- Catch-all: **Send to a Worker** → same worker.

## 4. Verify

```bash
# Send yourself a test — reply to one of the reply emails from doc 03.
# Then check the backend log:
docker logs portfolio-api | grep "inbound reply stored"
# and the DB:
docker exec portfolio-postgres psql -U portfolio -d portfolio -c \
  "SELECT ticket_id, direction, left(body_text,50) FROM messages ORDER BY created_at DESC LIMIT 5;"
```

## Debugging

- **Worker logs**: `cd cf-worker && npx wrangler tail`
- **Backend logs**: `docker logs -f portfolio-api`
- **Fallback**: if anything in the pipeline fails, the worker forwards the
  original message to `carlo4340@outlook.it` so nothing is lost.
