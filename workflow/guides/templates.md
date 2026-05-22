# Template Usage

`FlowForge` templates are reference defaults, not an automatic override system.

## Core rule

- If a project only needs the standard shape, use the default templates as-is.
- If a project needs project-specific wording, extra columns, or extra sections, copy the whole template or the relevant section files into the workspace-local template area and edit the copies directly.
- Do not rely on line-level or section-level merge semantics. Template customization is explicit copy-and-edit work.

## Recommended workspace-local location

Workspace-local template copies should live under the workspace docs root, for example:

- `docs/flowforge/_templates/`

This keeps project-specific template variants visible to the team without turning them into a hidden tool configuration.

## Project seed rules

Projects also receive an install-time seed rules bundle, usually under
`docs/flowforge/_rules/`.

- Use the bundle as the editable project-default working policy.
- Keep it separate from the core workflow guides so projects can tune analysis
  and writing defaults without forking the platform rules.
- If a project needs different behavior, edit the copied rules bundle directly
  instead of patching the core workflow first.
- Adapters should load the bundle according to
  `workflow/guides/rule-loading.md`.

## Model templates

The default model template is a single document, `model.md`.

Use the single-file template when:

- the project wants the standard model shape
- the model can be described in one coherent reading surface
- the project wants customization to stay explicit and easy to review

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
