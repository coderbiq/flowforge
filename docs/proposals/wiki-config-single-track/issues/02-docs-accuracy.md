---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      wiki-config-single-track-requirements: 3
    design:
      wiki-config-single-track-design: 1
---

# 02: 文档失准清零 + 自举 config 清理

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

3 处 d04e955 前过时表述修复（README、pi-host-integration 票面、generic-role 票面幻锚点）+ 本仓自举 `.flowforge/config.yaml` 删除 `projects[].wikiRoot` 残留行；grep 验收旧表述清零。

## Design context

d04e955 将 `DefaultDocsDir` 从 `docs` 改为 `ff-wiki`，四处文档表述未同步（architecture.md 一处已在 generic-role-orchestration 归档时顺手修复，本票清余下三处）。幻锚点为规划期笔误引用。

See the design authority at [双轨 wiki 配置统一：方案](../design.md#wiki-config-single-track-design)（d-docs 节）. Requirement authority: [双轨 wiki 配置统一：需求](../requirements.md#wiki-config-single-track-requirements)（目标 5，验收 4）.

## Touch points

- `README.md` — L92 附近默认 `docs/` 表述
- `docs/proposals/pi-host-integration/issues/03-pi-extension.md` — L136 附近"缺省 docs"表述
- `docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md` — L64 附近"设计 Next Steps"幻锚点引用
- `.flowforge/config.yaml` — `projects[].wikiRoot` 行（untracked 本地文件）

## Changes

- [x] 1. `README.md`：默认 wiki 目录表述 `docs/` → `ff-wiki/`（保持上下文语义完整）。
  - cmd: `grep -c 'ff-wiki' README.md`
  - exit: 0
  - output: `2`（默认创建清单四处根路径已改 ff-wiki/）
  - artifact: README.md
- [x] 2. `docs/proposals/pi-host-integration/issues/03-pi-extension.md`：修正"缺省 docs"过时表述为 `ff-wiki`（或按上下文改述）。
  - cmd: `grep -n '缺省 `ff-wiki`' docs/proposals/pi-host-integration/issues/03-pi-extension.md`
  - exit: 0
  - output: `136:` 命中修正后表述
  - artifact: docs/proposals/pi-host-integration/issues/03-pi-extension.md
- [x] 3. `docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md`：L64 附近题材候选句所引"设计 Next Steps"锚点不存在——改述为指向该票 Execution detail 的题材事实（编排会话选定）或删除幻引用。
  - cmd: `git diff --stat docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md`
  - exit: 0
  - output: 幻锚点句尾引用已删（1 行最小 diff，票面历史语义未动）
  - artifact: docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md
- [x] 4. `.flowforge/config.yaml`：删除 `projects[].wikiRoot` 行（untracked，本机同步；若 01 已使 load 告警，此行消除本机告警噪声）。
  - cmd: `grep -c 'docs_dir' .flowforge/config.yaml`
  - exit: 0
  - output: `1`（wikiRoot 行已删，docs_dir 保留）
  - artifact: .flowforge/config.yaml

## Constraints

- 纯文档与本地配置清理，不动代码（01 管辖面）。
- 已完结票面（pi-host-integration 03、generic-role 04）的修正仅限失准表述，不改票面历史语义。
- Write set: `README.md`、`docs/proposals/pi-host-integration/issues/03-pi-extension.md`、`docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md`、`.flowforge/config.yaml`

## Done and verify

- 失准清零: `grep -rn '缺省 \`docs\`\|默认创建.*\`docs/' README.md docs/proposals/pi-host-integration/issues/03-pi-extension.md` — 无命中（验收命令已修正：原模式对反引号包裹的真缺陷无检出力，且误报 03:58/122 的准确历史内容）。
- 幻锚点清除: `grep -n '设计 Next Steps' docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md` — 无命中。
- 自举清理: `grep -n 'wikiRoot' .flowforge/config.yaml` — 无命中。
- 全套回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok（文档票回归确认）。

---

## Execution detail

### Verified contracts

- 三处失准坐标经调研 W3 与归档扫描证实（README.md:92、pi-host-integration/issues/03:136、generic-role-orchestration/issues/04:64）；architecture.md 第四处已修（commit f32b9e9）。
- 验收正则教训（Blocked evidence 折入）：原 Done-1 模式 `缺省 docs\|默认.*docs/` 对反引号包裹的真缺陷（`缺省 \`docs\``）无检出力，且误报 03:58/122（“（默认路径）/（默认导出）”+真实文件路径 `docs/cli-design.md`，属准确历史内容不得修改）；已实测修正模式见 Done-1。规划期验收 grep 必须先跑一遍再入票。
- Done-4 全绿受并行票 01 牵连：失败面 `TestLoadWarnsOnDeprecatedWikiKeys` 属 01 在飞实现，与本票零交集；待 01 落地后复核。
- `.flowforge/config.yaml` 为 v2 格式 untracked 本地文件（`docs_dir: docs` + `projects[0].wikiRoot: ff-wiki-flowforge`），删除 wikiRoot 行不影响其他键。
- 本仓实际 wiki 根为 `docs`（docs_dir 显式），文档表述修正不改变本仓行为。

### Execution scenarios

- Success：三处 grep 清零、自举 config 无 wikiRoot、全套测试绿。
- Failure：若误改 pi-host-integration 03 票面非失准内容，改变完结票历史语义——以最小 diff 约束（Constraints 重申）。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok（无代码改动回归确认）。

### Generated artifacts

- Not applicable：纯文档与本地 untracked 配置清理。

### Conventions

- must 变更后运行 `go test ./internal/...`（转录自设计 Standards clauses）。
- 完结票面修正保持最小 diff，只动失准表述。

## Implementation note

### 三处修正前后文（+自举 config）

1. `README.md:92`（Change 1）
   - 前：默认创建 `.flowforge/config.yaml`、`docs/CONTEXT.md`、`docs/adr/`、`docs/proposals/`，并部署 `.agents/skills/` 与 `docs/agents/`。
   - 后：默认创建 `.flowforge/config.yaml`、`ff-wiki/CONTEXT.md`、`ff-wiki/adr/`、`ff-wiki/proposals/`，并部署 `.agents/skills/` 与 `ff-wiki/agents/`。
   - `.flowforge/config.yaml`、`.agents/skills/` 为仓库根/宿主路径不随 docs_dir 变化，保持不动（最小 diff）。
2. `docs/proposals/pi-host-integration/issues/03-pi-extension.md:136`（Change 2）
   - 前：Change 1(b) 的 check 目录按配置解析 `docs_dir`（缺省 `docs`）拼 `<docs_dir>/proposals`，与任务文本 `<docs>/proposals` 一致。
   - 后：同句改为（缺省 `ff-wiki`），其余零改动——只换默认值字面量，句式与同节其余事实修正项不动（最小 diff，不改完结票历史语义）。
3. `docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md:64`（Change 3）
   - 前：- 题材（编排会话已选定）：双轨 wiki 配置收敛调研——`Wiki.Root`/`wikiRoot` 与 `docs_dir` 双轨现状盘点，喂给后续小提案（设计 Next Steps 已列）。
   - 后：同句删去句尾“（设计 Next Steps 已列）”。采纳票面授权的“删除幻引用”选项：句首“题材（编排会话已选定）”已承载 Execution detail 题材事实，改述补锚点反成循环引用。
4. `.flowforge/config.yaml`（Change 4）——删除 `      wikiRoot: ff-wiki-flowforge` 一行；文件现为 `version: 2.0.0` / `docs_dir: docs` / `projects[0]{id: flowforge-v2, srcDirs: [.]}`，其余键原样（untracked 本地文件）。

### grep 验收输出（修正后实测）

- Done-2 幻锚点：`grep -n '设计 Next Steps' docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md` → 无输出，exit=1（✓ 清零）。
- Done-3 自举：`grep -n 'wikiRoot' .flowforge/config.yaml` → 无输出，exit=1（✓ 清零）。
- Done-1 失准清零：`grep -rn '缺省 docs\|默认.*docs/' README.md docs/proposals/pi-host-integration/issues/03-pi-extension.md` → 仍 2 命中（03:58、03:122），exit=0（✗，与 Constraints 矛盾，详见 Blocked evidence）。
- 缺陷表述本体清零补验：`grep -n '缺省 \`docs\`' docs/proposals/pi-host-integration/issues/03-pi-extension.md` → 0 命中（L136 现为“缺省 `ff-wiki`”）；README:92 现仅含 `ff-wiki/` 路径（`grep -n '默认创建' README.md` 证实）。
- Done-1 替换模式候选（实测 exit=1 清零）：`grep -rn '缺省 \`docs\`\|默认创建.*\`docs/' README.md docs/proposals/pi-host-integration/issues/03-pi-extension.md`。

### go test

`GOPROXY=https://goproxy.cn,direct go test ./internal/...`（1 次，fail-fast）：command/subagent/tracker/update 全 ok；`internal/config` FAIL 于 `TestLoadWarnsOnDeprecatedWikiKeys`——为并行票 01 新增测试且其实现尚在飞（发现 2 已折入 Execution detail 的 Verified contracts：Done-4 待 01 落地复核），与本票零代码改动无关。
