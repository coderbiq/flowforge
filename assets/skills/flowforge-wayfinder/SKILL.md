---
name: flowforge-wayfinder
description: Use ONLY when facing high uncertainty, greenfield architecture, unexplored tech stacks, or complex multi-phase systems where linear slicing is impossible. Builds a decision graph (MAP.md), manages exploration spikes, and advances through decision frontiers. Do NOT use for standard bug fixes or well-understood linear feature slices.
---

# flowforge-wayfinder (Decision Graph & Fog-of-War Exploration)

Navigate high uncertainty and complex multi-phase technical forks by maintaining an explicit Decision Map (`MAP.md`) and exploring along the decision frontier.

## ⛔ Absolute Rules

1. **DO NOT generate linear implementation plans under thick fog**: If the core architecture, data model, or key libraries are undecided, write Decision Tickets, not code slices.
2. **One Question per Decision Ticket**: Each ticket in `01-workspace/<proposal_id>/MAP.md` must resolve one specific fork or validate one technical spike.
3. **Advance only along the Frontier**: Only work on tickets whose dependencies are marked `[DONE]`.

## MAP.md Structure (`01-workspace/<proposal_id>/MAP.md`)

```markdown
# Exploration Decision Map: <Topic>

## 1. High-Level Vision & Hypotheses
- **North Star**: <What success looks like>
- **Core Unknowns**: <Top 2-3 risks to de-risk>

## 2. Decision Dependency Graph (Frontier)

- [ ] **TICK-1: <Title>** `[FRONTIER]`
  - **Question / Fork**: <Specific technical question>
  - **Prerequisites**: None
  - **Spike / Probe Action**: <Read-only probe, prototype script, or library benchmark>
  - **Outcome / Decision**: Pending

- [ ] **TICK-2: <Title>** `[BLOCKED by TICK-1]`
  - **Question / Fork**: <Dependent question>
  - **Prerequisites**: TICK-1

## 3. Settled Decisions Log
(Record resolved tickets here with rationales)
```

## Workflow

1. Identify core forks/unknowns and initialize `MAP.md`.
2. Pick the highest priority `[FRONTIER]` ticket.
3. Execute minimal spike/investigation with `flowforge-explore`.
4. Capture findings, record settled decisions in `MAP.md` and `README.md`.
5. Once the fog clears (core forks decided), graduate to `flowforge-align` / `flowforge-plan`.
