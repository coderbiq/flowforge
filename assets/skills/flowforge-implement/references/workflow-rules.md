# Workflow Rules

Use one primary mode per activation.

## Token-Aware Execution Loop

1. Get step context: `context feature --feature <id> --step <n>` (~400 tokens)
2. Implement the step as described in Approach
3. Handle Edge Cases listed in the step
4. Write tests per Verification criteria
5. Record progress: `card steps <id> --status done <n>` + `card log <id> --event "..."`
6. Validate: run tests + `flowforge validate all`
7. Next step or complete: `card evolve <id> --stage done` when all steps done

## Reading Strategy

| Need | Command | Tokens |
|------|---------|--------|
| Step execution context | `context feature --feature <id> --step <n>` | ~400 |
| Design rationale | `card read --section "Design.Key Decisions"` | ~150 |
| Only constraints | `card read --section "Constraints"` | ~100 |
| Feature overview | `card read --summary <id>` | ~100 |
| Proposal health | `proposal inspect <id>` | ~200 |

Never read the whole FEATURE card during step execution.

## Design Issue Protocol

If implementation reveals a design issue:
1. `card steps <id> --status blocked N --reason "..."`
2. `card log <id> --kind blocked --event "..."`
3. Update Open Questions in the FEATURE card
4. Switch to design mode: modify Design/Constraints
5. Update affected Implementation Plan steps
6. `card steps <id> --status in_progress N` to resume

## Finish

When the implementation is complete, report:
- changed files
- tests run
- validation result
- gaps or blockers
- one next step
