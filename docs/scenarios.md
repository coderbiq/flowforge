# FlowForge 场景导览

这 8 个场景每个回答三个问题：**解决什么痛点、你敲什么命令、feature 从哪个 proposal 来**。命令与角色名以当前二进制和 `assets/` 实物为准；新手主线（安装 → 说出需求 → frontier 循环）见 [README 快速开始](../README.md#快速开始)。场景对话均为虚构的通用示例。

## 场景 1：说出一个需求，看它变成 proposal

**痛点**：刚接触 FlowForge，不知道从哪下手——是不是要学一堆 Markdown 模板、建目录、写票？

**答案：都不用。** 你在 agent 会话里用业务语言说出方向，剩下的全部由 agent 完成。起步对话示例（虚构的通用场景）：

> “通知模块上次重构得比较乱，我整理了一版旧设计放在 docs/old-design.md。把消息模板管理拆出去单独立项，再开一个专题修复发送重试的 proposal。”

Agent 会先查代码与现有文档、只问你它答不了的产品取舍，然后建立 `<docs_dir>/proposals/<feature>/`、写下需求与设计事实、用 `flowforge-plan` 生成带 DAG 依赖的执行票，并跑：

```bash
flowforge check      # DAG、authority 与证据诊断（agent 自动跑）
flowforge frontier   # 展示无阻塞就绪批次
```

你推进的方式是说话：**“继续”**（最高频）、“同意，继续推进”、“不太看得明白这些票在做什么”。想自己看进度时才敲 `flowforge status`。

**涉及 feature**：对话驱动立项（import → align → solution-design → plan 全部由 agent 承担）；文件事实驱动的 DAG 检查与 frontier 投影；warning/gap/blocker 三档处置。

**来源**：[documentation-contract-refinement](proposals/documentation-contract-refinement/spec.md)

## 场景 2：从一份旧 PRD 起步

**痛点**：手里有一份几百行的旧 PRD，说不清哪段还算数、哪段是过时设计；直接"翻译成 ticket"会把噪音一起带进来。

**操作路径**：让 agent 进入 `flowforge-import`，它把来源材料分成五类——来源事实 / 需求候选 / 设计决定 / 交付与验证证据 / 未知与冲突——并保留文件与标题定位。需求候选交给 `flowforge-align` 裁决，模块与 seam 选择交给 `flowforge-solution-design`；import 本身不创建任何 authority。

配套命令：

```bash
flowforge assets verify    # 只读核对：项目内受管 skill 与二进制嵌入资产是否漂移
```

**涉及 feature**：外部材料分类交接（不机械转换文档）；受管资产漂移核对。

**来源**：[external-material-intake](proposals/external-material-intake/requirements.md)

## 场景 3：用便宜模型批量执行一个 proposal

**痛点**：整批 ticket 派给 flash 级模型省钱，但弱模型会"勾了 `[x]` 却什么都没做"、在一条失败命令上重试几百次、或把测试断言改掉混过验证。

**操作路径**：

1. `.flowforge/config.yaml` 钉执行者模型（如 `agents.models_by_name` 给 `flowforge-implementer` 配 flash），配 `agents.max_steps` 设迭代预算（默认 200，`-1` 关闭）。
2. `flowforge-plan` 产出三层信息结构的票（人类优先 / 共享执行契约 / agent 细节）；缺五段机器执行契约的票会被 frontier 标 gap，先用 `flowforge-refine-ticket` 补齐。
3. 执行走**轻量模式**：按顺序机械执行未勾选 Changes、机检自检、写 Implementation note 即停——不给弱模型自由发挥空间。
4. 勾选 Change 必须携带证据四元组 `cmd` / `exit` / `output` / `artifact`；`flowforge check --strict` 机械核对，缺证据直接失败。
5. 执行者被阻断时必须先写 `## Blocked evidence`（逐字错误、已试命令与退出码、下一步假设）再返回 `STATUS: BLOCKED`；`refine-ticket` 会把这段失败史转写成 Verified contracts，重派的新上下文不再踩同一个坑。
6. review 的实质发现开出 repair ticket（原票 `needs-repair`、下游等它修完），而不是往原票堆补丁。

`check` 还会对"同一命令非零退出 ≥3 次"报 `evidence-repeat-failure`，在图上点名失控重试。

**涉及 feature**：轻量执行契约；evidence 四元组门禁；steps 预算与 fail-fast；blocked evidence 持久化与消费；repair ticket 分流；refine-ticket 五段契约。

**来源**：[lightweight-execution-contract](proposals/lightweight-execution-contract/spec.md)、[fast-executor-reliability](proposals/fast-executor-reliability/requirements.md)、[ticket-refinement-contract](proposals/ticket-refinement-contract/requirements.md)、[executor-loop-hardening](proposals/executor-loop-hardening/requirements.md)、[blocked-evidence-persistence](proposals/blocked-evidence-persistence/requirements.md)

## 场景 4：角色工厂——一次派一批票

**痛点**：主会话上下文宝贵，不该亲自读 20 个文件做批量提取；每张票都该由干净上下文的执行者跑完再自动验证。

**操作路径**：

```bash
flowforge agents deploy            # 12 角色编译进四宿主
flowforge agents status            # 检查部署漂移（current/missing/drifted）
```

- 与流程无关的活按能力键派通用角色：`flowforge-batch-analyst`（批量提取对比）、`flowforge-scribe`（模板回填）、`flowforge-executor`（机械执行）、`flowforge-investigator`（有界调查）——flash 档可钉。
- PI 宿主一键整批：`flowforge frontier --pi-workflow` 直接输出 pi-subagents 可执行的 workflowScript，每张票一个干净上下文的 implementer，干完自动 `flowforge check`。PI 会话里还能直接调 `flowforge_frontier` / `flowforge_check` 原生工具。
- 有界调查的输出按"设计事实简报"契约交接，旗舰会话只在简报上做裁定。

**涉及 feature**：subagent 生命周期管理（deploy/remove/status）；通用能力角色；PI 一等宿主（原生工具 + 测试文件写拦截）；`--pi-workflow` 批次脚本。

**来源**：[subagent-lifecycle](proposals/subagent-lifecycle/requirements.md)、[generic-role-orchestration](proposals/generic-role-orchestration/requirements.md)、[pi-host-integration](proposals/pi-host-integration/requirements.md)

## 场景 5：限额时段换模型，不换流程

**痛点**：白天旗舰模型被限速，想整体切到便宜模型又怕"指针切了、部署产物没切"的中间态；换了一段时间后说不清到底省了钱还是只把问题推迟。

**操作路径**：

```bash
flowforge model-set list           # 列出命名模型方案
flowforge model-set use offpeak    # 切换：自动重部署 agent 文件，失败回滚
flowforge model-set show           # 查看当前方案
```

- 方案内容留在本机 `.flowforge/config.yaml` 的 `model_sets`：比如 offpeak 把四个决策/审查角色换成限额时段模型、执行者保持 flash。
- 调研、Standards 轴审查这类工作可以常驻 flash：`flowforge-investigator` 产出事实简报、`flowforge-reviewer-lite` 跑规范符合性轴，旗舰只做裁定与 Specification 轴。
- 价值用数据说话：executor 指标管线（`extract` / `report`）从本机会话库只读提取成本/时长/循环指标，按纪元与票难度分档聚合，按预登记的 G1–G4 规则给出"保持 flash / 回退旗舰"判定。

**涉及 feature**：命名模型方案一键切换（自动重部署 + 回滚）；flash 工作面扩展（investigator / reviewer-lite / explore）；执行者价值度量。

**来源**：[model-sets-switching](proposals/model-sets-switching/requirements.md)、[flash-worksurface-expansion](proposals/flash-worksurface-expansion/requirements.md)、[executor-value-measurement](proposals/executor-value-measurement/requirements.md)

## 场景 6：前端票自证"好看"

**痛点**：改一处样式，只能靠读代码推理颜色对不对；旗舰模型看不了图，规范里的对比度条款形同虚设。

**操作路径**：前端票交给 `flowforge-frontend-implementer`（视觉模型角色）：起 dev server → 截图 → 对照票内反例清单逐条批判 → 修复 → 再截图，截图与批判记录随票交付。复核交给 `flowforge-frontend-reviewer`：把交付截图与票内转录的前端条款逐条核对，量化输出（如对比度数值）与截图实际显示交叉验证。

前端规范作为一等公民贯穿全链：Align/Design 阶段把设计系统条款写进 authority，Plan 把条款转录进票（must/must not），执行与复核都由视觉模型角色承担。

**涉及 feature**：前端规范生命周期；截图自审循环；视觉条款级双轴复核。PoC 实测把一处暗底暗字从对比度 1.12:1 修到 18.43:1 才放行。

**来源**：[frontend-standards-lifecycle](proposals/frontend-standards-lifecycle/requirements.md)、[frontend-vision-loop-poc](proposals/frontend-vision-loop-poc/requirements.md)

## 场景 7：规范不靠人肉抄进每张卡

**痛点**：项目规范（分层依赖、错误处理约定）以前靠人抄进每张票，漏抄就漏执行；弱执行者更不会主动去翻规范文档。

**起步（一次性）**：`flowforge init` 部署的 `<docs_dir>/agents/standards.md` 默认只是通用模板。第一次使用前让 AI 重写它——一句话即可：

> “读一下这个项目的分层约定、测试规范和提交规范，把 standards.md 重写成能指导后续所有 ticket 的规范提取说明。”

从此规范进场全程自动：

1. `standards.guide`（默认 `agents/standards.md`）指向你重写后的规范提取说明；
2. `flowforge-align` 挑出与本次需求适用的规范条目交给 Design；
3. `flowforge-solution-design` 把规范转写为 `must` / `must not` 并决定其归属；
4. `flowforge-plan` 只做机械转写进票的 Constraints/Conventions；
5. implementer 开工前发现票里没有规范就退回，不裸奔。

配套的 dispatch 收敛：每个 skill 的 description 必须自带"触发短语 + 下游所有权 + 负边界"（NOT for X），你说"检查 review 发现"不会再误入设计或拆票。

**涉及 feature**：规范提取说明项目定制化（standards.md 重写 + `standards.guide` 指向）；规范注入链（Align 挑选 → Design 转写 → Plan 机械转录 → 缺规范退回）；description-driven dispatch 三段约束。

**来源**：[standards-injection](proposals/standards-injection/requirements.md)、[skill-routing-simplification](proposals/skill-routing-simplification/requirements.md)

## 场景 8：多机协作不互踩

**痛点**：换一台新电脑部署，子代理定义里没有 `model:` 字段，模型分配无法重现；同事 commit 了本机部署产物，你的模型配置被覆盖。

**操作路径**：

- 模型分配只存每台机器自己的 `.flowforge/config.yaml`（`agents.models_by_name` / `agents.models_by_host` 六级优先级链）；`agents deploy` 按配置向各宿主注入模型（PI 宿主也写出 `model:`）。
- `flowforge init` / `upgrade` 自动把部署产物与 config 加入 `.gitignore`（每机器文件）；已入库的存量产物给出 `git rm --cached` 指引，不自动改索引。
- 你手写进已部署 agent 文件的 `model:` 不会再被升级静默抹掉——会被保留并在 stderr 提示。
- wiki 根目录只由 `docs_dir` 一轨决定，旧的 `wiki.root` / `wikiRoot` 键加载时提示一行后被忽略。

**涉及 feature**：模型本地化（models_by_host 六级链）；部署产物 gitignore 自动化；手写 model 保留；wiki 配置单轨。

**来源**：[deploy-artifact-localization](proposals/deploy-artifact-localization/requirements.md)、[agent-model-preservation](proposals/agent-model-preservation/requirements.md)、[wiki-config-single-track](proposals/wiki-config-single-track/requirements.md)
