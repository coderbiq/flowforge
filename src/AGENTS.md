<!-- BEGIN FLOWFORGE v:0.14.0 profile:default -->

## FlowForge SKILL 路由

- 新需求、分析、设计、拆分任务 → `flowforge-design`
- 执行任务、继续推进 → `flowforge-implement`
- 归档、沉淀到 library → `flowforge-archive`
- 实施中发现问题、新认知 → `flowforge-feedback`
- 创建/修改 wiki 文档 → `flowforge-docs`
- 检查 library 健康 / 诊断过期断链 → `flowforge-library-doctor`
- 探索记录 / 优化重组 library 内容 → `flowforge-library-surgeon`
- 定期维护 / 持续监控 library → `flowforge-library-keeper`

## 任务操作规则

**所有任务操作通过 `flowforge task` CLI，严禁直接操作后端存储。**

- ❌ 禁止读写 `tasks.snapshot.md`（自动生成快照）
- ✅ 常用命令：`flowforge task status` 查看 | `ready/claim/done` 执行
- 📖 任务层级、完整命令和编写规范见 `.flowforge/guides/`

## CLI 入口

```bash
flowforge task status --proposal <CR-id>   # 任务状态
flowforge task ready --proposal <CR-id>    # 就绪任务
flowforge task claim --proposal <CR-id> <id>  # 认领
flowforge task done --proposal <CR-id> <id>   # 完成
```

---

以下动作后**必须**激活 `flowforge-progress`：

- 通过 `flowforge task` 完成任务操作
- 在 notes.md 中追加日志
- 创建、归档或移动 proposal 目录

### 会话收尾

1. 质量门禁通过（测试、lint、构建）
2. `git pull --rebase && git push`

<!-- END FLOWFORGE -->
