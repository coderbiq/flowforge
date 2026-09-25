# System Context & Domain Glossary

This document maintains the core domain glossary, bounded contexts, and system architecture invariants.

## Glossary

- **Generic capability role**: A subagent asset defined by capability (no flowforge workflow position); dispatched via the AGENTS.md capability table. Roster: `flowforge-batch-analyst`, `flowforge-scribe`, `flowforge-executor`. See [ADR 0001](adr/0001-hybrid-generic-subagent-roles.md).
- **Capability-first ordering**: `## Generic capability dispatch` precedes `## Subagent delegation` in AGENTS.md; any work first matches a capability key, process states then consult the process table.
- **workspace-write**: Semantic permission label with no compiler special-casing (only `read-only` is special-cased); write constraints live in the asset's Boundaries text.
- **Dual-track wiki config (historical)**: The `wiki.root` / `projects[].wikiRoot` config track was display-only dead code (zero production callers) and was removed in the wiki-config-single-track proposal; `docs_dir` is the single decision track. Legacy keys load with a stderr deprecation warning and are ignored. See [research note](research/2026-09-25-dual-track-wiki-config.md).
