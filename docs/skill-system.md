# mattpocock 方法论体系

FlowForge 全量采用 `mattpocock/skills` 的敏捷方法论规范。

---

## 1. Skill 双层架构

### User-Invoked（编排入轨）
直接由用户在对话中调起或按需调度的顶层工作流：
- `ask-matt`：路由与技能指引
- `triage`：持续性 Issue 分诊与状态机流转
- `grill-with-docs`：需求面试并当场建立领域模型与 ADR
- `to-spec`：将对话综合为功能规格（Spec）
- `to-tickets`：将 Spec 拆解为垂直切片（Tickets）并声明 blocking edges
- `implement`：在预定 seams 上驱动交付循环
- `wayfinder`：高不确定性复杂任务的迷雾决策地图
- `improve-codebase-architecture`：代码库架构深化扫描
- `handoff`：跨会话对话压缩与工作交接

### Model-Invoked（可复用纪律原语）
被编排技能自动调度的原子纪律：
- `grilling`：非阻塞的决策树前沿面试原语
- `domain-modeling`：当场创建/更新 `CONTEXT.md` 术语表与 ADR
- `tdd`：Red-Green-Refactor 驱动
- `code-review`：Standards 与 Spec 双轴隔离并行审查
- `diagnosing-bugs`：假设驱动的缺陷根因分析
- `research`：后台子 Agent 高信任源事实调查
- `prototype`：一次性技术探针验证
- `codebase-design`：模块深度与接口边界设计

---

## 2. 工件与流转规范

1. **工件物理化**：所有 Spec 与 Ticket 物理存在于 `<docs_dir>/proposals/<feature>/` 下的独立 Markdown 文件中。
2. **知识当场建**：在 grilling 与 triage 过程中即刻将领域名词写入 `<docs_dir>/CONTEXT.md`，重要架构决策写入 `<docs_dir>/adr/`，不再延后合流。
3. **依赖图驱动**：切片通过 `Blocked by:` 建立 DAG，由 `flowforge frontier` 驱动执行。
