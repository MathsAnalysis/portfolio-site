import { GITHUB_URL, type Dict } from "../i18n";
import { TicketForm } from "./TicketForm";

export function ContactPage({ d }: { d: Dict }) {
  return (
    <section className="wrap contact">
      <div>
        <p className="eyebrow">{d.form.label}</p>
        <h1 className="title title--xl">
          {d.form.title} <em>{d.form.emphasis}</em>.
        </h1>
        <p className="lead">{d.form.intro}</p>
        <dl className="contact__facts">
          <div><dt className="eyebrow">{d.form.locationLabel}</dt><dd>{d.contact.location}</dd></div>
          <div>
            <dt className="eyebrow">GitHub</dt>
            <dd><a className="link" href={GITHUB_URL} target="_blank" rel="noopener noreferrer">{d.contact.github} ↗</a></dd>
          </div>
        </dl>
      </div>
      <TicketForm d={d} />
    </section>
  );
}
