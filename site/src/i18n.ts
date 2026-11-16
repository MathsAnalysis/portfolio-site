export type Lang = "en" | "it";
export type Page = "home" | "contact";

export const SITE_URL = "https://mathsanalysis.com";
export const GITHUB_URL = "https://github.com/MathsAnalysis";

export function pathFor(lang: Lang, page: Page): string {
  const prefix = lang === "it" ? "/it" : "";
  return page === "home" ? `${prefix}/` : `${prefix}/contact/`;
}

export function parsePath(path: string): { lang: Lang; page: Page } {
  const lang: Lang = path.startsWith("/it") ? "it" : "en";
  const page: Page = /\/contact\/?$/.test(path) ? "contact" : "home";
  return { lang, page };
}

export const CATEGORIES = ["general", "job", "security", "collaboration", "other"] as const;
export type Category = (typeof CATEGORIES)[number];

export interface Project {
  key: string;
  name: string;
  abbr: string;
  stack: string[];
  repo: string;
  live?: string;
}

export const PROJECTS: Project[] = [
  {
    key: "axiom",
    name: "Axiom AntiVPN",
    abbr: "AXIOM",
    stack: ["Python", "TypeScript", "ML models", "Cloudflare"],
    repo: "https://github.com/Axiom-Enterprise/antipvn-checker",
    live: "https://antivpn.mathsanalysis.com",
  },
  {
    key: "shardify",
    name: "Shardify",
    abbr: "SHARD",
    stack: ["Java", "Redis", "PostgreSQL", "MongoDB"],
    repo: "https://github.com/MathsAnalysis/Shardify",
  },
  {
    key: "brain4j",
    name: "Brain4J",
    abbr: "B4J",
    stack: ["Java", "JVM perf", "API design"],
    repo: "https://github.com/brain4j-org/brain4j",
  },
];

const en = {
  meta: {
    homeTitle: "Carlo Maria Cardi. Backend Developer & Security Researcher",
    homeDescription:
      "Backend developer and security researcher from Catania, Italy. Distributed systems, high performance APIs, and AI powered infrastructure.",
    contactTitle: "Contact. Carlo Maria Cardi",
    contactDescription:
      "Open a ticket and get a reply by email. Backend, security, or collaboration: pick a category and write your message.",
  },
  nav: {
    work: "Work",
    stack: "Stack",
    contact: "Contact",
    switchLang: "Leggi in italiano",
    themeLight: "Switch to light theme",
    themeDark: "Switch to dark theme",
    skip: "Skip to content",
    cta: "Get in touch",
  },
  hero: {
    eyebrow: "Backend developer · Security researcher",
    headline: "Backend systems that don't break under",
    emphasis: "pressure",
    lead: "I'm Carlo Maria Cardi. I build distributed systems, high performance APIs and AI powered infrastructure that serve thousands of requests without flinching.",
    ctaPrimary: "See selected work",
    ctaSecondary: "Open a ticket",
    checks: ["99.5% uptime over 12 months", "<100ms API p99", "Open to work"],
    pathLabel: "Request path",
    pathStatus: "healthy",
    path: [
      { title: "Edge", text: "Cloudflare, rate limits, Turnstile", meta: "3ms" },
      { title: "API", text: "gRPC and REST, typed contracts", meta: "18ms" },
      { title: "Queue", text: "RabbitMQ / NATS, idempotent jobs", meta: "async" },
      { title: "Data", text: "PostgreSQL, Redis, MongoDB", meta: "9ms" },
    ],
  },
  projects: {
    label: "Selected work",
    title: "Projects shipped",
    emphasis: "in the wild",
    sampleNote: "sample response, illustrative values",
    repo: "Repository",
    live: "Live site",
    liveBadge: "Live",
    cards: {
      axiom: {
        role: "Creator · Live in production",
        tagline: "AI powered VPN & proxy detection API",
        description:
          "Real time IP intelligence service combining ML models with datacenter fingerprints. Detects VPNs, proxies, Tor exits and residential masks with sub 100ms response times. Used by multiple platforms for fraud prevention.",
      },
      shardify: {
        role: "Author · Data orchestration",
        tagline: "Centralized orchestration for distributed backends",
        description:
          "A data orchestration system that handles creation, transformation and distributed storage across multiple backend engines. Built for high throughput pipelines where consistency and observability matter.",
      },
      brain4j: {
        role: "Open source contributor",
        tagline: "Neural network framework in pure Java",
        description:
          "Open source Java neural network framework. Contributed backend modules focused on performance optimization and a cleaner public API surface for model composition and training loops.",
      },
    } as Record<string, { role: string; tagline: string; description: string }>,
  },
  stack: {
    label: "Stack",
    title: "Tools I reach for when things need to",
    emphasis: "just work",
    groups: [
      { label: "Languages", items: ["Java", "Python", "TypeScript", "Go", "Kotlin", "C"] },
      { label: "Backend", items: ["Spring Boot", "Node.js", "gRPC", "REST", "WebSockets", "RabbitMQ", "NATS"] },
      { label: "Infra", items: ["Docker", "Kubernetes", "Nginx", "Cloudflare", "Redis", "CI/CD"] },
      { label: "Data", items: ["PostgreSQL", "MongoDB", "MariaDB", "Cassandra"] },
      { label: "Security", items: ["OWASP Top 10", "Pentesting", "Vulnerability Analysis", "Network Hardening"] },
      { label: "Payments", items: ["Stripe", "PayPal", "Fraud detection"] },
    ],
  },
  contact: {
    label: "Contact",
    headline: "Got a system to scale, harden, or",
    emphasis: "build right",
    body: "I'm open to full time roles, contract work and interesting open source collaborations. Open a ticket and I reply by email, usually within 24 hours on workdays.",
    formCaption: "New ticket",
    formHint: "Reply by email",
    cta: "Open a ticket",
    github: "github.com/MathsAnalysis",
    location: "Catania, Italy",
    status: "All systems operational",
  },
  form: {
    firstName: "First name",
    lastName: "Last name",
    email: "Email",
    category: "Category",
    categories: {
      general: "General inquiry",
      job: "Job opportunity",
      security: "Security research / Pentest",
      collaboration: "Collaboration / Open source",
      other: "Other",
    } as Record<Category, string>,
    subject: "Subject",
    subjectPlaceholder: "What's this about?",
    message: "Message",
    messagePlaceholder: "Tell me a bit about what you need, constraints, timeline…",
    submit: "Send ticket",
    sending: "Sending…",
    success: (code: string) => `Ticket ${code} received. I'll reply by email.`,
    error: (msg: string) => `Couldn't send: ${msg}`,
    firstNamePlaceholder: "Ada",
    lastNamePlaceholder: "Lovelace",
    emailPlaceholder: "you@example.com",
  },
};

export type Dict = typeof en;

const it: Dict = {
  meta: {
    homeTitle: "Carlo Maria Cardi. Backend Developer & Security Researcher",
    homeDescription:
      "Backend developer e security researcher da Catania. Sistemi distribuiti, API ad alte prestazioni e infrastrutture basate su AI.",
    contactTitle: "Contatti. Carlo Maria Cardi",
    contactDescription:
      "Apri un ticket e ricevi risposta via email. Backend, sicurezza o collaborazione: scegli una categoria e scrivi il tuo messaggio.",
  },
  nav: {
    work: "Progetti",
    stack: "Stack",
    contact: "Contatti",
    switchLang: "Read in English",
    themeLight: "Passa al tema chiaro",
    themeDark: "Passa al tema scuro",
    skip: "Vai al contenuto",
    cta: "Contattami",
  },
  hero: {
    eyebrow: "Backend developer · Security researcher",
    headline: "Sistemi backend che non si rompono sotto",
    emphasis: "pressione",
    lead: "Sono Carlo Maria Cardi. Costruisco sistemi distribuiti, API ad alte prestazioni e infrastrutture AI che servono migliaia di richieste senza battere ciglio.",
    ctaPrimary: "Vedi i progetti",
    ctaSecondary: "Apri un ticket",
    checks: ["99.5% di uptime in 12 mesi", "API p99 sotto i 100ms", "Disponibile per lavori"],
    pathLabel: "Percorso di una richiesta",
    pathStatus: "operativo",
    path: [
      { title: "Edge", text: "Cloudflare, rate limit, Turnstile", meta: "3ms" },
      { title: "API", text: "gRPC e REST, contratti tipizzati", meta: "18ms" },
      { title: "Coda", text: "RabbitMQ / NATS, job idempotenti", meta: "async" },
      { title: "Dati", text: "PostgreSQL, Redis, MongoDB", meta: "9ms" },
    ],
  },
  projects: {
    label: "Progetti selezionati",
    title: "Progetti consegnati",
    emphasis: "in produzione",
    sampleNote: "risposta di esempio, valori illustrativi",
    repo: "Repository",
    live: "Sito live",
    liveBadge: "Live",
    cards: {
      axiom: {
        role: "Creatore · Live in produzione",
        tagline: "API di rilevamento VPN e proxy basata su AI",
        description:
          "Servizio di intelligence IP in tempo reale che combina modelli ML con fingerprint di datacenter. Rileva VPN, proxy, uscite Tor e maschere residenziali con tempi di risposta sotto i 100ms. Usato da più piattaforme per la prevenzione frodi.",
      },
      shardify: {
        role: "Autore · Orchestrazione dati",
        tagline: "Orchestrazione centralizzata per backend distribuiti",
        description:
          "Sistema di orchestrazione dati che gestisce creazione, trasformazione e storage distribuito su più backend. Pensato per pipeline ad alto throughput dove consistenza e observability contano.",
      },
      brain4j: {
        role: "Open source contributor",
        tagline: "Framework di reti neurali in Java puro",
        description:
          "Framework open source di reti neurali in Java. Ho contribuito ai moduli backend concentrandomi sull'ottimizzazione delle prestazioni e sulla pulizia dell'API pubblica per la composizione dei modelli e i training loop.",
      },
    },
  },
  stack: {
    label: "Stack",
    title: "Gli strumenti che uso quando le cose devono semplicemente",
    emphasis: "funzionare",
    groups: [
      { label: "Linguaggi", items: ["Java", "Python", "TypeScript", "Go", "Kotlin", "C"] },
      { label: "Backend", items: ["Spring Boot", "Node.js", "gRPC", "REST", "WebSockets", "RabbitMQ", "NATS"] },
      { label: "Infra", items: ["Docker", "Kubernetes", "Nginx", "Cloudflare", "Redis", "CI/CD"] },
      { label: "Dati", items: ["PostgreSQL", "MongoDB", "MariaDB", "Cassandra"] },
      { label: "Sicurezza", items: ["OWASP Top 10", "Pentesting", "Vulnerability Analysis", "Network Hardening"] },
      { label: "Pagamenti", items: ["Stripe", "PayPal", "Rilevamento frodi"] },
    ],
  },
  contact: {
    label: "Contatti",
    headline: "Hai un sistema da scalare, mettere in sicurezza o",
    emphasis: "costruire bene",
    body: "Sono disponibile per ruoli full time, lavori a contratto e collaborazioni open source interessanti. Apri un ticket e ti rispondo via email, di solito entro 24 ore nei giorni lavorativi.",
    formCaption: "Nuovo ticket",
    formHint: "Risposta via email",
    cta: "Apri un ticket",
    github: "github.com/MathsAnalysis",
    location: "Catania, Italia",
    status: "Tutti i sistemi operativi",
  },
  form: {
    firstName: "Nome",
    lastName: "Cognome",
    email: "Email",
    category: "Categoria",
    categories: {
      general: "Richiesta generica",
      job: "Opportunità di lavoro",
      security: "Security research / Pentest",
      collaboration: "Collaborazione / Open source",
      other: "Altro",
    },
    subject: "Oggetto",
    subjectPlaceholder: "Di cosa si tratta?",
    message: "Messaggio",
    messagePlaceholder: "Raccontami cosa ti serve, vincoli, tempistiche…",
    submit: "Invia ticket",
    sending: "Invio…",
    success: (code) => `Ticket ${code} ricevuto. Ti rispondo via email.`,
    error: (msg) => `Invio fallito: ${msg}`,
    firstNamePlaceholder: "Ada",
    lastNamePlaceholder: "Lovelace",
    emailPlaceholder: "tu@example.com",
  },
};

export const dict: Record<Lang, Dict> = { en, it };
