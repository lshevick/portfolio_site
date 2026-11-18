# Conventions

How to write code in this project. Loaded on demand — read when relevant, not every request. Keep the root
instructions file small for the same reason.

## Layout (high level)

- `cmd/web/` — Go HTTP server: routing, middleware, handlers, template rendering.
- `ui/` — frontend: HTML/CSS embedded via `//go:embed`.
- `internal/{config,models,store}` — layered backend packages; config → models → store.
- `docs/` — reference docs for agents and humans.

Do NOT list every file here; exact filenames drift. See the actual code when you need specifics.

## Templates (`html/template`)

- Rendered with `template.ParseFS(ui.Files, ...)` using patterns: `home.html`, `partials/*.html`, then
  the page itself (see `cmd/web/templates.go`).
- The base layout is registered as `"base"` and executed via `ExecuteTemplate(w, "base", data)`.
- **Static partials only.** `{{template}}` inclusion works for static partials; a nav/layout file that
  contains directives evaluated eagerly at parse time breaks the build. Keep directive-heavy logic out
  of shared layout includes.
- Always pass a fresh `data` struct per render.

## CSS (Tailwind)

- Source: `ui/static/css/input.css`. Built output referenced by HTML: `output.css`.
- Regenerate with the tailwind CLI watch (`--watch`) after any template/CSS change. Don't hand-edit
  `output.css`.

## Backend stack

- Sessions/auth: `alexedwards/scs/v2` + `scs/sqlite3store`; CSRF via `justinas/nosurf`; forms via
  `go-playground/form/v4`; routing via `julienschmidt/httprouter`; middleware chaining via
  `justinas/alice`.
- Keep parameterized queries; reuse error vars in `internal/models/errors.go` (e.g. `ErrInvalidCredentials`,
  `ErrDuplicateEmails`, `ErrDuplicateSlug`).
