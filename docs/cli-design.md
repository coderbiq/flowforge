# FlowForge CLI 设计规范

FlowForge CLI 是一个专为 `mattpocock/skills` 设计的轻量级辅助计算器。它严格遵循 **"文件负责内容，CLI 负责图计算"** 的原则。

---

## 1. 核心命令集

### `flowforge init [path]`
一键初始化目标项目：
- 部署 `docs/agents/issue-tracker.md`、`docs/agents/domain.md`、`docs/agents/triage-labels.md`
- 部署 `assets/skills/` 下的 18 个 mattpocock 官方 Skill 至 `.agents/skills/`
- 创建 `CONTEXT.md` 模板与 `.scratch/` 目录结构
- 注入 `AGENTS.md` 中的 skills 配置块

### `flowforge frontier [--dir .scratch] [--json] [--quiet]`
扫描 `.scratch/` 下的所有 issues，计算 DAG 依赖图，输出当前无阻塞且未被认领的就绪任务。
- `--json`: 输出结构化 JSON 格式
- `--quiet` (`-q`): 仅输出就绪任务的文件路径，便于脚本/Agent 直接读取

### `flowforge check [--dir .scratch]`
静态检查当前依赖图的健康状态：
- 循环依赖检测（Cycle Detection）
- 悬空依赖检测（Dangling Dependency，引用的前置 Ticket 不存在）
- 自依赖检测（Self Dependency）

### `flowforge status [--dir .scratch]`
汇总当前所有 Feature 和 Wayfinder 项目的完成度、阻塞任务数与就绪任务数。

### `flowforge version` & `flowforge upgrade`
- `flowforge version`: 输出当前版本与构建元数据
- `flowforge upgrade`: 自更新至最新 GitHub Release 版本

---

## 2. 设计铁律

1. **零内容 API**：禁止提供 `issue create --body "..."` 等需要通过命令行传长文本的接口。所有工件内容的创建与编辑由 Agent 直接通过 `write`/`edit` 工具操作 Markdown 文件。
2. **Token 精简**：CLI 默认输出必须精炼，避免大模型读取命令行输出消耗过多上下文。
3. **确定性优先**：所有拓扑排序、环检测由 Go 内部编译级算法完成，彻底避免大模型在纯文本上下文中的拓扑推导幻觉。
