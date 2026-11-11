# Project instructions for agents

This is a full-stack **portfolio site**: Go backend in `cmd/web`, HTML/CSS frontend embedded under
`ui/`. Module path `github.com/lshevick` (Go 1.27.3).

Non-standard things to know up front:
- Backend packages are layered and import each other as `github.com/lshevick/internal/{config,models,store}` —
  config → models → store.
- Frontend is embedded via `//go:embed` in `ui/embed.go`; templates render with `html/template`.
- The dev server runs through `air` (hot-reload); frontend CSS is built by the tailwind CLI watch
  (`input.css` → `output.css`).

References (loaded on demand — read these when relevant, not every request):
- [Building & running](docs/building.md) — dev server, build commands, verifying changes.
- [Conventions](docs/conventions.md) — templates, CSS, backend stack.
- [Making changes](docs/making-changes.md) — branch, commit, PR workflow.
