# Lifecycle

`FlowForge` uses a single lifecycle across all platforms:

1. `explore`
2. `propose`
3. `approve`
4. `apply`
5. `implement`
6. `archive`

## Phase intent

### Explore

- Capture the problem, context, evidence, unknowns, and candidate directions.
- Produce durable findings before any implementation starts.
- Start from the existing final knowledge base first: modules, architecture docs, and ADRs are the default reference corpus.
- Open a new exploration only for gaps, conflicts, or materially new questions that the final corpus does not already answer.
- Output lives under `docs/explorations/<slug>/`.

### Propose

- Convert a validated exploration into a decision-ready proposal.
- Define scope, success criteria, capabilities, constraints, and archive targets.
- Treat the archived knowledge base as the baseline; proposals should describe deltas against existing canonical docs, not rewrite the corpus from scratch.
- Output lives under `docs/proposals/<proposal-id>/`.

### Approve

- Lock the chosen approach before work starts.
- Confirm task backend and archive destinations.
- Proposal status changes from `proposed` to `approved`.

### Apply

- Create tasks from `task-map.md` and transition the proposal into execution.
- Task decomposition must follow `task-splitting.md`, including milestone boundaries and checkpoint rules for long-running work.
- Proposal status changes from `approved` to `active`.

### Implement

- Execute tasks, keep `notes.md` current, and write back scope changes into the proposal.
- Stop at declared checkpoints for verification when a proposal spans multiple sessions or days.
- Major scope changes require returning to exploration/proposal, not just implementation notes.

### Archive

- Verify tasks are complete.
- Update the primary archive target and any secondary targets.
- Treat archive as a knowledge-base maintenance pass, not a terminal dump.
- When existing final docs change, preserve the old fact in history or changelog sections so the corpus remains traceable.
- Keep the overview and linked subdocs in sync so readers can still navigate the full system from the canonical entry point.
- Close the proposal with status `archived`.

## Artifact map

- Exploration creates findings and candidate decisions.
- Proposal references exploration outputs.
- Task map bridges proposal capabilities to the task backend.
- Notes capture execution history.
- Archive updates modules, architecture views, and ADRs as needed.
- The archived knowledge base becomes the default target corpus for later explorations.
- Future proposals should start from the updated canonical corpus and record deltas against it instead of re-litigating already settled facts.
