# Skill mechanics

## Frontmatter

Every production Skill has YAML frontmatter with a stable `name` and a trigger-focused `description`. Dispatch is description-driven; there is no front-matter invocation gate.

## Invocation design

Put one distinct trigger branch in the description for each situation that should load the Skill. Collapse synonyms. Router descriptions select a path; the target Skill owns its detailed process.

Split separate invocation branches only when they require materially different instructions or hiding later steps prevents premature completion. Keep shared steps in one Skill and disclose branch-only reference material through a precise link.

## Description content convention

FlowForge has no central router. The model reads the flat skill `description` list and self-selects; there is no front-matter invocation gate, no advisory dispatcher. To make self-selection reliable and prevent sibling collisions, every production Skill `description` MUST carry three segments:

1. **Trigger phrases (触发短语)** — words the user or agent will actually say, e.g. "review" / "审查" / "复审" / "debug this". These are the only cues the model can evaluate at dispatch time.
2. **Downstream ownership (下游所有权)** — the follow-up the Skill owns, e.g. "review owns the fix-planning loop", "plan owns ticket slicing and frontier verification". Tells the caller where selecting this Skill advances the work.
3. **Negative boundary (负边界)** — `NOT for X`, distinguishing the Skill from its nearest sibling, e.g. review: "NOT for fix design — flowforge-review owns fix planning"; solution-design: "NOT for review-fix design". This is the direct remedy for the GLM-event root cause (review / solution-design / plan triangle collision).

Trigger phrases MUST NOT be body-internal methodology jargon — "genuine DAG edges", "deepening opportunities", "requirement-changing unknowns" are forbidden as dispatch cues: they are body terms the model cannot evaluate at dispatch time and would erode self-selection certainty.

`description` remains a single field; this is a content convention, not a front-matter schema split (still only `name` + `description`). When the model cannot decide from the description list, it asks the user rather than falling back to the AGENTS.md routing table — that table is human reference only, not an agent dispatch path.

## Deployment

Every relative reference must be packaged beside the Skill or at its declared shared path. A production Skill must not depend on a source-tree-only file.

