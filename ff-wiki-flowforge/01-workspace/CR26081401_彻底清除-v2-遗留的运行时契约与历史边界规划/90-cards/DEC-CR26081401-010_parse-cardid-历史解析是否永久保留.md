---
id: DEC-CR26081401-010
title: ParseCardID 历史解析是否永久保留
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081401
      relation: belongs_to
created: 2026-08-14T14:27:52.618973+08:00
updated: 2026-08-14T14:27:52.620699+08:00
source: CR26081401
---

# ParseCardID 历史解析是否永久保留

## Context

FIND-CR26081401-006 证明 `ParseCardID` 只被 `internal/core/naming_test.go` 命中，没有生产调用；实现仍支持普通 ID 与 TASK 特殊四/五段形状。当前 store 定位使用 `ParseFilename`/ID/PROP fallback，无法证明外部插件或未来历史读取不会依赖该 helper。

## Decision

用户已确认选择 B：立即删除 `ParseCardID` 及 TASK/subtask 专属解析支持；未来历史定位或迁移另立 Proposal。`ParseFilename` 同步收窄为只接受当前 v3 `{ID}_{slug}.md` 文件名契约，旧 TASK/旧模型文件名不再由运行时解析。

原文件、旧 ID、旧 links、body、frontmatter 和 path 保持原样，不通过 alias、fallback 或迁移写入恢复旧解析能力。

## Rationale

无生产调用、用户已明确不设兼容窗口，且旧 TASK/subtask 解析不属于当前 v3 契约；删除后历史定位/迁移责任必须进入独立 Proposal。

## Consequences

运行时 API 收窄到当前 v3 文件名解析；旧 TASK ID 仅作为不进入 current 域的原始历史数据存在。历史文件仍保持可审计的原始内容，但不再获得运行时解析兼容。

## Alternatives

拒绝保留旧 TASK 写入/查询能力；拒绝通过 alias、fallback 或迁移写入改变旧 ID、旧 links、body 或 path。

References: FEAT-CR26081401-002；FIND-CR26081401-006；PROP-CR26081401。

## Links

### Outgoing

- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划

### Incoming

- [FEAT-CR26081401-002](FEAT-CR26081401-002_运行时代码与历史解析边界清理规划.md) [feature] - 运行时代码与历史解析边界清理规划

