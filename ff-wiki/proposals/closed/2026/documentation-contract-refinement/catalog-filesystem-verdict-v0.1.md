# Artifact Catalog Filesystem Prototype Verdict

**Date:** 2026-08-25  
**Status:** Prototype result; production discovery remains unchanged

## Primary source

[catalog_fs.throwaway.go](parser-prototype/catalog_fs.throwaway.go)

Run:

```bash
go run ./docs/proposals/documentation-contract-refinement/parser-prototype/catalog_fs.throwaway.go \
  docs/proposals \
  /vol3/1000/develop/tangram-v2/ff-wiki-v5/proposals
```

## Verdict

The Artifact Catalog interface preserves real legacy ticket discovery and DAG behavior while removing false `spec.md` nodes.

### FlowForge proposal tree

```text
Markdown artifacts:          19
Current DiscoverIssues:       1
Catalog executable tickets:   0
Removed spec nodes:            1
Graph valid:                 yes
```

The current command treats this feature's design `spec.md` as an open task. Catalog projection correctly retains all documents as artifacts and produces no executable ticket.

### Tangram v5 proposal tree

```text
Markdown artifacts:          12
Current DiscoverIssues:      12
Catalog executable tickets:  11
Legacy tickets:              11
Removed spec nodes:            1
Graph valid:                 yes
Ready tickets:                 3
Blocked tickets:               1
Terminal tickets:              7
```

The executable ticket set matches the physical `issues/` directory exactly. Removing `spec.md` preserves a healthy graph and the expected three ready, one blocked, and seven terminal tickets.

## Refinements from real data

### Aggregate compatibility diagnostics

Emitting one `legacy-metadata` warning per ticket would turn every existing v5 frontier into “Ready with warnings,” even though the tickets remain fully compatible. This is high-volume, low-actionability noise.

Recommendation:

- Catalog records legacy status for each artifact internally.
- `check` aggregates it by feature: “11 legacy tickets use path-inferred role.”
- Default `frontier` does not attach this compatibility warning to each ticket.
- Strict policy may report or reject legacy artifacts.
- A legacy ticket with an actual parse fallback problem still receives a direct diagnostic.

Compatibility is therefore visible without reducing the signal of execution warnings.

### Infer safe conventional non-ticket roles

Metadata-free `spec.md` and `map.md` can be assigned legacy non-executable roles by filename. This improves catalog navigation without affecting safety.

Recommended inference:

- `spec.md` → legacy spec;
- `map.md` → legacy map;
- files under `issues/` → legacy ticket;
- all other metadata-free Markdown → unknown non-executable artifact.

Do not infer requirement/design/evidence from arbitrary names beyond approved conventions in schema v1.

### Catalog-level versus ticket-projected diagnostics

Not every diagnostic belongs on frontier tickets.

- Compatibility summary with no execution impact: catalog-level, shown by `check`.
- Invalid frontmatter fallback on one executable ticket: ticket-projected, shown by `frontier` and quiet stderr.
- Missing anchor in an unconsumed design area: catalog-level.
- Missing anchor in a consumed design area: projected to its consumers.

Diagnostic projection must be based on affected work, not merely the file where a diagnostic was found.

## Interface validation

The real trees validate the proposed separation:

```text
Artifact Catalog: 12 artifacts
Ticket projection: 11 tickets
Ticket DAG: healthy
```

If `spec.md` had to become an Issue before exclusion, the model would remain conceptually wrong. Keeping it as an Artifact and projecting only tickets is the deeper interface.

## Limits

- The filesystem prototype uses current legacy `ParseIssueFile` for ticket body fields after selecting ticket paths.
- It does not yet parse real schema v1 files because the proposal fixtures are intentionally non-production examples.
- It does not render `check/frontier` output streams.
- It validates one FlowForge and one Tangram proposal tree, not the full v1-v4 archive.

## Next work

The Catalog interface is sufficiently validated. Next design work should define coordinated Skill changes so new artifacts are authored consistently, then create implementation tickets that introduce Catalog in tracer slices while preserving current DAG behavior.

