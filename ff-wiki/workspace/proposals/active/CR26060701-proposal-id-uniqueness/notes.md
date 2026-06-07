# Notes: CR26060701-proposal-id-uniqueness

## 需求树

- Proposal ID 唯一性保障
  - ID 生成时的冲突检测
    - 检查 active/ 和 completed/ 目录下是否已有同名 CR-id 前缀的 proposal
    - 给出清晰的冲突错误提示，包含已有的 proposal 路径
  - NN 序号的自动化分配
    - 自动查询当日已有 proposal 的最大 NN，生成不冲突的序号
  - 校验闸门集成
    - 在 validate-proposal 中增加 ID 唯一性检查
    - 在 design-context 中增加 ID 冲突预警
  - SKILL 指令更新
    - flowforge-design SKILL 5.1 步骤增加查重指引

## 探索记录

### 2026-06-07 — ID 生成全链路分析

CR-id 由 Agent 手动从模板 `CR{YYMMDD}{NN}` 生成，无代码自动生成逻辑。关键风险：
- `{NN}` 序列号由 Agent 手动确定，同一日期可能冲突
- 目录由 Agent `fs.mkdirSync` 创建，无前置唯一性检查
- `findProposalById` 在 4 个脚本中重复实现但未用于创建前冲突检测
- `validate-proposal.js` 仅检查格式（正则），不检查唯一性
- `flowforge task init` 仅检查 beads epic，不检查文件系统

发现记录：
- `library/architecture/proposal-id-lifecycle.md` — 全链路架构事实
- `library/conventions/script-validation-patterns.md` — 校验模式约定

### 2026-06-07 — 设计决策

设计方案：最小侵入，不新增 CLI 命令。
- `design-context.js` 新增 `--check-id` / `--suggest-id` 模式
- `lib/config.js` 提取公用 `checkProposalId` / `suggestProposalId`
- `validate-proposal.js` 新增目录层级唯一性检查
- SKILL 5.1 增加查重步骤

详见 `design/proposal-id-uniqueness.md`。

## 实施日志

### 2026-06-07 — 实施完成 [实施 5/5]

- **flowforge-fj6.3.1** — `lib/config.js` 新增 `checkProposalId` 和 `suggestProposalId`，经验证 CR26060701 正确检测冲突，建议 ID 为 CR26060702
- **flowforge-fj6.3.2** — `design-context.js` 新增 `--check-id` / `--suggest-id` 模式，JSON 在末尾追加不干扰正常输出
- **flowforge-fj6.3.3** — `validate-proposal.js` 新增 ID 唯一性检查，排除自身目录仅报告外部冲突
- **flowforge-fj6.3.4** — `flowforge-design/SKILL.md` 阶段 5.1 步骤 1 更新查重流程
- **flowforge-fj6.3.5** — `tests/suite-context-output.js` 新增 Check 4/5，全量测试 533 passed 0 failed
