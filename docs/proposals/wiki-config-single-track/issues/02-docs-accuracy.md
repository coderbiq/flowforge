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
**Status:** open
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

- [ ] 1. `README.md`：默认 wiki 目录表述 `docs/` → `ff-wiki/`（保持上下文语义完整）。
- [ ] 2. `docs/proposals/pi-host-integration/issues/03-pi-extension.md`：修正"缺省 docs"过时表述为 `ff-wiki`（或按上下文改述）。
- [ ] 3. `docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md`：L64 附近题材候选句所引"设计 Next Steps"锚点不存在——改述为指向该票 Execution detail 的题材事实（编排会话选定）或删除幻引用。
- [ ] 4. `.flowforge/config.yaml`：删除 `projects[].wikiRoot` 行（untracked，本机同步；若 01 已使 load 告警，此行消除本机告警噪声）。

## Constraints

- 纯文档与本地配置清理，不动代码（01 管辖面）。
- 已完结票面（pi-host-integration 03、generic-role 04）的修正仅限失准表述，不改票面历史语义。
- Write set: `README.md`、`docs/proposals/pi-host-integration/issues/03-pi-extension.md`、`docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md`、`.flowforge/config.yaml`

## Done and verify

- 失准清零: `grep -rn '缺省 docs\|默认.*docs/' README.md docs/proposals/pi-host-integration/issues/03-pi-extension.md` — 无命中。
- 幻锚点清除: `grep -n '设计 Next Steps' docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md` — 无命中。
- 自举清理: `grep -n 'wikiRoot' .flowforge/config.yaml` — 无命中。
- 全套回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok（文档票回归确认）。

---

## Execution detail

### Verified contracts

- 三处失准坐标经调研 W3 与归档扫描证实（README.md:92、pi-host-integration/issues/03:136、generic-role-orchestration/issues/04:64）；architecture.md 第四处已修（commit f32b9e9）。
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
