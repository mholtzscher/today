# Spec: Configurable Database Path

**Type:** Feature  
**Effort:** S (<1 hour)  
**Status:** Ready for implementation

---

## Problem

Currently DB path is hardcoded to `~/today.db`. Users need to specify different paths for testing/isolation scenarios.

## Solution

Add `--db` flag with `TODAY_DB` env var fallback. Priority: flag > env > default. Empty values fall back to default.

---

## Deliverables

### D1: Add DB flag to root command

**Files:** `internal/cli/options.go`, `cmd/root.go`

1. Add to `internal/cli/options.go`:
   - `FlagDB = "db"` constant
   - `DBPath string` field in `GlobalOptions`

2. Add to `cmd/root.go`:
   - `&ufcli.StringFlag{Name: "db", EnvVars: []string{"TODAY_DB"}, Usage: "Database path", Value: defaultDBPath()}`
   - Resolve path: `flag > env > default`, empty falls back to default

### D2: Add testscript verification

**File:** `test/testscript/db-flag.txt`

Test:
- `today --db /tmp/test.db add "test"` creates DB at specified path
- `TODAY_DB=/tmp/test2.db today add "test"` uses env var
- `today --db "" add "test"` falls back to default

---

## Acceptance Criteria

- [ ] `today --db /custom/path.db add "x"` uses `/custom/path.db`
- [ ] `TODAY_DB=/custom/path.db today add "x"` uses env var
- [ ] `today --db "" add "x"` falls back to `~/today.db`
- [ ] Flag overrides env var when both set
- [ ] Tests pass

---

## Risks

| Risk | Mitigation |
|------|------------|
| Flag value not accessible in subcommands | Use root command's flag; urfave/cli/v3 propagates flags to subcommands |