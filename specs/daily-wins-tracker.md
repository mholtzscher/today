# Daily Wins Tracker

**Type:** Feature  
**Effort:** L (1-2 days)  
**Status:** Ready for implementation

## Problem

Need a frictionless CLI to capture daily accomplishments and retrieve them by date range. Quick-add is the critical path.

## Commands

```
today add <text>              # Quick add entry
today show                    # Show today's entries
today show --days 7           # Show past 7 days
today show 3                  # Show past 3 days (shorthand)
```

## Output Format

Grouped by day:

```
=== 2026-02-15 ===
• Fixed auth bug
• Shipped feature X

=== 2026-02-14 ===
• Refactored API layer
```

## Schema

```sql
-- +goose Up
CREATE TABLE entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  text TEXT NOT NULL,
  created_at TEXT DEFAULT (datetime('now'))
);

-- +goose Down
DROP TABLE entries;
```

## File Structure

```
internal/
  db/
    db.go              # Connection, auto-migrate
    migrations/
      fs.go            # embed.FS for migrations
      00001_init.sql   # Initial schema
  entry/
    entry.go          # Entry type, queries
cmd/
  add/add.go          # add subcommand
  show/show.go        # show subcommand
```

## Dependencies

- `modernc.org/sqlite` — pure Go SQLite (no CGO)
- `github.com/pressly/goose/v3` — migrations

## Trade-offs

| Choice | vs Alternative | Why |
|--------|-----------------|-----|
| Single command add | Interactive prompt | Faster, no friction |
| `~/today.db` | XDG dir | Simpler |
| Auto timestamp | Allow backdating | Simpler UX |
| Show today default | Show week | Less noise |
| Goose migrations | Inline schema | Versioned, reversible |
| Embedded migrations | File-based | Single binary, portable |

## Deliverables (Ordered)

1. **[S]** Add `modernc.org/sqlite` + `github.com/pressly/goose/v3` deps
2. **[M]** Create `internal/db/migrations/00001_init.sql` — entries table schema
3. **[S]** Create `internal/db/migrations/fs.go` — `//go:embed *.sql` + embed.FS
4. **[M]** Create `internal/db/db.go` — Open, goose.SetBaseFS, goose.Up
5. **[M]** Create `internal/entry/entry.go` — Entry struct, Insert, GetByDays
6. **[M]** Create `cmd/add/add.go` — add subcommand
7. **[M]** Create `cmd/show/show.go` — show subcommand, grouped output
8. **[S]** Wire into `cmd/root.go`, remove example command

## Risks

1. **DB file permissions** — ensure parent directory exists before creating
2. **Timezone handling** — SQLite `datetime('now')` is UTC; consider local time

## Open Questions

- [ ] Search/filter by text? (parked)
- [ ] Edit/delete commands? (parked)