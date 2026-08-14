# Classification Rules

Classify each discovery exactly once before routing it.

## Decision Tree

```
Discovery → reproducible deviation from an accepted contract? → bug
          → missing or ambiguous acceptance/design decision?    → design-gap
          → reusable convention, decision, or evidence?         → knowledge
          → otherwise                                           → finding
```

## Routes

### bug

Attach the evidence to the affected FEATURE with `card log --kind bug`. If the
fix is not already planned, create a FEATURE with `card init --type feature`,
link it with `card link`, and define executable steps before implementation.

### finding

Create or update a FIND card through the approved card workflow. Keep the
observation, source, evidence, impact, and open questions explicit. A finding
does not authorize implementation by itself.

### knowledge

Capture a stable reusable conclusion as a FIND or DEC and reference it from
the affected FEATURE. Use `library import` or `library promote` only through
the curation workflow after approval; never write library files directly.

### design-gap

Record the blocker in the affected FEATURE and Proposal Journal, mark the Step
blocked, and route it to `flowforge-design`. A design analyst must resolve the
scope, compatibility, migration, security, or ownership decision before a
new implementation Step is executed.

## Anti-patterns

- Do not turn every observation into an implementation request.
- Do not use deprecated command names as current examples.
- Do not silently reinterpret a historical card or wiki entry as a v3 route.
- Do not close a FEATURE while its verification or design decision is open.
