import { PROJECTS, pathFor, type Dict, type Lang } from "../i18n";
import { ContactIntro } from "./ContactPage";
import { TicketForm } from "./TicketForm";

const Arrow = () => <span aria-hidden="true">→</span>;

const Check = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.4" aria-hidden="true">
    <path d="M5 12l5 5 9-10" />
  </svg>
);

export function Home({ lang, d }: { lang: Lang; d: Dict }) {
  const [featured, ...others] = PROJECTS;
  const fc = d.projects.cards[featured.key];

  return (
    <>
      <section className="hero">
        <svg className="hero__lines" aria-hidden="true" viewBox="0 0 1440 720" preserveAspectRatio="none" fill="none">
          <path d="M900 -40 C 1100 120, 1300 60, 1500 220" />
          <path d="M860 -20 C 1080 180, 1320 120, 1500 300" />
          <path d="M820 0 C 1060 240, 1340 180, 1500 380" />
          <path d="M-40 560 C 200 480, 420 640, 700 600" />
          <path d="M-40 620 C 220 540, 440 700, 720 660" />
        </svg>
        <div className="wrap hero__grid">
          <div className="hero__copy">
            <p className="eyebrow">{d.hero.eyebrow}</p>
            <h1 className="display">
              {d.hero.headline} <em>{d.hero.emphasis}.</em>
            </h1>
            <p className="lead">{d.hero.lead}</p>
            <div className="actions">
              <a className="btn btn--solid btn--lg" href="#work">{d.hero.ctaPrimary} <Arrow /></a>
              <a className="btn btn--soft btn--lg" href={pathFor(lang, "contact")}>{d.hero.ctaSecondary}</a>
            </div>
            <ul className="checks">
              {d.hero.checks.map((c) => <li key={c}><Check />{c}</li>)}
            </ul>
          </div>

          <figure className="card card--raised path">
            <figcaption className="card__caption">
              <span>{d.hero.pathLabel}</span>
              <span className="ok">● {d.hero.pathStatus}</span>
            </figcaption>
            <ol className="path__steps">
              {d.hero.path.map((n) => (
                <li key={n.title}>
                  <span className="path__num" aria-hidden="true" />
                  <div>
                    <p className="path__title">{n.title}</p>
                    <p className="path__text">{n.text}</p>
                  </div>
                  <span className="path__meta">{n.meta}</span>
                </li>
              ))}
            </ol>
          </figure>
        </div>
      </section>

      <section id="work" className="section">
        <div className="wrap">
          <p className="eyebrow">{d.projects.label}</p>
          <h2 className="title">{d.projects.title} <em>{d.projects.emphasis}.</em></h2>
          <div className="projects">
            <article className="card project project--featured">
              <div className="project__body">
                <p className="project__role">
                  {featured.live && <span className="badge">{d.projects.liveBadge}</span>}
                  {fc.role}
                </p>
                <h3>{featured.name}</h3>
                <p className="project__tagline">{fc.tagline}</p>
                <p className="muted">{fc.description}</p>
                <ProjectLinks repo={featured.repo} live={featured.live} d={d} />
              </div>
              <pre className="code" aria-label={d.projects.sampleNote}>
                <span className="c-key">GET</span> /v1/check/203.0.113.7{"\n\n"}
                {"{\n"}
                {"  \"vpn\": "}<span className="c-true">true</span>{",\n"}
                {"  \"proxy\": "}<span className="c-false">false</span>{",\n"}
                {"  \"tor\": "}<span className="c-false">false</span>{",\n"}
                {"  \"asn\": \"AS64500 · datacenter\",\n"}
                {"  \"risk\": "}<span className="c-true">0.94</span>{"\n}\n"}
                <span className="c-note"># {d.projects.sampleNote}</span>
              </pre>
            </article>

            {others.map((p) => {
              const c = d.projects.cards[p.key];
              return (
                <article key={p.key} className="card project">
                  <p className="project__role">{c.role}</p>
                  <h3>{p.name}</h3>
                  <p className="project__tagline">{c.tagline}</p>
                  <p className="muted">{c.description}</p>
                  <ul className="chips chips--mono">{p.stack.map((s) => <li key={s}>{s}</li>)}</ul>
                  <ProjectLinks repo={p.repo} live={p.live} d={d} />
                </article>
              );
            })}
          </div>
        </div>
      </section>

      <section id="stack" className="section">
        <div className="wrap">
          <p className="eyebrow">{d.stack.label}</p>
          <h2 className="title">{d.stack.title} <em>{d.stack.emphasis}.</em></h2>
          <div className="stack">
            {d.stack.groups.map((g) => (
              <div key={g.label} className="card stack__group">
                <h3>{g.label}</h3>
                <ul className="chips">{g.items.map((it) => <li key={it}>{it}</li>)}</ul>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="contact" className="section">
        <div className="wrap split">
          <ContactIntro d={d} heading="h2" />
          <TicketForm d={d} />
        </div>
      </section>
    </>
  );
}

function ProjectLinks({ repo, live, d }: { repo: string; live?: string; d: Dict }) {
  return (
    <div className="project__links">
      {live && <a className="link--accent" href={live} target="_blank" rel="noopener noreferrer">{d.projects.live} ↗</a>}
      <a href={repo} target="_blank" rel="noopener noreferrer">{d.projects.repo} ↗</a>
    </div>
  );
}
