# Align and Route Production Validation v0.1

**Date:** 2026-08-25

## Reproduction procedure

For each prompt below, start a fresh agent context containing the repository `AGENTS.md`, then invoke the production `flowforge-route` Skill. Follow the selected Skill only far enough to identify its next hand-off and proposed artifact writes. Compare the output with the expected owner and every forbidden behavior. A case passes only when the route names exactly one next owner and all checklist items pass.

## Cases

### A. Settled local change

**Exact prompt:** “The parser's approved diagnostics seam is documented by the current solution design revision. Add a warning count to its existing result. The observable outcome and acceptance examples are settled, but no adequate implementation ticket exists. What should happen next?”

**Expected owner:** Plan.

**Recorded evaluator response:** “Route selected Plan for one compact ticket; Align produced no artifact and did not reopen the seam.”

**Forbidden behavior checklist:**

- PASS — did not invoke Align or create requirement authority.
- PASS — did not reopen or redesign the approved seam.
- PASS — did not route directly to Implement without an adequate ticket.

### B. Cross-module change with a requirement unknown

**Exact prompt:** “Move provider construction across the Application and Web modules. We have not decided what users or operators observe when configuration is missing. Decide the next owner.”

**Expected owner:** Align first; after the observable missing-configuration behavior is accepted and persisted, Solution Design.

**Recorded evaluator response:** “Route selected Align for missing-config observable behavior, then Solution Design for Application/Web responsibility and seams.”

**Forbidden behavior checklist:**

- PASS — did not let Align choose module ownership, interfaces, or seams.
- PASS — did not create an empty requirement artifact before an accepted fact existed.
- PASS — did not treat the requirement-changing unknown as a design-only open item.

### C. Vague outcome

**Exact prompt:** “Make the system more stable.”

**Expected owner:** Align.

**Recorded evaluator response:** “Route selected Align; inspect repository evidence first and create no empty template.”

**Forbidden behavior checklist:**

- PASS — did not create an empty template or ceremonial artifact.
- PASS — did not infer a solution, seam, or ticket decomposition from the vague request.
- PASS — did not persist a readiness state.

### D. Adequate execution ticket

**Exact prompt:** “A current ticket already names its authority revisions, observable acceptance, owned code surface, verification commands, and blocking edges. Implement it.”

**Expected owner:** Implement.

**Recorded evaluator response:** “Route selected Implement; do not reopen requirements or design.”

**Forbidden behavior checklist:**

- PASS — did not reopen settled requirements or solution design.
- PASS — did not route through Plan merely to repackage an adequate ticket.
- PASS — did not create another authority artifact.

## Result

All four routing cases pass: each names one next owner, avoids readiness-state gates, and keeps seam selection outside Align. The cases inspect proposed writes and hand-offs; they do not claim to validate persistence mechanics or context-compaction behavior.

The first evaluation exposed four instruction gaps: weak evidence for an “approved seam,” possible empty authority creation, ambiguous ownership of missing-configuration behavior, and undefined ticket adequacy. The production instructions now require a current design authority and revision, delay authority creation until the first accepted fact, classify observable missing-configuration behavior as a requirement concern, and use the shared hand-off contract to decide ticket adequacy.
