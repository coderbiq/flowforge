# workflow-core Design

## Current shape

`workflow-core` lives in `scripts/lib/flowforge.js` and is shared by the command wrappers under `scripts/`.

The runtime currently handles:

- project root detection
- tool root detection
- configuration loading and defaulting
- workspace enumeration and lookup
- proposal skeleton creation
- proposal validation
- archive target rendering

## Dependencies

- `workflow/guides/`
- `workflow/schema/`
- `workflow/templates/docs/`
- `docs/` as the canonical corpus

## Invariants

- `docs.default_workspace` must always resolve to a declared workspace, or be synthesized as `default`
- proposal metadata uses relative refs, not repo-absolute paths
- canonical corpus entries must refer to real documents in the workspace
- archive updates should preserve historical facts when replacing existing content
- workspace resolution should prefer explicit input over inferred scope

