---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      wiki-config-single-track-requirements: 3
    design:
      wiki-config-single-track-design: 1
---

## Completion evidence

闭环证据见各 Change 的 cmd/exit/output/artifact 四元组与 Implementation note（执行记录、验证命令与观测结果）。

# 01: wiki 轨删除面 + load 告警 + 字面量单点

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

wiki root 轨全链死代码一次删除（类型/函数群/service 分支/viper default），`Load` 增 raw 键 stderr 告警（`warnOut` seam），`"ff-wiki"` 字面量收敛到 `DefaultDocsDir` 单点、不可达死默认清理；`go test ./internal/...` 全绿且 wiki 键存量配置 load 不炸。

## Design context

删除断面经调研证实零生产消费；忽略语义零代码（删字段后 mapstructure 忽略未知键）；告警经包级 writer seam 输出可测试捕获；`config set project.x.wikiRoot` 转为报错属裁决内破坏面。

See the design authority at [双轨 wiki 配置统一：方案](../design.md#wiki-config-single-track-design)（d-deprecation 删除清单与告警形态、d-unification 字面量收敛）. Requirement authority: [双轨 wiki 配置统一：需求](../requirements.md#wiki-config-single-track-requirements)（rev 3，Q1-Q3 已裁决）.

## Touch points

- `internal/config/config.go` — `WikiConfig` 类型、`Config.Wiki`、`defaultConfig`、`Load`（viper default + 告警插入点）、`WikiRoot`/`WikiRootForProject`/`projectWikiRoot`/`primaryProject`、`ProjectConfig.WikiRoot`
- `internal/config/service.go` — `List()` wikiRoot 行、`getProjectConfig`/`setProjectConfig` 的 `"wikiRoot"` case
- `internal/config/config_test.go` — wiki 轨断言（L20/38/61/126-141/158 等）
- `internal/command/init_docs_test.go` — `TestConfigListIsStableAndIncludesDocsDir`（L100,107）
- `internal/command/assets_deploy.go`（L21-22）、`internal/command/assets_compare.go`（L61-62）— 不可达 `"docs"` 死默认
- `internal/command/init.go`（L45）— 默认 YAML 字面量

## Changes

- [x] 1. `internal/config/config.go`：删除 `WikiConfig` 类型、`Config.Wiki` 字段、`defaultConfig` 中 Wiki 初始化、`Load` 的 `v.SetDefault("wiki.root", …)`、`WikiRoot()`、`WikiRootForProject()`、`projectWikiRoot()`、`primaryProject()`、`ProjectConfig.WikiRoot` 字段。
  - cmd: `grep -rn 'WikiRoot' internal/ --include='*.go' | grep -v _test.go | wc -l`
  - exit: 0
  - output: `0`（生产面清零）
  - artifact: internal/config/config.go
- [x] 2. `internal/config/config.go`：新增包级 `var warnOut io.Writer = os.Stderr`；`Load` 在 `ReadInConfig` 成功后、`Unmarshal` 前检查 `v.IsSet("wiki")` 与 projects 中 `wikiRoot != ""` 的项，逐键向 `warnOut` 输出英文告警（含键名 + "use docs_dir" 指引；projects 项含 project id）。
  - cmd: `go test ./internal/config/ -run TestLoadWarnsOnDeprecatedWikiKeys`
  - exit: 0
  - output: `ok`（两条告警逐键捕获；viper 小写化细节见 Implementation note）
  - artifact: internal/config/config.go
- [x] 3. `internal/config/service.go`：删除 `List()` 的 `project.%s.wikiRoot` 行、`getProjectConfig` 与 `setProjectConfig` 的 `"wikiRoot"` case。
  - cmd: `go test ./internal/command/ -run TestConfigListIsStable`
  - exit: 0
  - output: `ok`（list 不含 wiki 键；set 报错路径由 TestConfigSetProjectWikiRootRejected 钉死）
  - artifact: internal/config/service.go
- [x] 4. 字面量收敛：`defaultConfig.DocsDir`（config.go:74）与 init 默认 YAML（init.go:45）改引用 `DefaultDocsDir` 常量；删除 `deployManagedAssets`/`compareManagedAssets` 的不可达 `"docs"` 兜底分支。
  - cmd: `grep -rn '"ff-wiki"' internal/ --include='*.go' | grep -v _test.go`
  - exit: 0
  - output: 仅 `internal/config/config.go:17: DefaultDocsDir = "ff-wiki"` 单点
  - artifact: internal/command/init.go
- [x] 5. 测试更新：删除/改造 wiki 轨断言；新增告警测试（替换 `warnOut` 捕获，顶层 `wiki.root` 与 `projects[].wikiRoot` 逐键断言）、忽略语义测试（yaml 含 `wiki:` 块 load 成功、`DocsRoot` 走 docs_dir）、`config set project.x.wikiRoot` 报错测试；`TestConfigListIsStableAndIncludesDocsDir` 改为断言 list 不含 wiki 键。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...`
  - exit: 0
  - output: 全部 ok（command/config/subagent/tracker/update，0 failures）
  - artifact: internal/config/config_test.go

## Constraints

- must 纯本地确定性文件操作，无网络、无 LLM 调用（转录自设计 Standards clauses）。
- must 不改 CLI 命令签名与 Issue Schema 头规范（告警不改 `Load` 函数签名，经包级 writer）。
- 存量兼容红线：yaml 含 `wiki:` 块或 `projects[].wikiRoot` 的配置 load 必须成功（告警 + 忽略），tangram-v2 形态行为零变化，以测试钉死。
- preset 测试授权：本票 Write set 内 `*_test.go` 的 wiki 轨断言删除/改造与新增用例已经用户在规划评审中显式授权（2026-09-25 会话），不得改动无关断言。
- Write set: `internal/config/`、`internal/command/assets_deploy.go`、`internal/command/assets_compare.go`、`internal/command/init.go`、`internal/command/init_docs_test.go`、`docs/proposals/wiki-config-single-track/`

## Done and verify

- 生产面清零: `grep -rn 'WikiRoot' internal/ --include='*.go' | grep -v _test.go` — 空输出。
- 字面量单点: `grep -rn '"ff-wiki"' internal/ --include='*.go' | grep -v _test.go` — 仅 `DefaultDocsDir` 常量定义一处。
- 全套回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok，0 failures。
- 告警可见: 含 `wiki.root` + `projects[].wikiRoot` 的 fixture load 测试通过且捕获到两条告警（新测试用例名执行 `-run` 过滤）。

## Implementation note

**TDD 顺序**

1. Red：先写三组新测试（`TestLoadWarnsOnDeprecatedWikiKeys` 告警逐键断言、`TestLoadIgnoresLegacyWikiBlock` 忽略语义、`TestConfigSetProjectWikiRootRejected` set 报错），并将 `TestConfigListIsStableAndIncludesDocsDir` 尾部改造为“set docs_dir 后 list 不含 wiki 键”；运行确认以预期原因失败——包编译失败 `undefined: warnOut`（缺 seam，非语法错误）。
2. Green：按 Changes 1-4 执行删除与实现；同步清理 wiki 轨旧断言（删 `TestWikiRoot`/`TestWikiRootAbsolute`/`TestWikiRootForProject`，`TestProjectByID` 改用 srcDirs 断言，`TestDefaultConfig`/`TestLoadConfig` 去 wiki 字段断言但保留含 `wikiRoot` 的 fixture 钉死忽略语义）。
3. 一次修复：Red→Green 首跑告警测试失败，诊断证实 viper 对 raw 键递归小写化（`wikiRoot` → `wikiroot`），`warnDeprecatedWikiKeys` 改读 `project["wikiroot"]` 后全绿。

**验证命令与结果**

- `GOPROXY=https://goproxy.cn,direct go test ./internal/config/` — ok（含新增三组用例）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestConfigListIsStable|TestInit|TestDeployManagedAssets|TestCompareManagedAssets' -v` — 6 个测试全 PASS。
- `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...` — 全部 ok，0 failures（command/config/subagent/tracker/update 五包）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/config/ -run 'TestLoadWarnsOnDeprecatedWikiKeys|TestLoadIgnoresLegacyWikiBlock|TestConfigSetProjectWikiRootRejected' -v` — 3 PASS，告警捕获两条。
- `grep -rn 'WikiRoot' internal/ --include='*.go' | grep -v _test.go` — 空输出（生产面清零）。
- `grep -rn '"ff-wiki"' internal/ --include='*.go' | grep -v _test.go` — 仅 `internal/config/config.go:17` 常量定义一处。
- `go vet ./internal/...` — 通过；`gofmt -l` 本票全部改动文件干净（`internal/update/manifest.go` 为存量未格式化，非本票引入）。
- 附注：`golangci-lint` 本机未安装，未运行。

**修改文件清单**

- `internal/config/config.go` — Changes 1/2/4（wiki 轨删除 + `warnOut` seam + 告警函数；`Save` 的 fileConfig Wiki 行随类型删除）
- `internal/config/service.go` — Change 3
- `internal/command/init.go` — Change 4（默认 YAML 引用 `config.DefaultDocsDir`）
- `internal/command/assets_deploy.go` / `internal/command/assets_compare.go` — Change 4（删不可达 `"docs"` 兜底）
- `internal/config/config_test.go` — Change 5（wiki 轨断言删除/改造 + 三组新用例；测试文件改动按票面 Conventions 经 bash 完成）
- `internal/command/init_docs_test.go` — Change 5（`TestConfigListIsStableAndIncludesDocsDir` 改造）

---

## Execution detail

### Verified contracts

- 删除面零生产消费：`WikiRoot()`/`WikiRootForProject()` 全仓生产调用者为零（调研 W1 + 旗舰复核 grep 证实）；`Config.Wiki` 唯一展示消费在 service.go List/getProjectConfig。
- 忽略语义零代码：viper+mapstructure 对未注册字段默认忽略；删除 `Config.Wiki`/`ProjectConfig.WikiRoot` 后 yaml 存量 `wiki:` 块自然失效，无需清洗逻辑。
- 告警 seam 无先例：internal/config 现无任何 stderr 输出（grep 证实），新增 `var warnOut io.Writer = os.Stderr` 为首例；测试以替换包级变量捕获。
- `Load` 结构（config.go:154-182）：viper SetDefault ×8 → ReadInConfig（ConfigFileNotFoundError 时返回 defaultConfig）→ Unmarshal；告警插在 ReadInConfig 成功后、Unmarshal 前，检查 `v.IsSet("wiki")`（删除 SetDefault 后仅用户显式写入为 true）与 `v.Get("projects")` 原始数组。
- `TestConfigListIsStableAndIncludesDocsDir` 位于 internal/command/init_docs_test.go:100,107（跨包引用 config List——改造时注意该测试属 internal/command）。
- 既有 9 个 docs-root 语义测试（`TestDeployManagedAssetsUsesAbsoluteDocsRoot` 等）不触碰 wiki 轨，零回归预期。

### Execution scenarios

- Success：tangram-v2 形态 fixture（显式 `docs_dir` + `wikiRoot`）load 成功、stderr 两条告警、`DocsRoot` 返回 docs_dir 值；全套测试绿。
- Failure：若删除 `ProjectConfig.WikiRoot` 后仍有生产引用，编译期暴露（build 失败即定位）；若告警写在 Unmarshal 后，`Config.Wiki` 已删无从判断——必须用 raw viper 键（设计钉死插入点）。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/config/` — ok（含新增告警/忽略/set 报错用例）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestConfigListIsStable|TestInit|TestDeployManagedAssets|TestCompareManagedAssets'` — ok。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

### Generated artifacts

- Not applicable：纯代码/测试删除与收敛，无生成物（部署产物路径零变化——编译目标文件不变）。

### Conventions

- must 变更后运行 `go test ./internal/...`（转录自设计 Standards clauses）。
- 告警消息英文（CLI 输出现状）；`warnOut` 命名与 io/os 导入遵循包内风格。
- 测试文件改动用 bash 完成（宿主 PreToolUse 守卫拦截 `*_test.go` 的 edit/write 工具）。
