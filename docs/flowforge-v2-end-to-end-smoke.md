# FlowForge v2 端到端 Smoke 清单

> 日期：2026-06-15  
> 用途：重建 `ff-wiki` 前后验证 FlowForge v2 是否具备完整自举闭环。

本文不是自动化测试脚本，而是一套固定的人工/Agent 执行清单。目标是验证：

- 参考资料可以先进入 library。
- proposal 可以从需求分析进入设计和任务拆分。
- 任务可以从 library 获取规范/模块上下文。
- 执行过程中产生的发现可以沉淀回 library。
- 所有卡片关系、内部导航和健康检查由 CLI 管理。

## 1. 前置约束

- 不直接编辑 `ff-wiki` 下的任何文件。
- 不手写内部卡片链接、frontmatter 或 wikilink。
- 内部导航只通过 `flowforge card refresh` 或其它 FlowForge CLI 生成。
- 手写 Markdown 链接只允许用于外部资料引用。
- 每个关键步骤后运行 `flowforge proposal inspect <proposal-id>` 查看 `Health Issues`。

## 2. 准备命令

```bash
make dev
flowforge --version
go test ./internal/...
git diff --check
```

预期：

- Go 测试通过。
- 没有 diff 空白错误。

## 3. 初始化项目

如果是在临时目录验证：

```bash
flowforge init --yes
flowforge project create default
flowforge project current
```

如果是在当前仓库验证，确认已有项目即可：

```bash
flowforge project current
```

预期：

- 当前 project 已设置。
- `ff-wiki` 目录结构由 CLI 生成。

## 4. 导入初始 Library 规范

先用结构化候选模拟“外部 source material 已被 SKILL 拆分”的结果：

```bash
flowforge proposal create "FlowForge v2 smoke"
flowforge proposal current
```

记录 proposal ID，下文用 `<proposal-id>` 表示。

创建一个可追溯的 source finding：

```bash
flowforge card create \
  --type finding \
  --title "Source material identifies service rule" \
  --body "## Finding\n\nService work must validate inputs before execution." \
  --proposal <proposal-id>
```

记录 finding ID，下文用 `<source-finding-id>` 表示。

导入 library convention：

```bash
flowforge library import \
  --type convention \
  --title "Validate inputs before service execution" \
  --body "## Rule\n\nService implementation must validate inputs before executing state changes.\n\n## Applies When\n\n- Implementing service behavior\n- Handling user supplied input" \
  --source-card <source-finding-id> \
  --tags layer:service,scenario:validation \
  --domain service
```

检查：

```bash
flowforge library facets
flowforge validate all
```

预期：

- library convention 写入 `02-library`。
- 新 convention `references -> <source-finding-id>`。
- `validate all` 通过。

## 5. 创建需求索引与需求卡

```bash
flowforge card create \
  --type requirement \
  --title "Smoke task can use library constraints" \
  --body "## Summary\n\nImplementation tasks should discover and link relevant library constraints.\n\n## Source\n\nSmoke validation.\n\n## Acceptance\n\n- Task context includes linked library convention.\n- Proposal inspect has no health issue for this requirement.\n\n## Scope\n\nFlowForge CLI workflow only.\n\n## Open Questions\n\nNone" \
  --proposal <proposal-id>
```

记录 requirement ID，下文用 `<req-id>` 表示。

```bash
flowforge structure add STR-<proposal-id>-REQ <req-id>
flowforge proposal inspect <proposal-id>
```

预期：

- requirement 已进入顶层需求索引。
- 没有 “requirement is not reachable from a requirement index”。

## 6. 创建设计与任务

创建 design：

```bash
flowforge card create \
  --type design \
  --title "Task context links library convention" \
  --body "## Goal\n\nEnsure implementation tasks carry precise library constraints.\n\n## Decision\n\nUse library suggest/classify to discover candidates, then link confirmed conventions to tasks.\n\n## Rationale\n\nTask execution should not load broad SKILL text.\n\n## Constraints\n\nConfirmed library cards only.\n\n## Impact\n\nTask context becomes precise and bounded.\n\n## Verification\n\ncontext task shows the linked convention.\n\n## Follow-up Tasks\n\nCreate implementation task." \
  --proposal <proposal-id> \
  --links <req-id>:designs
```

记录 design ID，下文用 `<des-id>` 表示。

创建 implementation task：

```bash
flowforge task create \
  --title "Verify task context includes convention" \
  --type i \
  --status not_ready \
  --body "## Goal\n\nVerify task context includes a linked library convention.\n\n## Inputs\n\n- Requirement\n- Design\n- Library convention\n\n## Deliverables\n\n- Smoke validation result\n\n## Acceptance\n\n- context task lists the convention in Stable Context Cards.\n\n## Out of Scope\n\n- Product code changes\n\n## Read Before Work\n\n- Requirement and design cards" \
  --links <des-id>:implements,<req-id>:satisfies
```

记录 task ID，下文用 `<task-id>` 表示。

刷新导航：

```bash
flowforge card refresh <req-id>
flowforge card refresh <des-id>
flowforge proposal inspect <proposal-id>
```

预期：

- requirement 的 `FlowForge Navigation` 展示 design 和 task。
- design 的 `FlowForge Navigation` 展示 implementation task。
- 没有 REQ/DES navigation stale/missing 的 Health Issue。

## 7. 查询 Library 并关联任务

```bash
flowforge library facets
flowforge library classify --for <task-id>
flowforge library suggest --for <task-id> --facet layer:service --facet scenario:validation --types convention,module
flowforge card read <convention-id> --summary
flowforge task link-add <task-id> <convention-id>:constrains
```

将任务改为 ready：

```bash
flowforge card update <task-id> --status ready
flowforge proposal inspect <proposal-id>
flowforge context task --task <task-id>
flowforge validate all
```

预期：

- `context task` 的 Stable Context Cards 包含 `<req-id>`、`<des-id>`、`<convention-id>`。
- ready implementation task 没有 “no linked convention constraints”。
- `validate all` 通过。

## 8. 记录过程并沉淀发现

```bash
flowforge log create \
  --kind progress \
  --title "Smoke task context verified" \
  --for <task-id> \
  --summary "Task context includes requirement, design, and library convention."
```

创建可复用 finding：

```bash
flowforge card create \
  --type finding \
  --title "Task context should stay card-scoped" \
  --body "## Finding\n\nTask execution context should be assembled from directly linked cards and backlink evidence, not broad library scans." \
  --proposal <proposal-id> \
  --links <task-id>:discovers
```

记录 finding ID，下文用 `<finding-id>` 表示。

沉淀到 library：

```bash
flowforge library promote <finding-id> \
  --type convention \
  --title "Keep task context card-scoped" \
  --tags context:task,scenario:execution
flowforge validate all
```

预期：

- library 中新增 convention。
- 新 convention `references -> <finding-id>`。
- 原 proposal finding 保留为追溯证据。

## 9. 最终检查

```bash
flowforge proposal inspect <proposal-id>
flowforge context proposal --proposal <proposal-id>
flowforge context task --task <task-id>
flowforge validate all
```

预期：

- `Health Issues` 为 `None`，或只剩明确可解释的非阻塞项。
- `context proposal` 可以从 requirement map 进入焦点。
- `context task` 能给出精准上下文，不需要读取整个 library。
- `validate all` 通过。

## 10. 失败处理

如果出现失败：

- 断链：优先修 CLI 写入路径，避免手工改卡片。
- 缺索引：使用 `flowforge structure add`。
- 缺导航：使用 `flowforge card refresh`。
- 缺规范：使用 `flowforge library suggest` 后确认并 `task link-add`。
- 缺沉淀：使用 `flowforge library promote` 或 `library import`。

不要在失败后直接修补 `ff-wiki` 文件。
