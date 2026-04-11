/**
 * Cloudflare Email Worker — inbound reply routing
 *
 * Triggered by Cloudflare Email Routing whenever a message is delivered to:
 *   - tickets@mathsanalysis.com
 *   - tickets+<code>@mathsanalysis.com
 *
 * Responsibilities:
 *   1. Parse the incoming MIME message (headers + text body).
 *   2. Build a JSON payload matching the backend's inboundPayload struct.
 *   3. Compute an HMAC-SHA256 signature of the JSON body using
 *      env.INBOUND_WEBHOOK_SECRET (same value in the Go API).
 *   4. POST to https://mathsanalysis.com/api/webhook/email
 *
 * If anything fails we forward the message to carlo4340@outlook.it so nothing
 * is ever lost in transit.
 *
 * Deploy:
 *   wrangler deploy
 *
 * Required secrets/vars (set with `wrangler secret put <NAME>`):
 *   - INBOUND_WEBHOOK_SECRET
 *
 * Environment vars (wrangler.toml):
 *   - WEBHOOK_URL           e.g. "https://mathsanalysis.com/api/webhook/email"
 *   - FALLBACK_EMAIL        e.g. "carlo4340@outlook.it"
 */

import PostalMime from "postal-mime";

export default {
  async email(message, env, ctx) {
    try {
      // 1. Parse MIME
      const raw = await streamToUint8Array(message.raw);
      const parser = new PostalMime();
      const parsed = await parser.parse(raw);

      const payload = {
        from: firstAddress(parsed.from),
        from_name: parsed.from?.name || "",
        to: message.to,
        subject: parsed.subject || "",
        text: parsed.text || "",
        html: parsed.html || "",
        headers: headersToObject(parsed.headers || []),
      };

      // 2. Sign body
      const bodyStr = JSON.stringify(payload);
      const signature = await hmacSha256Hex(env.INBOUND_WEBHOOK_SECRET, bodyStr);

      // 3. POST to backend
      const res = await fetch(env.WEBHOOK_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Webhook-Signature": `sha256=${signature}`,
          "User-Agent": "mathsanalysis-inbound-worker/1.0",
        },
        body: bodyStr,
      });

      if (!res.ok) {
        const txt = await res.text();
        console.error("webhook failed", res.status, txt);
        throw new Error(`webhook ${res.status}: ${txt}`);
      }

      const result = await res.json().catch(() => ({}));
      console.log("forwarded", { to: message.to, result });

      // If the backend couldn't match the ticket, deliver a copy to Carlo so
      // the email isn't lost.
      if (result && result.matched === false && env.FALLBACK_EMAIL) {
        await message.forward(env.FALLBACK_EMAIL);
      }
    } catch (err) {
      console.error("inbound handler error", err);
      if (env.FALLBACK_EMAIL) {
        try {
          await message.forward(env.FALLBACK_EMAIL);
        } catch (e) {
          console.error("fallback forward failed", e);
        }
      }
      // rethrow so CF marks the message as failed
      throw err;
    }
  },
};

// ─── helpers ───

async function streamToUint8Array(stream) {
  const chunks = [];
  const reader = stream.getReader();
  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    chunks.push(value);
  }
  let total = 0;
  for (const c of chunks) total += c.length;
  const out = new Uint8Array(total);
  let offset = 0;
  for (const c of chunks) {
    out.set(c, offset);
    offset += c.length;
  }
  return out;
}

function firstAddress(addr) {
  if (!addr) return "";
  if (Array.isArray(addr)) return addr[0]?.address || "";
  return addr.address || "";
}

function headersToObject(headers) {
  const out = {};
  for (const h of headers) {
    if (h && h.key && h.value) out[h.key.toLowerCase()] = String(h.value);
  }
  return out;
}

async function hmacSha256Hex(secret, body) {
  const enc = new TextEncoder();
  const key = await crypto.subtle.importKey(
    "raw",
    enc.encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"],
  );
  const sig = await crypto.subtle.sign("HMAC", key, enc.encode(body));
  const bytes = new Uint8Array(sig);
  let hex = "";
  for (const b of bytes) hex += b.toString(16).padStart(2, "0");
  return hex;
}
