---
name: flowforge-research
description: Investigate a question against high-trust primary sources and capture the findings as a Markdown file in the repo. Use when the user wants a topic researched, docs or API facts gathered, or reading legwork delegated to a background agent. Spins up a background agent and persists the cited findings as a Markdown file in the repo. NOT for bug diagnosis — that's flowforge-diagnose; NOT for design decisions — that's flowforge-solution-design.
---

Spin up a **background agent** to do the research, so you keep working while it reads.

Its job:

1. Investigate the question against **primary sources** (official docs, source code, specs, first-party APIs), not a secondary write-up of them. Follow every claim back to the source that owns it.
2. Write the findings to a single Markdown file, citing each claim's source.
3. Save it where the repo already keeps such notes; match the existing convention, and if there is none, put it somewhere sensible and say where.

## Design fact brief

When the caller (typically `flowforge-architect` preparing a `flowforge-solution-design` run) needs facts before a design ruling, produce a **design fact brief** (设计事实简报) instead of a free-form research note. The brief gathers facts; it never makes the design ruling itself — that stays with solution design.

Write the brief with this four-part structure:

1. **Question** — the design decision context in one sentence (what ruling this brief feeds).
2. **Facts** — each relevant fact with a verifiable citation: `file:line` for repository facts, or the command output for environment/behavior facts. No fact without a citation; no interpretation mixed in.
3. **Constraint surfaces** — the facts that bound the design space, each with its citation: interface signatures, data models, and compatibility boundaries as separate entries.
4. **Open items** — the questions that require architect judgment (trade-offs with multiple credible answers, unsettled responsibilities or seams). This list is the ruling checklist for solution design; state each as a question, not a conclusion.

Save the brief alongside the proposal it serves — inside the proposal's directory (`<docs_dir>/proposals/<proposal-id>/`) or at the caller-specified location — and reference that path in your reply, so solution design can add it to its reading list without re-reading the codebase.
