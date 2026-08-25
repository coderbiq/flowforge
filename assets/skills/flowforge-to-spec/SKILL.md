---
name: flowforge-to-spec
description: Create an optional compact navigation and review baseline over existing requirement, design, ticket, verification, and gap authorities. Use for multi-session work, external review, or when a feature needs one entry point.
disable-model-invocation: true
---

# Synthesize navigation

Use the shared contract's [authority roles](../_shared/ARTIFACT-CONTRACT.md#roles-and-authority), [hand-offs](../_shared/ARTIFACT-CONTRACT.md#hand-offs), and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value). Read [schema v1](../_shared/SCHEMA-V1.md) when publishing metadata.

This Skill synthesizes accepted content without interviewing or deciding. The result is optional navigation: requirement and design authorities remain authoritative, tickets remain executable, and evidence owns results.

## Process

### 1. Decide whether an overview earns its cost

Create one for a multi-session feature, external review package, or feature whose separate authorities need a stable entry point. Skip it when one compact ticket is the clearest reading interface. State the skip reason and link the compact artifact.

### 2. Resolve sources and gaps

Read the effective requirement and solution-design authorities, current tickets, verification/evidence navigation, and scoped open items. Inspect the repository only to resolve current artifact locations; do not invent a missing requirement, interface, seam, migration decision, or verification strategy.

If a missing decision prevents an accurate overview, stop synthesis and return its owning Skill. If the authorities already record a scoped open item, link that gap and its affected area rather than silently completing it.

### 3. Publish the compact entry point

Write `<docs_dir>/proposals/<feature>/spec.md` with schema v1 `role: spec` and consumed authority revisions. State near the title that it is non-executable navigation, not requirement or design authority.

Include only information that helps a reader navigate:

- one concise feature purpose and scope boundary;
- semantic links to requirement and design authorities with their reviewed revisions;
- a short list of key decisions, each expressed as a local reading cue and authority link rather than copied rationale;
- links to ticket/DAG and verification or evidence entry points;
- linked remaining gaps and their affected areas.

Omit empty roles, extensive user-story inventories, raw conversation history, and duplicated authority sections. Human link labels state meaning; machine IDs stay in metadata.

### 4. Validate optionality

Run Catalog/DAG checks. Confirm the spec is not executable and does not change issue count or frontier. Apply the deletion test: removing only this overview must leave requirement/design authority, ticket consumption, DAG edges, and evidence traceability intact.

## Completion

Return the spec link and reviewed authority revisions, or the explicit skip/owner result. A new reader can reach scope, key decisions, execution, verification, and current gaps without treating the overview as a second authority.
