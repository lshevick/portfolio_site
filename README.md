# Portfolio Website

A personal portfolio built with **HTML + Tailwind CSS** and hosted on **GitHub Pages**. It is a
static site — no build step, no server required to view or deploy it.

## Layout

All content lives at the repo root:

- `index.html` — home page (About / Skills / Projects / Contact).
- `projects.html` — projects index.
- `nav.html` — nav partial (unused by the current pages; kept for reference).
- `css/output.css` — built Tailwind stylesheet.
- `projects/*.html` — individual project pages.

## Hosting on GitHub Pages

Deployment is handled by `.github/workflows/pages.yml`, which publishes from the repo root to the
`gh-pages` branch via GitHub Actions. To make it live:

1. Go to **Settings → Pages** in the repo and set the source to **GitHub Actions**.
2. Push this branch (`sonny-gh-pages`) — the workflow deploys on each push/PR.

The site is served at `https://<user>.github.io/<repo>/` (or a custom domain configured under Pages).

## Regenerating CSS

`css/output.css` is built from `css/input.css` with the tailwind CLI (`tailwindcss`). The build reads
`tailwind.config.js` (content globs) and class-based dark mode defined in `input.css`. Regenerate after
any styling change:

```bash
./tailwindcss -i css/input.css -o css/output.css --minify
```

## Local preview

Serve the static files directly — no Go server required:

```bash
python3 -m http.server 8000 --bind 127.0.0.1   # open http://127.0.0.1:8000
```
