# Workflow Rules

Use one mode per activation.

## Mode A: External Import — Extract

1. Read the source file(s) completely.
2. Identify standalone knowledge assertions: self-contained, singular, addressable.
3. Restate in your own words. Do not copy-paste.
4. For card type, see `extraction-guide.md`.
5. Record source location (section/paragraph) for each unit.

## Mode B: Proposal Archive — Extract

1. Scan proposal cards with `flowforge proposal inspect <id>`.
2. Filter reusable candidates: `finding`, `decision`, `design` cards. Skip `log`, `requirement`, `task`, `ROOT`, `STR`.
3. Evaluate each candidate's knowledge type for the library.

## Shared: Cluster and Plan

1. **Cluster** knowledge units by concept (not by source). Each cluster becomes one STR index card.
2. **Generate review plan** as text output — do NOT write any cards yet. Include:
   - Source info (file path or proposal ID)
   - Proposed STR index cards with titles
   - Proposed atomic cards (type, title, 2-3 sentence summary, target STR)
   - Duplicate/merge candidates
   - Warnings for oversized or vague units
3. **Wait for user review** — do not proceed until user confirms.

## Batch Execution (After User Approval)

1. Create the plan card:
   ```
   flowforge library import --type finding --title "Curation Plan: <source>" --status active --tags "curation-plan" --body "<plan>" --source-card <source-card>
   ```
   The body must list all items in batches of 5-10 with `- [ ]` checkboxes.

2. Execute one batch (5-10 items). For each item:
   - `create`: `flowforge library import --type <type> --title "<title>" --status draft --body "<body>" --source-card <plan-card-id> --links <plan-card-id>`
   - `merge`: `flowforge card read <target> --summary`, then `flowforge card update <target>` with appended content
   - `skip`: record reason only

3. Create STR index cards: `flowforge card create --type structure --title "<title>" --status active`

4. Add cards to STR indexes: `flowforge structure add <str-id> <card-id>`

5. Link related cards: `flowforge card link <from> <to> --relation references`

6. Update plan card: mark completed items as `- [x]` and update progress count.

7. Report: `Batch N/M complete. Processed: X/Y. Say "continue" to process next batch.`

8. When all batches done: `flowforge index rebuild`

## Mode B Only: Wrap Up

After all batches complete: `flowforge proposal archive <proposal-id>`