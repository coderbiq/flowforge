---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      deploy-artifact-localization-requirements: 1
    design:
      deploy-artifact-localization-design: 2
---

# 03: 产物本地化（gitignore 自动化 + 存量迁移指引）+ AGENTS 模板约定

**Blocked by:** None
**Status:** open
**Mode:** lightweight

## Delivery

init/upgrade 自动管理 `.gitignore`（受管部署产物路径 + `.flowforge/config.yaml` 单文件，幂等、同名视为已满足）；部署前 `git ls-files` 检测已跟踪受管产物并输出可复制的 `git rm --cached -r` 指引（绝不自动改索引）；assets/AGENTS.md 与自举根增补每机器文件约定句。

## Design context

无 --untrack flag（用户确认）：一次性迁移不值新旗标，自动改索引违反最小惊讶。init/upgrade 的 `syncProjectAssets` 同通道守恒。

See the design authority at [部署产物本地化与模型注入方案](../design.md#deploy-artifact-localization-design)（d-artifact-localization 节）. Requirement authority: [部署产物本地化与模型注入需求](../requirements.md#deploy-artifact-localization-requirements).

## Touch points

- `internal/command/init.go` — gitignore 管理（**现状：internal/command 无任何 gitignore 处理，全新行为**）+ ls-files 检测
- `internal/command/upgrade.go` — `syncProjectAssets`（L95 起）同通道 + 检测
- `internal/command/assets_deploy.go` / 部署入口 — 受管路径清单来源
- `assets/AGENTS.md` 与仓根 `AGENTS.md` — 约定句
- `internal/command/init_test.go` / `agents_test.go` 相关用例

## Changes

- [ ] 1. 新增受管路径 gitignore 管理（init 与 upgrade 的 `syncProjectAssets` 共用）：向项目根 `.gitignore` 追加受管部署产物目录条目 + `.flowforge/config.yaml` 单文件条目；幂等（已有同名条目视为已满足、不重复写、不改用户其余内容）；无 `.gitignore` 时创建。
- [ ] 2. 部署前存量检测：`git ls-files --error-unmatch <受管路径>` 逐条（仅本地查询，无网络；无 git 仓时静默跳过）；命中向 stderr 输出逐条可复制 `git rm --cached -r <path>` 指引 + 一句原因（"部署产物是每机器文件，不应随 git 旅行"），绝不自动改动 git 索引。
- [ ] 3. `assets/AGENTS.md` 与仓根 `AGENTS.md` 增补一句：部署产物与 `.flowforge/config.yaml` 是每机器文件（init 自动 gitignore）；用户自写且想入库的 agent 用 `git add -f` 例外。
- [ ] 4. TDD 测试：gitignore 幂等（二次运行零追加、已有条目识别）；config 单文件条目落点；ls-files 命中输出指引、未命中静默、无 git 环境跳过；AGENTS 模板锚点断言扩展。

## Constraints

- must 纯本地确定性：git 检测仅本地 `git ls-files` 查询，绝不写 git 索引（Standards clause 转录）。
- must `assets/` 只放部署内容（AGENTS 约定句属部署内容）。
- 不提供任何 untrack flag；指引文本可复制即用。
- preset 测试授权：init_test.go/agents_test.go 新增用例与锚点断言扩展已经用户规划评审授权（2026-09-25）；测试文件改动用 bash。
- Write set: `internal/command/init.go`、`internal/command/upgrade.go`、`internal/command/assets_deploy.go`、`internal/command/init_test.go`、`internal/command/agents_test.go`、`assets/AGENTS.md`、`AGENTS.md`、`docs/proposals/deploy-artifact-localization/`

## Done and verify

- 幂等: 临时 fixture 二次 init 后 `.gitignore` 行数不变 — 测试 ok。
- 单文件忽略: `grep -n '.flowforge/config.yaml' <fixture>/.gitignore` — 命中且 `.flowforge/` 整目录条目不存在（自定义源仍入库）。
- 指引: 模拟已跟踪受管产物的 fixture 输出含 `git rm --cached -r` 且索引未变 — 测试 ok。
- 模板: `grep -n 'per-machine\|每机器' assets/AGENTS.md AGENTS.md` — 命中。
- 全套: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

---

## Execution detail

### Verified contracts

- **gitignore 全新行为**：`grep -rn gitignore internal/command/*.go`（非测试）零命中——无既有管理逻辑，无冲突面。
- init 部署入口与 upgrade `syncProjectAssets`（upgrade.go:95 起）是两个调用面，共用新函数；受管路径清单从部署目标推导（宿主目录 + config 路径）。
- `git ls-files --error-unmatch` 退出码语义：命中 0/未命中非 0——检测用；仓库不存在时 git 报 fatal，需优雅跳过（非 git 环境合法场景）。
- `.gitignore` 单文件条目 `.flowforge/config.yaml` 与目录条目 `.flowforge/` 语义对立——设计钉死只写单文件（subagents/ 自定义源仍入库）。
- AGENTS 模板锚点测试 `TestAgentRulesDescribeGenericCapabilityDispatch` 为 Contains 断言，加句零破坏。

### Execution scenarios

- Success：新 init 项目 `.gitignore` 含受管路径 + config 单文件；已跟踪存量项目 deploy 前输出指引、索引不动。
- Failure：若 gitignore 追加非幂等，二次 init 测试暴露；若误写 `.flowforge/` 目录条目，单文件断言暴露。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestGitignore|TestInitDeploys|TestUpgradeSync'` — ok。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

### Generated artifacts

- producer init/upgrade → consumer 项目根 `.gitignore`（每机器文件清单）；stderr 指引文本。

### Conventions

- must 变更后运行 `go test ./internal/...`。
- must 纯本地确定性文件操作，无网络（转录自设计 Standards clauses）。
- 测试文件改动用 bash（宿主守卫）。

