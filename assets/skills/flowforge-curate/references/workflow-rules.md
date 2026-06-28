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
   flowforge library import --type finding --title "Curation Plan: <source>" --status active --tags "curation-plan" --body - --source-card <source-card> <<'EOF'
   ## 来源
   ...
   ## 计划条目
   ### 批次 1（条目 1-8）
   - [ ] CONV / title / STR-NAME / create
   ...
   EOF
   ```
   The body must list all items in batches of 5-10 with `- [ ]` checkboxes.

2. Generate a batch YAML manifest for the current batch, then execute via heredoc:
    ```
    flowforge card batch - -o json <<'EOF'
    cards:
      - ref: "str1"
        type: structure
        title: "Index Card Title"
        status: active
        body: |
          STR index card body.
        links:
          - "FIND-xxx:references"
      - type: convention
        title: "Convention Title"
        status: draft
        body: |
          Atomic card body here.
        links:
          - "FIND-xxx:references"
          - "@str1:indexes"
    EOF
    ```
   - `ref` creates a batch-local name for cross-references.
   - `@ref:indexes` links to a batch-local STR and automatically performs `structure add`.
   - Cards are pre-validated; all pass or none are written.
   - Use `-o json` to capture created card IDs.

3. After batch creation, update the plan card's progress section:
   ```
   flowforge card update <plan-card-id> --section "批次 1" --body - <<'EOF'
   - [x] CONV-xxx / title / STR-xxx / create
   ...
   EOF
   ```

4. Merge/skip items:
   - `merge`: `flowforge card read <target> --summary`, then `flowforge card update <target> --section "<section>" --body -` with heredoc
   - `skip`: record reason only

5. When all batches done: `flowforge index rebuild`

6. Report: `Batch N/M complete. Processed: X/Y. Say "continue" to process next batch.`

## Mode B Only: Wrap Up

After all batches complete: `flowforge proposal archive <proposal-id>`