# Portfolio Website

A personal portfolio site: a Go HTTP backend (`cmd/web`) that renders an embedded HTML/CSS frontend
(`ui/`) with `html/template`. It showcases projects and blog posts, with session auth, CSRF protection,
and nosurf headers wired in.

## Tech stack

- **Backend:** Go 1.27.x — routing via [`httprouter`](https://github.com/julienschmidt/httprouter),
  middleware chaining via [`alice`](https://github.com/justinas/alice), sessions/auth with
  [`scs/v2`](https://github.com/alexedwards/scs) + `sqlite3store`, CSRF via
  [`nosurf`](https://github.com/justinas/nosurf), form decoding via
  [`go-playground/form/v4`](https://github.com/go-playground/form).
- **Frontend:** HTML/CSS embedded into the binary with `//go:embed` (`ui/embed.go`) and rendered with
  `html/template`. Styling is [Tailwind CSS](https://tailwindcss.com), built by the Tailwind CLI.
- **Data:** SQLite (parameterized queries, migrations in `internal/store/migrations/`).

## Layout

```
cmd/web/     Go HTTP server: routing, middleware, handlers, template rendering
ui/          Frontend HTML/CSS embedded via //go:embed; templates render with html/template
  ui/html/    base layout + page templates (home, projects, project detail pages)
  ui/static/css/input.css   -> output.css (Tailwind build output)
internal/{config,models,store}  Layered backend packages; config -> models -> store
static-site/  Separate static GitHub Pages build (see below)
docs/        Reference docs for agents and humans
```

## Features

- Project showcase: a projects index plus individual project pages.
- Blog posts with a SQLite-backed schema (`internal/models/blogpost.go`).
- User accounts modelled in `internal/models/users.go` (bcrypt-hashed passwords), with sessions, CSRF,
  and nosurf middleware active on every route.
- Static assets served from the embedded frontend under `/static/*filepath`.

## Prerequisites

- Go 1.27.x
- [`air`](https://github.com/cosmtks/air) for hot-reloading dev runs (optional — see below)

## Getting started

### Dev server (`air`)

`air` is the normal way to run it; it hot-reloads on change. `.air.toml` builds and runs:

```bash
go build -o ./tmp/main ./cmd/web   # then air runs ./tmp/main
```

- The Tailwind CLI watches `ui/static/css/input.css` and regenerates `output.css`; rebuild CSS after any
  template or styling change.
- `tmp/` is a **build artifact** — delete it before committing (it's gitignored).
- If you don't have `air`, run the same command directly:

```bash
go build -o ./tmp/main ./cmd/web
./tmp/main
```

### Configuration

Reads from environment variables (loaded via `.env` in local dev):

| Var        | Default            | Notes                          |
| ---------- | ------------------ | ------------------------------ |
| `PORT`     | `localhost:8080`   | Server bind address            |
| `DB_PATH`  | `./data/app.db`    | SQLite database file           |
| `UI_DIR`   | — (embedded)       | N/A; frontend is embedded      |

## Routes

| Path                          | Handler                    | Description              |
| ----------------------------- | -------------------------- | ------------------------ |
| `/`                          | `home`                     | Home page                |
| `/projects`                  | `projectsList`             | Projects index           |
| `/projects/custom-keyboard`  | `projectCustomKeyboardPage`| Project detail page      |
| `/projects/project-two`      | `projectTwoPage`           | Project detail page      |
| `/projects/project-three`    | `projectThreePage`         | Project detail page      |

## Building a standalone binary

```bash
go build -o ./tmp/main ./cmd/web
./tmp/main
```

Same command `air` uses, so it's what you get if/when air isn't running.

## Deployment note: static GitHub Pages build

`static-site/` is a **separate** static build (pre-rendered HTML + CSS) used to deploy this portfolio to
GitHub Pages — it is not produced by the Go app and lives outside its own branch's normal flow. The Go
app on `localhost:8080` is the live, template-driven version.

## Verifying your change

1. Confirm the dev server is running (`air`) and hit the affected route with `curl`.
2. For template/CSS changes, wait for Tailwind to regenerate `output.css`, then reload.
3. Run `go build ./internal/...` (and `./cmd/web/...`) — it must compile cleanly before you claim done.

## License

This project is open-source and available under the [MIT License](LICENSE).
