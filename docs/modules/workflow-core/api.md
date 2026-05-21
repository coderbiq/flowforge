# workflow-core API

`workflow-core` is consumed through the CLI wrappers in `scripts/`.

## Proposal commands

- `scripts/flowforge-create-proposal.js`
  - create a proposal skeleton
  - accept archive targets and canonical corpus entries
- `scripts/flowforge-approve-proposal.js`
  - move a proposal from `draft` or `proposed` to `approved`
- `scripts/flowforge-apply-proposal.js`
  - prepare execution notes and transition to `active`
- `scripts/flowforge-validate-proposal.js`
  - validate proposal metadata and task map consistency
- `scripts/flowforge-proposal-status.js`
  - summarize proposal state and backend health
- `scripts/flowforge-archive-proposal.js`
  - update archive targets and close the proposal

## Public behaviors

- proposal ids use `CRYYMMDDNN`
- `archive_targets[].key` is the stable reference for task mapping
- `canonical_corpus` records the final docs reviewed as the baseline
- `scope` distinguishes `workspace`, `cross-workspace`, and `monorepo`

