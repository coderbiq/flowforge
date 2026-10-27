# frontend-standards-lifecycle — Observations

> 2026-09-29 端到端演练验收(验收 1)。载体:tangram-v2 的
> frontend-enforcement-toolchain proposal(6 票,06:41 立项 → 09:14 全闭)。

## 验收结论:通过

| 验收项 | 结果 |
|---|---|
| 端到端演练(五阶段) | ✓ align→design→plan→refine→三波执行→三轴 review→Fix 循环→收口,全程未离场 |
| 条款四级可追溯 | ✓ requirement → design(Standards clauses,d- 前缀域)→ ticket Constraints(带语义锚)→ review 报告(Fix 挂账带编号) |
| 模型矩阵落地 | ✓ 13+ runs 全部与矩阵一致(见下) |
| 守卫实战 | ✓ 编排首试 workflowScript 三路裸脚本 fan-out → block → 改三次具名派发 |
| 视觉轴 | ✓ 04/05 双票触发,gemini 判读 16+16 项全 PASS + 量化对照 + 残余风险识别;06 无截图 → 正确不启用 |
| 普适性(验收 2) | ⏳ 未做(需非 tangram-v2 项目,如 GIIS 前端) |

## 模型矩阵实录

| 环节 | Run | 模型 | 符合矩阵 |
|---|---|---|---|
| planner(6 票 DAG) | 07035d85 | glm-5.3 | ✓ |
| refine(六票 Tier 3) | a0815fd9 | deepseek-v4.1-flash | ✓ |
| 01/02/03/04/05/06 实现 + Fix | ×7 | deepseek-v4.1-flash | ✓ |
| 双轴 review(Spec 轴) | ×5 | glm-5.3 | ✓ |
| 视觉判读(04/05) | ×2 | gemini-3.8-flash-high | ✓ |

## 流程韧性实录(演练中的真实事件)

1. **三层独立纠错**:planner 核验 Touch points → refine 纠 planner 数据(catalog 26→实测 23/30 键)→ review 纠票面(DoV 与实测偏差)。
2. **三次 supervisor 裁决**:axe 白名单口径(B+C 合体不掏空规则)、vitest/e2e 域分离(授权扩写集)、outputDir 计划重叠(预接线消解)。全部"先补权威再裁决"。
3. **并行写集互斥**:02 见 03 规则序 44 error 未越写集补偿,归因实证后等 03 自修。
4. **故障恢复**:网络瞬断致 05 视觉判读 run failed(报告残缺 769B)→ 核实不可采信 → resume 续写完整(16 项 PASS)。06 run "complete 但契约未履行" → resume → 挖两层根因,拦截率 0/3→3/3。
5. **review 真咬**:06 六项 Fix(六票唯一非零,证据口径类)→ Round 2 收口。

## 已知偏差(待改进,记录于 design 或后续票)

1. review 双轴在 reviewer 单角色内串行(子代理不可嵌套)→ Standards 轴未用 flash lite,成本优化机会:编排会话拆派 reviewer-lite。
2. 快速视觉判读曾派 frontend-implementer(模型正确但角色应为 frontend-reviewer)。
3. async 完成唤醒丢失一次(需人工推动续跑)— pi 侧观察项。

## gemini-3.8-flash-high 质量评估(2026-09-29)

全历史 3 runs(NUC+本地):PoC 修复 54msg PASS / 04 判读 PASS(量化对照+只读纪律) / 05 判读(网络断→resume 完整)。**3/3 约束遵循,0 模型性违反,判定不换**(用户确认)。亮点:hex×阈值量化对照习惯、跨 OS AA 残余风险识别。

## 结论

五阶段规范贯穿 + 前后端模型切分 + 三轴 review 的完整体系在真实 proposal 上闭环跑通,
零人工干预实现(推动仅 3 次:指令投递/唤醒丢失/评估要求)。验收 2(普适性)待 GIIS
或其他项目实测后关闭本 proposal。
