---
name: flowforge-roadmap
description: Manage the intent layer above proposals — roadmap direction files with an unlimited node hierarchy that records coarse requirements and hatches mature nodes into proposals — roadmap owns long-horizon planning and aggregation, NOT for one-proposal work — that's flowforge-align, and NOT for decision-battle maps — that's flowforge-wayfinder. Use when the user says "长期规划"/"roadmap"/"大方向"/"需求树"/"孵化"/"对账".
---

# Roadmap direction files

Own the intent layer: direction files under `<docs_dir>/roadmap/` hold coarse, deliberately-unrefined requirements as an unlimited node tree. Proposal directories keep owning every settled requirement, design, and execution contract — the tree never copies their truth. Use the shared contract's [packaging](../_shared/ARTIFACT-CONTRACT.md#packaging) when a hatched node's proposal publishes artifacts.

## Layer discipline (non-negotiable)

- The tree stores intent only. Requirement truth, acceptance behavior, and interface detail live exclusively in the hatched proposal's authorities.
- Node state is always derived (`flowforge status` computes it from artifacts and ticket statuses). Never write status, percentages, or checkboxes into nodes.
- Nodes carry no IDs and no dependencies between each other. Ordering constraints belong to proposals via their existing `Blocked by`; cross-proposal references already resolve globally.
- One node links at most one proposal. Multi-link is a reported degenerate state; fix it by splitting the node, not by adding links.
- `roadmap/` sits outside the proposals scan root: it has no execution semantics, and a ticket's reading list never needs the tree.

## Direction file format

One file per direction at `<docs_dir>/roadmap/<direction>.md`, frontmatter `role: map`. Completed or abandoned directions move to `roadmap/closed/` (no year tier; directions accumulate far slower than proposals).

```markdown
---
flowforge:
  schema: 1
  role: map
  id: <direction>-map
  revision: 1
---
# <Direction title>

## 意图
Two or three lines: what the world looks like when this direction is walked to the end.

## 边界
- Out-of-scope statements at direction level (optional).

## 树
- **<Title>**: <one-sentence summary>
  <optional description paragraph: background, fog, hatch conditions, constraints
   from sibling nodes. Records intent and constraints only — never acceptance
   behavior or interface detail.>
  → &lt;slug-link to ../proposals/&lt;slug&gt;/README.md&gt;  ← the hatch link line (arrow + markdown link to the proposal README)
```

Increment `revision` only on tree-structure changes (nodes added/removed/re-parented); prose polish does not increment.

## Actions

### Chart (建图)

A loose idea too big for one proposal arrives. Ask at most three questions to pin 意图 (destination) and 边界, then write the direction file with its initial tree. **Refuse over-building**: if the idea fits one proposal, do not create a direction file — route directly to `flowforge-align`.

### Grow (养树)

Add, split, or re-parent nodes; maintain summaries and description paragraphs under the format rules above. A description paragraph that outgrows a few lines is the signal to hatch the node or split it — not to keep writing.

### Hatch (孵化)

The core hand-off action. Preconditions: the node's summary and description are coherent, and the user has confirmed the scope.

1. Scan sibling nodes and already-hatched proposals' description paragraphs; collect cross-node constraints into a hatch note.
2. Create `proposals/<slug>/README.md` as a stub: title, one sentence, back-link to the direction file and node title, and one line "requirement authority 由 align 建立". **Do not write requirements.md** — align creates the requirement authority only when the first accepted requirement fact exists.
3. Write the link line onto the node: a line starting with `→` followed by a markdown link whose target is the proposal README, e.g. arrow, slug label in brackets, then `../proposals/<slug>/README.md` in parens.
4. Hand off to `flowforge-align` with the hatch note (node summary + collected constraints + stub location) as input.

Slug naming is free-form; specific conventions are set per project.

### Reconcile (对账)

Run `flowforge status` and read the tree view plus its reports:

- dead-link nodes (link target missing) → fix the link (grow) or the proposal location (`flowforge proposal close/reopen`);
- orphan proposals (live, unlinked) → report only; adopt or close at the user's pace;
- multi-link nodes → split the node (grow);
- direction complete or abandoned → move the file to `roadmap/closed/` (fix the proposal stub back-links by hand; direction lifecycle is year-scale).

`flowforge check` continuously reports `roadmap-dead-link` warnings between reconciles.

## Proposal lifecycle commands

- `flowforge proposal close <slug>` — moves a fully-terminal proposal to `proposals/closed/<current-year>/` and rewrites every relative link to it across `<docs_dir>`. Refuses when non-terminal tickets exist.
- `flowforge proposal reopen <slug>` — moves it back and re-rewrites. Refuses on live name collision or ambiguous archive matches.

Both take `-d <proposals-dir>` like `status`. Manual `mv` is forbidden — it breaks links.

## Boundaries

This skill MUST NOT clarify requirements (align), design solutions (solution-design), slice tickets (plan), write node states, or create node dependencies. It routes to those skills at the hatch boundary and reads their artifacts for derived aggregation only.
