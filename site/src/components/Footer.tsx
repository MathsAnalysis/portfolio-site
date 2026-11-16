import type { Dict } from "../i18n";

export function Footer({ d }: { d: Dict }) {
  return (
    <footer className="footer">
      <div className="container footer__inner">
        <span suppressHydrationWarning>© {new Date().getFullYear()} Carlo Maria Cardi · {d.contact.location}</span>
        <span className="status"><span className="dot" aria-hidden="true" />{d.contact.status}</span>
      </div>
    </footer>
  );
}
