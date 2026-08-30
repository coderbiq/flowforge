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
- **Standards guide**: where project coding standards live and how to extract applicable standards per ticket

This is a prompt-driven skill, not a deterministic script. Explore, present what you found, confirm with the user, then write.

## Process

### 1. Explore

Look at the current repo to understand its starting state. Read whatever exists; don't assume:

- `.flowforge/config.yaml`: check if custom `docs_dir` is configured
- `AGENTS.md` and `CLAUDE.md` at the repo root: does either exist? Is there already an `## Agent skills` section in either?
- `<docs_dir>/CONTEXT.md`
- `<docs_dir>/adr/`
- `<docs_dir>/agents/`: does this skill's prior output already exist?
- `<docs_dir>/proposals/`: active proposal directories
- Is the `flowforge-triage` skill installed? This decides whether Section B runs at all.

### 2. Present findings and ask

Summarise what's present and what's missing. Then take the sections in order. One section, one answer, then the next.

**Section A: Issue tracker.**
Default: Local-First Markdown Wiki. Issues and specs live in `<docs_dir>/proposals/<feature-slug>/` and DAG dependencies are computed via `flowforge frontier`. Record this in `<docs_dir>/agents/issue-tracker.md`.

**Section B: Triage label vocabulary.** Skip this section entirely if the `triage` skill isn't installed.

If it is installed, ask:
> Do you want to keep the default triage labels? (recommended: **yes**)

The defaults are: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`.

**Section C: Domain docs.** Default to **single-context** (one `<docs_dir>/CONTEXT.md` + `<docs_dir>/adr/`). Record this in `<docs_dir>/agents/domain.md`.

**Section D: Standards guide.** The file `<docs_dir>/agents/standards.md` was deployed by `flowforge init` as a managed asset with a default generic extraction guide. Show the user its contents and ask whether they want to customise it now. The default version provides a working generic heuristic (derive layer/module from Touch points and Write set, look up governing standards, extract `must`/`must not` clauses). Encourage the user to edit it to match this project's actual standard documents and extraction logic—where standards live, how to match them to a ticket. The default is already usable; customisation can happen later at any time by editing the file directly. The standards guide path is configured at `standards.guide` (default `agents/standards.md`).

### 3. Confirm and edit

Show the user a draft of:
- The `## Agent skills` block to add to `AGENTS.md` / `CLAUDE.md`
- The contents of `<docs_dir>/agents/issue-tracker.md`, `<docs_dir>/agents/domain.md`, `<docs_dir>/agents/standards.md`, and `<docs_dir>/agents/triage-labels.md`

Let them edit before writing.

### 4. Write

Write the docs files:
- `<docs_dir>/agents/issue-tracker.md`: local-markdown wiki with FlowForge DAG engine
- `<docs_dir>/agents/domain.md`: domain doc consumer rules + layout
- `<docs_dir>/agents/standards.md`: already deployed by `flowforge init`; guide the user to customise it if they chose to in Section D
- `<docs_dir>/agents/triage-labels.md`: label mapping (only if `flowforge-triage` is installed)

### 5. Done

Tell the user the setup is complete. FlowForge is configured for local-first issue tracking and DAG dependency calculation.
