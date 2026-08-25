---
name: flowforge-setup
description: "Configure this repo for FlowForge engineering skills: set up local wiki issue tracking, triage label vocabulary, and domain doc layout. Run once before first use of the other engineering skills."
disable-model-invocation: true
---

# Setup FlowForge Skills

Scaffold the per-repo configuration that the engineering skills assume:

- **Issue tracker**: local-first markdown wiki under `<docs_dir>/proposals/` managed by FlowForge DAG engine
- **Triage labels**: the strings used for the five canonical triage roles
- **Domain docs**: where `CONTEXT.md` and ADRs live, and the consumer rules for reading them

This is a prompt-driven skill, not a deterministic script. Explore, present what you found, confirm with the user, then write.

## Process

### 1. Explore

Look at the current repo to understand its starting state. Read whatever exists; don't assume:

- `.flowforge/config.yaml`: check if custom `docs_dir` is configured
- `AGENTS.md` and `CLAUDE.md` at the repo root: does either exist? Is there already an `## Agent skills` section in either?
- `docs/CONTEXT.md` (or `<docs_dir>/CONTEXT.md`)
- `docs/adr/` (or `<docs_dir>/adr/`)
- `docs/agents/`: does this skill's prior output already exist?
- `docs/proposals/`: active proposal directories
- Is the `flowforge-triage` skill installed? This decides whether Section B runs at all.

### 2. Present findings and ask

Summarise what's present and what's missing. Then take the sections in order. One section, one answer, then the next.

**Section A: Issue tracker.**
Default: Local-First Markdown Wiki. Issues and specs live in `<docs_dir>/proposals/<feature-slug>/` and DAG dependencies are computed via `flowforge frontier`. Record this in `docs/agents/issue-tracker.md`.

**Section B: Triage label vocabulary.** Skip this section entirely if the `triage` skill isn't installed.

If it is installed, ask:
> Do you want to keep the default triage labels? (recommended: **yes**)

The defaults are: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`.

**Section C: Domain docs.** Default to **single-context** (one `<docs_dir>/CONTEXT.md` + `<docs_dir>/adr/`). Record this in `docs/agents/domain.md`.

### 3. Confirm and edit

Show the user a draft of:
- The `## Agent skills` block to add to `AGENTS.md` / `CLAUDE.md`
- The contents of `docs/agents/issue-tracker.md`, `docs/agents/domain.md`, and `docs/agents/triage-labels.md`

Let them edit before writing.

### 4. Write

Write the docs files:
- `docs/agents/issue-tracker.md`: local-markdown wiki with FlowForge DAG engine
- `docs/agents/domain.md`: domain doc consumer rules + layout
- `docs/agents/triage-labels.md`: label mapping (only if `flowforge-triage` is installed)

### 5. Done

Tell the user the setup is complete. FlowForge is configured for local-first issue tracking and DAG dependency calculation.
