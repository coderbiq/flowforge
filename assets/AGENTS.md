<!-- FLOWFORGE:START -->
## FlowForge (v4 Working Memory & Living Docs)

FlowForge is an engineering AI collaboration harness powered by a **Multi-Tier Working Memory System**, **Conversational Agility Skills**, and **Living Documentation Synthesis**.

### Working Memory System
- **Tier 1: Global Memory (`docs/CONTEXT.md`)**: Ubiquitous language, architectural constraints, and active proposals.
- **Tier 2: Proposal Scratchpad (`01-workspace/<proposal_id>/README.md`)**: Core objective, grilling consensus, explored facts, open questions, and actionable slices. Supports seamless cross-session recovery.
- **Tier 3: Slice Context**: Dynamic, minimal context for executing a single Tracer Bullet work item.

### Skills
| Phase | Skill | Role & Responsibility |
|:---|:---|:---|
| **Align** | `flowforge-align` | Conversational Grilling: clarify boundaries, challenge assumptions, capture consensus into Proposal Scratchpad |
| **Explore** | `flowforge-explore` | Fact-Finding: investigate existing codebase/data nuances, inject evidence (`file:line`) into Scratchpad |
| **Plan** | `flowforge-plan` | Decomposition: slice consensus into 3-6 minimal Tracer Bullets with mandatory test bindings |
| **Implement** | `flowforge-implement` | TDD Delivery: Red-Green-Refactor cycle against bound automated tests |
| **Diagnose** | `flowforge-diagnose` | Root-Cause Analysis: hypothesis-driven bug and regression diagnosis protocol |
| **Review** | `flowforge-review` | Non-blocking Adversarial Review: architectural drift, security, and cognitive load reduction |
| **Curate** | `flowforge-curate` | Living Docs Synthesis: extract ADRs, patch domain docs (`docs/domains/`), and archive completed proposals |

### CLI Support
- `flowforge memory init [--proposal <id>]` to initialize or refresh global/proposal working memory
- `flowforge context slice [--proposal <id>] --slice <n>` to extract minimal context for a single slice
- `flowforge curate diff [--proposal <id>]` to preview domain documentation diff before merge
- `flowforge status` for lightweight progress visibility

### Subagent Orchestration
- Use subagents when delegating specialized tasks.
<!-- FLOWFORGE:END -->
