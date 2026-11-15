---
flowforge:
  schema: 1
  role: requirement
  id: roadmap-management-requirements
  revision: 1
---
<a id="roadmap-management-requirements"></a>
# Roadmap 需求树管理

## 问题

1. **大需求无机制承载**：FlowForge 的工作单元是 proposal（`<docs_dir>/proposals/<slug>/`），其生命周期假设"先澄清、再设计、再执行"。长期规划（一个大方向拆成多个需求渐进完成）装不进单 proposal：要么硬塞把 authority 撑爆，要么在系统外人工拆成多个平级 proposal。本仓 `ff-wiki/proposals/` 的 27 个 proposal 即人工拆分的活样本：互引全靠散文、已出现死链（如 `backend-arch-refinement`）、无任何汇聚视图。
2. **拆分时机错位**：真实流程是"先粗粒度路线图，节点成熟后再深入分析设计"，proposal 模型要求"全想清楚再立项"，导致拆分被提前到 align 之前、在系统外完成。
3. **proposal 膨胀**：长期项目 proposals 根无限平铺。bytes-ai-space 项目手工维护 close/reopen 目录迁移补丁，未产品化。
4. **wayfinder 未覆盖交付**：fog-of-war 决策地图止步于"路线清晰"，节点从决策走向交付（孵化成 proposal）无衔接。

## 目标

1. **意图层与合同层分离**：新增 roadmap 需求树（无限层级、每节点概要级信息），只管理意图与拆分；proposal 保持"完全独立完整的需求交互"不变，继续承载 align→design→plan→implement 合同。
2. **节点→proposal 单向孵化**：树节点成熟后孵化为 proposal（hatch 动作），节点留链接；需求真相只在 proposal，树不复制真相、不手写状态。
3. **状态派生**：节点/方向进度由链接目标 proposal 的 artifact 存在性与票状态计算，禁止持久化手写状态。
4. **proposal 生命周期命令**：`close` 迁移到 `closed/<年份>/` 归档、`reopen` 迁回，自动重写全库引用链接。
5. **向上汇聚**：`status` 呈现 roadmap 树状汇聚视图（live/closed 分开），供人看整个规划走到哪。

## 边界

- **不修改**：proposal 内部管线（align/solution-design/plan/implement/review 的职责与 artifact 契约）、Issue Schema 头规范、catalog/DAG/frontier 的现有语义。
- **不引入**：树节点间依赖/阻塞（顺序约束放 proposal 层用现有 `Blocked by`，跨 proposal 已支持）；节点 ID；跨 codebase 需求管理。
- **roadmap/ 不进 catalog 扫描**：树是纯意图层，无执行语义；执行上下文（叶子票阅读清单）不感知树的存在。
- wayfinder 保持现状，roadmap 立项后按使用体感另行裁决两者关系。

## 关键裁决（2026-11-14 与用户对齐）

1. **路线**：不动 proposal 内部，外部增加 roadmap 机制。否决"统一 issue 层级模型"（拆掉角色分离承重墙）与"proposal 层级化"（重定义已验证对象语义）。
2. **布局**：proposals 根下不设 active/，工作 proposal 现状不动（存量零迁移）；`closed/` 按年份归档；roadmap 每大方向一文件，完结方向进 `roadmap/closed/`（不按年，膨胀速度差两个量级）。
3. **节点格式**：标题+一句话概要（必填）、描述段（可选，记意图/背景/雾/孵化条件/兄弟约束；红线：不记验收行为与接口细节）、链接行（孵化后 hatch 写入，形式固定 `→ [slug](../proposals/...)`）。无手写状态、无节点 ID、一节点一链接。
4. **孵化**：hatch 建 README 存根（标题、一句话、回链方向文件+节点标题）+ 节点写链接，**不写 requirements.md**（align 自己的规矩：第一个被接受的需求事实出现前不建 authority），交接给 flowforge-align，输入=节点概要+孵化备注（含扫描兄弟节点收集的跨节点约束）。
5. **slug**：自由命名，具体格式在项目中约定。
6. **roadmap skill**：新建 `flowforge-roadmap`（现有 skill 体系），四动作：建图 chart、养树 grow、孵化 hatch、对账 reconcile（死链/孤儿 proposal/多链节点，只报告不强制）。

## 可观察验收

1. roadmap 方向文件（`role: map`）+ 节点链接可被 `status` 解析并呈现树状汇聚；节点状态与链接目标 proposal 实际 artifact/票状态一致（抽测含已交付、实施中、未孵化三类）。
2. `proposal close` 后：目标目录在 `closed/<当年>/`；全库（roadmap 文件 + 各 proposal 正文）对该 slug 的相对链接被重写指向新位置；`flowforge check` 无新增悬空诊断。
3. `proposal reopen` 将其迁回 proposals 根；live 侧同名 slug 已存在时拒绝并报错。
4. `flowforge-roadmap` skill 部署后（`flowforge init`），chart/grow/hatch/reconcile 四动作可被 agent 按说明执行；hatch 产物为 README 存根（无 requirements.md）+ 节点链接。
5. 存量 25 个已 resolved proposal 迁入 `closed/2026/` 后，`flowforge status`/`check`/`frontier` 语义不变（票数、依赖 resolve、诊断集不因迁移增减）。
6. `go test ./internal/...` 全绿，新增命令与 status 视图有测试覆盖。
