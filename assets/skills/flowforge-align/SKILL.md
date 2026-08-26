---
name: flowforge-align
description: Clarify and persist why a feature exists, its observable outcomes, scope, scenarios, constraints, terminology, and requirement-changing unknowns. Use before solution design when different requirement answers would change the solution space.
disable-model-invocation: true
---

# Align requirement authority

Own requirement truth without choosing implementation modules, interfaces, seams, migration, or ticket slices.

Use the shared contract's [authority roles](../_shared/ARTIFACT-CONTRACT.md#roles-and-authority), [source intake and semantic rewrite](../_shared/ARTIFACT-CONTRACT.md#source-intake-and-semantic-rewrite), [hand-offs](../_shared/ARTIFACT-CONTRACT.md#hand-offs), and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value). Read [schema v1](../_shared/SCHEMA-V1.md) only when promoting requirements into an independently consumed artifact.

## Process

### 1. Locate the authority

Find the feature's current requirement authority and applicable project context/ADRs. Do not create a new authority file before the first accepted requirement fact exists. For a small change whose requirement can remain locally understandable in one compact ticket, keep it there. Create or promote `requirements.md` only when requirement information needs independent review, multiple consumers, a longer lifecycle, or protection from context compaction.

### 2. Resolve repository facts first

Inspect code, configuration, existing behavior, and durable project authority before asking the user. Ask only for product intent, trade-offs, or external facts the repository cannot answer. For example, whether missing configuration fails startup or selects a default is a requirement question when it changes observable behavior; the internal error-handling seam remains Solution Design's choice.

### 3. Work the requirement frontier

Invoke `flowforge-grilling` for the current branching requirement decisions. Invoke `flowforge-domain-modeling` only when a term needs durable glossary authority or the choice qualifies for an ADR.

Resolve or explicitly scope:

- the problem and evidence;
- observable outcomes and acceptance behavior;
- in-scope and out-of-scope behavior;
- representative and boundary scenarios;
- externally meaningful constraints and terminology;
- unknowns whose answers change the solution space.

These are completeness checks, not mandatory headings. Omit empty roles and merge overlapping prose.

### 4. Persist accepted decisions immediately

After each settled requirement decision, edit the authority directly. Replace superseded meaning instead of appending interview history. Increment its semantic revision only when meaning changes. Keep machine identity in metadata and human-readable meaning in prose and semantic links.

When this edit creates or revises schema requirement authority, run `flowforge check --dir <feature-dir> --strict` before reporting publication. Repair diagnostics introduced by this edit; this is a document check, not a readiness phase or a reason to create tickets.

An unresolved requirement fact uses a scoped [schema v1 open item](../_shared/SCHEMA-V1.md). Name affected requirement/design areas; do not create a feature-wide readiness state.

### 5. Hand off by unresolved owner

- A local change that clearly reuses a seam identified by current design authority and revision may proceed to Plan as a compact ticket.
- Responsibility, interface, seam, cross-module flow/order, migration, or verification-strategy decisions go to `flowforge-solution-design` with the requirement authority link/revision.
- Missing primary-source facts may detour through Research, then return to the same requirement question.

## Completion

Return the authority link and revision, requirement decisions changed this run, scoped unknowns/waivers, affected areas, and the next owner. Every requirement question that changes the solution space is resolved, explicitly scoped, or waived with a reason.

Keep `CONTEXT.md` glossary-only and ADRs sparse. MUST NOT persist `requirements-ready`, choose implementation architecture, create ticket decomposition, or pad authority with template prose.
