# Authoring Rules

## General

- Prefer concise, decision-oriented writing over transcript dumps.
- Every durable statement should be attributable to a finding, decision, or explicit assumption.
- Use relative links between artifacts to keep references navigable in Git.

## Explorations

- `index.md` is the reading surface.
- `journal/` preserves chronology.
- `findings/` contains atomic statements worth reusing.
- `decisions/` contains candidate decisions with status.
- Declare `ownership`, `expected_size_class`, and `reusable_rules` in `index.md` once the question is scoped. These propagate into the resulting proposal.
- Mirror the ownership graph in human-readable form inside `index.md`: explicitly summarize owning modules, system or architecture targets, and reusable conventions instead of relying only on raw `type:target` lines.
- Use the archived knowledge base as the first source of truth for exploration research, including `docs/conventions/`.
- Treat proposals and explorations as delta records against the existing final corpus, not as replacements for it.
- Do not mix implementation logs into exploration files.

## Proposals

- `meta.yaml` is the proposal bundle manifest.
- Each Markdown artifact carries its own YAML frontmatter for Obsidian-style indexing and doc-local routing.
- `proposal.md` answers why and what, and surfaces `size_class`, `ownership`, and any promoted `reusable_rules` in its frontmatter and body.
- Human-readable proposal docs must summarize the ownership graph explicitly:
  - which module docs this work belongs to
  - which system or architecture docs it affects
  - which reusable conventions it introduces or updates
- Design lives in either `design.md` or a `design/` directory depending on `size_class`. See `workflow/guides/sizing.md`.
- Template customization follows the reference-copy rule in `workflow/guides/templates.md`: projects may copy the whole template or the relevant part files and edit the copies directly.
- `task-map.md` is authoritative for task decomposition.
- `task-map.md` must follow `task-splitting.md` for deliverable-first decomposition, milestone boundaries, and checkpoint rules.
- `notes.md` is for execution history only.
- Proposals should begin by reviewing the canonical corpus and then describe the delta from that corpus.
- `canonical_corpus` in `meta.yaml` records which final docs were reviewed as the baseline for the proposal.
- Proposal creation may infer `canonical_corpus` from the declared archive targets and same-type final docs in the workspace; explicit overrides are allowed for broader baselines.
- Any manually supplied `canonical_corpus` entry must resolve to an existing document in the corresponding workspace.
- When a proposal changes existing final docs, describe the merge surface explicitly:
  - what section is updated in place
  - what facts are appended as new material
  - what facts are replaced or deprecated
  - what history note or changelog entry preserves the old fact
- For large modules, keep one canonical overview doc as the reader entry point and split dense details into subdocs such as `design.md`, `api.md`, `history.md`, or feature-specific pages.
- Do not duplicate the same fact across multiple final docs unless one of them is clearly marked as a pointer or historical record.
- If a proposal redistributes knowledge across docs, update the linked docs in one archive pass so the reader path stays coherent.

### Design surface by size class

- `small`: single-file `design.md`. Optional sections may be omitted.
- `medium`: default is single-file `design.md`. The proposal may opt into the `design/` directory layout when the change spans multiple concerns. When the proposal introduces two or more business models, a `model/` directory is mandatory regardless of whether `design.md` stays single-file.
- `large`: `design/` directory is mandatory and must contain at minimum `README.md`, `architecture.md`, `model.md`, and `lifecycle.md`. The `model/` directory is mandatory and must contain one document per core business model.

### Business model documents

- One document per core business model when the `model/` directory is in use.
- Each model document describes: data structure, responsibilities, lifecycle, validation, referenced conventions, and links to related models.
- `model/README.md` lists every model in the proposal and groups them by role (core configuration, lifecycle, view-facing helper, etc).
- Model docs should state their owning module and any related convention targets as part of the model identity block.
- The model template is split into readable parts so projects can customize the data-structure section, including project-specific columns such as `Master table`, by copying the template or the relevant part file.

### Convention authoring

- Conventions live in `docs/conventions/<topic>.md` and are not embedded only inside module or architecture docs.
- A proposal that introduces a convention must declare a matching `ownership` entry of type `convention` and a matching `archive_target` of type `convention`.
- A convention document follows the canonical template under `workflow/templates/docs/conventions/convention.md`.

## Decisions

- Use ADRs only for stable, high-cost decisions with meaningful alternatives.
- Draft decisions belong in explorations or proposals until they are accepted.

## Archive targets

- Every proposal must declare at least one primary archive target.
- Primary target is where a future reader should start.
- Secondary targets exist to preserve alternate reading paths.
- Ownership entries and archive targets must align: every ownership entry should have a corresponding archive target.
- The archived target corpus is the default baseline for future explorations, so archive targets should be kept navigable and up to date.
