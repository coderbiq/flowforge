# Design: <Proposal Title>

This directory is the design entry point for large proposals and for medium proposals that opt into a split layout.

The files in this directory are reference sections, not a template engine. A project may:

- use the directory as-is
- copy one section file and adjust it
- copy the whole `design/` directory into a workspace-local template area and tailor it as a unit

If a proposal only needs a small change, keep the single-file `design.md` layout instead.

## Reading order

1. [Architecture](./architecture.md)
2. [Model overview](./model.md) and detailed models under [`../model/`](../model/README.md)
3. [Lifecycle](./lifecycle.md)
4. [Flow](./flow.md)
5. [API](./api.md)
6. [Constraints](./constraints.md)
7. [Tradeoffs](./tradeoffs.md)

## Canonical corpus reviewed

<Canonical corpus reviewed>

## Chosen approach

Short statement of the intended implementation at the level of "what is built and why".

## Major decisions

### Decision 1

- Choice:
- Reason:
- Alternatives rejected:

## Knowledge impact

<Knowledge impact>

## Canonical corpus maintenance

- Canonical entry point: <which doc remains the first read>
- In-place updates: <which existing sections will be edited rather than duplicated>
- Historical trace: <where replaced facts or deprecated assumptions will be preserved>
- Sync set: <which linked docs must be updated together>

## Milestones and checkpoints

- Milestone 1: <phase boundary and intermediate deliverable>
- Milestone 2: <next phase boundary>
- Checkpoint rules: <when work must pause for verification or review>

## Risks and mitigations

- Risk:
- Mitigation:

## When to customize this directory

Customize the design directory when the proposal needs:

- architecture details that differ from the default explanation
- a different flow order or additional flow sections
- project-specific lifecycle language
- extra constraints or tradeoff sections
- section-level wording that should be visible to agents as a reference

If several sections need adjustment, copy the whole directory and keep the section headings aligned so the agent can still follow the reading order.
