<!-- BEGIN FLOWFORGE v:0.8 profile:default -->

## FlowForge SKILL 路由

- 新需求、分析、设计、拆分任务 → `flowforge-design`
- 执行任务、继续推进 → `flowforge-implement`
- 归档、沉淀到 library → `flowforge-archive`
- 实施中发现问题、新认知 → `flowforge-feedback`
- 创建/修改 wiki 文档 → `flowforge-docs`

## 任务操作规则

**proposal 任务的创建和状态变更必须通过 FlowForge SKILL**（`flowforge-design` / `flowforge-implement` / `flowforge-feedback`），禁止直接用 `bd create`。

**实施过程中发现的 bug、遗漏、改进**也属于 proposal 的一部分，通过 `flowforge-feedback` 回流，不要脱离 FlowForge 用 `bd create`。

`bd create / update / close` 仅限与任何 proposal 无关的独立事务（如环境配置、工具脚本、临时调研）。不确定是否相关时，默认走 FlowForge。

知识持久化用 `bd remember`。

## 路径约定

脚本 → `.flowforge/scripts/`，schema → `.flowforge/schema/`，guides → `.flowforge/guides/`。

---

以下动作后**必须**激活 `flowforge-progress`：

- 修改 proposal 的 `meta.yaml` status
- 在 task-map.yaml / notes.md 中标记任务或追加日志
- 创建、归档或移动 proposal 目录

### 会话收尾

1. 质量门禁通过（测试、lint、构建）
2. `task-sync.js --check` 确认一致
3. `git pull --rebase && bd dolt push && git push`

<!-- END FLOWFORGE -->
