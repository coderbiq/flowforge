# Standards extraction guide

This document tells Plan and Review how to find the project standards that apply to a given ticket. Edit it to match your project's actual standard documents and extraction logic.

## Where project standards live

<!-- Replace the list below with your project's actual standard documents. Each entry should name the file and what it covers. -->

- `CODING_STANDARDS.md` — general coding conventions (if your repo has one).
- `CONTRIBUTING.md` — contribution and review conventions (if your repo has one).
- Architecture or layering rules — wherever your project keeps them (e.g. `docs/` directory, ADRs, or a dedicated standards folder).

## How to determine which standards apply to a ticket

<!-- This section describes the extraction logic. The default below is a generic heuristic: derive the layer or module from the ticket's Touch points and Write set, then look up standards that govern that layer or module. Replace it with whatever logic fits your project. -->

1. Read the ticket's **Touch points** and **Write set**. They name the concrete files and directories the ticket will modify.
2. Identify the layer, module, or area those files belong to. Common signals:
   - path segments (e.g. `domain/`, `application/`, `infrastructure/`, `api/`, `ui/`)
   - module or package names
   - the kind of change (new endpoint, data model change, test, config, etc.)
3. Consult the standard documents listed above. For each document, determine whether any rule inside it governs the layer, module, or area identified in step 2.
4. For each applicable rule, extract a `must` or `must not` statement that captures what the implementer must do or avoid. Keep it specific to this ticket's work — do not copy entire sections.
5. Link each extracted statement to its source with a relative path and anchor, e.g. `../docs/dependency-rules.md#section-anchor`.

## Where to write extracted standards in the ticket

Extracted standards go into the existing ticket tiers, based on the nature of each rule:

- **Hard invariants** (violation means failure, affects completion) → `## Constraints`, alongside ticket-specific invariants.
- **Conventions** (the implementer must follow but softer, not a direct completion gate) → Tier 3 `### Conventions`.

Each extracted standard uses this format:

```
- must <required behavior> — <source document link>
- must not <forbidden behavior> — <source document link>
```

## Extraction status markers

After extraction, the ticket must carry one of these states so the implementer's pre-flight check can proceed:

- `must`/`must not` statements present → extraction done, proceed.
- `standards: none found per guide` → extraction attempted but no applicable standards found for this ticket; proceed.
- `standards: pending` → extraction not yet done; implementer returns the ticket to Plan.

Tickets without a Write set (pure documentation changes) do not require extraction markers.

## Customizing this guide

<!-- Edit the two sections above to match your project. You can organise standards by layer, module, scenario, dependency graph, or any structure that fits. FlowForge does not prescribe the internal structure of this guide — it only requires that Plan can read it and extract `must`/`must not` statements with source links. -->
