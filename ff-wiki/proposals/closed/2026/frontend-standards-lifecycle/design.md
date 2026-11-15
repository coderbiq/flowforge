# frontend-standards-lifecycle — Design

## 总原则(三条,全部有实证依据)

1. **智能分层**:理解规范 → 旗舰(设计/Plan);视觉判读 → 视觉模型(执行自审/
   review 视觉轴);机械遵循 → 票内条款(任何执行者)。(依据:weak-model
   调研 — 结构化输出合规率 100% vs <40%;转录链 SKILL 原文)
2. **证据机器可校验优先**:视觉结论必须落到可 grep 的量化输出(PASS/FAIL+数值),
   截图是证据附件而非唯一判据。(依据:PoC 的对比度量化输出形态)
3. **确定性兜底**:缺槽位/缺能力/缺入口 → 显式 BLOCKED 上报,不静默降级。
   (依据:治理调研 — 提示≠强制)

## 生命周期数据流

```
槽位 manifest(项目声明:本仓有前端规范体系)
   │
   ▼
[Align] requirement.md
   frontend scope 标记 + manifest 引用(识别,不读全文)
   │
   ▼
[Solution Design] design.md
   ## Standards clauses += 前端条款(编号锚点转录自槽位切片,强模型完成)
   验证策略声明 verification-entry
   │
   ▼
[Plan] ticket(frontend / backend 分票,DAG 边连接)
   Constraints += 前端 must/must not(现有转录链,零新机制)
   Mode: lightweight(默认)+ Write set
   │
   ▼
[Implement]
   前端票 → frontend-implementer(视觉模型)
     遵循票内条款 → 按锚点查槽位切片 → 反例×截图批判循环 → 量化输出
   后端票 → implementer(文本 flash)
   │
   ▼
[Review] 三轴
   Standards 轴(两类型通用)→ reviewer-lite(flash:lint/build/机械条款)
   Spec 轴·后端 → reviewer(旗舰)
   Spec+视觉轴·前端 → frontend-reviewer(视觉模型:截图×票内条款核对)
```

## 角色与模型矩阵

| 阶段/轴 | 角色 | 模型档 | 依据 |
|---|---|---|---|
| Align/Design/Plan | analyst/architect/planner | 旗舰(文本) | 规范理解与转录 |
| 执行·后端 | implementer | flash(文本) | 已验证 |
| 执行·前端 | frontend-implementer | **flash(视觉)** gemini 系 | PoC 验证 |
| Review·Standards 轴 | reviewer-lite | flash(文本) | 机械可 grep |
| Review·Spec 轴·后端 | reviewer | 旗舰(文本) | 契约核对 |
| Review·Spec+视觉·前端 | **frontend-reviewer(新)** | **flash(视觉)** | 截图判读 |

关键裁决:**Spec 轴前端 review 不给旗舰** — glm-5.3 纯文本无法看图;给视觉
flash 是智能匹配(判读截图+对照票内条款是低理解负荷、高视觉负荷任务)。
旗舰若未来获得视觉能力,config 一行切换,角色不变。

## 槽位契约 v2(声明式,项目自证)

```
.agents/references/frontend/manifest.md     ← 唯一入口(声明下列槽位存在)
  design-tokens        token 源(推荐 DTCG,不强制)
  component-entry      组件白名单 + 自检命令(项目自证,不规定 lint 实现)
  counter-examples     编号反例清单(任意形态,批判对照用)
  verification-entry   可重复视觉验证命令 + 可 grep 的 PASS/FAIL 量化输出契约
```

- skill 不假设 Playwright/antd/React;web-only;非 web 为 adapter 扩展点
- 执行者读槽位的方式:**按票内条款锚点查切片**,不全文通读
- 第 1 步前移三项自检:槽位齐备 / 运行时图像能力 / verification-entry 存在 —
  任一缺失立即 BLOCKED(修 PoC 的"延迟失败"瑕疵)

## Review 视觉轴的工作方式(frontend-reviewer)

1. 读票(条款+量化输出)+ 槽位反例切片
2. 取执行者交付的截图(路径在票的 Verification 段)
3. 逐条款视觉核对(票内条款 × 截图),发现 → Fix: Change 回票
4. 量化输出交叉核对(数字是否与截图实际一致 — 防执行者伪造量化)
5. 产出:逐条款判定 + 截图引用 + 可 grep 结论

## 混合需求拆票规则(Plan)

- Write set 命中前端目录(JSX/样式/组件包)→ 前端票;其余 → 后端票
- 同一 feature 的前后端票用 DAG 边连接(接口契约为依赖锚)
- 单一前端票且无后端耦合 → 不拆

## 降级语义(零配置新项目)

沿用已确认矩阵:机制层全运行;执行层在检查点显式 BLOCKED 并报告缺口;
后端票零影响。启用 = 三件声明(槽位/入口/pin)+ deploy。

## 变更资产清单

| 资产 | 变更 |
|---|---|
| assets/skills/flowforge-frontend-implement | v2:槽位契约化+条款优先+前置自检 |
| assets/subagents/flowforge-frontend-reviewer.md | 新增(视觉 pin) |
| assets/skills/flowforge-align | 增补:frontend scope 识别+manifest 引用 |
| assets/skills/flowforge-solution-design | 增补:前端条款转录义务+verification-entry 声明 |
| assets/skills/flowforge-review | 增补:前端票三轴编排(视觉轴指向 frontend-reviewer) |
| assets/AGENTS.md | 路由表 + 生命周期参与说明 |
| (无 Go 变更) | 具名 pin + 现有转录链已覆盖 |

## 风险与对策

1. 视觉 flash 判读质量波动 → 量化输出兜底(数字是硬证据)+ review 视觉轴独立复核
2. 条款转录成本(每张前端票)→ 条款切片表复用,Plan 只引用编号
3. 角色数膨胀(frontend-reviewer 是第 12 个)→ 换来模型智能匹配,且 config
   可按项目裁剪(hosts/models_by_name 不配即不激活)
4. review 视觉轴与执行自审重复 → 自审是修复循环(执行者内部),review 是独立
   核对(不同上下文),参照"对抗式 review 用新上下文"实证

## 开放决策点(留给 align 阶段裁决)

1. frontend scope 标记进 requirement.md 的具体 schema(字段名/位置)
2. 视觉轴是否对 Standards 轴也开放(现状:Standards 轴保持纯机械,视觉条款
   归 Spec 轴 — 还是允许 lite 也看图?)
3. frontend-reviewer 是否同时承担"品味类"目标(英文直出/信息密度,调研确认
   无自动化门的那类)— 建议首期 No(人审保留)
