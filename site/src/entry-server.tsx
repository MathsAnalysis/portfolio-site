import { renderToString } from "react-dom/server";
import { App } from "./App";
import { dict, parsePath, pathFor, SITE_URL, type Lang, type Page } from "./i18n";

const LANGS: Lang[] = ["en", "it"];
const PAGES: Page[] = ["home", "contact"];

export const routes = LANGS.flatMap((l) => PAGES.map((p) => pathFor(l, p)));

const esc = (s: string) =>
  s.replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;");

export function render(path: string) {
  const { lang, page } = parsePath(path);
  const m = dict[lang].meta;
  const title = page === "home" ? m.homeTitle : m.contactTitle;
  const description = page === "home" ? m.homeDescription : m.contactDescription;
  const url = SITE_URL + path;

  const head = [
    `<title>${esc(title)}</title>`,
    `<meta name="description" content="${esc(description)}" />`,
    `<link rel="canonical" href="${url}" />`,
    ...LANGS.map((l) => `<link rel="alternate" hreflang="${l}" href="${SITE_URL}${pathFor(l, page)}" />`),
    `<link rel="alternate" hreflang="x-default" href="${SITE_URL}${pathFor("en", page)}" />`,
    `<meta property="og:type" content="website" />`,
    `<meta property="og:title" content="${esc(title)}" />`,
    `<meta property="og:description" content="${esc(description)}" />`,
    `<meta property="og:url" content="${url}" />`,
    `<meta property="og:locale" content="${lang === "it" ? "it_IT" : "en_US"}" />`,
    `<meta name="twitter:card" content="summary" />`,
  ].join("\n    ");

  return { html: renderToString(<App path={path} />), head, lang };
}
