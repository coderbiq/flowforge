---
description: 创建提案和任务 Epic。用户提到"提案"、"proposal"时触发。
allowed-tools: Skill(tg-proposal)
---

Use the `tg-proposal` skill to create a proposal and task Epic.

## 执行流程

1. 生成提案编号 `CR{YYMMDD}{序号}`
2. 创建提案目录结构
3. 从探索笔记提取内容填充 proposal.md
4. 询问是否创建任务 Epic
5. 存储长期记忆

## 触发场景

- 用户要创建提案
- 探索笔记已定稿

Arguments: $ARGUMENTS
