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
- Output lives under `docs/explorations/<slug>/`.

### Propose

- Convert a validated exploration into a decision-ready proposal.
- Define scope, success criteria, capabilities, constraints, and archive targets.
- Output lives under `docs/proposals/<proposal-id>/`.

### Approve

- Lock the chosen approach before work starts.
- Confirm task backend and archive destinations.
- Proposal status changes from `proposed` to `approved`.

### Apply

- Create tasks from `task-map.md` and transition the proposal into execution.
- Proposal status changes from `approved` to `active`.

### Implement

- Execute tasks, keep `notes.md` current, and write back scope changes into the proposal.
- Major scope changes require returning to exploration/proposal, not just implementation notes.

### Archive

- Verify tasks are complete.
- Update the primary archive target and any secondary targets.
- Close the proposal with status `archived`.

## Artifact map

- Exploration creates findings and candidate decisions.
- Proposal references exploration outputs.
- Task map bridges proposal capabilities to the task backend.
- Notes capture execution history.
- Archive updates modules, architecture views, and ADRs as needed.
