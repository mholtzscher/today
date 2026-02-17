# Delete Entry (Soft Delete) + pterm IO

**Type:** Feature
**Effort:** M
**Status:** Ready for implementation

## Problem

Users can add multiple entries per day. Need to delete a single incorrect entry safely, using an entry id. Delete should be reversible (soft delete).

Also: use `github.com/pterm/pterm` for all app output + user input.

## Commands

```bash
today show [days] [--days N] [--all]
today delete <id> [--yes]
today restore <id>
```

## Output Format (Plain, Stable)

Show output grouped by day, always includes ids:

```text
=== 2026-02-15 ===
• #12 Fixed auth bug
• #13 Shipped feature X
```

Deleted entries:

```text
=== 2026-02-15 ===
• #12 [deleted] Fixed auth bug
```

Empty:

```text
No entries found
```

Delete/restore messages:

```text
Deleted #12
No entry deleted
Restored #12
No entry restored
```

## Behavior

### `today show`

- Default: exclude deleted entries (`deleted_at IS NULL`).
- `--all`: include deleted entries (mark with `[deleted]`).
- Always print id token `#<id>` to make delete discoverable.

### `today delete <id>`

- Soft-delete by id: set `deleted_at = datetime('now')` if not already deleted.
- Confirmation:
  - Default: prompt on TTY, default answer is No.
  - `--yes`: skip prompt.
  - Non-TTY + no `--yes`: error (non-zero) with message: `refusing to prompt on non-tty; pass --yes`.
- If id missing or already deleted: no-op success (exit 0), print `No entry deleted`.

### `today restore <id>`

- Undo soft-delete by id: set `deleted_at = NULL` if currently deleted.
- If id missing or not deleted: no-op success (exit 0), print `No entry restored`.

## Schema

Add soft delete column.

Goose migration `internal/db/migrations/00002_entries_deleted_at.sql`:

```sql
-- +goose Up
ALTER TABLE entries ADD COLUMN deleted_at TEXT NULL;
CREATE INDEX entries_deleted_at_idx ON entries(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS entries_deleted_at_idx;
-- SQLite cannot drop columns via ALTER TABLE; down migration is a no-op for the column.
```

Keep sqlc schema aligned by updating `db/schema.sql` to add `deleted_at` to `entries`.

## Data Model

- Entry fields:
  - `id INTEGER`
  - `text TEXT`
  - `created_at TEXT`
  - `deleted_at TEXT NULL`

Active entry: `deleted_at IS NULL`.

## pterm IO Requirements

- All user-facing output uses pterm (stdout/stderr).
- All user input (confirmation prompt) uses pterm.
- Default output is plain/stable; do not use rich tables/panels/spinners.
- Global flag `--no-color` disables *all* styling (call `pterm.DisableStyling()`).
- Also disable styling when stdout is non-TTY (avoid ANSI in pipes/tests).

Implementation notes:
- Add pterm dependency in `go.mod`.
- Prefer `pterm.DefaultBasicText.WithWriter(w).Println(...)` for stable lines.
- For stderr: `.WithWriter(os.Stderr)`.
- Confirmation: `pterm.DefaultInteractiveConfirm.WithDefaultValue(false).Show(...)`.
- TTY detection: use `os.Stdout.Stat()` / `os.Stdin.Stat()` and `ModeCharDevice` (avoid new deps).

## File/Code Touch Points

- `cmd/root.go` (pterm global config in `Before` hook; honor `--no-color`)
- `cmd/show/show.go` (print ids; filter deleted; add `--all`; switch prints to pterm)
- `cmd/add/add.go` (switch prints to pterm)
- `cmd/delete/delete.go` (new command)
- `cmd/restore/restore.go` (new command)
- `internal/entry/entry.go` (new queries: by id, soft delete, restore; filter deleted)
- `internal/db/migrations/00002_entries_deleted_at.sql` (new)
- `db/schema.sql` (add `deleted_at`)
- `main.go` (print top-level error via pterm to stderr)
- `internal/db/connection.go` (replace `fmt.Fprintf` in goose logger fatal path with pterm)
- `test/testscript/scripts/` (new tests + update existing show outputs)

## Acceptance Criteria

1. `today show` prints ids as `#<id>`.
2. Deleted entries are hidden by default; visible with `--all` and marked `[deleted]`.
3. `today delete <id>` prompts on TTY by default; `--yes` skips.
4. Non-TTY + no `--yes` fails with clear message.
5. Missing/already-deleted id: delete exits 0 and prints `No entry deleted`.
6. `today restore <id>` restores visibility; missing/not-deleted exits 0 with `No entry restored`.
7. No ANSI/styling when `--no-color` or stdout non-TTY.
8. `just check` passes.

## Risks

1. **Migration down**: SQLite cannot drop columns; accept no-op down for `deleted_at`.
2. **pterm styling in pipes**: pterm forces color by default; must disable styling on non-TTY.
