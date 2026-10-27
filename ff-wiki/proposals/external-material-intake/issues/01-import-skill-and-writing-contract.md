---
flowforge:
  schema: 1
  role: ticket
  id: import-skill-and-writing-contract
  revision: 1
  consumes:
    requirements:
      external-material-intake-requirements: 1
    design:
      external-material-intake-design: 1
---

# 01: Add external-material intake and a shared writing contract

**Blocked by:** None

**Status:** closed

**Delivery:** An agent can use a deployed `flowforge-import` Skill to classify any local source material and hand only its relevant, traceable facts to the existing authority owners in the target language.

## Design context

Implement [external-material intake](../design.md#external-material-intake-design) for the outcomes in [the requirement authority](../requirements.md#external-material-intake-requirements). Import is an input interpreter: Align remains the owner of requirements and Solution Design remains the owner of module/interface decisions.

## Touch points

`assets/skills/flowforge-import/`; `assets/skills/_shared/ARTIFACT-CONTRACT.md`; packaged-skill reference tests in `internal/command/assets_deploy_test.go`.

## Changes

1. Create the Import Skill with its source-path, feature, optional-language input contract; classify retained material as source fact, requirement candidate, design decision, delivery/verification evidence, or unknown/conflict.
2. Define compact source attribution and the promotion threshold for a `role: research` `source-notes.md` without copying source prose.
3. Add the target-language semantic-rewrite and information-value rules to the shared contract, including stable code identifiers and terminology.
4. Add the new Skill and its references to the packaging contract tests.

## Constraints

- Do not create a legacy-proposal migration lifecycle, issue files, or a readiness state.
- Do not let Import publish requirement, design, ticket, or evidence authority in another Skill's name.
- Keep machine IDs/revisions as traceability metadata; source locations and semantic prose are the human reading interface.

## Done and verify

- A mixed local PRD/old-proposal fixture produces the five classifications and a compact hand-off to Align/Solution Design without copied source paragraphs.
- Packaged-skill tests resolve every Import reference in source and deployed layouts: `go test ./internal/command/...`.

## Completion evidence

- Added the deployed [Import Skill](../../../../assets/skills/flowforge-import/SKILL.md), which classifies local material before handing it to the authority owner and records target language without converting sources into authority.
- Added the shared [source-intake and semantic-rewrite contract](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#source-intake-and-semantic-rewrite); the [mixed-source walk-through](../scenario-fixtures/mixed-source-import.example.md) exercises all five classifications in Chinese without copied source text or ticket creation.
- Added package coverage for the new contract pointers; `go test ./internal/command/...`, `GOPROXY=https://goproxy.cn,direct go test ./internal/...`, `flowforge check --dir docs/proposals/external-material-intake --strict`, `flowforge frontier --dir docs/proposals/external-material-intake --strict`, and `git diff --check` passed.
- Dual-axis review: Standards found no violations or smell concerns. Specification review found missing scenario proof and ambiguous target-language hand-off; both were corrected by the walk-through and explicit completion contract.
