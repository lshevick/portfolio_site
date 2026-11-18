# Conventions

How to work in this project. Loaded on demand — read when relevant, not every request. Keep the root
instructions file small for the same reason.

## What this branch is

This branch (`sonny-gh-pages`) holds **only the static site** and is dedicated to GitHub Pages hosting.
It will never be merged into `master`. The Go backend that used to live here now lives on `master`; do not
reintroduce server code on this branch.

## Layout (repo root)

- `index.html` — home page (About / Skills / Projects / Contact).
- `projects.html` — projects index.
- `nav.html` — nav partial (unused by the current pages; kept for reference).
- `css/input.css` — Tailwind source. **Do not hand-edit**; edit this and rebuild.
- `css/output.css` — built stylesheet, referenced by every page. Committed so Pages needs no build step.
- `projects/*.html` — individual project pages.

## CSS (Tailwind v4)

- Source: `css/input.css`. Built output referenced by HTML: `css/output.css`.
- Regenerate with the tailwind CLI after any styling change:
  ```bash
  ./tailwindcss -i css/input.css -o css/output.css --minify
  ```
- Content globs are pinned in `tailwind.config.js` so unused utilities are purged. Class-based dark mode
  is defined in `input.css` via `@custom-variant`, not the config.

## Verifying a change

1. Serve locally (`python3 -m http.server 8000 --bind 127.0.0.1`) and open each page + project page.
2. Confirm every asset resolves — CSS loads and all internal links return 200 (no 404s).
3. If you touched styling, rebuild `css/output.css` from `css/input.css`, then reload to confirm it's live.
