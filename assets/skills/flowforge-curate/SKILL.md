# flowforge-curate

Use when the user wants to import knowledge from external documents or archive a completed proposal into the library. Do not activate for single-card creation, proposal design, or task execution.

## Start

Check for an in-progress plan card (FIND with tag `curation-plan` and status `active`). If one exists, resume batch execution. Otherwise, determine the mode:

- **Mode A (external import)**: User provided a file path → read the source file(s).
- **Mode B (proposal archive)**: User provided a proposal ID → run `flowforge proposal inspect <id>` and `flowforge context proposal --proposal <id>`.

## Workflow

### Mode A: Extract Knowledge Units

1. Read the source file(s) completely.
2. Identify standalone knowledge assertions. Each must be: self-contained, singular (one assertion), addressable (can be independently referenced).
3. Restate each in your own words. Do NOT copy-paste from the source.
4. Determine the card type for each unit:
   - Constraint on future work → `convention`
   - Accepted technical/product decision → `decision`
   - Reusable fact or caveat → `finding`
   - Guiding principle → `convention`
   - Reusable design pattern → `design`
   - Module/system knowledge → `module`
5. Record the source location (section/paragraph) for each unit.

### Mode B: Extract Knowledge Units

1. Scan proposal cards using `flowforge proposal inspect <id>`.
2. Filter reusable candidates: `finding` cards, `decision` cards, `design` cards with reusable patterns. Skip: `log`, `requirement`, `task`, `ROOT`, `STR`.
3. Evaluate each candidate's knowledge type for the library.

### Shared: Cluster and Plan

1. **Cluster** knowledge units by concept (NOT by source structure). Each cluster becomes one STR index card.
2. **Generate review plan** — output as text, do NOT write any cards yet. Include:
   - Source info (file path or proposal ID)
   - Proposed STR index cards with titles
   - Proposed atomic cards (type, title, 2-3 sentence summary, target STR)
   - Duplicate/merge candidates (cards that may already exist)
   - Warnings for oversized or vague units
3. **Wait for user review** — do not proceed until user confirms.

### Batch Execution

After user approval:

1. **Create the plan card**: `flowforge library import --type finding --title "Curation Plan: <source>" --status active --tags "curation-plan" --body "<plan>" --source-card <source-card>`. The body must list all items in batches of 5-10, with `- [ ]` checkboxes.
2. **Execute one batch** (5-10 items). For each item:
   - `create`: `flowforge library import --type <type> --title "<title>" --status draft --body "<body>" --source-card <plan-card-id> --links <plan-card-id>` 
   - `merge`: `flowforge card read <target> --summary`, then `flowforge card update <target>` with appended content
   - `skip`: record reason only
3. Create STR index cards: `flowforge card create --type structure --title "<title>" --status active`
4. Add cards to STR indexes: `flowforge structure add <str-id> <card-id>`
5. Link related cards: `flowforge card link <from> <to> --relation references`
6. Update plan card: mark completed items as `- [x]` and update progress count
7. Report: `Batch N/M complete. Processed: X/Y. Say "continue" to process next batch.`
8. When all batches done: `flowforge index rebuild`

### Mode B Only: Wrap Up

After all batches complete: `flowforge proposal archive <proposal-id>`

## Hard Rules

- CLI is the only read/write path for cards.
- Never read wiki files directly (except source files for Mode A).
- Never hand-write card files, frontmatter, wikilinks, or internal card links.
- Always create `status: draft` cards. Promote to `active` only after review.
- Always generate a review plan before writing any cards.
- Always batch execution: 5-10 items per activation.
- Never skip the review step.
- Always use single quotes for `--body` content containing special characters.
- Use `--source-card` to link each created card to the plan card.
- The plan card tracks progress; read it on each activation to resume.

## Output

Report batch progress: completed items, remaining items, created card IDs, and next step ("continue" or "done").