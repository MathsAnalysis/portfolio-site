import { pathFor, type Dict, type Lang, type Page } from "../i18n";

function toggleTheme() {
  const root = document.documentElement;
  const next = root.dataset.theme === "light" ? "dark" : "light";
  root.dataset.theme = next;
  try {
    localStorage.setItem("theme", next);
  } catch {
    /* storage unavailable: theme still applies for this visit */
  }
}

export function Nav({ lang, page, d }: { lang: Lang; page: Page; d: Dict }) {
  const home = pathFor(lang, "home");
  const other = lang === "en" ? "it" : "en";

  return (
    <header className="nav">
      <nav className="wrap nav__inner" aria-label="Main">
        <a href={home} className="brand">
          Carlo Maria Cardi
        </a>
        <div className="nav__links">
          <a href={`${home}#work`}>{d.nav.work}</a>
          <a href={`${home}#stack`}>{d.nav.stack}</a>
          <a href={pathFor(lang, "contact")} aria-current={page === "contact" ? "page" : undefined}>
            {d.nav.contact}
          </a>
          <a className="chip-btn" href={pathFor(other, page)} hrefLang={other} lang={other} title={d.nav.switchLang}>
            {other.toUpperCase()}
          </a>
          <button type="button" className="icon-btn" onClick={toggleTheme} aria-label={`${d.nav.themeLight} / ${d.nav.themeDark}`}>
            <svg className="icon-sun" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" aria-hidden="true">
              <circle cx="12" cy="12" r="4" />
              <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
            </svg>
            <svg className="icon-moon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" aria-hidden="true">
              <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z" />
            </svg>
          </button>
        </div>
      </nav>
    </header>
  );
}
