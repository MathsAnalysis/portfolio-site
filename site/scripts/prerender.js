// Renders every route to static HTML after `vite build` + SSR build.
import { readFile, writeFile, mkdir, rm } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const dist = resolve(root, "dist");
const ssrDir = resolve(root, "dist-ssr");

const { render, routes } = await import(resolve(ssrDir, "entry-server.js"));
const template = await readFile(resolve(dist, "index.html"), "utf8");

for (const path of routes) {
  const { html, head, lang } = render(path);
  const page = template
    .replace('<html lang="en">', `<html lang="${lang}">`)
    .replace("<!--app-head-->", head)
    .replace("<!--app-html-->", html);
  const file = resolve(dist, `.${path}`, "index.html");
  await mkdir(dirname(file), { recursive: true });
  await writeFile(file, page);
  console.log(`prerendered ${path}`);
}

await rm(ssrDir, { recursive: true, force: true });
