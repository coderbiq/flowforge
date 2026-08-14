---
id: DEC-CR26081401-011
title: 旧 CardType 常量与 Prefix 处置边界
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081401
      relation: belongs_to
created: 2026-08-14T14:27:52.740174+08:00
updated: 2026-08-14T14:27:52.740837+08:00
source: CR26081401
---

# 旧 CardType 常量与 Prefix 处置边界

## Context

FIND-006 证明 current ID 生成仍使用 `Prefix()`，而 `CardType` 的旧常量/Prefix/FromPrefix 分支仍参与 legacy 解析过滤、目录映射、校验和 STR control-plane 相关路径。普通域已经排除旧卡/STR，但不能整体删除旧常量而不先拆 metadata loader 与历史边界。

## Decision

用户已确认：删除旧 `CardType` 常量和公开 `Prefix`/prefix 规则；不保留 legacy CardType/Prefix 公开兼容层。仅保留 Proposal control-plane 对 `STR-<proposal>-REQ` 的内部识别与读取，普通 CardStore、list/search/index 和 current query 不识别 STR。

当前 v3 新卡所需的 ID 生成继续由内部 current-v3 实现承担，但不得重新暴露旧类型常量、旧 Prefix 映射或旧类型解析 API。旧 ID 统一按非当前资源处理；原文件、旧 ID、旧 links、body、frontmatter 和 path 不变。

## Rationale

用户决策区分 current-v3 ID 生成与旧公开契约：current-v3 只保留必要的内部实现，旧 CardType 常量和公开 Prefix 规则不再成为 API 或历史解析入口。

## Consequences

运行时减少旧类型和公开 Prefix 的维护面；Proposal inspect/traceability 必须通过独立 STR control-plane loader 工作，不能依赖普通 CardType 域。所有旧文件仍保持原样。

## Alternatives

拒绝保留旧 CLI/旧类型写入；拒绝把 STR control-plane 内部识别扩展回普通 Card 域；拒绝修改历史文件以适配新的解析边界。

References: FEAT-CR26081401-002；FIND-CR26081401-006；PROP-CR26081401。

## Links

### Outgoing

- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划

### Incoming

- [FEAT-CR26081401-002](FEAT-CR26081401-002_运行时代码与历史解析边界清理规划.md) [feature] - 运行时代码与历史解析边界清理规划

