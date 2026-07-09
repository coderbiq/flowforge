# FlowForge Agent 配置

> 本文档约束 Agent 如何开发 FlowForge。设计细节参考 docs/ 下的专项文档。

## Commands

- Build (当前平台): `make dev` 或 `go build -trimpath -o bin/flowforge ./cmd/flowforge`
- Build (所有平台): `make build` 或 `./scripts/build.sh <version> all`
- Test: `go test ./internal/...`
- Lint: `golangci-lint run ./...`
- Release: `make release` 或 `./scripts/release.sh <version>`

## 发布流程

每次发布时按以下步骤操作，不要遗漏：

1. **确认版本号**：每次发布必须有新的版本号。Bug 修复递增 patch（`v3.0.1` → `v3.0.2`），功能变更递增 minor（`v3.0.0` → `v3.1.0`）。**禁止重用已发布的 tag**——客户端靠版本号差异检测更新，同版本号不会触发升级。
2. **提交所有变更**：`git add -A && git commit -m "<message>"`。commit message 必须描述本次发布的变更内容。
3. **打 tag**（格式 `v<major>.<minor>.<patch>[-alpha.N]`）：
   ```bash
   git tag -a v3.0.0-alpha.1 -m "v3.0.0-alpha.1: <release summary>"
   ```
4. **推送代码和 tag**：
   ```bash
   git push origin main --tags
   ```
5. **验证**：确认 tag 出现在 GitHub Releases 页面，对应的构建产物自动生成。确认 `flowforge upgrade --dry-run` 在已安装旧版本的客户端上能看到新版本。

tag 命名规范：
- 正式版：`v3.0.0`、`v3.1.0`
- Bug 修复递增 patch：`v3.0.0` → `v3.0.1` → `v3.0.2`
- 功能变更递增 minor：`v3.0.0` → `v3.1.0`

## Boundaries

- ✅ **Always**: 变更前先读现有代码模式；变更后运行 `go test ./internal/...`；保持 SKILL 文件 < 200 tokens；错误必须显式处理（不用 `_` 忽略）
- ⚠️ **Ask first**: 添加新依赖、修改卡片 schema、变更 CLI 命令签名、修改 `internal/version/` 的注入逻辑
- 🚫 **Never**: 在 `assets/` 中放不部署的内容；使用 `panic()` 处理可恢复错误；直接编辑目标项目文件（必须通过 CLI）；在 `internal/` 中写公开 API（用 `pkg/` 或独立模块）；在方案确认前直接实施代码变更（必须先讨论方案，得到明确同意后再动手）

## `assets/` 是部署边界

`assets/` 下的所有文件都会部署到目标项目。不部署的内容放在 `assets/` 之外。

| assets/ 目录 | 部署目标 |
|--------------|---------|
| `assets/skills/` | `.agents/skills/` |
| `assets/templates/` | `.flowforge/templates/` |
| `assets/wiki/` | 项目 wiki 根目录 |
| `assets/AGENTS.md` | 目标项目 `AGENTS.md` |

添加文件前先问：**"这个文件会部署到目标项目吗？"** 不会就不放 `assets/`。

## 项目结构

```
flowforge/
├── cmd/flowforge/     ← CLI 入口（package main）
├── internal/          ← 私有业务逻辑（Go 编译器保护）
│   ├── command/       ← Cobra 命令实现
│   ├── config/        ← 配置加载
│   ├── core/          ← 核心业务（卡片 CRUD、上下文聚合）
│   ├── update/        ← 自更新引擎
│   ├── daemon/        ← 守护进程管理
│   └── version/       ← 版本注入
├── assets/            ← 部署制品（复制到目标项目）
├── docs/              ← 开发文档（不部署）
├── scripts/           ← 构建、安装脚本（不部署）
├── tests/             ← 集成测试（不部署）
├── go.mod             ← Go module 定义
└── Makefile           ← 构建命令
```

## SKILL 编写原则

SKILL 是本项目的核心产出物。编写或审查 SKILL 时对照以下原则：

| 原则 | 说明 |
|------|------|
| 单一职责 | 每个 SKILL 只做一件事 |
| 薄适配器 | SKILL 委托给 CLI，不内联所有内容 |
| 自洽命中 | 靠 description 让模型准确识别激活时机 |

### Description 审查清单

新增或修改 SKILL 时必须通过：

1. 能否 3 秒内说出"用户说了什么话，这个 SKILL 就该激活"？
2. 与相邻 SKILL 的 description 是否互不冲突？
3. 反例（不该激活的场景）是否明确？
4. description 是为模型写的，还是为人写的？

**禁止**：描述实现细节、使用抽象术语、缺少反例、与其他 SKILL 重叠。

## Agent 工作流驱动设计

设计顺序：**SKILL → 工作流模拟 → 实现 → 文档**。

禁止：
- 没有 SKILL 入口设计就写规则文档
- 写"描述性"规则而非"可执行"规则
- 先定义 artifact 结构再思考 Agent 如何使用

## 测试要求

- 变更后必须执行 `go test ./internal/...`
- 新增 SKILL、修改 CLI 参数、改变 context 输出格式时，测试必须同步更新

## Go 编码规范

### 错误处理

- 所有错误必须显式处理，禁止 `_ = someFunc()` 忽略错误
- 使用 `fmt.Errorf("context: %w", err)` 包装错误，保留错误链
- 自定义错误类型放在 `internal/` 对应包中

### 包结构

- `cmd/` 只放 `package main`，保持薄层（< 50 行），只做依赖注入和启动
- `internal/` 是真正逻辑所在，Go 编译器阻止外部导入
- 每个命令一个文件：`internal/command/init.go`、`internal/command/upgrade.go`
- 跨包共享的类型放在 `internal/` 顶层或专门的子包

### 依赖管理

- 优先使用标准库（`os`、`io`、`net/http`、`encoding/json`）
- 添加新依赖前检查是否已有类似功能
- CLI 框架：`github.com/spf13/cobra` + `github.com/spf13/viper`
- 版本比较：`github.com/Masterminds/semver/v3`

### 版本注入

- 版本通过 `-ldflags` 注入到 `internal/version.injected`
- 不要修改 `internal/version/version.go` 的结构，除非理解 ldflags 机制
- 本地构建使用 `make dev`，自动注入 git 信息

### 跨平台编译

- 使用 `CGO_ENABLED=0` 静态编译
- Windows 二进制加 `.exe` 后缀
- 打包格式：Linux/macOS 用 `.tar.gz`，Windows 用 `.zip`

## 设计文档索引

| 文档 | 说明 |
|------|------|
| [v3 重构方案](docs/proposal-v3/) | 卡片模型 v3、CLI 规格、SKILL 方法论、实现计划 |
| [架构设计](docs/architecture.md) | 项目定位、核心设计决策 |
| [CLI 设计](docs/cli-design.md) | 命令体系、init/upgrade/uninstall |
| [知识卡片系统](docs/knowledge-system.md) | 卡片模型、ID 规范、目录结构、索引系统 |

## 语言偏好

使用中文进行对话和文档编写。

<!-- FLOWFORGE:START -->
## FlowForge

Use `card init --type feature` to create cards; then edit the `.md` file directly for body content.
Use CLI for structured operations: `card link`, `card evolve`, `card log`, `card steps`.

### CLI
- `card init --type feature --title "..." --proposal <id>` to create a FEATURE card skeleton
- `card evolve <id> --stage designed|planned|done` for stage transitions (CLI enforces gates)
- `card log <id> --event "..." [--kind progress|bug|blocked]` to append to History
- `card steps <id> --status done|in_progress|blocked <n>` to update step status
- `context feature --feature <id> --step <n>` for minimal execution context
- `proposal inspect <id>` for auto-generated Feature Map and health checks
- `--body 'content\nwith\nnewlines'` for inline multi-line content
- Use single quotes for --body and --manifest to protect backticks, $, ! from shell expansion
- Never use shell redirects (`2>&1`, `<<`, `|`, `>`) with flowforge CLI — they trigger agent permission prompts
- `-o json` for machine-readable output
- `task`, `structure`, `log create` are DEPRECATED — use FEATURE-based commands instead

### Skills
| When | Skill |
|------|-------|
| Design / decompose proposal | `flowforge-design` |
| Execute implementation task | `flowforge-implement` |
| Report bug / finding / gap | `flowforge-feedback` |
| Import docs / archive proposal | `flowforge-curate` |
<!-- FLOWFORGE:END -->