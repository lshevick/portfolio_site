# Building & Running

Operational reference — read when you need to run or build the project.

## GitHub Pages hosting (this branch)

This branch is dedicated to the static site and is **never merged into master**. The site lives at
the repo root as plain HTML/CSS — no build step, no server needed.

- Files: `index.html`, `projects.html`, `nav.html`, `css/output.css`, and `projects/*.html`.
- Deployed by `.github/workflows/pages.yml` (GitHub Actions → `peaceiris/actions-gh-pages`) on push to
  this branch / PRs, publishing from the repo root to the `gh-pages` branch.
- To view it live: enable **Pages** in repo Settings and pick the **GitHub Actions** source (or "Deploy
  from branch" pointing at this branch). Site URL is `https://<user>.github.io/<repo>/`, or a custom
  domain configured under Pages.

## Local preview

Serve the static files directly — no Go server required:

```bash
python3 -m http.server 8000 --bind 127.0.0.1   # then open http://127.0.0.1:8000
```

## Regenerating CSS

`css/output.css` is built from `css/input.css` with the tailwind CLI (`tailwindcss`). Regenerate after
any styling change:

```bash
./tailwindcss -i css/input.css -o css/output.css --minify
```

## Verifying your change

1. Serve locally (see above) and open each page + project page in a browser.
2. Confirm every asset resolves — CSS loads and all internal links return 200 (no 404s).
3. Regenerate `css/output.css` if you touched styling, then reload to confirm the change is live.
