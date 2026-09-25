# System Context & Domain Glossary

This document maintains the core domain glossary, bounded contexts, and system architecture invariants.

## Glossary

- **Generic capability role**: A subagent asset defined by capability (no flowforge workflow position); dispatched via the AGENTS.md capability table. Roster: `flowforge-batch-analyst`, `flowforge-scribe`, `flowforge-executor`. See [ADR 0001](adr/0001-hybrid-generic-subagent-roles.md).
- **Capability-first ordering**: `## Generic capability dispatch` precedes `## Subagent delegation` in AGENTS.md; any work first matches a capability key, process states then consult the process table.
- **workspace-write**: Semantic permission label with no compiler special-casing (only `read-only` is special-cased); write constraints live in the asset's Boundaries text.
- **Dual-track wiki config (historical)**: The `wiki.root` / `projects[].wikiRoot` config track was display-only dead code (zero production callers) and was removed in the wiki-config-single-track proposal; `docs_dir` is the single decision track. Legacy keys load with a stderr deprecation warning and are ignored. See [research note](research/2026-09-25-dual-track-wiki-config.md).
- **Per-machine deploy artifacts**: Subagent/AGENTS deploy outputs and `.flowforge/config.yaml` are per-machine files; init/upgrade auto-gitignore them. Legacy tracked artifacts get copyable `git rm --cached -r` guidance — never automatic index changes. See [ADR 0002](adr/0002-deploy-artifact-localization-and-model-chain.md).
- **models_by_host**: Nested per-host model overrides (host → agent-name | profile-key → model). The per-host layer beats global `models`/`models_by_name` wholesale; within a layer, name keys beat profile keys. `models_by_name` accepts agent names only (inert keys are config corruption signals).
- **Six-level model precedence**: `by_host[name] > by_host[profile] > by_name[name] > models[profile] > preserve-merge backfill > host default`. pi subagent frontmatter gains `model:` when the chain resolves a value (omitted when empty — unpinned output is byte-identical).
