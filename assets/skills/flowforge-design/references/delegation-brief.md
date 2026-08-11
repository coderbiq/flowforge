# Delegation Brief

The Design Analyst publishes work; the Coordinator schedules it; an Investigator writes only the designated FIND.

## Required work-item fields

- `workId`: stable within the analysis cycle and revision.
- `question`: one answerable question.
- `scope`: included and excluded investigation boundaries.
- `role`: normally `flowforge-investigator`.
- `sources`: allowed code, Library, document, log, or network surfaces.
- `skill`: required specialist skill, or `None`.
- `inputs`: persisted FEATURE, DEC, FIND, paths, and commands needed to start.
- `evidenceTarget`: FIND ID created before dispatch.
- `doneWhen`: observable completion condition.
- `dependsOn` and `parallelGroup`: explicit scheduling relationships.
- `required`: whether re-entry waits for this result.
- `budget`: bounded calls, sources, or time proxy.

## Investigator brief

State the exact question, known facts, exclusions, allowed sources, designated FIND, completion condition, and return statuses. The Investigator may edit only `Evidence`, `Source`, `Impact`, and `Open Questions` in that FIND. It must distinguish fact, interpretation, assumption, unknown, and risk; cite recoverable sources; and return `COMPLETED`, `INCONCLUSIVE`, `EVIDENCE_CONFLICT`, or `BLOCKED`.

Coordinator must not invent work, broaden scope, accept evidence, or revise design. A returned suggestion beyond scope becomes a candidate for the next Analyst revision, never an automatic follow-up.
