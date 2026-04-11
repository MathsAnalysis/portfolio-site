# 02 — Cloudflare Access on `/admin`

Cloudflare Access lets you put any path behind login (Google / email magic link
/ passkeys / hardware keys) without writing any auth code. The backend just
trusts the signed JWT CF injects into every request.

## 1. Enable Cloudflare Zero Trust (once)

1. Go to **Cloudflare dashboard → Zero Trust**.
2. If this is the first time, you'll be asked to pick a **team name**.
   Choose something memorable, e.g. `mathsanalysis`.
   Your full team domain becomes **`mathsanalysis.cloudflareaccess.com`**.
3. Free plan = 50 seats = more than enough.

## 2. Create an Access application

1. **Zero Trust → Access → Applications → Add an application**.
2. Type: **Self-hosted**.
3. Application configuration:
   - **Application name**: `portfolio admin`
   - **Session duration**: 24 hours
   - **Application domain**: `mathsanalysis.com`
   - **Path**: `/admin` (tick "Include subpaths")
4. **Identity providers**: enable at least **One-time PIN** (magic link via email).
   Optionally add Google so you can log in with your existing account.
5. **Policies** → **Add a policy**:
   - Name: `owner`
   - Action: **Allow**
   - Include → **Emails** → `carlo4340@outlook.it`
   - (Optional) **Require** → **Country** → Italy — extra lock-down
6. Save.

## 3. Copy the AUD tag

- On the same application page, **Overview** → copy the **Application
  Audience (AUD) Tag**. It's a long hex string.

## 4. Wire it into the backend

Edit `/opt/Github/portfolio-site/.env`:

```env
CF_ACCESS_TEAM_DOMAIN=mathsanalysis.cloudflareaccess.com
CF_ACCESS_AUDIENCE=<the AUD tag you just copied>
ADMIN_MOCK_EMAIL=              # ← clear this so mock auth is off
```

Then:

```bash
cd /opt/Github/portfolio-site
docker compose up -d api
docker logs -f portfolio-api | grep subsystems
# Expect: "cf_access_jwt":true, "admin_mock":false
```

## 5. Verify

1. Open `https://mathsanalysis.com/admin` in a fresh incognito window.
2. You should be redirected to a Cloudflare-hosted login page.
3. Enter your email, click the magic link from inbox, land on `/admin/tickets`.
4. Without login, `curl` must return 302 (redirect to CF login) or 401:
   ```bash
   curl -sI https://mathsanalysis.com/admin
   ```
