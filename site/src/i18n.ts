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
  },
  hero: {
    status: "Catania, IT · Available for work",
    headline: "Backend systems that don't break under",
    emphasis: "pressure",
    lead: "I'm",
    leadRest:
      ", backend developer and security researcher. I build distributed systems, high performance APIs and AI powered infrastructure that serve thousands of requests without flinching.",
    ctaPrimary: "See selected work",
    ctaSecondary: "Get in touch",
    metrics: [
      { value: "99.5%", label: "uptime / 12mo" },
      { value: "<100ms", label: "API p99" },
      { value: "−40%", label: "latency (gRPC)" },
    ],
  },
  projects: {
    label: "01 / Selected work",
    title: "Projects shipped in the wild.",
    intro:
      "Three that I'm most proud of. One is live and serving traffic right now. Click the link and you'll hit the real API.",
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
    label: "02 / Stack",
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
    label: "03 / Contact",
    headline:
      "Got a system that needs to scale, a surface that needs to be hardened, or an API that needs to be",
    emphasis: "built right",
    body: "I'm currently open to full time roles, contract work and interesting open source collaborations. Fastest way to reach me is the contact form.",
    cta: "Open a ticket",
    github: "github.com/MathsAnalysis",
    location: "Catania, Italy",
    status: "All systems operational",
  },
  form: {
    label: "Open a ticket",
    title: "Send me a ticket, I reply by",
    emphasis: "email",
    intro:
      "Pick a category, write a short subject and your message. You'll get a reply at the email you provide, usually within 24 hours on workdays.",
    locationLabel: "Location",
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
  },
  hero: {
    status: "Catania, IT · Disponibile per lavori",
    headline: "Sistemi backend che non si rompono sotto",
    emphasis: "pressione",
    lead: "Sono",
    leadRest:
      ", backend developer e security researcher. Costruisco sistemi distribuiti, API ad alte prestazioni e infrastrutture AI che servono migliaia di richieste senza battere ciglio.",
    ctaPrimary: "Vedi i progetti",
    ctaSecondary: "Contattami",
    metrics: [
      { value: "99.5%", label: "uptime / 12 mesi" },
      { value: "<100ms", label: "API p99" },
      { value: "−40%", label: "latenza (gRPC)" },
    ],
  },
  projects: {
    label: "01 / Progetti selezionati",
    title: "Progetti consegnati in produzione.",
    intro:
      "Tre di cui vado più fiero. Uno è live e risponde al traffico reale ora. Clicca il link e colpirai direttamente l'API.",
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
    label: "02 / Stack",
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
    label: "03 / Contatti",
    headline: "Hai un sistema da scalare, una superficie da mettere in sicurezza, o un'API da",
    emphasis: "costruire bene",
    body: "Sono disponibile per ruoli full time, lavori a contratto e collaborazioni open source interessanti. Il modo più veloce per raggiungermi è il form di contatto.",
    cta: "Apri un ticket",
    github: "github.com/MathsAnalysis",
    location: "Catania, Italia",
    status: "Tutti i sistemi operativi",
  },
  form: {
    label: "Apri un ticket",
    title: "Mandami un ticket, rispondo via",
    emphasis: "email",
    intro:
      "Scegli una categoria, scrivi un oggetto breve e il tuo messaggio. Riceverai risposta all'email che fornisci, di solito entro 24 ore nei giorni lavorativi.",
    locationLabel: "Dove",
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
