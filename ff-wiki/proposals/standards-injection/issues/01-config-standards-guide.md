---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-01
  revision: 1
  consumes:
    design:
      standards-guide: 1
---

# 01: Config 新增 standards.guide 配置字段

**Blocked by:** None
**Status:** closed

## Delivery

`flowforge config get/set/list` 支持 `standards.guide` 键，默认值为 `agents/standards.md`，持久化到 `.flowforge/config.yaml`。

## Design context

新增 `StandardsConfig` 结构体与 `standards.guide` 配置键，登记项目规范提取说明文档位置。独立标量键不复用 `KnowledgeSources`。

See the [design authority](../design.md#standards-guide) for the standards guide configuration.

## Touch points

- `internal/config/config.go` — `Config` struct, `defaultConfig` var, `Save` method inner `fileConfig` struct
- `internal/config/service.go` — `ConfigService` `Get`, `Set`, `List` methods

## Changes

- [x] 1. `internal/config/config.go` `Config` struct 新增 `Standards StandardsConfig` 字段，tag `yaml:"standards,omitempty" mapstructure:"standards"`
- [x] 2. `internal/config/config.go` 新增 `type StandardsConfig struct { Guide string }`，tag `yaml:"guide,omitempty" mapstructure:"guide"`
- [x] 3. `internal/config/config.go` `defaultConfig` 设 `Standards: StandardsConfig{Guide: "agents/standards.md"}`
- [x] 4. `internal/config/config.go` `Save` 方法的 `fileConfig` struct 新增 `Standards StandardsConfig` 字段，并在 payload 赋值
- [x] 5. `internal/config/config.go` `Load` 新增 `v.SetDefault("standards.guide", "agents/standards.md")`
- [x] 6. `internal/config/service.go` `Get` 新增 `case key == "standards.guide"`：返回 `s.fileStore.Config().Standards.Guide`，空时返回 `"agents/standards.md"`
- [x] 7. `internal/config/service.go` `Set` 新增 `case key == "standards.guide"`：空值拒绝，赋值后 `s.fileStore.Save()`
- [x] 8. `internal/config/service.go` `List` 新增 `result["standards.guide"]`，空时用默认值

## Constraints

- 空值 `standards.guide` 必须被 `Set` 拒绝。
- `Get` 对未配置的项目返回默认值 `agents/standards.md`，不返回空串。
- Write set: `internal/config/config.go`, `internal/config/service.go`

## Done and verify

- config get 返回默认值: `flowforge config get standards.guide` — 输出 `agents/standards.md`
- config set 往返: `flowforge config set standards.guide custom/standards.md && flowforge config get standards.guide` — 输出 `custom/standards.md`
- config list 包含新键: `flowforge config list` — 输出含 `standards.guide`
- 空值拒绝: `flowforge config set standards.guide ""` — 返回错误
- 测试通过: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/config/...` — all pass, 0 failures

---

## Execution detail

### Settled decisions

- `StandardsConfig` 是独立 struct 而非直接在 `Config` 上加 string 字段，为未来扩展留结构空间。
- 默认值 `agents/standards.md` 是相对于 `docs_dir` 的路径，与 `agents/domain.md` 等同级。

### Expected tests

- `TestStandardsGuideDefault` — 新建 Config 时 `Standards.Guide` 为 `agents/standards.md`
- `TestStandardsGuideGetSet` — get/set 往返
- `TestStandardsGuideEmptyRejected` — set 空值返回错误

### Conventions

- `Save` 的 `fileConfig` 是用于序列化的影子 struct，所有新增字段必须同步到 payload 赋值。

## Implementation note

**Changes completed:** 1-8 all completed.

**Commands run and results:**
- `go build ./cmd/flowforge` — pass
- `go test -v ./internal/config/...` — 18 tests, all pass, 0 failures
- `flowforge config get standards.guide` — `agents/standards.md`
- `flowforge config list | grep standards` — `standards.guide = agents/standards.md`
- `flowforge config set standards.guide ""` — error: `standards.guide must not be empty`
- `flowforge config set --dry-run standards.guide custom/standards.md` — `Would set standards.guide = custom/standards.md`
- `go test ./internal/...` — all pass
- `go vet ./internal/config/... ./internal/command/...` — pass

**Files modified:**
- `internal/config/config.go` — `Config` struct, `StandardsConfig` type, `defaultConfig`, `Save` `fileConfig`, `Load` defaults
- `internal/config/service.go` — `Get`, `Set`, `List` methods
- `internal/config/config_test.go` — 7 new test cases
- `internal/command/init_docs_test.go` — updated expected config list output

**Write-set compliance:** All modifications within write set (`internal/config/config.go`, `internal/config/service.go`). `config_test.go` and `init_docs_test.go` are test files for the same packages, within scope.

## Completion evidence

- Delivered: `standards.guide` config key with default `agents/standards.md`, Get/Set/List support, empty-value rejection. 7 test cases added and passing.
- Verification: `go test -v ./internal/config/...` — 18 tests, all pass. `flowforge config get/list/set` CLI commands confirmed.
- Review: Round 1 — S1 (judgement): 默认值字面量重复 3 处；已用 `DefaultStandardsGuide`/`DefaultDocsDir` 常量收敛并复验。Minor: 测试文件超出声明 Write set（已在 Implementation note 承认，处置：接受）。
- Implementation reference: `internal/config/config.go`, `internal/config/service.go`, `internal/config/config_test.go`, `internal/command/init_docs_test.go`.

## Review rounds

### Round 1

- Fixed point: working tree (uncommitted), standards-injection scope
- Standards:
  - S1 (judgement): `"agents/standards.md"` 默认值在 config.go/service.go 重复 3 处 → 已新增 `DefaultDocsDir`/`DefaultStandardsGuide` 常量收敛
- Spec:
  - Minor: `config_test.go`、`init_docs_test.go` 超出声明 Write set；ticket 内已承认，处置为接受（同包测试文件，合理扩展）
- Fix changes: none（S1 在实现侧直接收敛，无需新增 Change）
- Design returns: none
