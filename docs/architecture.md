# FlowForge 架构与权威模型

FlowForge 是 **Local-First Markdown Artifact Store + Engineering Skills + Deterministic DAG Engine**。目标不是用 CLI 管理长文本或用状态机代替判断，而是让需求、设计和执行上下文完整持久化，同时把确定性计算交给程序。

## 三个组成部分

```text
用户与 Agent
    │
    ▼
flowforge-* Skills ── 分配问题所有权、维护语义权威、实施与审查
    │
    ▼
<docs_dir>/
├── CONTEXT.md       领域词汇与系统不变量
├── adr/             持久架构决策
├── agents/          项目内 Agent 规则
└── proposals/
    └── <feature>/
        ├── requirements.md / design.md / spec.md（按信息价值选用）
        ├── issues/*.md                         （唯一可执行工件）
        └── evidence/*.md                       （需要独立生命周期时）
    │
    ▼
FlowForge CLI ── 工件发现、语义诊断、DAG 校验和 frontier 投影
```

Markdown 是内容接口。Agent 直接读写文件；CLI 不接受大段需求、设计或 ticket 正文。CLI 解析少量机器元数据和人类可见的 `Blocked by`、`Status` 字段，输出可复验的图与诊断结果。

## 权威分工

| 工件角色 | 唯一负责的含义 | 不负责 |
|---|---|---|
| Requirement | 问题、可观察结果、范围、场景、外部约束、术语、需求未知项 | 模块、接口、执行切片 |
| Design | 模块责任、接口/seam、信息流、迁移、替代方案、验证策略 | 改写需求、发布 ticket |
| Ticket | 一个可独立验证的交付增量、局部上下文、动作、约束、验证 | 现场选择架构 |
| Evidence | 已交付行为、真实验证结果、审查处置、偏差、实现引用 | 未来计划或复制 ticket |

正文和语义链接供人阅读并具有权威性。Schema v1 的 ID、revision、area、`consumes`、open item 和 waiver 负责机器追踪。一个语义含义只在一个权威位置维护；下游仅保留当前工作需要的摘要。

## 自适应工件，而非固定阶段

局部、已确定且复用现有 seam 的需求可由一张 compact ticket 承载需求、设计增量、执行和完成证据。跨模块或需要独立评审的内容再提升为独立 authority 文件。`spec.md` 是多权威入口和阅读导航，不凌驾于 requirement/design authority，也不是每个需求必经步骤。

系统不持久化 readiness phase。当前可执行性由这些事实推导：

- DAG edge：真实前置 ticket；
- warning：风险可见，默认允许继续；
- gap：受影响区域缺少决定，默认不进入 frontier，可显式 `--include-gaps`；
- blocker：依赖或外部事实确实阻止执行，override 不能解除；
- exact waiver：只豁免一个诊断和目标，保留原始诊断及理由。

一个区域的 gap 不污染已解决区域。状态字段只表示 ticket 执行生命周期，不表示需求或设计质量。

## 实现边界

- `internal/config`：项目根、`docs_dir` 与兼容配置解析。
- `internal/tracker`：Markdown ticket 解析、Artifact Catalog、语义诊断、DAG 与 frontier。
- `internal/command`：Cobra 命令、策略投影、初始化和受管资产部署。
- `internal/update`：CLI 更新与同版本资产同步。
- `assets/skills`、`assets/agents`、`assets/AGENTS.md`：编译进二进制并部署到目标项目的生产资产。

`docs_dir` 默认为 `docs`，可为相对项目根目录或绝对路径。命令从任意子目录向上定位 `.flowforge/config.yaml`；发现损坏配置时返回错误，不静默退回其他目录。没有 FlowForge 配置的普通目录仍使用 `docs/proposals` 作为兼容默认值。

## 完成不变量

Implement 只有在交付可观察、验证实际运行、Standards/Specification 两轴发现均处理且 evidence 已写入后才关闭 ticket。Catalog 对 closed ticket 缺少非空 `Completion evidence` 发出 `missing-completion-evidence` warning；严格检查将其视为失败。CLI 负责发现事实，Skill 负责理解证据是否真实充分。
