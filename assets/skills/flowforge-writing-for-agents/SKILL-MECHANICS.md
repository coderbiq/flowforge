# Skill mechanics

## Frontmatter

Every production Skill has YAML frontmatter with a stable `name` and a trigger-focused `description`. Use `disable-model-invocation: true` when invocation must be explicit or controlled by another Skill.

## Invocation design

Put one distinct trigger branch in the description for each situation that should load the Skill. Collapse synonyms. Router descriptions select a path; the target Skill owns its detailed process.

Split separate invocation branches only when they require materially different instructions or hiding later steps prevents premature completion. Keep shared steps in one Skill and disclose branch-only reference material through a precise link.

## Deployment

Every relative reference must be packaged beside the Skill or at its declared shared path. A production Skill must not depend on a source-tree-only file.

