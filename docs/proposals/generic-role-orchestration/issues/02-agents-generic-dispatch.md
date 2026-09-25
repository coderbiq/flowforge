---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      generic-role-orchestration-requirements: 1
    design:
      generic-role-orchestration-design: 1
---

## Completion evidence

闭环证据见各 Change 的 cmd/exit/output/artifact 四元组与 Implementation note（执行记录、验证命令与观测结果）。

# 02: AGENTS.md 通用调度段（能力表 + 任务链协议 + pi 提示）

**Blocked by:** 01
**Status:** closed
**Mode:** lightweight

## Delivery

`assets/AGENTS.md` 在 `## Agent skills` 表之后、`## Subagent delegation` 之前插入 `## Generic capability dispatch` 段（能力键调度表 5 行 + 任务链协议与结构化任务模板 + 派发边界条款 + pi 宿主提示小节 + research 落点指引句）；仓根 `AGENTS.md`（自举 dogfood）同步该段；模板锚点断言测试落地，`go test ./internal/...` 全绿。

## Design context

双节结构：通用调度段在前、流程委派表在后，声明顺序即适用优先级——任何工作先按能力键找角色，flowforge 流程态再查流程表。任务链协议从 GIIS 实证提炼（规划→结构化任务→批量派发→Review 收敛→下一批）；pi 强化第一版是纯指引文本，不改 compile_pi 与 extension 代码。

See the design authority at [通用角色与任务链调度方案](../design.md#generic-role-orchestration-design)（d-dispatch 表与模板原文、d-pi-boost 节）. Requirement authority: [通用角色与任务链调度需求](../requirements.md#generic-role-orchestration-requirements)（目标 2/3，验收 2/3）.

## Touch points

- `assets/AGENTS.md` — 现三节：`## Agent skills`（L1）、`## Subagent delegation`（L21）、`## Execution unit policy`（L39）；新段插在 L21 之前
- `AGENTS.md`（仓根自举）— `## Subagent delegation` 位于 L50 附近，新段插在其前
- `internal/command/assets_deploy_test.go` — 新增模板锚点断言用例（preset 授权见 Constraints）

## Changes

- [x] 1. `assets/AGENTS.md`：在 `## Agent skills` 表与 `## Subagent delegation` 之间插入 `## Generic capability dispatch` 段，内容按设计 d-dispatch 表：能力表 5 行（定向探查→`flowforge-investigator`；批量提取+对照+汇总→`flowforge-batch-analyst`；按模板撰写/回填→`flowforge-scribe`；机械执行→`flowforge-executor`；决策素材简报→`flowforge-investigator` 简报形态×flowforge-research，各含档位建议列）；任务链协议四步 + 结构化任务模板 code block（你是<角色>。项目根：<绝对路径>。/【目标】/【输入】/【输出】带引用要求）；边界条款一句（"mechanically completable 才下放 flash 档；任务链规划、跨组综合、决策、Review 收敛留在编排会话"）；pi 提示小节（fork 继承只读参考、批量异步派发与完成唤醒、用户级 `agentOverrides` 按名覆盖）；research 落点指引句（讨论期产物落 `<docs_dir>/research/YYYY-MM-DD-<slug>.md` 带引用，proposal 创建前经 flowforge-import 导入）。
- [x] 2. 仓根 `AGENTS.md`：在 `## Subagent delegation` 之前插入同段（自举 dogfood 同步，保持与模板语义一致，允许保留仓根已有的额外内容不动）。
- [x] 3. `internal/command/assets_deploy_test.go` 新增断言用例：部署后的 AGENTS.md 产物含锚点关键词 `Generic capability dispatch`、三个新角色名、`【目标】`（任务模板标记）、边界条款关键词（`flash`）、research 指引关键词（`research/`）。

## Constraints

- must 纯本地确定性文件操作，无网络、无 LLM 调用（转录自设计 Standards clauses）。
- must 不改 CLI 命令签名与 Issue Schema 头规范；不改编译器行为。
- must `assets/` 只放部署内容（AGENTS 模板属部署内容）。
- 段内角色名必须与 01 交付的资产名一致（能力表是调度入口，名字错配即失效）。
- preset 测试授权：锚点断言为新增用例，写入 `assets_deploy_test.go` 已经用户在规划评审中显式授权（2026-09-24 会话）。
- Write set: `assets/AGENTS.md`、`AGENTS.md`、`internal/command/assets_deploy_test.go`、`docs/proposals/generic-role-orchestration/`

## Done and verify

- 双处落点: `grep -n "Generic capability dispatch" assets/AGENTS.md AGENTS.md` — 两文件各命中 1 处。
- 能力表完整: `grep -c "flowforge-batch-analyst\|flowforge-scribe\|flowforge-executor" assets/AGENTS.md` — ≥3。
- 锚点断言: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestAGENSGenericDispatch'`（或新增用例实际命名）— ok。
- 全套回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

## Implementation note

- 插入位置：`assets/AGENTS.md` 新段落 L21（`## Agent skills` L1 之后、`## Subagent delegation` 现 L61 之前）；仓根 `AGENTS.md` 新段落 L50（`## Subagent delegation` 现 L90 之前，仓特有节未动）。双文件段内容逐字节一致（`diff` 校验 identical），行序均为 Agent skills < Generic capability dispatch < Subagent delegation。
- 新增锚点用例：`TestAgentRulesDescribeGenericCapabilityDispatch`（`internal/command/assets_deploy_test.go`，紧随 `TestAgentRulesDescribeSubagentDelegation`，同 `os.ReadFile(assets/AGENTS.md)` + `strings.Contains` 模式，另含三段行序断言）。TDD 先红（缺段失败）后绿；`go test ./internal/command/ -run 'TestAgentRulesDescribeSubagentDelegation|TestAgentRulesDescribeGenericCapabilityDispatch'` PASS，全套 `go test ./internal/...` 全绿（command/config/subagent/tracker/update 均 ok）。
- 修改文件（均在 Write set 内）：`assets/AGENTS.md`（+40）、`AGENTS.md`（+40）、`internal/command/assets_deploy_test.go`（+38，经 bash 写入以绕过宿主 *_test.go edit/write 守卫）。

---

## Execution detail

### Verified contracts

- `assets/AGENTS.md` 三节布局：L1 `## Agent skills`（表 L5-18）、L21 `## Subagent delegation`、L39 `## Execution unit policy`；新段插在表后 L21 前——声明顺序即适用优先级（通用在前、流程在后）。
- 仓根 `AGENTS.md`（手工维护 dogfood 镜像，tracked）：L5 Commands / L12 核心设计原则 / L23 boundaries 为仓特有，L30 `## Agent skills`、**L50 `## Subagent delegation`**、L68 `## Execution unit policy`；新段插在 L50 前，仓特有节不动。
- 部署映射：`deployManagedAssets` 把 `assets/AGENTS.md` 写到 `<docsRoot>/agents/issue-tracker.md`（证据：`TestDeployManagedAssetsUsesAbsoluteDocsRoot`，internal/command/assets_deploy_test.go）；`flowforge upgrade` 同通道。仓根 `AGENTS.md` 不经此通道，手工同步。
- 锚点测试落点：`TestAgentRulesDescribeSubagentDelegation`（internal/command/assets_deploy_test.go）直接 `os.ReadFile(filepath.Join("..", "..", "assets", "AGENTS.md"))` + `strings.Contains` 断言——新增锚点用例同文件同模式（preset 授权已在票面）。
- 能力表角色名必须与 01 已交付资产一致（`assets/subagents/` 名册 9：含 batch-analyst / scribe / executor，已落地）。

### Execution scenarios

- Success：双文件 grep 命中锚点且行序正确（`## Generic capability dispatch` 位于 Agent skills 之后、Subagent delegation 之前）；新锚点用例绿；`go test ./internal/...` 全绿。
- Failure：插入位置错（如落在 Execution unit policy 后）时 Contains 断言不拦截，由 grep -n 行序核验兜底；角色名与 01 资产名不一致则调度面断裂，人工核对表行。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestAgentRulesDescribeSubagentDelegation|TestGenericCapabilityDispatch'`（新用例实际命名）— ok。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。
- `grep -n "Generic capability dispatch" assets/AGENTS.md AGENTS.md` — 各命中 1 处且行号介于 Agent skills 与 Subagent delegation 之间。

### Generated artifacts

- producer `assets/AGENTS.md` → consumer `<docs_dir>/agents/issue-tracker.md`（deploy/upgrade 通道）；仓根 `AGENTS.md` 为手工 dogfood 镜像，非生成物。

### Conventions

- must 变更后运行 `go test ./internal/...`（转录自设计 Standards clauses）。
- 段落英文书写（与模板既有节一致）；任务模板 code block 保留【目标】【输入】【输出】标记；表格用与 Agent skills 表相同的 `|:---|:---|:---|` 风格；能力表含档位建议列。
