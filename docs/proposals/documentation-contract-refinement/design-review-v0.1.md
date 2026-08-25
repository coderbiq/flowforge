# Documentation Contract Refinement Design Review v0.1

**Date:** 2026-08-25  
**Fixed point:** `HEAD` (`d329cc3`); proposal files reviewed as scoped additions  
**Scope:** `docs/proposals/documentation-contract-refinement/`

## Standards

Three hard findings were corrected:

1. `solution-design-interface-v0.2.md` contained stale out-of-scope and next-work statements that contradicted completed schema and scenario work.
2. `spec.md` still described approved and validated contracts as pending discussion.
3. `solution-design-interface-v0.2.md` said production Skill design had not occurred although the coordination design already defined it.

Two judgement findings were addressed:

- Pre-validation v0.1 authorities are moved under `history/` so directory-level discovery does not expose competing current authorities.
- Solution-design guardrails now lead with positive ownership and handoff behavior; `MUST NOT` remains only for concrete boundary failures.

## Spec

Three contract findings were corrected:

1. The physical schema now defines exact, reasoned, persistent waivers. Waivers retain the original diagnostic and cannot act as blanket skips.
2. The physical schema now defines scoped `open_items`, allowing the Catalog to reconstruct declared gaps and blockers from current files without session memory or heuristic prose judgment.
3. Strict mode now changes policy outcome only; diagnostic severity remains a fact and is not promoted from warning to error.

No material scope creep was found.

## Verdict

The proposal is ready for implementation ticket planning after validation commands confirm links, fixtures, and repository tests remain sound.
