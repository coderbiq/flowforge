# Adapter Contract

Adapters 和 agent wrappers 必须保留 `task-splitting.md` 里定义的 canonical
task-splitting 规则。

平台集成在适配 `FlowForge` 时，不能重新定义 workflow behavior。

## Adapter 责任

- 暴露平台原生的 commands 或 prompts
- 加载 workflow guides 和 templates
- 按照 `scripts/flowforge-rules-context.js` 和
  `workflow/guides/rule-loading.md` 加载 project rules
- 加载 project configuration
- 更新本地状态快照
- 调用 task 和 memory provider
- 让 install 和 upgrade 的入口与同一套 managed payload 规则保持一致

Adapters 可以提供不同的入口面：

- Claude Code 和 OpenCode：repo-local commands，加上 hooks/plugins
- Codex：project `AGENTS.md` 加上 workflow scripts

## Adapter 不能做的事

- fork lifecycle definitions
- 发明平台专用的 proposal states
- 硬编码 project tags 或 archive rules
- 重新定义 document directory structures

## 必备能力

### Workflow adapter

- `explore(topic)`
- `propose(exploration_path, title)`
- `upgrade(target_project)`
- `apply(proposal_id)`
- `archive(proposal_id)`
- `status(proposal_id)`
- `list()`
- `notes(proposal_id, entry)`

### Task backend adapter

- `create_epic(proposal)`
- `create_tasks(task_map)`
- `query_by_proposal(proposal_id)`
- `close_epic(epic_id)`

### Memory provider adapter

- `store(memory)`
- `search(query, filters)`
- `list_due_reviews()`
- `supersede(memory_id)`

## 配置查找顺序

Adapters 应该按下面顺序解析配置：

1. environment overrides
2. project-local `.flowforge/config.json`
3. user-level adapter config
4. built-in defaults

## Upgrade 行为

- Upgrade 必须保留项目拥有的 `.flowforge/config.json`
- Upgrade 必须保留 `.flowforge/state/`
- 平台 command surface 应该和底层 scripts 共享同样的 upgrade boundary
