# Readiness Gates

Stage readiness is semantic, not only structural.

## Designed

A FEATURE may evolve to `designed` when:

- Objective includes the intended result and non-goals.
- Current Understanding separates accepted facts from assumptions.
- Key Decisions explain why the selected design wins.
- Complex mode Evidence records a support state for decisive claims.
- Working Design defines behavior, boundaries, ownership, and failure handling.
- Rejected or Revised Assumptions records changes or `None`.
- Open Questions contains no unresolved gate-affecting item.
- Next Investigation is `None`.
- `flowforge analysis status` shows no active required work, conflict, blocker, or user decision for the Proposal.

## Planned

Additionally:

- Each Step is executable without new architecture or product decisions.
- Files, Approach, Edge Cases, Dependencies, Parallel, and Verification are explicit.
- FEATURE Verification maps user-visible outcomes and design risks to tests or inspection.
- Dependencies and accepted risks are visible.
- The active analysis plan is completed or absent.

## Decisions and exceptions

The Analyst cannot accept product, compatibility, migration, security, or conflicting-goal risk on the user's behalf. Record `user.decision_required`; after resolution, persist the choice in a DEC and update affected FEATUREs. `None` is valid for sections with no content; `TBD` is never ready.

Run `proposal inspect`, `analysis validate`, `validate all`, then use `card evolve`. A rejected gate produces Artifact fixes or a new bounded revision, not a manual stage edit.
