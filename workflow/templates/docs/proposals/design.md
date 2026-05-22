# Design: <Proposal Title>

This single-file design is intended for `small` and `medium` proposals.
For `large` proposals, use the `design/` directory layout instead and delete this file.
For `medium` proposals that opt into the split layout, also use the `design/` directory.

> Section guidance per size:
> - `small`: sections marked OPTIONAL may be omitted.
> - `medium`: all sections should be present, but may be brief.

## Canonical corpus reviewed

<Canonical corpus reviewed>

## Chosen approach

Describe the intended implementation at the right level of abstraction.

## Major decisions

### Decision 1

- Choice:
- Reason:
- Alternatives rejected:

## Architecture <!-- OPTIONAL for small -->

- Module boundary, layering, and cross-module impact
- Allowed and forbidden dependency directions

## Data and interfaces

- Data structures (inline for small; link to `model/` for medium when 2+ models exist)
- APIs or commands
- State transitions

## Lifecycle <!-- OPTIONAL for small -->

- State model
- Validation gates
- Audit and history
- Enable / disable semantics

## Flow <!-- OPTIONAL for small -->

- Primary flows and sequences

## Constraints

- Hard constraints
- Non-goals
- Preconditions

## Tradeoffs <!-- OPTIONAL for small -->

- Chosen approach vs alternatives
- Cost of each alternative
- Conditions under which the tradeoff should be revisited

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
