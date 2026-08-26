---
name: flowforge-solution-design
description: Design how approved requirements will be realized when work changes module responsibility, an interface or seam, cross-module information flow or ordering, migration compatibility, or has multiple credible solutions. Route a local change that clearly reuses an existing seam directly to flowforge-plan.
---

# Solution design

Turn approved requirement authority into the implementation decisions Plan needs, without choosing ticket slices or changing production code.

Use the shared contract's [authority roles](../_shared/ARTIFACT-CONTRACT.md#roles-and-authority), [source intake and semantic rewrite](../_shared/ARTIFACT-CONTRACT.md#source-intake-and-semantic-rewrite), [packaging rules](../_shared/ARTIFACT-CONTRACT.md#packaging), [hand-offs](../_shared/ARTIFACT-CONTRACT.md#hand-offs), [diagnostics](../_shared/ARTIFACT-CONTRACT.md#diagnostics), and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value).

## Inputs

Resolve the requirement authority and revision, applicable project context and ADRs, current repository responsibilities and seams, and the design area the caller needs advanced. If only a conversation request exists, label the output provisional and return creation or approval of requirement authority to `flowforge-align`; do not invent an ID or revision and do not prevent reversible fact discovery.

Return a question to Align only when different answers change the externally observable outcome, scope, scenario, constraint, or meaning. Choosing a module owner, interface shape, seam, internal flow, migration order, or verification seam remains solution design work.

## Process

### 1. Establish the decision frontier

Identify which responsibility, interface, seam, flow, ordering, migration, or verification facts are unsettled. Do not inventory the whole codebase. A local change at an already approved seam returns directly to Plan with the reason this Skill adds no value.

### 2. Resolve facts before choices

Inspect code and project authority for the current model. Invoke `flowforge-research` only for a named unavailable primary-source fact. Invoke `flowforge-prototype` only for one behavior hard to judge on paper. Resume the same decision after the verdict and retain only evidence that changed it.

Use `flowforge-domain-modeling` only for durable vocabulary or an ADR. Invoke `flowforge-codebase-design` whenever this work selects or moves a module responsibility, interface, or seam; it advises and does not own the artifact. Proposal fixtures are examples, never requirement or design authority for the current feature.

### 3. Compare credible designs

For consequential seams, sketch at least two materially different responsibility boundaries. Compare caller knowledge, hidden complexity, failure handling, compatibility order, test seam, and change locality. Select one and retain the decision-rich reason alternatives lost.

When repository facts are settled but a real product or engineering trade-off leaves a branching decision frontier, invoke `flowforge-grilling` with the alternatives, consequences, and decision owner. Resume after the choice and write it to design authority; do not use grilling to compensate for missing facts or requirement ambiguity.

### 4. Maintain authority incrementally

After each settled decision, update design authority. Treat current/target responsibilities, interfaces, information flow, migration, scoped constraints, verification, and decision-rich alternatives as a completeness checklist—not mandatory headings. Omit any section with no independent information.

Use a stable semantic area and increment only its semantic revision. Formatting or link repair does not increment it. Read [schema v1](../_shared/SCHEMA-V1.md) when writing machine metadata. Read [adaptive design packaging](DESIGN-PACKAGING.md) only when deciding whether the design should remain compact or split.

After creating or revising schema design authority, run `flowforge check --dir <feature-dir> --strict`; repair diagnostics caused by the edit before reporting it published. The command validates present document facts and does not create a state gate.

### 5. Record scoped unresolved facts

An unresolved required fact uses the [schema v1 open-item shape](../_shared/SCHEMA-V1.md) with a semantic ID, diagnostic code, `gap` or `blocker`, explanation anchor, and exact affected authority areas; add ticket targets later when Plan creates consumers. Split areas narrowly enough that an unresolved branch does not contaminate settled work. One incomplete area never becomes a global not-ready state. A missing exact command is a warning when the observable seam is settled; it is a gap when no feasible observation method has been selected.

### 6. Derive planning coverage

For every affected area, report only status, authority link/revision, and any diagnostic: resolved, warning, gap, or blocker. The design authority—not the coverage report—owns the decision details. Planning proceeds for resolved areas. Override keeps the fact visible and does not resolve it.

### 7. Compress

Delete synonymous headings, implementation commonplaces, exploration that changed no choice, copied project constraints, and repeated requirement prose. Every retained paragraph carries a design fact, decision, constraint, alternative, verification method, or named unknown.

### 8. Classify implementation returns

- Apply a factual correction to design authority, record its source, and revise semantics only when meaning changed.
- Let a ticket refine a local detail only when it preserves approved responsibility, interfaces, seams, DAG order, constraints, and acceptance behavior.
- Resume solution design when the discovery changes any of those approved facts; update affected authority and coverage before planning or implementation continues.

## Completion

Return authority links, scoped diagnostics, areas Plan may consume, supporting verdict links, semantic titles/links for decisions added or changed this run, and material compression actions. Do not repeat decision bodies or assign a synthetic quality score. Every affected area has an explicit coverage result; every area recommended to Plan passes the responsibility/interface/seam/migration/verification completeness checklist applicable to it.

MUST NOT persist `design-ready`, publish tickets, change production code, or silently resolve requirement ambiguity.
