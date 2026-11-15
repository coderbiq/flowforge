---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      documentation-refresh-requirements: 4
    design:
      documentation-refresh-design: 2
---

# 06: 文档刷新终检收口

**Blocked by:** 01, 02, 03, 04, 05
**Status:** closed
**Mode:** lightweight

## Delivery

对 docs-refresh 全部交付做一次机械终检：gap 三清单归属项清零复核、跨文档命令一致性抽检、链接有效性、`check --strict` 通过，并记录复核证据。

## Design context

收口票：依赖 01–05 全部交付后才有可检对象（真实 DAG 边——每张上游票消解自己的 gap 子集，本票验证总和为零）。

See the design authority at [文档刷新方案](../design.md#documentation-refresh-design)（d-dag 节）. Requirement authority: [文档刷新需求](../requirements.md#documentation-refresh-requirements)（可观察结果 1–4）.

## Touch points

- 只读复核 + `docs/research/2026-09-29-doc-refresh/` 内追加复核记录；无生产写集（发现问题时开 Fix/回上游票，不在本票直接改文档）

## Changes

 - [x] 1. current-surface-gaps.md A(17)/B(12)/C(4) 逐条复核：每条确认"已消除/已裁决/归属外"，结论追加为 gap 文档附页。
  - cmd: `./bin/flowforge check --dir docs/proposals/documentation-refresh --strict` 无关；实质复核见附页 §一（逐条 doc 证据行）
  - exit: 0（F1 修复后复验）
  - output: "复验（refine 消费后）：33 条 = 已消除 28 / 归属外 3（A12/A16/A17）/ 待修复 0；B2/C1 经票 05 Round 1 修复后旧表述 0 命中（grep 'docs/proposals` 作为兼容默认值' → 0），附页行已转「已消除」并附消费记录"
  - artifact: docs/research/2026-09-29-doc-refresh/final-review-appendix.md
- [x] 2. 跨文档命令一致性抽检：README、scenarios.md、cli-design.md 中出现的同一命令表述逐字一致（抽 8 条）。
  - cmd: `python3 /tmp/cmd_consistency.py`（抽检脚本，见附页 §二）
  - exit: 0
  - output: "8/8 命令（init / agents deploy / model-set use / frontier --pi-workflow / check --strict / status / assets verify / upgrade）在共现文件中逐字一致，且与 ./bin/flowforge <cmd> --help 相符；附页记 1 观察项（README:71 只列 PI 扩展两项能力，非矛盾）"
  - artifact: docs/research/2026-09-29-doc-refresh/final-review-appendix.md
- [x] 3. 链接有效性：六份文档中相对链接逐条可解析。
  - cmd: `python3 /tmp/links.py`（逐条解析 + os.path.exists）
  - exit: 0（F2/F3 修复后复验）
  - output: "复验（refine 消费后）：scenarios.md 21 条相对链接 broken=0；六份文档全量相对链接复扫均 NONE（python3 逐条 os.path.exists，从各文档自身目录解析）"
  - artifact: docs/research/2026-09-29-doc-refresh/final-review-appendix.md
- [x] 4. `./bin/flowforge check --dir docs/proposals/documentation-refresh --strict` exit 0。
  - cmd: `./bin/flowforge check --dir docs/proposals/documentation-refresh --strict`
  - exit: 0
  - output: "Checked 6 issues in docs/proposals/documentation-refresh / ✓ Dependency graph is healthy. No cycles or dangling references found.（**复核时点**：写入 `## Blocked evidence` 之前；本票被阻断后 strict 因自身新增的 `blocked-evidence-present` 诊断 exit 1，属诊断 #7 的预期机制，非 proposal 图缺陷——默认 `check` 仍 exit 0）"
  - artifact: docs/research/2026-09-29-doc-refresh/final-review-appendix.md

## Implementation note（收口轮补充）

- 首轮执行（flowforge-implementer，fresh context）：Change 2/4 过、Change 1/3 按票面 Failure scenario 阻断并留 Blocked evidence；附页落盘 33 条结论与两项编排裁决。
- 上游修复：票 05 Round 1（F1，architecture.md 兼容默认句）、票 02 Round 4（F2/F3，链接基线与目标文件）。
- refine 消费（编排会话）：失败事实转写为 Verified contracts、Blocked evidence 原文移除、Change 1/3 复验勾选、附页行翻转 + 消费记录；本轮全程未改六份文档正文（约束遵守）。

## Completion evidence

- 交付行为：33 条 gap 逐条复核结论落附页（已消除 28 / 归属外 3 / 待修复 0，含消费记录与两项编排裁决原文）；8 条跨文档命令一致性 8/8；六份文档全量相对链接复扫 0 断链；blocked evidence 按 refine 契约消费（转写 Verified contracts、移除原文）。
- 验证方法与观测：Change 1–4 证据四元组齐备（1/3 含修复前后双时点）；`./bin/flowforge check --dir docs/proposals/documentation-refresh --strict` exit 0（本 evidence 写入后）。
- 双轴与发现处置：首轮阻断发现 3 项上游缺陷全部经 Fix 轮修复并复验（票 05 Round 1 / 票 02 Round 4）；无 waiver。
- 偏差：无（终检票只写附页与本票，六份文档正文零改动，Constraints 遵守）。
- 实现参考：本 commit（docs: finalize documentation-refresh verification）。

## Constraints

- 本票不改六份文档正文；发现缺陷 → 追加 `Fix:` Change 到对应上游票并重开（repair 语义）。

## Done and verify

四条 Changes 全部通过并在 gap 附页留证；strict check exit 0。

---

## Execution detail

### Verified contracts

- 复核基线：current-surface-gaps.md 三清单 33 条，每条含双侧引用（gap 文档结构实测）。
- strict 语义：`check --strict` 使未豁免 warning/gap 失败（`check --help` 实测）。

### Execution scenarios

- Success：33 条全处置、strict 通过、附页留证。
- Failure：任一条目未消除且无裁决记录 → 回上游票修复，本票不直接改。

### Expected tests

- `./bin/flowforge check --dir docs/proposals/documentation-refresh --strict` — exit 0。
- 附页存在且 33 条计数完整。

### Generated artifacts

- `docs/research/2026-09-29-doc-refresh/` 附页复核记录。

### Conventions

- must 只读复核 + 附页留证（源：Constraints）。
- 变更后运行 `./bin/flowforge check --dir docs/proposals/documentation-refresh`（源：AGENTS.md）。

---

## Implementation note

轻量模式执行（票面 `**Mode:** lightweight`）。只读复核；**未改动六份文档正文**；仅新建附页 `docs/research/2026-09-29-doc-refresh/final-review-appendix.md` 并写本票。

- 完成：Change 2（命令一致性 8/8 逐字一致）、Change 4（`check --strict` exit 0）。
- 未完成：Change 1（33 条中 27 已消除 / 3 归属外 —— 但 B2/C1 同处 `architecture.md:66` 未消除且无裁决记录）、Change 3（`docs/scenarios.md` 20 条来源链接断链）。
- 命令：`./bin/flowforge check --dir docs/proposals/documentation-refresh --strict`（复核时点 exit 0；写入 `## Blocked evidence` 后 strict 因诊断 #7 `blocked-evidence-present` exit 1，默认 check 仍 exit 0）；`python3 /tmp/cmd_consistency.py`（8/8 一致）；`python3 /tmp/links.py`（可解析 19、断链 20）。
- 文件修改：`docs/research/2026-09-29-doc-refresh/final-review-appendix.md`（新建）、本票。
- write-set 合规：本票声明「无生产写集，仅 research 附页」——所有修改在声明范围内；六份文档正文零改动。

## Blocked evidence（已消费）

按 refine-ticket 契约消费：F1/F2/F3 失败事实已转写为下列 Verified contracts 并处置（上游票 05 Round 1 / 票 02 Round 4 修复 + 附页消费记录），本小节原始逐字证据已按"瞬态即焚"移除，复验命令与结果见 Change 1/3 更新后的证据四元组。

### Verified contracts（消费转写）

- `docs/architecture.md` 兼容默认句现为 `<startDir>/ff-wiki/proposals`（真值 `internal/config/config.go:17,225` `DefaultDocsDir="ff-wiki"` + `ResolveProposalsDir`）；旧表述 `docs/proposals` 作为兼容默认值 复验 0 命中；与 `docs/cli-design.md` 一致。
- `docs/scenarios.md` 来源链接基线为 `proposals/<name>/...`（同级相对路径）；documentation-contract-refinement 与 lightweight-execution-contract 两处指向各自实际存在的 `spec.md`；21 条相对链接 broken=0。
- 链接验证正确姿势（教训转写）：从文档自身目录解析相对路径并 `os.path.exists` 目标文件——只验证目录名存在性会漏掉基线错误与目标文件缺失（本轮 F2/F3 根因）。
- 附页 final-review-appendix.md 已含消费记录与 B2/C1 行结论翻转；strict check 于消费后复跑见 Change 4 补验。

**STATUS: COMPLETED（消费后收口）**
