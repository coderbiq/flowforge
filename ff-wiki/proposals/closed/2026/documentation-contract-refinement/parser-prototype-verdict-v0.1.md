# Schema v1 Parser Prototype Verdict

**Date:** 2026-08-25  
**Status:** Prototype result; production parser is unchanged

## Question

Can the existing Go tracker adopt Minimal Physical Markdown Schema v0.1 with its current YAML dependency while preserving legacy tickets, excluding non-ticket artifacts safely, and producing scoped revision/link diagnostics?

## Primary source

[schema_v1.throwaway.go](parser-prototype/schema_v1.throwaway.go)

Run from the repository root:

```bash
go run ./docs/proposals/documentation-contract-refinement/parser-prototype/schema_v1.throwaway.go
```

## Verdict

Yes. The prototype passes 17/17 scenarios using the repository's existing `gopkg.in/yaml.v3` dependency.

No new parsing library or content database is needed. Directory-safe discovery, YAML role metadata, legacy body fields, area revisions, consumer dependencies, and deterministic diagnostics can coexist.

## Scenarios passed

### Discovery and compatibility

- `spec.md` is excluded.
- Legacy `issues/*.md` remains executable with `legacy-metadata` warning.
- Schema v1 ticket under `issues/` is executable.
- Ticket role outside `issues/` remains non-executable.
- Evidence role inside `issues/` remains non-executable.
- Invalid YAML falls back to a legacy issue in default mode.
- Invalid YAML fails safely in strict mode.
- Unsupported future schema warns while the issues-directory safety rule still applies.

### Authority structure

- A valid area anchor resolves.
- A missing area anchor warns.
- Duplicate feature-local semantic identity warns.

### Revision and semantic links

- Matching area revision produces no diagnostic.
- Only the changed area warns its consumer.
- Missing human link warns.
- Consumed revision newer than authority warns.
- A semantic link without a machine dependency warns as `untracked-upstream`.

### Human DAG fields

- `Blocked by: 01, 03` parses as two identifiers.
- `Blocked by: None` parses as an empty blocker set with visible human meaning.

## Production seam found

The current tracker combines three responsibilities:

1. discover executable files;
2. parse a file as an `Issue`;
3. build the ticket DAG.

It also deliberately discovers `spec.md`, defaults every parsed file to `type: task`, and represents both artifact role and work kind through the overloaded `Issue.Type` field.

Schema v1 needs a deeper seam:

```text
proposal files
  -> artifact catalog + diagnostics
  -> executable ticket projection
  -> ticket DAG
```

Recommended conceptual interface:

```text
DiscoverArtifacts(root) -> Catalog

Catalog:
  Artifacts
  Tickets
  Diagnostics
```

`BuildGraph` consumes only `Catalog.Tickets`. Requirement, design, spec, evidence, research, and map artifacts remain available for semantic dependency checks without becoming graph nodes.

This interface centralizes safety, compatibility, and diagnostics. Adding frontmatter parsing directly to `ParseIssueFile` while retaining `DiscoverIssues` would leave non-ticket authority invisible and continue conflating role with work kind.

## Required model separation

The production model should distinguish:

- `ArtifactRole`: requirement, design, spec, ticket, evidence, research, map;
- `WorkKind`: bug, task, refactor, or other ticket classification;
- `TicketStatus`: existing lifecycle compatibility values;
- `Diagnostic`: severity, code, artifact, authority area, and explanation;
- `AuthorityArea`: semantic identity, revision, anchor;
- `ConsumedAuthority`: kind, qualified identity, consumed revision.

This does not require every field to land in the first tracer bullet. It defines the seam so later diagnostics do not force another parser rewrite.

## Diagnostics must reach frontier output

Default-mode invalid frontmatter deliberately falls back to an executable legacy ticket. This is compatible and non-blocking only if the warning is visible to the user and execution model.

Therefore:

- discovery diagnostics cannot be logged and discarded;
- `check` should report structural and semantic diagnostics;
- `frontier` should attach relevant warnings/gaps to each ticket;
- quiet output needs a defined policy so automation does not silently lose risk information;
- strict mode filters or fails according to policy, not by changing persisted ticket state.

Without diagnostic propagation, safe degradation becomes silent degradation.

## Frontmatter/body interaction

The current `ParseIssueFile` begins body parsing at the first non-bold line. A leading `---` would end its metadata-header scan before it reaches `**Status:**` and `**Blocked by:**`.

Production parsing must first split YAML frontmatter from Markdown body, then pass the body to the legacy human-field parser. This preserves the decision that YAML does not duplicate status or blockers.

## Cross-artifact discovery cost

Duplicate semantic IDs, anchors, consumed revisions, and human-link checks require reading metadata from non-ticket artifacts. Discovery should scan proposal Markdown files for a small FlowForge envelope while keeping executable projection restricted to `issues/`.

It does not need to parse every Markdown section into a database. The catalog can retain file paths, role metadata, identified areas, and the body needed for anchor/link checks.

## Limits of the prototype

- It uses in-memory fixtures rather than the filesystem walker.
- Human-link checking is intentionally simple and does not implement a full Markdown parser.
- It does not render CLI diagnostics.
- It does not implement cross-feature identity resolution.
- It does not test concurrent file mutation.
- It does not decide whether invalid frontmatter fallback is allowed under every project policy.
- It does not modify current parser or prove backward compatibility against every historical proposal.

These are implementation and tracer-test concerns; none invalidates the physical schema shape.

## Refinements recommended before implementation

1. Add an artifact-catalog seam rather than growing `DiscoverIssues` into a mixed parser.
2. Split frontmatter before applying legacy body-field parsing.
3. Propagate diagnostics through `check` and `frontier`; warning-only fallback must never be silent.
4. Separate ArtifactRole from WorkKind in the model.
5. Keep DAG construction ticket-only while cross-artifact checks operate on the catalog.
6. Define quiet-output diagnostic behavior before relying on it for automation.

## Next design decision

The schema itself is viable. The next design work should define the artifact-catalog interface and diagnostic projection into `check` and `frontier`, then validate it against real legacy proposal fixtures before production implementation tickets are created.

