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

# 02: AGENTS.md 通用调度段（能力表 + 任务链协议 + pi 提示）

**Blocked by:** 01
**Status:** open
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

- [ ] 1. `assets/AGENTS.md`：在 `## Agent skills` 表与 `## Subagent delegation` 之间插入 `## Generic capability dispatch` 段，内容按设计 d-dispatch 表：能力表 5 行（定向探查→`flowforge-investigator`；批量提取+对照+汇总→`flowforge-batch-analyst`；按模板撰写/回填→`flowforge-scribe`；机械执行→`flowforge-executor`；决策素材简报→`flowforge-investigator` 简报形态×flowforge-research，各含档位建议列）；任务链协议四步 + 结构化任务模板 code block（你是<角色>。项目根：<绝对路径>。/【目标】/【输入】/【输出】带引用要求）；边界条款一句（"mechanically completable 才下放 flash 档；任务链规划、跨组综合、决策、Review 收敛留在编排会话"）；pi 提示小节（fork 继承只读参考、批量异步派发与完成唤醒、用户级 `agentOverrides` 按名覆盖）；research 落点指引句（讨论期产物落 `<docs_dir>/research/YYYY-MM-DD-<slug>.md` 带引用，proposal 创建前经 flowforge-import 导入）。
- [ ] 2. 仓根 `AGENTS.md`：在 `## Subagent delegation` 之前插入同段（自举 dogfood 同步，保持与模板语义一致，允许保留仓根已有的额外内容不动）。
- [ ] 3. `internal/command/assets_deploy_test.go` 新增断言用例：部署后的 AGENTS.md 产物含锚点关键词 `Generic capability dispatch`、三个新角色名、`【目标】`（任务模板标记）、边界条款关键词（`flash`）、research 指引关键词（`research/`）。

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

---

## Execution detail

### Verified contracts

- <filled by flowforge-refine-ticket>

### Execution scenarios

- <filled by flowforge-refine-ticket>

### Expected tests

- <filled by flowforge-refine-ticket>

### Generated artifacts

- <filled by flowforge-refine-ticket>

### Conventions

- must 变更后运行 `go test ./internal/...`（转录自设计 Standards clauses）。
- must 通用角色 description 按能力书写、不含流程术语（能力表行同理）。
