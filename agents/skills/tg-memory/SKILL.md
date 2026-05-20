---
name: tg-memory
description: |
  Memory and work-restoration skill for tg-workflow.

  Load this skill when the user asks about previous decisions, current progress, pending review items, or storing reusable experience.
---

# tg-memory

`tg-memory` manages two distinct layers and must keep them separate.

## Layer 1: local work-restoration state

Purpose:

- recover the current working context quickly
- track the active proposal, touched files, current focus, and next actions

Storage:

- `.workflow/state/active-session.json`
- `.workflow/state/sessions/*.json`
- `.workflow/state/workstreams/*.json`

Rules:

- this layer is operational, not semantic
- do not send routine progress into the external memory provider
- this layer is updated automatically by hooks/plugins

## Layer 2: reusable experience memory

Purpose:

- store decisions, architecture insight, debugging lessons, workflow conventions, and explicit preferences

Provider model:

- use a provider interface
- `Memory MCP` is a default implementation, not a required hardcoded dependency
- project tags come from `workflow/config.json`, not from code constants

Allowed memory types:

- `decision`
- `architecture`
- `debugging`
- `workflow`
- `preference`

## Configuration

Read `workflow/config.json` for:

- `project.id`
- `project.slug`
- `paths.state_root`
- `memory_provider.*`

If `memory_provider.tags` is absent, derive tags from the configured project identity.

## Review workflow

- reviewable decisions must include `review-pending`
- delayed review checks happen at session start
- completing a review should supersede or update the prior memory record

## Boundaries

Do:

- update local session state automatically
- query the memory provider for reusable knowledge
- store only durable, reusable experience in the provider

Do not:

- hardcode project tags
- mix implementation progress with semantic experience
- treat local state snapshots as long-term design memory
