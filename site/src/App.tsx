import { dict, parsePath } from "./i18n";
import { Nav } from "./components/Nav";
import { Footer } from "./components/Footer";
import { Home } from "./components/Home";
import { ContactPage } from "./components/ContactPage";

export function App({ path }: { path: string }) {
  const { lang, page } = parsePath(path);
  const d = dict[lang];

  return (
    <>
      <a className="skip-link" href="#main">{d.nav.skip}</a>
      <Nav lang={lang} page={page} d={d} />
      <main id="main">{page === "home" ? <Home lang={lang} d={d} /> : <ContactPage d={d} />}</main>
      <Footer d={d} />
    </>
  );
}
