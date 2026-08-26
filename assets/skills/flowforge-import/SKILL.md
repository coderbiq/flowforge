---
name: flowforge-import
description: Interpret local PRDs, old proposals, briefs, or notes into traceable candidate facts before requirement or solution-design work. Use when an existing local document is the starting point; it does not convert documents mechanically.
disable-model-invocation: true
---

# Import external material

Turn local source material into a compact, traceable hand-off for the Skill that owns the next unsettled question. Import owns source interpretation; it does not own new requirement, design, ticket, or evidence authority.

Use the shared contract's [source intake and semantic rewrite](../_shared/ARTIFACT-CONTRACT.md#source-intake-and-semantic-rewrite), [authority roles](../_shared/ARTIFACT-CONTRACT.md#roles-and-authority), and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value).

## Inputs

Resolve the supplied local source paths, target feature, optional target language, and current project authority. Read only the material needed to establish the current request; retain the source path and nearest heading for each fact that survives.

## Classify before hand-off

Classify retained content as one of:

- **Source fact**: background, current behavior, or historical fact that informs later authority.
- **Requirement candidate**: an outcome, scope, scenario, constraint, or term that Align must accept or change.
- **Design decision**: a responsibility, interface, seam, flow, migration, or verification choice that Solution Design must review.
- **Delivery/verification evidence**: completed behavior or observed proof; it remains evidence rather than becoming a new request.
- **Unknown/conflict**: incompatible sources or an unresolved claim; name the affected candidate and return it to its owner.

Discard copied templates, repeated history, and source prose that changes no decision. A short source summary belongs beside the receiving authority. Promote a `source-notes.md` artifact with `role: research` only when several sources need independent review, a substantive conflict needs a durable rationale, or later work must revisit the classification.

## Hand off

Send requirement candidates and requirement-changing unknowns to `flowforge-align`. Send accepted requirement work that still needs responsibility, interface, seam, flow, migration, or verification decisions to `flowforge-solution-design`. Evidence stays with the relevant evidence or source record. Return the source locations, classification, target language, and selected next owner; do not publish a ticket or declare the source superseded.

## Completion

The hand-off names every retained fact's category and source location, target language, and next owner; it separates completed evidence from future work and contains only the meaning the receiver needs. The receiving Skill writes any new authority in that target language as a semantic rewrite and retains stable identifiers and terms.
