# Template Usage

`FlowForge` templates are reference defaults, not an automatic override system.

## Core rule

- If a project only needs the standard shape, use the default templates as-is.
- If a project needs project-specific wording, extra columns, or extra sections, copy the whole template or the relevant part files into the workspace-local template area and edit the copies directly.
- Do not rely on line-level or section-level merge semantics. Template customization is explicit copy-and-edit work.

## Recommended workspace-local location

Workspace-local template copies should live under the workspace docs root, for example:

- `docs/flowforge/_templates/`

This keeps project-specific template variants visible to the team without turning them into a hidden tool configuration.

## Model templates

The model template is intentionally split into parts so projects can customize the data-structure section, add model-specific notes such as master-table columns, or copy the entire model template and adjust it as a unit.

Use the split parts when:

- the project only needs to tweak one section
- the agent needs a clear explanation of what each section does

Copy the whole template when:

- several sections need adjustment
- the project wants a heavily tailored model document shape
- the resulting document should read as a project-specific reference rather than a generic default

## Design templates

The split `design/` layout follows the same reference-copy rule as the model template.

Projects may:

- use the default `design/` files as-is
- copy one design section file and adjust it
- copy the whole `design/` directory and tailor it as a project-specific design surface

This is useful when the proposal needs project-specific wording for architecture, lifecycle, flow, API, constraints, or tradeoffs without changing the workflow core.

## Agent-facing guidance

Every template directory should include readable explanatory text so an agent can tell:

- what the template is for
- which parts are standard
- which parts are intended for project-specific adjustment
- when to copy the whole template instead of modifying a single part

Avoid silent template behavior. If the shape changes, explain it in the template README or the file header.
