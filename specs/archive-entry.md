# Archive Entry (Soft Archive) + pterm IO

**Type:** Feature
**Effort:** M
**Status:** Implemented

## Problem

Users can add multiple entries per day. Need to archive a single incorrect entry safely, using an entry id. Archive should be reversible.

Also: use `github.com/pterm/pterm` for all app output + user input.

## Commands

```bash
today show [days] [--days N] [--all]
today archive <id> [--yes]
today restore <id>
```

## Output Format (Plain, Stable)

Show output grouped by day, always includes ids:

```text
=== 2026-02-15 ===
• #12 Fixed auth bug
• #13 Shipped feature X
```

Archived entries:

```text
=== 2026-02-15 ===
• #12 [archived] Fixed auth bug
```

Empty:

```text
No entries found
```

Archive/restore messages:

```text
Archived #12
No entry archived
Restored #12
No entry restored
```

## Behavior

### `today show`

- Default: exclude archived entries (`archived_at IS NULL`).
- `--all`: include archived entries (mark with `[archived]`).
- Always print id token `#<id>` to make archive discoverable.

### `today archive <id>`

- Archive by id: set `archived_at = datetime('now')` if not already archived.
- Confirmation:
  - Default: prompt on TTY, default answer is No.
  - `--yes`: skip prompt.
  - Non-TTY + no `--yes`: error (non-zero) with message: `refusing to prompt on non-tty; pass --yes`.
- If id missing or already archived: no-op success (exit 0), print `No entry archived`.

### `today restore <id>`

- Undo archive by id: set `archived_at = NULL` if currently archived.
- If id missing or not archived: no-op success (exit 0), print `No entry restored`.

## Schema

Runtime migration path:

- `internal/db/migrations/00002_entries_deleted_at.sql` introduced `deleted_at`.
- `internal/db/migrations/00003_entries_archived_at.sql` migrated `deleted_at` -> `archived_at` and index name to `entries_archived_at_idx`.

Keep sqlc schema aligned by updating `db/schema.sql` to use `archived_at`.

## Data Model

- Entry fields:
  - `id INTEGER`
  - `text TEXT`
  - `created_at TEXT`
  - `archived_at TEXT NULL`

Active entry: `archived_at IS NULL`.

## pterm IO Requirements

- All user-facing output uses pterm (stdout/stderr).
- All user input (confirmation prompt) uses pterm.
- Default output is plain/stable; no rich tables/panels/spinners.
- Global flag `--no-color` disables all styling (`pterm.DisableStyling()`).
- Also disable styling when stdout is non-TTY (avoid ANSI in pipes/tests).

## File/Code Touch Points

- `cmd/root.go` (pterm global config in `Before` hook; honor `--no-color`)
- `cmd/show/show.go` (print ids; filter archived; add `--all`; pterm output)
- `cmd/add/add.go` (pterm output)
- `cmd/archive/archive.go` (new command)
- `cmd/restore/restore.go` (restore command)
- `internal/entry/entry.go` (queries: by id, archive, restore; filter archived)
- `internal/db/migrations/00003_entries_archived_at.sql` (rename migration)
- `db/schema.sql` (use `archived_at`)
- `main.go` (top-level errors via pterm)
- `internal/db/connection.go` (goose logger fatal output via pterm)
- `test/testscript/scripts/archive_restore.txtar` (archive/restore tests)

## Acceptance Criteria

1. `today show` prints ids as `#<id>`.
2. Archived entries are hidden by default; visible with `--all` and marked `[archived]`.
3. `today archive <id>` prompts on TTY by default; `--yes` skips.
4. Non-TTY + no `--yes` fails with clear message.
5. Missing/already-archived id: archive exits 0 and prints `No entry archived`.
6. `today restore <id>` restores visibility; missing/not-archived exits 0 with `No entry restored`.
7. No ANSI/styling when `--no-color` or stdout non-TTY.
8. `just check` passes.

## Risks

1. **Migration transitions**: upgrading across `deleted_at` and `archived_at` must preserve data.
2. **pterm styling in pipes**: pterm forces color by default; must disable styling on non-TTY.
