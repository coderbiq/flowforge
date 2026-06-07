---
proposal_id: CR26060701
title: 加强 Proposal ID 唯一性校验
status: draft
created: 2026-06-07
updated: 2026-06-07
author: Sisyphus
---

# CR26060701: 加强 Proposal ID 唯一性校验

## 问题

当前 CR-id（`CR{YYMMDD}{NN}`）中的 `{NN}` 序列号由 Agent 手动确定，没有任何机制检查两个 proposal 是否会生成相同的 ID。同一天创建两个 proposal 时，如果都选了相同的 NN（如 `01`），第二个会因目录已存在的文件系统错误而失败——但这个错误发生在 `mkdir` 阶段，而不是在 ID 生成阶段被提前预防。

**根因分析**：
- `findProposalById` 查找逻辑已在 4 个 context 脚本中实现，但仅用于"已有 CR-id 时查找对应目录"，从未在创建新 proposal 时用于冲突检测
- `flowforge task init` 的 `hasTaskSpace` 只检查 beads epic 存在性，不检查文件系统
- `flowforge validate-proposal` 仅校验 ID 格式（正则），不检查唯一性

## 方案

采用**最小侵入**策略，在现有脚本上增加 ID 检查能力，不引入新脚本或新 CLI 命令。

1. **`design-context.js`** 新增 `--check-id <CR-id>` 模式：扫描所有 project 的 active/ 和 completed/ 目录，检测 CR-id 前缀冲突
2. **`design-context.js`** 新增 `--suggest-id` 模式：自动计算当日可用的下一个 NN 序号
3. **`lib/config.js`** 提取公用函数 `checkProposalId` / `suggestProposalId`，消除 4 个脚本中的重复 `findProposalById` 实现
4. **`validate-proposal.js`** 新增目录层级唯一性检查：在格式校验之外，检测同一 CR-id 是否被其他 proposal 占用
5. **`flowforge-design/SKILL.md`** 阶段 5.1 步骤增加 `--check-id` 查重指引

详见 `design/proposal-id-uniqueness.md`。

## 影响范围

| 文件 | 变更 |
|------|------|
| `src/cli/scripts/design-context.js` | 新增 `--check-id` / `--suggest-id` 参数处理 |
| `src/cli/scripts/lib/config.js` | 新增 `checkProposalId` / `suggestProposalId` |
| `src/cli/scripts/validate-proposal.js` | 新增 ID 唯一性校验 |
| `src/agents/flowforge-design/SKILL.md` | 5.1 增加查重步骤 |
| `tests/suite-context-output.js` | 新增参数处理验证 |

## 不变部分

- `flowforge` CLI 入口不变（不新增子命令）
- `design-context.js` 正常输出格式不变（`--check-id` / `--suggest-id` 在末尾追加 JSON）
- `proposal_id` 模板 `CR{YYMMDD}{NN}` 不变
- meta.yaml 格式不变
