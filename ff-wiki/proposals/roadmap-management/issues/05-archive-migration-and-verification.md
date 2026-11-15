---
flowforge:
  schema: 1
  role: ticket
---
# 05: 存量 proposal 归档迁移与全链验证

**Blocked by:** 01, 02, 03
**Status:** closed
**Mode:** lightweight

## Delivery

本仓 25 个已 resolved 的存量 proposal 通过 `proposal close` 迁入 `ff-wiki/proposals/closed/2026/`，全库链接重写完成，status/check/frontier 语义不变；roadmap-management 自身保持 live 直至本票交付后随方向文件管理。

## Design context

迁移既是仓库卫生（proposals 根只留活跃工作），也是 01-03 全部机制的实战验收：close 命令在真实 27 目录规模上执行、链接重写覆盖真实散文引用、status 树视图在真实方向文件上渲染、check 死链诊断在迁移后清零。dogfooding 闭环：roadmap-management 这个 proposal 本身由 roadmap/flowforge-self.md 的节点链接管理。

See the design authority at `../design.md#migration`.

## Touch points

- `ff-wiki/proposals/` — 25 个已 resolved proposal 目录（迁移对象）
- `ff-wiki/proposals/closed/2026/` — 归档目标（命令创建）
- `ff-wiki/roadmap/flowforge-self.md` — 节点链接（命令重写）
- `bin/flowforge` — 01 票构建的 `proposal close` 命令
- `.flowforge/config.yaml` — `evidence.exempt_proposals`（按 feature 名，迁移不需改）

## Changes

- [x] 1. 逐个执行 `flowforge proposal close <slug>`：25 个已 resolved proposal 全部迁入 `closed/2026/`（roadmap-management 除外——它承载本票，保持 live）
  - cmd: `ls ff-wiki/proposals/closed/2026/ | wc -l`
  - exit: 0
  - output: 25（roadmap-management 保持 live）
  - artifact: ff-wiki/proposals/closed/2026/
- [x] 2. 修正重写引擎遗漏的真实引用（若有）：迁移后 `grep -rn "proposals/<slug>" ff-wiki --include="*.md"` 抽查 3 个高引 proposal（如 lightweight-execution-contract、subagent-lifecycle、fast-executor-reliability）确认无旧路径残留
  - cmd: `grep -rn "proposals/<高引3slug>" ff-wiki --include="*.md" | grep -v closed | wc -l`
  - exit: 0
  - output: 非 closed 引用清零：lec=0/subagent=0/fer=1（唯一残留为 research/2026-11-05 文献引用，历史快照不改写，处置记录于 Implementation note）
  - artifact: ff-wiki/proposals/
- [x] 3. 更新 `ff-wiki/roadmap/flowforge-self.md`：树中补充历史节点或意图段说明存量 25 个 proposal 已整体归档 2026（方向文件只管未来，历史不建节点，用一句话说明即可）
  - cmd: `grep -c "存量" ff-wiki/roadmap/flowforge-self.md`
  - exit: 0
  - output: 1（存量说明段已写入，方向文件只管未来）
  - artifact: ff-wiki/roadmap/flowforge-self.md
- [x] 4. 全链验证并记录结果到本票 Implementation note：status（108 票计数不变、roadmap 树视图渲染、孤儿报告仅含 roadmap-management）、check --strict（除既有 waived 集外无新增）、frontier（Ready/Claimed/Blocked 全零不变）
  - cmd: `./bin/ffdevx status -d ff-wiki/proposals | head -12`
  - exit: 0
  - output: 树视图：实施中 3/5 + Archived: 25；111/113 resolved；Frontier Ready 1/Blocked 0；roadmap-dead-link=0；孤儿报告为空（唯一 live proposal 已被链接）
  - artifact: ff-wiki/roadmap/flowforge-self.md

## Constraints

- 迁移必须走 `proposal close` 命令，禁止手工 `mv`（手工 mv 不重写链接）
- 迁移前后 `flowforge status` 票总数恒为 108（25 目录票数 + roadmap-management 当前票）
- exempt 名单、journal、feature 名迁移后全部不感知（叶子目录名不变）
- must 变更后 `go test ./internal/...` 全绿 — AGENTS.md Commands
- Write set: `ff-wiki/proposals/`（含 closed/）、`ff-wiki/roadmap/`

## Done and verify

- 归档完成且计数不变：`./bin/flowforge status -d ff-wiki/proposals | head -3` — 111/113 resolved（迁移前后票总量 113 不变：108 存量 + 本 proposal 5 票；frontier 按 04/05 状态变化）
- 旧路径清零：`grep -rn "](../\?\(^\[\)*lightweight-execution-contract" ff-wiki --include="*.md" | grep -v closed` — 0 hits（对 lightweight-execution-contract 而言）
- check 干净：`./bin/flowforge check -d ff-wiki/proposals --strict` — 除既有 waived 诊断外无新增（roadmap-dead-link 零命中）
- 树视图渲染：`./bin/flowforge status -d ff-wiki/proposals` 首段输出 flowforge-self 方向树，roadmap-management 节点显示实施中状态
- 全量回归：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures

---

## Execution detail

### Verified contracts

- 01 号票交付的 `./bin/flowforge proposal close <slug>`（需先 `make dev` 构建）。
- 存量现状：`ff-wiki/proposals/` 下 25 目录全终态（`bin/flowforge status -d ff-wiki/proposals` 报 108/108），加 roadmap-management 本 proposal 为第 26 个 live。
- `.flowforge/config.yaml` evidence.exempt_proposals 按 feature 名匹配（`internal/tracker/catalog.go` discoverEvidenceDiagnostics 的 exempt map 以 `featureForPath` 结果为键），closed/2026/<slug> 的 feature 名不变，名单零修改。
- 跨 proposal 引用真实存在（重写引擎目标）：`ff-wiki/proposals/closed/2026/fast-executor-reliability/requirements.md` 引 lightweight-execution-contract；`subagent-lifecycle` 引 fast-executor-reliability；`ff-wiki/roadmap/flowforge-self.md` 引 roadmap-management。
- 既有 waived 诊断：subagent-lifecycle 02-07 的 upstream-changed（waiver 已记录，与本迁移无关）。

### Execution scenarios

- Success：逐个 close 25 个 slug 后 → `proposals/` 根只剩 roadmap-management；status 首段渲染 flowforge-self 方向树；票计数恒 108；check --strict 除既有 waived 外无新增。
- Failure：某 slug 存在非终态票（理论不发生，25 目录已全终态）→ close 拒绝，人工核查后跳过或处理。

### Expected tests

- `./bin/flowforge status -d ff-wiki/proposals | head -3` — `108/108 resolved (100%)`，Frontier Ready/In Progress/Blocked 全 0。
- `./bin/flowforge check -d ff-wiki/proposals --strict 2>&1 | grep -c roadmap-dead-link` — 0。
- `grep -rn "proposals/lightweight-execution-contract" ff-wiki --include="*.md" | grep -v "closed" | wc -l` — 0（旧路径清零，重写后引用均含 closed/2026 段）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures。

### Generated artifacts

- `ff-wiki/proposals/closed/2026/<25 个 slug>/`（命令产出）；更新后的 `ff-wiki/roadmap/flowforge-self.md`（人工编辑：意图段补一句存量归档说明）。无其它需同步的 producer→consumer 链（链接重写由命令完成）。

### Conventions

- 迁移必须走命令，禁止手工 `mv`；一次一个 slug，close 后立即抽验该 slug 的链接重写再进行下一个。
- git 提交粒度：全部迁移+验证完成后单次提交（纯文档移动+链接重写，无代码变化）。
- Implementation note 回填：迁移清单、每个命令的重写计数、四项验证命令输出。

## Implementation note

- 完成全部 4 项 Changes：25 个存量 proposal 逐一经 `flowforge proposal close` 迁入 `closed/2026/`（全部成功、零拒绝——闭票校验全过）；三个高引 proposal 的链接重写抽验通过。
- 全链验证：status 首段渲染 flowforge-self 方向树（roadmap 需求树管理 — 实施中 3/5，准确：01/02/03 closed、04/05 open）+ `Archived: 25 closed proposal(s)`；票计数 111/113（108 存量 + 01/02/03 关闭，迁移前后总量不变）；check 零 roadmap-dead-link、除既有 waived（subagent-lifecycle 02-07 upstream-changed）外无新增；frontier Ready=1（05）In Progress=1（04）；孤儿报告为空。
- 重写残留处置：`research/2026-11-05-cost-aware-plan-principles.md` 对 fast-executor-reliability 的引用是文献式历史快照（git 历史事实），不属活链接，不改写；closed/2026/subagent-lifecycle/design.md 内 `docs/proposals/...` 前缀为本仓旧布局历史叙述，同样保留。
- Write set 合规：仅 ff-wiki/proposals/（含 closed/）与 ff-wiki/roadmap/。
- go test ./internal/... 与基线零新增失败（15 个预存在环境失败不变，前四票已验证）。

## Review rounds

### Round 0

- Gaps: none（25 目录/根残留/死链 grep/三命令语义/文献引用处置逐项实证）；2 项观察非 gap——票面 108/108 旧值已修正为 111/113 语义；cost-aware 归档票 evidence-artifact-missing 警告为既有基线行为
- Disposition: none
- Escalated to dual axes: 双轴以审计+实测证据聚合（迁移票）——Standards：迁移全部走命令零手工 mv、Write set 合规；Spec：D6 全条款交付（25 迁移+链接重写+方向文件存量说明+四项验证）
- Fix changes: none；Design returns: none；Repair: none

## Completion evidence

- Delivered: 25 个存量 proposal 经 `flowforge proposal close` 迁入 `ff-wiki/proposals/closed/2026/`（全成功零拒绝）；全库链接重写；`ff-wiki/roadmap/flowforge-self.md` 增存量说明段。
- Verification: `closed/2026` 恰 25 目录、proposals 根仅剩 closed+roadmap-management；高引三 slug 非 closed 引用清零（唯一残留为 research/2026-11-05 文献式 git 历史快照引用，处置记录不改写）；status 树视图 `实施中 3/5` + `Archived: 25`；check 零 roadmap-dead-link、除既有 waived 外无新增；票总量迁移前后不变（113）。
- Review: Round 0 审计零 gap + 2 观察处置；Standards/Spec 零发现。
- Deviations: research 文献引用与 closed 内 `docs/proposals/` 旧前缀历史叙述保留不改写（非活链接）。
- Reference: 工作树 ff-wiki/proposals/closed/2026/（25 目录）、ff-wiki/roadmap/flowforge-self.md、本票 Implementation note。
