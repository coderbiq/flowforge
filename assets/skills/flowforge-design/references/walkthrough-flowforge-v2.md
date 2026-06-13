# FlowForge v2 Walkthrough

This walkthrough uses FlowForge itself as the design target.
It is an example, not a new rule document.

## Scenario

User input:

> Start designing the `flowforge-design` skill so an agent can move from a vague request to a requirement index tree, analysis tasks, design cards, and implementation tasks.

## Step 1: Inspect the work surface

```bash
flowforge project current
flowforge proposal current
flowforge proposal inspect <proposal-id>
flowforge context proposal --proposal <proposal-id>
```

What to look for:
- root card
- top-level STR / requirement index
- active analysis tasks
- open questions
- not_ready or blocked implementation tasks

## Step 2: Seed the requirement index tree

Create or update the top-level STR with a few themes:
- trigger boundary
- requirement indexing
- analysis task flow
- library discovery
- design and implementation handoff

If one theme grows too large, split it into a child index instead of stuffing more entries into one card.

## Step 3: Split atomic requirements

Examples:
- the skill must use CLI, not direct file reads
- the skill must start from proposal inspect and context proposal
- the skill must end each turn with cards, relations, gaps, and next step

Each requirement should be testable and small enough to verify on its own.

## Step 4: Create an analysis task for uncertainty

Example analysis task:

```markdown
# Decide the minimum CLI surface for the design skill

## Goal
Confirm which commands are required for the first deployable version.

## Inputs
- Requirement index
- Current docs
- Existing CLI behavior

## Investigation Plan
- inspect proposal and context commands
- inspect library discovery commands
- compare task readiness rules

## Expected Outputs
- design update
- finding
- implementation task

## Done When
The minimum command set is stable enough to write the skill body.
```

## Step 5: Discover library context

```bash
flowforge library suggest --for REQ-xxx --types convention,module,design,finding --limit 10
flowforge card search "thin adapter design skill" --scope library --type convention,design --limit 10
flowforge card read CONV-001 --section Rules
```

Use the results to decide which cards are actually relevant.
Do not read the whole library.

## Step 6: Create a design card

Example design card:

```markdown
# flowforge-design uses a thin adapter and a card-growth workflow

## Goal
Keep the skill short while still driving proposal analysis.

## Decision
The skill starts with inspect/context, then grows requirement, analysis, design, and implementation cards through CLI.

## Rationale
This avoids long proposal text and keeps the work traceable.

## Constraints
- CLI only
- no direct library file reads
- no long proposal docs in the skill body

## Impact
- skill body
- reference templates
- library discovery flow

## Verification
Run a design turn on FlowForge itself and confirm the card network grows.

## Follow-up Tasks
- create the skill body
- create card templates
- create library discovery rules
```

## Step 7: Create an implementation task only when ready

Ready implementation task example:

```markdown
# Create the deployable flowforge-design skill assets

## Goal
Produce the deployable skill body and reference files.

## Inputs
- requirement cards for trigger boundary and output format
- the design card above
- library discovery constraints

## Deliverables
- `assets/skills/flowforge-design/SKILL.md`
- `assets/skills/flowforge-design/references/card-templates.md`
- `assets/skills/flowforge-design/references/library-discovery.md`
- `assets/skills/flowforge-design/references/walkthrough-flowforge-v2.md`

## Acceptance
- skill body stays short
- templates stay in references
- library access stays CLI-only
- walkthrough covers the full turn

## Out of Scope
- Go CLI implementation
- docs/ changes

## Read Before Work
- requirement cards
- design card
- library discovery findings
```

If the design is still missing constraints, keep the task `not_ready`.

## Step 8: End the turn

The final response should list:
- cards added or updated
- relations added
- open gaps
- next step

That output is the user-facing proof that the design turn is complete.
