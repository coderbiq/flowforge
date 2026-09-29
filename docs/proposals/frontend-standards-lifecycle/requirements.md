# frontend-standards-lifecycle — Requirements

> Type: proposal(完整立项,承接 frontend-vision-loop-poc 的通过结论)
> Status: draft
> Prior art: docs/research/ 三轮调研(ui-enforcement-×2、frontend-skill-generality、
> weak-model-execution-accuracy)+ PoC 实证(AP-19 修复 1.12:1→18.43:1)

## 问题

前端规范(设计系统/视觉标准)目前在 flowforge 生命周期中**只在执行末端出现**,
且依赖执行者(视觉 flash 模型)自行读取理解 — 违背本仓已实证的两条原则:

1. **规范前置转录链**(Plan 从 design.md 的 Standards clauses 转录,执行者只读票)— 前端条款不在链内;
2. **智能分层**(理解规范归旗舰,视觉判读归视觉模型,机械遵循归票内条款)— 现状让弱模型承担规范理解。

同时,review 阶段没有前后端模型切分:旗舰 reviewer 是纯文本模型,无法做前端
Spec 轴的视觉核对;PoC 的视觉自审只覆盖执行者自身,缺少独立视觉复核。

## 目标(端到端)

前端规范作为一等公民贯穿五阶段,每一步可追溯(条款编号从 requirement 一路
带到 review 报告),且各阶段由智能匹配的模型承担:

1. **Align**:识别前端涉及面 — 需求的可观察结果涉及 UI/交互/视觉时,requirement
   标记 frontend scope 并引用项目槽位 manifest(本项目存在前端规范体系的事实入口)。
2. **Solution Design**:遵循规范 — design.md 的 Standards clauses 含从槽位切片
   转录的适用前端条款(编号锚点引用);前端变更的验证策略必须声明
   verification-entry。
3. **Plan**:复用现有转录链 — 前端票 Constraints 带前端 must/must not 条款;
   混合变更拆分为前端票/后端票(DAG 边连接);前端票默认 `Mode: lightweight`。
4. **Implement**:按 Write set 切分角色 — 前端票 → frontend-implementer(视觉
   模型,只遵循票内条款 + 按锚点查切片 + 反例×截图批判循环);后端票 →
   implementer(文本 flash)。
5. **Review**:按轴 × 前后端切模型 — Standards 轴始终 reviewer-lite(flash 机械);
   Spec 轴后端票归 reviewer(旗舰);**前端票归 frontend-reviewer(视觉模型,
   对照截图与票内条款做视觉核对)**。

## 范围(In)

- skill v2:`flowforge-frontend-implement` 槽位契约化(design-tokens /
  component-entry / counter-examples / verification-entry,项目自证、web-only、
  adapter 扩展点)
- 新角色:`flowforge-frontend-reviewer`(视觉模型 pin)
- `flowforge-align` / `flowforge-solution-design` skill 增补:前端涉及面识别、
  前端条款转录义务
- AGENTS.md 模板:路由表与生命周期参与说明
- 模型矩阵与 config pin 约定
- 执行第 1 步前移视觉能力自检(能力缺失即 BLOCKED,不延迟到批判步)

## 范围(Out)

- tangram-v2 项目层(805/反例/白名单/lint 链/CI 基线)— 归 tangram-v2 proposal
- 非 web 前端宿主(仅留 adapter 槽位协议)
- Go 侧 tier 体系变更(具名 pin 已够)
- 组件元数据 MCP 化(待静态白名单拦截率数据后再议)

## 验收

1. 端到端演练:一个含前后端的需求 → align 标记 → design 带条款 → 拆双票 →
   双角色执行 → 三轴 review(其中前端 Spec 轴由视觉模型完成),全程条款编号
   可追溯(requirement → design → ticket → review 报告四级一致)。
2. 普适性:槽位契约在 ≥1 个非 tangram-v2 项目以零 skill 改动声明并通过降级语义
   检查(缺槽位 → 第 1 步 BLOCKED 并报告缺口)。
3. 模型切分落地:执行/Review 各轴实际调用模型与矩阵一致(会话记录可审计)。
