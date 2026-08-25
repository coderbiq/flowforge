# Proposal schema v1

Read this reference only when authoring or revising FlowForge machine metadata.

```yaml
---
flowforge:
  schema: 1
  role: design
  id: feature-design
  revision: 1
  areas:
    semantic-area:
      revision: 1
      anchor: semantic-area
  open_items:
    - id: unresolved-fact
      diagnostic: design-decision-missing
      severity: gap
      affects: [semantic-area]
      anchor: unresolved-fact
---
```

Roles are `requirement`, `design`, `spec`, `ticket`, `evidence`, `research`, and `map`. Only files under `issues/` may be executable, and only when their role is `ticket` or they are compatible legacy tickets.

Use a positive whole-artifact revision while the authority changes together. Use `areas` when consumers independently consume parts of one authority. Every area has a positive semantic revision and explicit same-file anchor.

A ticket records reviewed authorities without replacing human links:

```yaml
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      feature-requirements: 2
    design:
      semantic-area: 3
```

An `open_items` entry requires `id`, `diagnostic`, `severity` (`gap` or `blocker`), non-empty `affects`, and an explanation `anchor`. Prefer affected authority areas before Plan creates tickets.

A persistent waiver requires one exact diagnostic, target, and non-empty reason:

```yaml
waivers:
  - diagnostic: upstream-changed
    target: semantic-area
    reason: "Reviewed the new revision; this consumer uses unchanged behavior."
```

Do not encode ticket status or `Blocked by` in YAML. Keep those human fields immediately after the ticket title for legacy compatibility.

