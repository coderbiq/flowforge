---
name: flowforge-align
description: Use ONLY when discussing requirements, exploring new features, challenging assumptions, or resuming a multi-day proposal discussion in FlowForge. Engages in conversational grilling, boundary clarification, domain modeling, and maintains working memory (Flat or Hierarchical Mode). Do NOT use for code implementation, test running, or final documentation curation.
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

## Active DDD & Seam Challenges

- **Domain Glossary Conflicts**: If user uses vague/conflicting terms against `docs/CONTEXT.md`, challenge immediately (*"You said 'account', do you mean 'Customer' or 'User'?"*).
- **Physical Seam Boundaries**: For refactoring/new modules, clarify where the seam lies (*"Will `module-a` call `module-b` directly, or via a newly extracted SPI in `contracts/`?"*).
- **Out of Scope (Non-goals)**: Explicitly push for what is *NOT* being built in this proposal.

## Working Memory Modes: Flat vs. Hierarchical

Choose mode based on task complexity (from `flowforge-triage`):

### Mode A: Flat Mode (Level 2 Vertical Feature)
Store in single file `01-workspace/<proposal_id>/README.md`:
- `## 1. Objective`: 1-2 sentences.
- `## 2. Ubiquitous Language`: Clarified terms table.
- `## 3. Explored Facts`: Direct `path:line` pointers found by explore.
- `## 4. Key Decisions & Consensus`: Confirmed decisions so far.
- `## 5. Open Questions`: Current frontier.

### Mode B: Hierarchical Mode (Level 3 Refactor / Level 4 Greenfield)
When touching $\ge 2$ modules, maintain top-level master `README.md` and module specs:
- `01-workspace/<proposal_id>/README.md`: Overall objective, topology diagram, global invariants, and module index.
- `01-workspace/<proposal_id>/modules/<module>.md`:
  - **Module Purpose & Seams**: Exact public interfaces and visibility rules (`internal` vs `public`).
  - **Physical Move Matrix (for Refactor)**: `SourceClass.kt -> TargetClass.kt` mapping.
  - **Dependency Changes**: Gradle/Maven/Package manifest modifications.

## 🔄 Mandatory Dual-Way Memory Anchoring (Every Turn)

At the end of **every conversation turn**:
1. **Anchor Artifacts**: If any file was created or investigated in the workspace, append its path and a 1-line summary to `README.md` under `## 3. Explored Facts & Artifacts Index`.
2. **Ingest Consensus**: Record confirmed decisions under `## 4. Key Decisions & Consensus`.
3. **Keep Resume-Ready**: The working memory must reflect the cumulative state so any subsequent agent session can resume seamlessly.
