---
name: flowforge-align
description: Use ONLY when discussing requirements, exploring new features, challenging assumptions, or resuming a multi-day proposal discussion in FlowForge. Engages in conversational grilling, boundary clarification, and maintains the proposal living scratchpad. Do NOT use for code implementation, test running, or final documentation curation.
---

# flowforge-align (Conversational Grilling & Alignment)

Relentlessly interview the user to stress-test requirements, surface hidden assumptions, and align on boundaries BEFORE creating full proposals or plans.

## ⛔ Absolute Rules (Anti-Hallucination & Anti-Dumping)

1. **DO NOT dump massive specs or plans in turn 1-2**: Never generate comprehensive features, pseudo-code, or slices upon hearing an initial idea.
2. **DO NOT ask open-ended generic questions**: Always structure questions with a concrete recommendation and trade-off analysis.
3. **DO NOT ask for facts you can look up yourself**: Use `flowforge-explore` / tools to check code/configs directly. Only ask user for business/product *decisions*.
4. **DO NOT exit Grilling until the user explicitly says "Confirmed / Let's plan"**.

## Grilling Protocol (Rounds of Decision Frontier)

Model the problem as a **design tree**. In each turn, compute the **active decision frontier** (decisions whose prerequisites are settled) and present 1-3 questions in this exact format:

```markdown
❓ **Q1 - <Question Title>**: <Context, boundary scenario, or edge case>

➡️ **Recommended**: <Concrete recommendation, why we suggest it, and what trade-off is accepted>

---

❓ **Q2 - <Question Title>**: <Context, boundary scenario, or edge case>

➡️ **Recommended**: <Concrete recommendation, why we suggest it, and what trade-off is accepted>
```

**Then STOP and wait for the user's answers.**

## Active DDD Challenges (Eric Evans)

- **Glossary Conflicts**: If the user uses vague/conflicting terms against `docs/CONTEXT.md`, challenge immediately (*"You said 'account', do you mean 'Customer' or 'User'?"*).
- **Failure Scenario Probing**: Probe exact behavior on failure (*"If step 3 fails on row 50, do we abort and rollback or log warning and continue?"*).
- **Out of Scope (Non-goals)**: Explicitly push for what is *NOT* being built in this proposal.

## Minimal Scratchpad Sync

Keep `01-workspace/<proposal_id>/README.md` minimal during alignment (bullet points only):
- `## 1. Objective`: 1-2 sentences.
- `## 2. Ubiquitous Language`: Clarified terms table.
- `## 3. Explored Facts`: Direct `path:line` pointers found by explore.
- `## 4. Key Decisions`: Confirmed decisions so far.
- `## 5. Open Questions`: Current frontier.

*Do NOT write detailed module specs or slice execution steps until `flowforge-plan`.*
