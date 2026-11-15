---
flowforge:
  schema: 1
  role: ticket
---
# 01: proposal close/reopen 生命周期命令与链接重写引擎

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`flowforge proposal close <slug>` 将已完结 proposal 迁入 `proposals/closed/<当年年份>/` 并重写全库指向它的相对链接；`flowforge proposal reopen <slug>` 迁回 proposals 根，live 侧同名冲突时拒绝。

## Design context

意图层与合同层分离下的 proposal 生命周期收口：归档不破坏依赖图（`featureForPath` 取叶子目录名，feature 不随迁移变化），但相对链接会断，所以 mv 与全库链接重写是同一个命令的原子两步。

See the design authority at `../design.md#commands`.

## Touch points

- `internal/command/root.go` — `cmd.AddCommand(...)` 注册新命令组
- `internal/command/proposal.go` — 新文件：`newProposalCmd()`、close/reopen 子命令
- `internal/command/link_rewrite.go` — 新文件：链接重写引擎（扫描、匹配、改写）
- `internal/command/proposal_test.go` — 新文件：命令测试
- `internal/command/link_rewrite_test.go` — 新文件：重写引擎测试

## Changes

- [x] 1. 在 `internal/command/link_rewrite.go` 实现 `RewriteReferences(docsDir, oldRelPath, newRelPath string) (rewrites int, err error)`：扫描 docsDir 下全部 `.md`，将相对链接中指向 oldRelPath 的路径段改写为按各文件位置重算的 newRelPath（`../` 前缀按文件所在深度重算），返回重写文件数
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run "TestRewriteReferences|TestMatchProposalPath" -count=1`
  - exit: 0
  - output: ok flowforge/internal/command — 7/7 pass
  - artifact: internal/command/link_rewrite.go
- [x] 2. 在 `internal/command/link_rewrite.go` 实现路径匹配辅助 `matchProposalPath(linkTarget, fromFile, proposalsRoot string) (slugPath string, ok bool)`：判定一个 markdown 链接目标是否解析到某个 proposal 目录（解析 `../` 与 `./` 段，容忍 `README.md`/`requirements.md` 等文件后缀，取到目录粒度）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestMatchProposalPath -count=1`
  - exit: 0
  - output: ok — 后缀容忍/目录粒度断言通过
  - artifact: internal/command/link_rewrite.go
- [x] 3. 在 `internal/command/proposal.go` 实现 `newProposalCmd()` 命令组
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestProposalCmdRegistered -count=1`
  - exit: 0
  - output: ok — 命令组注册断言通过
  - artifact: internal/command/proposal.go
- [x] 4. 实现 close 子命令定位 `proposals/<slug>/`（不存在报错）→ 用 `tracker.DiscoverIssues` 校验全部票终态（非终态列出并拒绝）→ `mv` 到 `proposals/closed/<当年>/`（按需建目录）→ 调用 `RewriteReferences` 重写全库 → 输出迁移与重写计数
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run "TestProposalCloseAndRewrite|TestProposalCloseRejectsNonTerminal|TestProposalCloseRejectsMissingProposal" -count=1`
  - exit: 0
  - output: ok — close 迁移/重写/拒绝路径 3/3
  - artifact: internal/command/proposal.go
- [x] 5. 实现 reopen 子命令在 `proposals/closed/*/<slug>/` 查找唯一匹配（零或多处匹配均报错列出）→ live 侧同名存在则拒绝 → `mv` 回 `proposals/<slug>/` → `RewriteReferences` → 输出计数
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run "TestProposalReopenAndRewrite|TestProposalReopenRejectsLiveCollision|TestProposalReopenRejectsMultipleMatches|TestProposalReopenMissing" -count=1`
  - exit: 0
  - output: ok — reopen 4/4
  - artifact: internal/command/proposal.go
- [x] 6. 在 `internal/command/proposal_test.go` 写测试临时目录构造含 roadmap 方向文件与相邻 proposal 正文链接的 fixture，覆盖 close 迁移+重写、close 拒绝非终态、reopen 迁回+重写、reopen 同名拒绝、reopen 多处匹配报错五条路径
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run "TestProposalClose|TestProposalReopen" -count=1`
  - exit: 0
  - output: ok — 7 用例
  - artifact: internal/command/proposal_test.go
- [x] 7. 在 `internal/command/link_rewrite_test.go` 写测试同深浅不同深度文件的 `../` 前缀重算、文件后缀容忍、非 proposal 链接不动
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run "TestRewriteReferences" -count=1`
  - exit: 0
  - output: ok — 深度重算/后缀容忍/他链不动/文字保留 5 用例
  - artifact: internal/command/link_rewrite_test.go

## Constraints

- close 只迁移 live 目录（`proposals/<slug>/`），不重复归档已 closed 的
- 重写只动链接路径，不改链接文字与正文其它内容
- 不引入 CLI 传长文本接口；slug 是唯一位置参数
- must 构建走 `make dev` 或 `go build -trimpath`，变更后 `go test ./internal/...` 全绿 — AGENTS.md Commands
- must not 引入通过 CLI 传长文本的接口；文件负责内容，CLI 负责计算 — AGENTS.md 核心设计原则
- Write set: `internal/command/`（新增 proposal.go、proposal_test.go、link_rewrite.go、link_rewrite_test.go；修改 root.go）

## Done and verify

- close 后目标目录在 `closed/<当年>/` 且全库链接指向新位置：临时 fixture 上 `go test ./internal/command/ -run TestProposalClose` — all pass
- 非终态票拒绝 close：`go test ./internal/command/ -run TestProposalCloseReject` — all pass, 错误信息列出未终态票
- reopen 全路径与冲突拒绝：`go test ./internal/command/ -run TestProposalReopen` — all pass
- 全量回归：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures

## Implementation note

- 完成全部 7 项 Changes：`internal/command/proposal.go`（close/reopen 子命令，`-d` flag 与 status 一致）、`internal/command/link_rewrite.go`（`RewriteReferences`+`matchProposalPath`）、两份测试共 14 用例、`internal/command/root.go` 注册命令组。
- 验证：`go test ./internal/command/ -run 'TestProposalClose|TestProposalReopen|TestRewriteReferences|TestMatchProposalPath' -count=1` — all pass（14/14）；`go vet` 干净；包级 build 绿。
- 端到端 smoke（真实二进制 + 临时 fixture）：close 迁入 `closed/2026/`，roadmap 方向文件与邻居 proposal 正文链接均按各自深度重算改写；reopen 迁回并反向重写，两向均报 `links rewritten in 2 file(s)`。
- Write set 合规：仅动 `internal/command/` 下票内四个新增文件与 `root.go`。
- 环境注记：本机 go 最终链接到 `/tmp` 的 exe 临时文件报 permission denied（trimacl 环境怪癖），`GOTMPDIR` 指到项目内即正常；与代码无关。
- 未尽事项：全量 `go test ./internal/...` 因子代理超时未跑完最终统计，留给 review 收口时补跑。

---

## Execution detail

### Verified contracts

- `internal/command/root.go:18-23` — `cmd.AddCommand(newCheckCmd(), newStatusCmd(), ...)` 是新命令组的注册点。
- `internal/tracker/parser.go:177` — `DiscoverIssues(root string) ([]*Issue, error)`；`Issue.Status` 终态判定在 `internal/tracker/model.go:22-24`（`IsTerminal()` = resolved/closed/wontfix）。
- `internal/config/config.go:238` — `ResolveProposalsDir(startDir)` 返回 proposals 根；docsDir = `filepath.Dir(proposalsRoot)`。
- `internal/tracker/catalog.go:422` — `featureForPath` 取叶子目录名（剥离 issues/evidence/design 父层），迁移后 feature 不变是本票不破坏依赖图的前提。
- 命令 flag 风格：`internal/command/status.go:105` 的 `cmd.Flags().StringVarP(&statusDir, "dir", "d", ...)`。

### Execution scenarios

- Success：fixture 含 `roadmap/d.md`（节点链接 `→ [p1](../proposals/p1/README.md)`）与相邻 proposal 正文引用；对全终态 p1 执行 close → p1 位于 `closed/2026/p1/`，两处链接改写指向新路径，输出迁移+重写计数。
- Failure：p1 存在 `Status: open` 票时 close → 非零退出，错误列出未终态票，目录未移动。
- Failure：reopen 时 live 侧已有同名 slug / `closed/*/` 下多处同名 → 非零退出并列出冲突路径。

### Expected tests

- `go test ./internal/command/ -run 'TestProposalClose|TestProposalReopen|TestRewriteReferences' -v` — all pass（close 迁移+重写、close 拒绝非终态、reopen 迁回+重写、reopen 同名拒绝、reopen 多处匹配、`../` 深度重算、后缀容忍、非 proposal 链接不动）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures。

### Generated artifacts

- `proposals/closed/<YYYY>/` 目录与重写后的 markdown 文件由命令在用户仓库产生；本票不产生需 producer→consumer 同步的仓库内新 artifact（消费方：05 号票在真实仓库执行这些命令）。

### Conventions

- 测试统一用 `t.TempDir()` 构造 docs 根，不依赖本仓真实 ff-wiki。
- 错误信息走 `fmt.Errorf` 返回给 cobra（非 os.Exit），保持命令层错误风格一致。
- 年份取 `time.Now().UTC().Year()`。

## Review rounds

### Round 0

- Gaps: 1 — contradicts 1（Done and verify「全量回归 all pass」与预存在环境失败矛盾；git stash 基线对照确认 FAIL 集合与改动零 diff）
- Disposition: 处置记录而非 Fix Change——15 个预存在失败为环境性（HEAD 同样失败），已在 Implementation note 环境注记中记录；后续票面的全量回归预期统一按「scoped 全绿 + 全量与基线零新增失败」解释
- Escalated to dual axes: no — gap 为证据措辞问题，非交付缺口
### Round 1

- Fixed point: working tree @ 06f6cf4（Round 0 后增量修复，未提交）
- Standards: P2×4 — 死代码三处（rel/恒nil error/悬空注释）、裸目录死链判定分歧（tracker vs command）、"proposals" 魔法串散布（记录不修：与现有代码库风格一致，等真实多根需求出现再收敛）、双份正则（票面 Constraints 已授权，驳回）
- Spec: P2×3 — 零票时树视图不渲染（Delivery 直接条款，已修）、空 Completion evidence 段误判已交付（已修）、CRLF frontmatter 失效（已修）
- Fix changes: 见上
- Design returns: none
- Repair: none
- 验证：scoped 测试全绿（command 11 用例 + tracker 4 用例）；go vet ./internal/... 干净；gofmt 干净；全量 go test ./internal/... 与基线（git stash 对照）零新增失败（15 个预存在环境失败不变）
## Completion evidence

- Delivered: `flowforge proposal close <slug>` / `reopen <slug>` 生命周期命令 + 全库相对链接重写引擎（`internal/command/proposal.go`、`link_rewrite.go`、root.go 注册）。
- Verification: `go test ./internal/command/ -run "Proposal|Rewrite|MatchProposal" -count=1` — ok（14 用例）；`go vet` 干净；真实二进制端到端 smoke：close 迁入 `closed/2026/` + 双向链接重写（roadmap 方向文件与邻居 proposal 正文各按深度重算），reopen 迁回 + 反向重写；非终态拒绝、同名冲突拒绝、多处匹配拒绝均有测试。
- Review: Round 0 — 1 contradicts（全量回归措辞 vs 预存在环境失败；git stash 基线对照 FAIL 集合零 diff，处置记录）；Round 1 — Standards P2×4 / Spec P2×3 中涉本票的均已处置（"proposals" 魔法串记录不修，与库内风格一致）。
- Deviations: 本机 /tmp go exe 链接 ACL 怪癖以 GOTMPDIR 绕开（Implementation note 环境注记）；委派子代理超时后由主会话按票面合同收口验证。
- Reference: 工作树 `internal/command/proposal.go` + `link_rewrite.go` + `proposal_test.go` + `link_rewrite_test.go` + `root.go`（提交随 proposal 收口）。
