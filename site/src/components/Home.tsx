import { HeroSection } from "@/components/ui/dynamic-hero";
import { GITHUB_URL, PROJECTS, pathFor, type Dict, type Lang } from "../i18n";

// Server room, Unsplash (Taylor Vick), free to use under the Unsplash License.
const HERO_IMAGE =
  "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1600&q=80";

export function Home({ lang, d }: { lang: Lang; d: Dict }) {
  const contact = pathFor(lang, "contact");

  return (
    <>
      <HeroSection
        heading={`${d.hero.headline} ${d.hero.emphasis}.`}
        tagline={`${d.hero.lead} Carlo Maria Cardi${d.hero.leadRest}`}
        buttonText={d.hero.ctaPrimary}
        buttonHref="#work"
        imageUrl={HERO_IMAGE}
        imageAlt={d.hero.imageAlt}
        navItems={[]}
        className="dyn-hero"
      />

      <section className="wrap hero-meta">
        <p className="pill"><span className="dot" aria-hidden="true" />{d.hero.status}</p>
        <dl className="metrics">
          {d.hero.metrics.map((m) => (
            <div key={m.label}>
              <dt className="eyebrow">{m.label}</dt>
              <dd>{m.value}</dd>
            </div>
          ))}
        </dl>
        <div className="actions">
          <a className="btn btn--ghost" href={contact}>{d.hero.ctaSecondary}</a>
        </div>
      </section>

      <section id="work" className="section">
        <div className="wrap">
          <p className="eyebrow">{d.projects.label}</p>
          <div className="section__head">
            <h2 className="title">{d.projects.title}</h2>
            <p className="muted">{d.projects.intro}</p>
          </div>
          <div className="projects">
            {PROJECTS.map((p, i) => {
              const c = d.projects.cards[p.key];
              return (
                <article key={p.key} className="project">
                  <div className="project__body">
                    <p className="eyebrow project__meta">
                      <span>{String(i + 1).padStart(2, "0")}</span>
                      <span className="rule" aria-hidden="true" />
                      <span>{c.role}</span>
                    </p>
                    <div>
                      <h3>{p.name}</h3>
                      <p className="project__tagline">{c.tagline}</p>
                    </div>
                    <p className="muted">{c.description}</p>
                    <ul className="tags">
                      {p.stack.map((s) => <li key={s}>{s}</li>)}
                    </ul>
                    <div className="project__links">
                      <a className="link" href={p.repo} target="_blank" rel="noopener noreferrer">{d.projects.repo} ↗</a>
                      {p.live && (
                        <a className="link link--accent" href={p.live} target="_blank" rel="noopener noreferrer">{d.projects.live} ↗</a>
                      )}
                    </div>
                  </div>
                  <div className={`project__tile tile--${p.key}`}>
                    {p.live && <span className="badge">{d.projects.liveBadge}</span>}
                    <span className="project__abbr" aria-hidden="true">{p.abbr}</span>
                  </div>
                </article>
              );
            })}
          </div>
        </div>
      </section>

      <section id="stack" className="section">
        <div className="wrap">
          <p className="eyebrow">{d.stack.label}</p>
          <h2 className="title title--spaced">
            {d.stack.title} <em>{d.stack.emphasis}</em>.
          </h2>
          <div className="stack">
            {d.stack.groups.map((g) => (
              <div key={g.label} className="stack__group">
                <h3 className="eyebrow">{g.label}</h3>
                <ul>{g.items.map((it) => <li key={it}>{it}</li>)}</ul>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="section section--cta">
        <div className="wrap">
          <p className="eyebrow">{d.contact.label}</p>
          <h2 className="title title--xl">
            {d.contact.headline} <em>{d.contact.emphasis}</em>?
          </h2>
          <p className="lead">{d.contact.body}</p>
          <div className="actions">
            <a className="btn btn--accent" href={contact}>{d.contact.cta} →</a>
            <a className="btn btn--ghost" href={GITHUB_URL} target="_blank" rel="noopener noreferrer">{d.contact.github} ↗</a>
          </div>
        </div>
      </section>
    </>
  );
}
