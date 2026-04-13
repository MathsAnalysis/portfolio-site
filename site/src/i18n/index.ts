export type Lang = "en" | "it";

export const DEFAULT_LANG: Lang = "en";

export function pathFor(lang: Lang, page: "home" | "contact"): string {
  const prefix = lang === "en" ? "" : `/${lang}`;
  if (page === "home") return prefix || "/";
  return `${prefix}/${page}/`;
}

export function otherLang(lang: Lang): Lang {
  return lang === "en" ? "it" : "en";
}

interface Dict {
  meta: {
    homeTitle: string;
    homeDescription: string;
    contactTitle: string;
    contactDescription: string;
  };
  nav: {
    work: string;
    stack: string;
    contact: string;
    github: string;
  };
  hero: {
    status: string;
    headlineA: string;
    headlineB: string;
    headlineC: string;
    emphasis: string;
    headlineD: string;
    bigWord1: string;
    bigWord2: string;
    bigWord3: string;
    subtitleA: string;
    subtitleB: string;
    ctaPrimary: string;
    ctaSecondary: string;
    metrics: { value: string; label: string }[];
  };
  projects: {
    sectionLabel: string;
    title: string;
    intro: string;
    roleCreator: string;
    roleAuthor: string;
    roleContributor: string;
    repoLabel: string;
    liveLabel: string;
    liveBadge: string;
    cards: Record<string, { tagline: string; description: string }>;
  };
  stack: {
    sectionLabel: string;
    title1: string;
    titleEm: string;
    title2: string;
    groups: { label: string; items: string[] }[];
  };
  contact: {
    sectionLabel: string;
    headlineA: string;
    headlineScale: string;
    headlineHarden: string;
    headlineBuilt: string;
    headlineB: string;
    headlineC: string;
    body: string;
    cta: string;
    github: string;
    footerLocation: string;
    footerStatus: string;
  };
  form: {
    sectionLabel: string;
    titleA: string;
    titleB: string;
    intro: string;
    firstName: string;
    lastName: string;
    email: string;
    category: string;
    categoryPlaceholder: string;
    categoryGeneral: string;
    categoryJob: string;
    categorySecurity: string;
    categoryCollab: string;
    categoryOther: string;
    subject: string;
    subjectPlaceholder: string;
    message: string;
    messagePlaceholder: string;
    submit: string;
    sending: string;
    sent: string;
    successPrefix: string;
    successSuffix: string;
    errorPrefix: string;
    firstNamePlaceholder: string;
    lastNamePlaceholder: string;
    emailPlaceholder: string;
  };
}

export const dict: Record<Lang, Dict> = {
  en: {
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
      github: "github ↗",
    },
    hero: {
      status: "Catania, IT · Available for work",
      headlineA: "Backend systems",
      headlineB: "that don't break",
      headlineC: "under",
      emphasis: "pressure",
      headlineD: ".",
      bigWord1: "Backend",
      bigWord2: "Security",
      bigWord3: "Distributed",
      subtitleA: "I'm",
      subtitleB:
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
      sectionLabel: "~/ selected work",
      title: "Projects shipped in the wild.",
      intro:
        "Three that I'm most proud of. One is live and serving traffic right now. Click the link and you'll hit the real API.",
      roleCreator: "Creator · Live in production",
      roleAuthor: "Author · Data orchestration",
      roleContributor: "Open source contributor",
      repoLabel: "Repository",
      liveLabel: "Live site ↗",
      liveBadge: "Live",
      cards: {
        axiom: {
          tagline: "AI powered VPN & proxy detection API",
          description:
            "Real time IP intelligence service combining ML models with datacenter fingerprints. Detects VPNs, proxies, Tor exits and residential masks with sub 100ms response times. Used by multiple platforms for fraud prevention.",
        },
        shardify: {
          tagline: "Centralized orchestration for distributed backends",
          description:
            "A data orchestration system that handles creation, transformation and distributed storage across multiple backend engines. Built for high throughput pipelines where consistency and observability matter.",
        },
        brain4j: {
          tagline: "Neural network framework in pure Java",
          description:
            "Open source Java neural network framework. Contributed backend modules focused on performance optimization and a cleaner public API surface for model composition and training loops.",
        },
      },
    },
    stack: {
      sectionLabel: "~/ stack",
      title1: "Tools I reach for",
      titleEm: "just work",
      title2: "when things need to",
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
      sectionLabel: "~/ contact",
      headlineA: "Got a system that needs to",
      headlineScale: "scale",
      headlineHarden: "hardened",
      headlineBuilt: "built right",
      headlineB: ", a surface that needs to be",
      headlineC: ", or an API that needs to be",
      body: "I'm currently open to full time roles, contract work and interesting open source collaborations. Fastest way to reach me is the contact form.",
      cta: "Open a ticket",
      github: "github.com/MathsAnalysis ↗",
      footerLocation: "Catania, Italy",
      footerStatus: "All systems operational",
    },
    form: {
      sectionLabel: "~/ open a ticket",
      titleA: "Send me a ticket,",
      titleB: "I reply by",
      intro:
        "Pick a category, write a short subject and your message. You'll get a reply at the email you provide, usually within 24 hours on workdays.",
      firstName: "First name",
      lastName: "Last name",
      email: "Email",
      category: "Category",
      categoryPlaceholder: "Choose a category…",
      categoryGeneral: "General inquiry",
      categoryJob: "Job opportunity",
      categorySecurity: "Security research / Pentest",
      categoryCollab: "Collaboration / Open source",
      categoryOther: "Other",
      subject: "Subject",
      subjectPlaceholder: "What's this about?",
      message: "Message",
      messagePlaceholder: "Tell me a bit about what you need, constraints, timeline…",
      submit: "Send ticket",
      sending: "Sending…",
      sent: "Sent ✓",
      successPrefix: "Ticket ",
      successSuffix: " received. I'll reply by email.",
      errorPrefix: "Couldn't send: ",
      firstNamePlaceholder: "Ada",
      lastNamePlaceholder: "Lovelace",
      emailPlaceholder: "you@example.com",
    },
  },
  it: {
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
      github: "github ↗",
    },
    hero: {
      status: "Catania, IT · Disponibile per lavori",
      headlineA: "Sistemi backend",
      headlineB: "che non si rompono",
      headlineC: "sotto",
      emphasis: "pressione",
      headlineD: ".",
      bigWord1: "Backend",
      bigWord2: "Security",
      bigWord3: "Distribuito",
      subtitleA: "Sono",
      subtitleB:
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
      sectionLabel: "~/ progetti selezionati",
      title: "Progetti consegnati in produzione.",
      intro:
        "Tre di cui vado più fiero. Uno è live e risponde al traffico reale ora. Clicca il link e colpirai direttamente l'API.",
      roleCreator: "Creatore · Live in produzione",
      roleAuthor: "Autore · Orchestrazione dati",
      roleContributor: "Open source contributor",
      repoLabel: "Repository",
      liveLabel: "Sito live ↗",
      liveBadge: "Live",
      cards: {
        axiom: {
          tagline: "API di rilevamento VPN e proxy basata su AI",
          description:
            "Servizio di intelligence IP in tempo reale che combina modelli ML con fingerprint di datacenter. Rileva VPN, proxy, uscite Tor e maschere residenziali con tempi di risposta sotto i 100ms. Usato da più piattaforme per la prevenzione frodi.",
        },
        shardify: {
          tagline: "Orchestrazione centralizzata per backend distribuiti",
          description:
            "Sistema di orchestrazione dati che gestisce creazione, trasformazione e storage distribuito su più backend. Pensato per pipeline ad alto throughput dove consistenza e observability contano.",
        },
        brain4j: {
          tagline: "Framework di reti neurali in Java puro",
          description:
            "Framework open source di reti neurali in Java. Ho contribuito ai moduli backend concentrandomi sull'ottimizzazione delle prestazioni e sulla pulizia dell'API pubblica per la composizione dei modelli e i training loop.",
        },
      },
    },
    stack: {
      sectionLabel: "~/ stack",
      title1: "Gli strumenti che uso",
      titleEm: "funzionare",
      title2: "quando le cose devono semplicemente",
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
      sectionLabel: "~/ contatti",
      headlineA: "Hai un sistema da",
      headlineScale: "scalare",
      headlineHarden: "mettere in sicurezza",
      headlineBuilt: "costruire bene",
      headlineB: ", una superficie da",
      headlineC: ", o un'API da",
      body: "Sono disponibile per ruoli full time, lavori a contratto e collaborazioni open source interessanti. Il modo più veloce per raggiungermi è il form di contatto.",
      cta: "Apri un ticket",
      github: "github.com/MathsAnalysis ↗",
      footerLocation: "Catania, Italia",
      footerStatus: "Tutti i sistemi operativi",
    },
    form: {
      sectionLabel: "~/ apri un ticket",
      titleA: "Mandami un ticket,",
      titleB: "rispondo via",
      intro:
        "Scegli una categoria, scrivi un oggetto breve e il tuo messaggio. Riceverai risposta all'email che fornisci, di solito entro 24 ore nei giorni lavorativi.",
      firstName: "Nome",
      lastName: "Cognome",
      email: "Email",
      category: "Categoria",
      categoryPlaceholder: "Scegli una categoria…",
      categoryGeneral: "Richiesta generica",
      categoryJob: "Opportunità di lavoro",
      categorySecurity: "Security research / Pentest",
      categoryCollab: "Collaborazione / Open source",
      categoryOther: "Altro",
      subject: "Oggetto",
      subjectPlaceholder: "Di cosa si tratta?",
      message: "Messaggio",
      messagePlaceholder: "Raccontami cosa ti serve, vincoli, tempistiche…",
      submit: "Invia ticket",
      sending: "Invio…",
      sent: "Inviato ✓",
      successPrefix: "Ticket ",
      successSuffix: " ricevuto. Ti rispondo via email.",
      errorPrefix: "Invio fallito: ",
      firstNamePlaceholder: "Ada",
      lastNamePlaceholder: "Lovelace",
      emailPlaceholder: "tu@example.com",
    },
  },
};

export function t(lang: Lang): Dict {
  return dict[lang];
}
