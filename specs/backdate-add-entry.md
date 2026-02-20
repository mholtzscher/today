# Backdate Entry On Add

**Type:** Feature
**Effort:** M
**Status:** Ready for implementation

## Problem Statement

**Who:** CLI journal users who forgot to log a win on a prior day.

**What:** Allow adding a new entry assigned to a specific prior calendar day.

**Why it matters:** Keeps daily log accurate; avoids manual DB edits; preserves fast add flow.

## Proposed Solution

Extend `today add` with `--date YYYY-MM-DD` (local timezone).

- If `--date` omitted: existing behavior unchanged (timestamp defaults in SQLite).
- If `--date` provided:
  - Parse strictly as `YYYY-MM-DD`.
  - Enforce past-only: reject future dates (relative to local date).
  - Store `created_at` as the local start-of-day (00:00:00) for that date.

Also render day headings in `today show` using local timezone so entries group under the expected day.

## Scope & Deliverables

| Deliverable | Effort | Depends On |
|------------|--------|------------|
| [D1] Add `--date` flag + validation to `today add` | M | - |
| [D2] Add sqlc query to insert with explicit `created_at` | M | - |
| [D3] Add Store method for backdated insert | S | D2 |
| [D4] Show groups by local date | S | - |
| [D5] Stabilize ordering for same timestamps (`ORDER BY created_at DESC, id DESC`) | S | - |
| [D6] Testscript coverage + README update | M | D1-D5 |

## Non-Goals

- Relative date inputs (`yesterday`, `2d`) or interactive date picking.
- New `show --date` filtering or calendar-day range semantics.

## Data Model

No schema changes.

- `entries.created_at` remains `INTEGER` unix epoch seconds.

## API/Interface Contract

### CLI

`today add [--date YYYY-MM-DD] <text>`

`--date` behavior:

- Parsing: strict `2006-01-02`.
- Timezone: local (`time.Local`) for date boundaries.
- Validation: must be `<=` today (local date) or return error (exit 1).
- Stored timestamp: local start-of-day (00:00:00) converted to unix epoch seconds.

### DB/sqlc

Add a new query (do not change existing `CreateEntry` signature):

```sql
-- name: CreateEntryAt :exec
INSERT INTO entries (text, created_at)
VALUES (?, ?);
```

### Store

Add a new method used by `cmd/add` when `--date` is present:

- `CreateEntryAt(ctx, text, createdAt time.Time) error`

## Output Behavior

### `today show`

- Day headers use local timezone:
  - `e.CreatedAt.In(time.Local).Format("2006-01-02")`

### Ordering

- Update list queries to include a deterministic tiebreaker:
  - `ORDER BY created_at DESC, id DESC`

Rationale: multiple backdated entries will share the same midnight timestamp.

## Acceptance Criteria

- [ ] `today add --date 2000-01-01 "X"` succeeds.
- [ ] `today show --days 20000` includes `=== 2000-01-01 ===` and `X`.
- [ ] `today add --date 2999-01-01 "X"` fails (exit 1) with a clear future-date error.
- [ ] `today add --date not-a-date "X"` fails (exit 1) with a clear format/parse error.
- [ ] `today add "X"` behavior unchanged.
- [ ] `today show` groups by local date.
- [ ] `just check` passes.

## Test Strategy

- Add `test/testscript/scripts/backdate.txtar`:
  - Success: add with far-past date; verify via large `show --days`.
  - Failure: future date rejected.
  - Failure: invalid format rejected.

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Backdated midnight may fall outside some recent windows | Medium | Low | Document; user can increase `--days` |
| Local timezone grouping differs from prior UTC-based output | Medium | Medium | Make behavior explicit; keep consistent across add/show |

## Trade-offs Made

| Chose | Over | Because |
|------|------|---------|
| `today add --date` | new subcommand | minimal UX surface area |
| strict `YYYY-MM-DD` | relative parsing | scriptable + unambiguous |
| new sqlc query | changing `CreateEntry` | avoids breaking existing callers |
| midnight timestamp | keep current time-of-day | matches “for that day” semantics |

## Open Questions

- None
