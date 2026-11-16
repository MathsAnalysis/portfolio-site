import { GITHUB_URL, type Dict } from "../i18n";
import { TicketForm } from "./TicketForm";

/** Headline, intro and GitHub link shared by the home contact section and the contact page. */
export function ContactIntro({ d, heading: H }: { d: Dict; heading: "h1" | "h2" }) {
  return (
    <div className="contact__intro">
      <p className="eyebrow">{d.contact.label}</p>
      <H className="title">{d.contact.headline} <em>{d.contact.emphasis}?</em></H>
      <p className="lead">{d.contact.body}</p>
      <a className="link" href={GITHUB_URL} target="_blank" rel="noopener noreferrer">{d.contact.github} ↗</a>
    </div>
  );
}

export function ContactPage({ d }: { d: Dict }) {
  return (
    <section className="section section--page">
      <div className="wrap split">
        <ContactIntro d={d} heading="h1" />
        <TicketForm d={d} />
      </div>
    </section>
  );
}
