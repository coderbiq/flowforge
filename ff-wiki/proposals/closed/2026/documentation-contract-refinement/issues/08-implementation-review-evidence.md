# 08: Close work with implementation and review evidence

**Blocked by:** 07

**Status:** closed

**What to build:** Implement, TDD, and Review consume the effective linked specification, return design discoveries to the owning Skill, and close a ticket only with concise verification and review disposition evidence.

See [Evidence interface](../artifact-handoff-interface-v0.2.md#evidence-interface) and [effective specification](../artifact-handoff-interface-v0.2.md#effective-specification).

## Touch points

Implement, TDD, Review, and Handoff Skills; inline and promoted evidence examples; review fixed-point handling.

## Changes

1. Resolve applicable project, requirement, design, ticket, and waiver authorities before execution or review.
2. Distinguish factual corrections from responsibility/seam changes and return the latter to solution design.
3. Record delivered behavior, actual commands/results, review dispositions, deviations, and implementation reference after execution.
4. Promote evidence only when it has independent audit, lifecycle, size, or cross-ticket value.

## Constraints

- Checking boxes alone is not completion evidence.
- Evidence must not copy the ticket, dump an unfiltered terminal transcript, or silently redefine authority.
- Standards and Spec review findings remain separate axes.

## Done and verify

- Scenario tests cover stale authority, design return, factual correction, inline evidence, promoted evidence, and reasoned finding disposition.
- A ticket cannot be closed through the new flow without observable verification evidence.

## Completion evidence

- Implement now resolves effective linked authority and diagnostics, classifies discoveries by owner, drives TDD at approved seams, invokes fixed-scope dual-axis review, and records evidence before close.
- TDD reuses current pre-agreed seams and returns absent/conflicting/moved seams as design gaps; Handoff transports authority/evidence links plus context delta.
- Review accepts committed or working-tree scope from a resolvable fixed point, evaluates Standards and effective Specification independently, and never writes evidence or closes tickets.
- [Implementation, Review, and Evidence Production Validation v0.1](../implementation-production-validation-v0.1.md) covers stale authority, design return, factual correction, inline/promoted evidence, reasoned dispositions, and the no-evidence/no-close invariant with exact prompts and preserved evidence shapes.
- Catalog emits `missing-completion-evidence` for a closed ticket without a non-empty `Completion evidence` section; tracker and command tests prove normal warning projection and `check --strict` rejection.
