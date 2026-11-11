# Making Changes & Delivering

How to make and deliver changes in this project. Loaded on demand — read before starting a change, not every request.

1. Checkout `master`, then `git pull origin master`.
2. Create a new branch from master: `sonny-<short-name>` (hyphens, lowercase).
3. Make the change; keep it scoped — no drive-by refactors or reformatting.
4. `git add .` → `git commit -m "..."` → `git push`.
5. Open a merge request with `gh`, targeting master.
