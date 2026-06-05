<!-- BEGIN FLOWFORGE v:0.10.0 profile:default -->

## FlowForge SKILL 路由

- 新需求、分析、设计、拆分任务 → `flowforge-design`
- 执行任务、继续推进 → `flowforge-implement`
- 归档、沉淀到 library → `flowforge-archive`
- 实施中发现问题、新认知 → `flowforge-feedback`
- 创建/修改 wiki 文档 → `flowforge-docs`

## 任务操作规则

**所有任务操作必须通过 `flowforge task` CLI，严禁直接读写任务文件。**

- ❌ **禁止** 读取 `tasks.snapshot.md` —— 这是自动生成的只读快照，供人类 git diff 审查
- ❌ **禁止** 读取 `task-map.yaml` —— v0.9 已废弃，任务数据在 beads 后端
- ❌ **禁止** 直接用 `bd create/update/close` 操作 proposal 任务
- ✅ **必须** 使用 `flowforge task status/ready/claim/done` 等命令
- ✅ `bd create/update/close` 仅限与任何 proposal 无关的独立事务
- ✅ 知识持久化用 `bd remember`

任务查询命令：

```bash
flowforge task status --proposal <id>      # 全部任务状态（含 byType 分组）
flowforge task ready --proposal <id>       # 就绪任务列表
flowforge task blocked --proposal <id>     # 阻塞任务列表
```

## CLI 入口

项目根目录 `flowforge` 是统一入口。常用命令：

```bash
flowforge task ready --proposal <CR-id>     # 就绪任务
flowforge task claim --proposal <CR-id> <id> # 认领任务
flowforge task done --proposal <CR-id> <id>  # 完成任务
flowforge task status --proposal <CR-id>     # 状态概览
flowforge implement-context [CR-id]           # 加载实施上下文
flowforge design-context [CR-id]              # 加载设计上下文
flowforge task --help                         # 任务管理帮助
```

---

以下动作后**必须**激活 `flowforge-progress`：

- 修改 proposal 的 `meta.yaml` status
- 通过 `flowforge task` 完成任务操作
- 在 notes.md 中追加日志
- 创建、归档或移动 proposal 目录

### 会话收尾

1. 质量门禁通过（测试、lint、构建）
2. `git pull --rebase && bd dolt push && git push`

<!-- END FLOWFORGE -->
