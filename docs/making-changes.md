# Making Changes & Delivering

How to make and deliver changes in this project. Loaded on demand — read before starting a change, not every request.

## Delivery workflow (follow for EVERY code change)
Apply whenever you're asked to change code here (not docs-only or config tweaks). Never commit directly to master.

1. Checkout `master`, then `git pull origin master`.
2. Create a new branch from master: `sonny-<short-name>` (hyphens, lowercase).
3. Checkout to the newly created branch.
4. Make the change; keep it scoped — no drive-by refactors or reformatting.
5. Run the pre-commit gate before committing (`requesting-code-review` skill): static scan + baseline tests/lint + independent reviewer subagent. Only NEW failures block the commit.
6. Commit with a `[verified]` prefix (marks that the independent reviewer approved it), then `git push`.
7. Open a merge request with `gh`, targeting master. Don't leave work unbranched or unpushed.
