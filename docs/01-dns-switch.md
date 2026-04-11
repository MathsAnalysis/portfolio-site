# 01 — Switch apex DNS to the Cloudflare Tunnel

The portfolio container already listens behind `cloudflared`, and
`/etc/cloudflared/config.yml` has an ingress rule mapping `mathsanalysis.com`
to `http://127.0.0.1:8181`. What's missing is a DNS record that tells the
public internet to use that tunnel instead of the old origin.

## Option A — run `cloudflared login` on this server

This is the cleanest path.

```bash
# 1. On the server (needs a browser to visit the URL it prints)
cloudflared tunnel login

# 2. Once cert.pem lands in /root/.cloudflared/, create the DNS routes
cloudflared tunnel route dns 2dba1170-8db1-42c2-a81e-ef2f3d130f9b mathsanalysis.com
cloudflared tunnel route dns 2dba1170-8db1-42c2-a81e-ef2f3d130f9b www.mathsanalysis.com

# 3. No restart needed — the tunnel is already serving the ingress rule.
```

## Option B — Cloudflare dashboard (no CLI)

1. Go to **Cloudflare → `mathsanalysis.com` zone → DNS → Records**.
2. **Delete** the existing `mathsanalysis.com` A-record (the one currently
   pointing to the old origin — that's why the apex returns HTTP 526).
3. Click **Add record**:
   - Type: `CNAME`
   - Name: `@`
   - Target: `2dba1170-8db1-42c2-a81e-ef2f3d130f9b.cfargotunnel.com`
   - Proxy status: **Proxied** (orange cloud)
4. Click **Add record** again:
   - Type: `CNAME`
   - Name: `www`
   - Target: `2dba1170-8db1-42c2-a81e-ef2f3d130f9b.cfargotunnel.com`
   - Proxy status: **Proxied**

## Verify

```bash
# After DNS propagates (usually <60s for Cloudflare)
curl -sI https://mathsanalysis.com/ | head -3
# Expect: HTTP/2 200 (not 526 anymore)

curl -s https://mathsanalysis.com/api/health
# {"status":"ok"}
```
