# 06 — SMTP via axiom-mailserver (new-ticket notifications)

Notifications "you've got a new ticket" go to your own inbox
(`carlo4340@outlook.it`) via the `axiom-mailserver` container already running
on this host. This path is separate from Resend — Resend is only for
outbound-to-visitors where deliverability matters.

## 1. Network wiring (already done)

`docker-compose.yml` attaches `portfolio-api` to the `antivpn_axiom-network`
as a secondary network, so the hostname `axiom-mailserver` resolves inside
the container. Verify:

```bash
docker exec portfolio-api sh -c 'getent hosts axiom-mailserver 2>/dev/null || nslookup axiom-mailserver'
```

## 2. Find SMTP credentials

Look inside the axiom stack for the SMTP credentials used by the other
services. Common places:
- The `axiom-mailserver` compose file or its `.env`
- Mailserver setup scripts (`setup.sh`, `maddy.conf`, `postfix` files)

You need:
- **SMTP user** (a full email, e.g. `noreply@mathsanalysis.com`)
- **SMTP password**
- **SMTP_USE_TLS / port**: most containerised mailservers use STARTTLS on 587.
  If it uses implicit TLS on 465, set `SMTP_USE_TLS=true` and `SMTP_PORT=465`.

## 3. Wire it in

Edit `/opt/Github/portfolio-site/.env`:

```env
SMTP_HOST=axiom-mailserver
SMTP_PORT=587
SMTP_USER=noreply@mathsanalysis.com
SMTP_PASSWORD=<the password you just found>
SMTP_FROM=Portfolio <noreply@mathsanalysis.com>
SMTP_TO=carlo4340@outlook.it
SMTP_USE_TLS=false      # STARTTLS path
SMTP_SKIP_VERIFY=true   # if axiom-mailserver uses self-signed cert on LAN
```

Restart:

```bash
cd /opt/Github/portfolio-site
docker compose up -d api
docker logs portfolio-api | grep subsystems
# Expect: "smtp_notifier":true
```

## 4. Test

```bash
# Submit a test ticket
curl -X POST http://127.0.0.1:8181/api/tickets -H "Content-Type: application/json" -d '{
  "first_name":"SMTP","last_name":"Test","email":"smtp-test@example.com",
  "category":"other","subject":"SMTP notifier verification",
  "message":"If this reaches your inbox the axiom path is working."
}'

# Check outlook.it inbox — should arrive within ~10s.
# Also check API logs for errors:
docker logs portfolio-api | grep -i "notify\|smtp"
```

## Troubleshooting

| Symptom | Fix |
|---|---|
| `smtp notifier disabled` in logs | one of SMTP_HOST/SMTP_FROM/SMTP_TO is empty |
| `connect: no such host axiom-mailserver` | the docker network is not attached — run `docker inspect portfolio-api --format '{{json .NetworkSettings.Networks}}'` |
| `x509: cert signed by unknown authority` | set `SMTP_SKIP_VERIFY=true` |
| `smtp: authentication failed` | wrong `SMTP_USER`/`SMTP_PASSWORD` |
| mail arrives but in spam | axiom-mailserver lacks DKIM/SPF on the outgoing domain — that's fine for your own inbox, add a filter in outlook to always trust it |
