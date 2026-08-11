# Evidence Rules

## Evidence classes

Label material conclusions:

- **Fact**: directly observed in a recoverable source.
- **Interpretation**: explanation consistent with facts but requiring Analyst judgment.
- **Assumption**: provisional statement used to continue design.
- **Unknown**: unresolved information that may affect a gate.
- **Risk**: possible negative outcome with trigger and impact.

## Support state

Each decisive claim in a complex FEATURE is `accepted`, `rejected`, `conflicting`, or `inconclusive`. Cite a FIND ID or an independently recoverable code/document source. A FIND is not accepted merely because work returned; the Analyst must copy the conclusion and support state into FEATURE `Evidence`.

No-hit results remain Unknown. Conflicting sources stay visible and trigger synthesis or user decision. Do not silently choose the convenient source. Rejected evidence remains in the FIND and is summarized under `Rejected or Revised Assumptions` with the reason and revision.

## FIND write-back

Create the FIND skeleton through CLI, add `<!-- analysis-mode: complex -->` and `<!-- analysis-work-id: <revision-work-id> -->`, then edit only:

- `Evidence`: observations and labeled interpretations.
- `Source`: paths, card IDs, commands, versions, timestamps, or URLs sufficient to reproduce.
- `Impact`: consequence for the published question, not a final design decision.
- `Open Questions`: remaining uncertainty or `None`.

Run `flowforge validate all` after the edit. Never use Agent memory, chat summaries, or SQLite rows as evidence.
