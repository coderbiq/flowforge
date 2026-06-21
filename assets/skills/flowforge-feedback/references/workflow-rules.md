# Workflow Rules

This reference defines the complete feedback SKILL turn loop, from receiving a
discovery to routing it to the correct output (task card, log card, or library
card).  Every step references the classification rules in
`classification-rules.md` — do not duplicate classification logic here.

## Turn Loop

Use one primary mode per turn.  Do not mix discovery, classification, and
routing in a single step.

```
1. Receive discovery
2. Classify (classification-rules.md)
3. Route (this document)
4. Record (log card)
5. Verify card state
```

## Step 1 — Receive Discovery

A discovery arrives as any of the following:

- A failed test or check reported during task execution.
- A behavior deviation noticed by the implementor.
- A cognitive update: a new constraint, convention, or decision learned from
  the codebase or from a human.
- A stakeholder request or review comment that reveals a gap.
- A code review finding from `flowforge-feedback` activation.

Collect the minimum context needed for classification:

- What was expected vs. what was observed.
- Which task, requirement, or design card is relevant.
- Whether the item is reproducible.

Do not write a card yet.  Do not attempt to classify before the full discovery
text has been read.

## Step 2 — Classify

Apply the decision tree in `classification-rules.md §Decision Tree`.  Every
discovery must be assigned exactly one of the five types.  If it genuinely
matches more than one type, pick the most actionable type first and record the
ambiguity in the log card.

The card type is set via the CLI `--type` flag when creating the card.

## Step 3 — Route

Routing is the mapping from type to card creation command and relation.

### 3.1 bug → task card (not_ready)

```bash
flowforge card create \
  --type task \
  --title "<concise bug title>" \
  --status not_ready \
  --body - <<'EOF'
# <bug title>

## Goal
What the fix must achieve.

## Inputs
- Link to the source discovery or failed task.

## Deliverables
- Changed file(s) or configuration.
- Updated test that verifies the fix.

## Acceptance
- The reproducing check passes after the change.

## Out of Scope
- Adjacent improvements not related to the root cause.

## Read Before Work
- <source card ID>
EOF
  --links "<source-card-id>:records"
```

After creation, add the `requires` link from the new task to the relevant
requirement or design card if not already set by the batch manifest.

### 3.2 finding → finding card (draft)

```bash
flowforge card create \
  --type finding \
  --title "<finding title>" \
  --status draft \
  --body - <<'EOF'
# <finding title>

## Summary
One-paragraph description of what was observed.

## Source
<where the finding came from: task, test, review>

## Evidence
- Concrete data points, log snippets, or test output.

## Impact
Why this matters and what it might affect.

## Open Questions
What is still unknown.
EOF
  --links "<source-card-id>:records"
```

### 3.3 knowledge → library card (import / promote)

```bash
# Option A — import from scratch (new knowledge not in any proposal card)
flowforge library import \
  --type convention \
  --title "<knowledge title>" \
  --body - <<'EOF'
## Summary

## Rule or Finding

## Applies When

## Source Evidence
EOF
  --source-card "<source-card-id>"

# Option B — promote an existing stable proposal card to library
flowforge library promote "<proposal-card-id>"
```

The library card type must be one of: `convention`, `decision`, `module`,
`finding`, `design`.  Do not create a library card without `--source-card` or
`--links`.

### 3.4 missing-requirement → requirement card (draft)

```bash
flowforge card create \
  --type requirement \
  --title "<requirement title>" \
  --status draft \
  --body - <<'EOF'
# <requirement title>

## Summary
What is missing and why it matters.

## Source
<discovery context>

## Acceptance
- Testable condition 1
- Testable condition 2

## Scope
What is in scope and what is out of scope.

## Open Questions
None (or list open questions)
EOF
  --links "<source-card-id>:records"
```

After creation, add the requirement to the proposal STR index:

```bash
flowforge structure add "<STR-index-id>" "<new-requirement-id>"
```

### 3.5 design-flaw → requirement card (design change request, draft)

```bash
flowforge card create \
  --type requirement \
  --title "<design flaw: short title>" \
  --status draft \
  --body - <<'EOF'
# <design flaw title>

## Summary
What is wrong with the current design and why.

## Source
<affected design card ID>

## Acceptance
- The updated design must resolve the structural risk.
- No new regression is introduced.

## Scope
What aspects of the design are affected.

## Open Questions
None (or list open questions)
EOF
  --links "<affected-design-card-id>:references"
```

Route the new requirement through the design SKILL to produce an updated design
card before any implementation task is created.

## Step 4 — Record (log card)

Every meaningful discovery produces a log card.  The log is process evidence,
not a replacement for requirement, design, or finding cards.

```bash
flowforge log create \
  --kind feedback \
  --title "<category>: <short description>" \
  --for "<source-task-or-discovery-card-id>" \
  --summary "<one-line summary of the discovery and its routing>"
```

`--for` links the log to the card that triggered the feedback.  Use `--kind
feedback` so the log can be filtered in `context proposal` output.

Log kinds and when to use them:

| Kind | Use when |
|------|----------|
| `feedback` | All 5-type discoveries |
| `progress` | Task state change (claim / done / block) |
| `analysis` | Analysis task conclusion |
| `decision` | Design decision recorded |

## Step 5 — Verify Card State

After routing, confirm the card network is consistent:

```bash
flowforge validate all
```

Required checks:

1. The new card passes validation (outbound link, frontmatter, body).
2. If a requirement was created, it is indexed in the proposal STR.
3. If a task was created, it links to at least one requirement via `satisfies`
   or `requires`.
4. No dangling `@ref` placeholders remain in any card body.

If validation fails, fix the card before moving to the next discovery.

## Batch Mode

When handling multiple discoveries from a single task execution, group them into
a single batch manifest:

```bash
flowforge card batch /tmp/feedback-batch.yaml
```

Batch manifest rules:

- Set `proposal` to the current proposal ID.
- Group by type: put all bug tasks together, all findings together, etc.
- Use `ref` for cross-references within the same batch.
- The `indexes` relation on a requirement card automatically performs
  `structure add`.

## Evidence Checklist (end of turn)

Before closing the feedback turn, confirm each of the following is present:

- [ ] Every discovery has a card (task / finding / requirement / library card).
- [ ] Every card has at least one typed outbound frontmatter link.
- [ ] Every discovery has a log card with `--kind feedback`.
- [ ] `flowforge validate all` passes with zero errors.
- [ ] No wiki files or `02-library/` files were edited directly.

If any item is unchecked, do not close the turn — fix it first.
