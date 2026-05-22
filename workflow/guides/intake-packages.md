# Intake Packages

An input package is a project-local, pre-exploration artifact that collects the
initial durable material for a topic before an exploration exists.

## Goals

- capture the first readable version of the request
- allow incremental updates as new information appears
- keep project-specific intake shape lightweight and extensible
- give exploration a deterministic bundle to read before it writes a skeleton

## Recommended layout

```text
docs/intake/<slug>/
├── index.md
├── references.md
├── questions.md
└── assets/
```

Projects may add or remove files when their intake needs differ. The workflow
core only requires that the package be discoverable and readable.

## Reading expectations

- `index.md` is the primary reading surface.
- `references.md` lists links, screenshots, and external resources.
- `questions.md` records open questions that should be answered during
  exploration.
- `assets/` holds images or other non-Markdown material.

## Update behavior

- input packages may evolve during exploration
- updates should be reflected in timestamps or revision notes inside the
  package
- exploration must re-read the updated material rather than relying on the
  previous snapshot

## Helper commands

- `scripts/flowforge-create-intake.js` creates the package scaffold.
- `scripts/flowforge-intake-context.js` materializes the current package into a
  deterministic context block.
- `scripts/flowforge-explore-context.js` combines the project rules bundle and,
  when present, the intake package into the exploration seed context.

## Relationship to rules

- project rules define how to interpret the package
- the intake package itself carries the request evidence
- exploration uses both to generate the initial skeleton and cite its sources
