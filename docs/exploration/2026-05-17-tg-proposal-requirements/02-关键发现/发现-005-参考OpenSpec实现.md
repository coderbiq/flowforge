# 发现-005：参考 OpenSpec 实现

**发现时间**: 2026-05-17 15:30

---

## 核心洞察

tg-proposal 的各阶段行为可以参考 OpenSpec 的实现模式，但需要调整：不使用独立 CLI，用 Skill 实现。

---

## OpenSpec 核心架构

- **Change + Artifacts**：一次变更的容器
- **Schema 驱动**：工件结构和依赖由 schema 定义
- **CLI 作为单一真相源**：所有状态查询通过 `openspec` CLI

---

## 各阶段核心行为

### Propose 模式

| 行为 | OpenSpec | tg-proposal 调整 |
|------|----------|-----------------|
| 获取用户意图 | AskUserQuestion | 同样主动询问 |
| 创建提案目录 | CLI: `openspec new change` | Skill 直接创建 |
| 获取工件依赖 | CLI: `openspec status --json` | 硬编码依赖顺序 |
| 按依赖生成 | 循环创建直到完成 | 分步引导填写 |

### Apply 模式

| 行为 | OpenSpec | tg-proposal 调整 |
|------|----------|-----------------|
| 选择 Change | CLI 参数 | 提案编号参数 |
| 获取上下文 | CLI: `openspec instructions` | 直接读取文件 |
| 实现任务 | 按顺序执行 | 按 Capabilities 拆解 |
| 更新状态 | 更新 tasks.md | 更新 Beads 任务 |

### Archive 模式

| 行为 | OpenSpec | tg-proposal 调整 |
|------|----------|-----------------|
| 检查完成 | 检查 tasks.md | 检查 Beads 任务状态 |
| 归档动作 | 移动到 archive/ | 移动到 completed/ |
| 额外操作 | 无 | 更新模块文档 |

---

## 关键设计要点

1. **Schema 驱动**：工件结构可配置，非硬编码
2. **Fluid Workflow**：非阶段锁定，可随时切换模式
3. **Context 分离**：Skill 的约束不写入工件
4. **用户确认优先**：关键操作都需用户确认

---

## 为什么重要

复用成熟方案的设计模式，避免重新发明轮子，同时根据 tg-workflow 的特点进行调整。
