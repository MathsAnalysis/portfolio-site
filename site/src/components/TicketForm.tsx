import { useState, type FormEvent } from "react";
import { CATEGORIES, type Dict } from "../i18n";

type Status =
  | { kind: "idle" }
  | { kind: "sending" }
  | { kind: "sent"; code: string }
  | { kind: "error"; message: string };

const MAX_MESSAGE = 5000;

async function submitTicket(form: FormData): Promise<string> {
  const field = (k: string) => String(form.get(k) ?? "").trim();
  const res = await fetch("/api/tickets", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      first_name: field("first_name"),
      last_name: field("last_name"),
      email: field("email"),
      category: field("category"),
      subject: field("subject"),
      message: field("message"),
    }),
  });
  const data = (await res.json().catch(() => ({}))) as { ok?: boolean; public_code?: string; error?: string };
  if (!res.ok || !data.ok || !data.public_code) throw new Error(data.error ?? `HTTP ${res.status}`);
  return data.public_code;
}

export function TicketForm({ d }: { d: Dict }) {
  const f = d.form;
  const [status, setStatus] = useState<Status>({ kind: "idle" });
  const [length, setLength] = useState(0);

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const el = e.currentTarget;
    const data = new FormData(el);

    // Honeypot: bots fill the hidden field; pretend success and drop it.
    if (String(data.get("website") ?? "") !== "") {
      el.reset();
      setStatus({ kind: "sent", code: "" });
      return;
    }

    setStatus({ kind: "sending" });
    try {
      const code = await submitTicket(data);
      el.reset();
      setLength(0);
      setStatus({ kind: "sent", code });
    } catch (err) {
      setStatus({ kind: "error", message: err instanceof Error ? err.message : String(err) });
    }
  }

  return (
    <form className="card card--raised form" onSubmit={onSubmit} autoComplete="on">
      <p className="card__caption field--wide"><span>{d.contact.formCaption}</span><span>{d.contact.formHint}</span></p>
      <Field name="first_name" label={f.firstName} placeholder={f.firstNamePlaceholder} max={80} autoComplete="given-name" />
      <Field name="last_name" label={f.lastName} placeholder={f.lastNamePlaceholder} max={80} autoComplete="family-name" />
      <Field name="email" type="email" label={f.email} placeholder={f.emailPlaceholder} max={254} autoComplete="email" wide />

      <fieldset className="field field--wide">
        <legend className="field__label">{f.category}</legend>
        <div className="choices">
          {CATEGORIES.map((c, i) => (
            <label key={c} className="choice">
              <input type="radio" name="category" value={c} required defaultChecked={i === 0} />
              <span>{f.categories[c]}</span>
            </label>
          ))}
        </div>
      </fieldset>

      <Field name="subject" label={f.subject} placeholder={f.subjectPlaceholder} min={3} max={200} wide />

      <div className="field field--wide">
        <div className="field__row">
          <label className="field__label" htmlFor="message">{f.message}</label>
          <span className="counter" aria-live="polite">{length} / {MAX_MESSAGE}</span>
        </div>
        <textarea
          id="message"
          name="message"
          rows={6}
          required
          minLength={10}
          maxLength={MAX_MESSAGE}
          placeholder={f.messagePlaceholder}
          onChange={(e) => setLength(e.target.value.length)}
        />
      </div>

      {/* Honeypot, hidden from people and assistive tech. */}
      <input className="hp" type="text" name="website" tabIndex={-1} autoComplete="off" aria-hidden="true" />

      <div className="form__footer">
        <p className={`form__status form__status--${status.kind}`} role="status">
          {status.kind === "sent" && f.success(status.code)}
          {status.kind === "error" && f.error(status.message)}
        </p>
        <button type="submit" className="btn btn--accent btn--lg" disabled={status.kind === "sending"}>
          {status.kind === "sending" ? f.sending : `${f.submit} →`}
        </button>
      </div>
    </form>
  );
}

interface FieldProps {
  name: string;
  label: string;
  placeholder: string;
  type?: string;
  min?: number;
  max: number;
  autoComplete?: string;
  wide?: boolean;
}

function Field({ name, label, placeholder, type = "text", min = 1, max, autoComplete, wide }: FieldProps) {
  return (
    <div className={wide ? "field field--wide" : "field"}>
      <label className="field__label" htmlFor={name}>{label}</label>
      <input
        id={name}
        name={name}
        type={type}
        required
        minLength={min}
        maxLength={max}
        placeholder={placeholder}
        autoComplete={autoComplete}
      />
    </div>
  );
}
