---
flowforge:
  schema: 1
  role: design
  id: artifact-catalog
  revision: 1
---

<a id="artifact-catalog"></a>
# Artifact Catalog Interface v0.1

**Date:** 2026-08-25  
**Status:** Approved design baseline; production implementation not started  
**Purpose:** Separate proposal artifact discovery and diagnostics from executable ticket projection and DAG calculation

## Module interface

```text
DiscoverArtifacts(root) -> Catalog, error

Catalog:
  Artifacts
  Tickets
  Diagnostics
```

The module hides filesystem walking, YAML/body splitting, legacy parsing, role inference, identity indexing, anchor/link indexing, revision comparison, and safe executable projection.

Callers learn one interface: a trustworthy catalog of proposal artifacts, the executable tickets derived from it, and every recoverable diagnostic found while building it.

## Error contract

`DiscoverArtifacts` returns a Go error only when it cannot construct a trustworthy catalog, including:

- proposal root cannot be inspected for a reason other than non-existence;
- filesystem walk aborts;
- a read failure prevents safe classification of a candidate executable file;
- an internal invariant makes the catalog result unreliable.

Recoverable file problems become diagnostics:

- invalid or unsupported frontmatter;
- role/location conflicts;
- legacy metadata;
- invalid anchor or revision;
- duplicate semantic identity;
- dangling semantic dependency;
- human-link/dependency mismatch.

A missing proposal root produces an empty Catalog, preserving current behavior.

## Catalog model

Conceptual Go shape:

```go
type Catalog struct {
    Artifacts   []Artifact
    Tickets     []*Ticket
    Diagnostics []Diagnostic
}
```

`Tickets` is a safe projection of `Artifacts`; it is not independently discovered.

### Artifact

```go
type Artifact struct {
    Path          string
    Feature       string
    Role          ArtifactRole
    Schema        SchemaVersion
    Identity      string
    Revision      int
    Title         string
    Areas         map[string]AuthorityArea
    Consumes      []ConsumedAuthority
    SemanticLinks []SemanticLink
    Ticket        *Ticket
}
```

The exact Go syntax may change during implementation. The interface semantics are fixed:

- common artifact identity and role are separate from ticket lifecycle;
- only ticket artifacts contain a Ticket projection;
- non-ticket authority remains available for semantic validation;
- the public model exposes indexed semantic information, not every raw Markdown line.

### Ticket

Ticket preserves current DAG behavior while separating overloaded concepts:

```go
type Ticket struct {
    ID        string
    Slug      string
    Feature   string
    Path      string
    Title     string
    Status    TicketStatus
    WorkKind  WorkKind
    BlockedBy []string
    Assignee  string
    Labels    []string
}
```

Temporary compatibility may retain body content for existing JSON consumers, but full body is not part of the desired Catalog interface.

### AuthorityArea

```go
type AuthorityArea struct {
    QualifiedID string
    Revision    int
    Anchor      string
    Source      SourceLocation
}
```

`QualifiedID` is feature-local identity resolved as `feature/id`. Human display normally uses the unqualified semantic ID.

### ConsumedAuthority

```go
type ConsumedAuthority struct {
    Kind        requirement | design
    QualifiedID string
    Revision    int
    Source      SourceLocation
}
```

### SemanticLink

The catalog indexes explicit Markdown links whose target fragment resolves to an identified authority area. It does not treat every Markdown link as a semantic dependency.

```go
type SemanticLink struct {
    TargetPath   string
    TargetAnchor string
    LinkText     string
    Source       SourceLocation
}
```

## Diagnostic model

```go
type Diagnostic struct {
    Code             DiagnosticCode
    Severity         warning | gap | blocker
    Artifact         string
    Area             string
    Message          string
    RelatedArtifacts []string
    Source           SourceLocation
    Waiver           *AppliedWaiver
}
```

### Stable machine contract

- `Code`, `Severity`, `Artifact`, and `Area` are machine fields.
- `Message` is human-facing and MUST NOT be parsed by consumers.
- Source location points to a file and, when available, a line.
- Related artifacts make upstream/downstream effects explainable.
- An applied waiver is structured policy metadata; it does not change `Severity` or remove the diagnostic.

### Policy separation

Catalog severity describes a fact and does not change under strict mode.

- Catalog returns facts.
- Command policy decides grouping, filtering, and exit status.
- Strict mode does not rebuild the catalog or mutate diagnostics.

Unrecoverable Go errors remain outside the diagnostic model.

## Build pipeline

### 1. Walk proposal Markdown files

Scan Markdown files under the proposal root. Directory and filename conventions provide feature and placement context but do not by themselves turn non-issue documents into tickets.

### 2. Split frontmatter and body

If YAML frontmatter exists:

- parse the `flowforge` envelope;
- retain unsupported/unknown-schema facts as diagnostics;
- never allow parse failure outside `issues/` to create an executable artifact.

Pass the remaining body to legacy human-field parsing so `Status`, `Type`, and `Blocked by` remain visible authority.

### 3. Determine role safely

- Valid schema role is recorded.
- Legacy `issues/*.md` infers ticket role with compatibility warning.
- Metadata-free `spec.md` and `map.md` infer legacy non-executable roles.
- Other metadata-free Markdown remains unknown and non-executable.
- Explicit non-ticket role under `issues/` is excluded.
- Ticket role outside `issues/` is non-executable.
- Other metadata-free Markdown is catalogued only when needed for legacy authority/link compatibility; it is never executable.

### 4. Parse artifact-specific metadata

- Requirement/design identity and revision.
- Design/requirement areas and anchors.
- Ticket consumed authorities.
- Evidence ticket reference.
- Existing ticket human fields.

### 5. Build feature-local identity index

- Resolve bare IDs within the current feature.
- Resolve `feature/id` explicitly across features.
- Detect duplicates before comparing revisions.
- Project invariants and ADR paths remain external references in schema v1.

### 6. Build semantic-link index

- Parse explicit Markdown links.
- Resolve target file and anchor.
- Match only identified authority areas.
- Do not infer dependency from prose similarity or unmarked terms.

### 7. Compare authority consumption

Emit scoped diagnostics for:

- missing authority;
- stale consumed revision;
- future consumed revision;
- machine dependency without human semantic link;
- semantic link without consumed revision;
- missing or duplicate anchor/identity.

### 8. Project executable tickets

A ticket enters `Catalog.Tickets` only when:

- it resides under `issues/`;
- it is a legacy issue or has valid/usable ticket role under default compatibility policy;
- no role/location safety rule excludes it.

Warnings do not automatically remove it. Gaps and blockers are attached through diagnostics for command policy to project.

## Diagnostic scope projection

Each diagnostic must be attributable to one or more affected tickets when possible.

### Direct diagnostic

A problem in a ticket applies directly to that ticket.

Examples:

- invalid frontmatter fallback;
- missing human link;
- stale consumed revision.

### Authority diagnostic

A problem in a requirement/design area applies to consumers of that area.

Examples:

- missing anchor;
- duplicate identity;
- invalid authority revision.

If no consumer exists yet, the diagnostic remains catalog-level and appears in `check` but not on an unrelated frontier ticket.

### Declared open-item diagnostic

The Catalog deterministically projects each valid `open_items` entry from the physical schema. Its code, severity, and affected scopes come from metadata; its human explanation comes from the linked authority prose. The Catalog does not infer architectural gaps from prose quality or session memory.

- A `gap` applies only to named consumers or tickets.
- A `blocker` excludes only named work and cannot be command-overridden.
- Missing scope, missing explanation anchor, or unknown target produces a structure diagnostic instead of a global gate.

### Waiver matching

The Catalog matches `waivers` by exact diagnostic code and exact target after diagnostics are built. A match attaches the waiver and reason to the original diagnostic. It never removes the fact, changes severity, marks work ready, or matches future diagnostics by wildcard.

### Compatibility diagnostic

Legacy path-inferred tickets are recorded per artifact but aggregated by feature for presentation. Default frontier does not attach `legacy-metadata` to every compatible ticket. Strict policy may report or reject legacy artifacts.

### Global diagnostic

Filesystem or feature-level ambiguity that prevents trustworthy classification is returned as Go error or, when recoverable, as a catalog diagnostic with no ticket projection.

## DAG interface

The graph module consumes only ticket facts:

```text
BuildGraph(catalog.Tickets) -> TicketGraph
```

It remains responsible for:

- blocking edges;
- cycles;
- dangling ticket blockers;
- self-blocking;
- terminal and claimed ticket behavior.

It does not parse Markdown, resolve artifact roles, compare semantic revisions, or decide content diagnostics.

Deleting Artifact Catalog would force discovery, compatibility, identity, and semantic diagnostics into every command and future UI; the module therefore passes the deletion test.

## `check` projection

Default `flowforge check` combines Catalog and DAG results.

### Exit behavior

- DAG cycle, dangling ticket blocker, or self-blocking: non-zero.
- Catalog blocker: non-zero.
- Gap: displayed, default exit zero.
- Warning: displayed, default exit zero.
- `--strict`: any gap or warning produces non-zero.

Display severity and exit behavior must agree. A diagnostic shown as a warning cannot be described as a fatal error while silently using different semantics.

### JSON

JSON always includes:

- catalog diagnostics;
- DAG diagnostics/results;
- validity under the selected policy;
- artifact and ticket counts separately.

Human messages remain informational; automation consumes stable diagnostic fields.

## `frontier` projection

Text output groups:

```text
READY
READY WITH WARNINGS
GAPS
CLAIMED
BLOCKED
```

### Classification

- **Ready:** DAG executable and no warning/gap projected to ticket.
- **Ready with warnings:** DAG executable with warning diagnostics.
- **Gaps:** DAG executable but an authority Skill explicitly recorded a required design/verification fact as missing for that ticket.
- **Claimed:** current claimed compatibility behavior.
- **Blocked:** unfinished DAG dependency or non-overridable external blocker.

### Overrides and strict policy

- `--include-gaps` includes gap tickets as explicit override while preserving their diagnostics.
- It does not include blocked tickets.
- `--strict` treats warning and gap tickets as not executable.
- No command writes a readiness status back to Markdown.

## Quiet and JSON behavior

### Quiet

- stdout contains executable ticket paths only.
- diagnostics are written to stderr.
- default quiet includes Ready and Ready with warnings.
- default quiet excludes Gaps and Blocked.
- `--include-gaps` includes gap paths and writes gap explanations to stderr.
- `--quiet --strict` emits only clean Ready paths.
- strict policy produces non-zero when diagnostics violate policy.

This preserves pipeline compatibility without making safe degradation silent.

### JSON

JSON carries all frontier groups and structured diagnostics. It never requires parsing rendered messages.

## Compatibility adapter

Retain `DiscoverIssues(root)` temporarily as a wrapper:

```text
DiscoverIssues(root):
  catalog = DiscoverArtifacts(root)
  return catalog.Tickets
```

Limitations are explicit:

- wrapper callers do not receive diagnostics;
- production `check` and `frontier` migrate directly to Catalog;
- wrapper exists only to stage internal migration;
- new call sites MUST use `DiscoverArtifacts`;
- wrapper is deleted after remaining callers migrate.

## Performance and persistence

- No persistent catalog database in schema v1.
- Walk and parse on each command, matching current local-first behavior.
- Parse frontmatter for proposal Markdown files.
- Parse body fields for tickets and body links/anchors for identified authority artifacts and consumers.
- Avoid full semantic Markdown interpretation.
- Optimization requires measured evidence; do not add caching before a real performance problem.

## Verification seam

Test through the Catalog interface using filesystem fixtures. Each fixture directory is a miniature proposal root.

Required cases:

1. Empty/missing root.
2. Current legacy tickets and statuses.
3. `spec.md` excluded from tickets.
4. Schema v1 role/location safety.
5. YAML/body splitting preserves human blocker parsing.
6. Invalid/future schema under default and strict policy.
7. Feature-local and cross-feature identity resolution.
8. Duplicate identity and missing anchor.
9. Scoped revision warning.
10. Human link/dependency mismatch.
11. Artifact diagnostic projected only to affected consumers.
12. Catalog Tickets produce the current DAG result.
13. Check default/strict exit policy.
14. Frontier five-group classification.
15. Quiet stdout and diagnostic stderr separation.
16. Legacy `DiscoverIssues` wrapper compatibility.

## Rejected interfaces

### Extend `DiscoverIssues` with more fields

Rejected because non-ticket authority would remain invisible or be misrepresented as Issue nodes.

### Let each command parse artifacts independently

Rejected because safety, legacy compatibility, identity, and diagnostics would diverge across commands.

### Store strict/default severity inside Catalog

Rejected because policy would contaminate parsing facts and require rebuilding for each caller.

### Build a persistent index now

Rejected because current scale has no measured parsing problem and persistence creates invalidation complexity.

## Out of scope

- Production code changes.
- Exact Cobra flag naming beyond the approved semantics.
- GUI projection.
- Automated Markdown rewriting.
- Content-quality scoring.
- Cross-repository authority resolution.

## Next work

1. Build a filesystem-level Catalog prototype against real legacy proposal trees.
2. Validate diagnostic projection for `check`, `frontier`, quiet, and JSON.
3. Refine this interface only where prototype evidence requires it.
4. Then design coordinated Skill changes and implementation tickets.

The real-tree validation is captured in [Artifact Catalog Filesystem Prototype Verdict](catalog-filesystem-verdict-v0.1.md). It confirmed ticket/DAG preservation and refined compatibility diagnostic aggregation.
