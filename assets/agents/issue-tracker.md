# Issue tracker: Local Markdown (with FlowForge DAG Engine)

Issues, specs, and wayfinder maps live as markdown files in `.scratch/`.

## Conventions
- One feature per directory: `.scratch/<feature-slug>/`
- The spec is `.scratch/<feature-slug>/spec.md`
- Implementation tickets are at `.scratch/<feature-slug>/issues/<NN>-<slug>.md`
- State is recorded as a `Status:` line (`open`, `claimed`, `resolved`)
- Blocking is recorded as a `Blocked by: NN, NN` line

## Publishing & Editing (Use File Tools Directly)
- **To publish a spec/ticket**: Use your file writing tool to create the `.md` file directly.
- **To update status/answer**: Use your file editing tool to update `Status:` or append `## Comments`.

## Frontier & Dependency Calculation (Use FlowForge CLI)
- **To find the next unblocked ticket**: Run `flowforge frontier`. It instantly resolves DAG dependencies across `.scratch/` and outputs the ready tickets, saving context tokens and eliminating graph hallucinations.
- **To validate ticket plans**: Run `flowforge check` to ensure no circular dependencies exist.
