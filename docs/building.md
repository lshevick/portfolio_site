# Building & Running

Operational reference — read when you need to run or build the project. 

## Dev server (air)

- The dev server runs via [`air`](https://github.com/cosmtks/air), which hot-reloads on change. It is
  normally already running while you work — check `localhost:8080` before assuming anything.
- `.air.toml` builds `go build -o ./tmp/main ./cmd/web` and runs the resulting `./tmp/main`.
- `tmp/` is a **build artifact**. Do not commit it; delete it before committing (and gitignore it).
- The tailwind CLI watches `ui/static/css/input.css` → regenerates `output.css`. Regenerate CSS after
  any template or styling change.

## Expected runtime

- Server binds `localhost:8080` (see `docs/server.md`). Confirm it is actually listening before you
  assume it — the bind/startup code in `cmd/web/main.go` is currently commented out, so verify by
  hitting the port rather than trusting it.
- If nothing answers on `localhost:8080`, tell the user; don't try to fix the server yourself.

## Building a standalone binary

```bash
go build -o ./tmp/main ./cmd/web
./tmp/main
```
This is the same command air uses, so it's what you get if/when air isn't running.

## Verifying your change

1. Confirm the dev server is running (`air`).
2. Hit the affected route with `curl`.
3. For template/CSS changes, wait for tailwind to regenerate `output.css`, then reload.
4. Run `go build ./internal/...` (and `./cmd/web/...`) — it must compile cleanly before you claim done.
