# Solution-design Production Validation v0.1

**Date:** 2026-08-25
**Subject:** `assets/skills/flowforge-solution-design/SKILL.md`

## Independent forward test

An independent agent applied the production Skill to a cross-module configured-provider request with one missing requirement: whether absent/unknown configuration fails startup or selects a default.

Observed behavior:

- Correctly invoked solution design for responsibility, seam, information-flow, and migration changes.
- Returned only the externally observable missing-configuration behavior to Align.
- Continued application/web responsibility, neutral construction seam, explicit selection, expand–migrate–contract order, and dependency-inversion verification.
- Scoped the requirement gap to provider selection, missing-configuration verification, final cutover, and legacy-path removal rather than blocking the whole design.
- Invoked codebase design; correctly skipped research, prototype, and domain modeling because their branch conditions were absent.
- Kept fixtures as examples rather than authority and produced area-level coverage without tickets or a readiness state.

The run exposed provisional-authority ambiguity, an overly broad Align return, checklist-as-template risk, repeated coverage details, ticket-biased open-item scope, and fixture contamination risk. Production instructions were revised for each observed failure mode.

## Required scenario matrix

A second independent evaluator ran all ten inputs directly against the revised production Skill and its conditionally loaded references. It reported 10/10 PASS and audited persisted readiness, ticket publication, production changes, silent requirement assumptions, and duplicate authority for every case.

| Scenario | Actual production behavior and verdict |
|---|---|
| Local feature at an existing seam | Step 1 returns directly to Plan with the reason solution-design ceremony adds no value. |
| Cross-module responsibility change | Steps 1–4 own module responsibility, interface, seam, flow, migration, and verification authority. Confirmed by forward test. |
| Requirement ambiguity | Inputs return only outcome/scope/scenario/constraint/meaning questions to Align; unaffected areas continue. Confirmed by forward test. |
| Missing primary-source fact | Step 2 invokes Research for one named fact and resumes the same decision from its verdict. |
| Behavioral/interface uncertainty | Step 2 invokes Prototype for one falsifiable question and retains only its verdict. |
| Partial planning | Steps 5–6 use narrow authority-area scopes and return per-area coverage without a global state. Confirmed by forward test. |
| Adaptive packaging | Step 4 loads `DESIGN-PACKAGING.md` only when deciding whether independent authority earns a split. |
| Semantic revisions | Step 4 increments only meaning-changing areas and preserves identity across packaging changes. |
| Information-value compression | Step 7 deletes repetition, commonplaces, discarded exploration, and copied upstream prose. |
| Implementation discovery return | Role/hand-off boundaries keep responsibility/interface/seam changes in solution design while factual corrections remain downstream. |

## Verdict

The production Skill covers all required interaction branches. The independent complex run confirmed the highest-risk partial-planning and requirement-return behavior; the ten-input evaluation passed every route and forbidden-behavior audit. Deterministic deployment tests protect its references and shared-contract anchors. Schema metadata is progressively disclosed through the shared schema v1 reference rather than recreated by each Skill.

## Reproduction protocol

Run a fresh agent with only the production `SKILL.md` and references it conditionally selects. For each row below, ask for actual route, produced/updated authorities, coverage output, and an audit of persisted readiness, ticket publication, production mutation, silent requirement assumption, and duplicate authority. A pass requires the stated observable result and zero forbidden behaviors.

1. “Add one validation branch inside an already approved parser seam; the requirement and verification are settled.” → direct Plan route.
2. “Move configured provider construction out of Application while preserving dependency inversion and gradual migration.” → solution design plus codebase-design; design authority updated.
3. “The requirement does not say whether missing provider configuration fails or falls back.” → Align only for that observable branch; unaffected areas continue.
4. “The design depends on an unavailable upstream API lifecycle guarantee.” → one named Research detour, then resume.
5. “Two interface state models remain hard to evaluate on paper.” → one-question Prototype plus codebase-design, then resume.
6. “One selection-policy question remains open while responsibility and migration are settled.” → scoped open item and partial planning, no global gate.
7. “A design hub has independently reviewed areas with unrelated revisions.” → load adaptive packaging and lazily split only earned authorities.
8. “The design contains repeated requirements, exploration transcripts, and generic implementation advice.” → compression removes filler while retaining decisions and constraints.
9. Present separately: corrected current-code fact; local implementation refinement preserving approved contracts; discovery moving a seam. → update authority fact, keep ticket-local detail, or resume solution design respectively.
10. “Plan asks whether design-ready is true.” → return area status/link/revision/diagnostics and current authority facts, never a persisted Boolean.

The recorded evaluator returned PASS for all ten inputs. Cases 3, 4, and 7 require semantic judgment, but their branch criteria were sufficient and did not authorize a forbidden behavior.
