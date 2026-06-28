# FlowForge Agent 配置

> 本文档约束 Agent 如何开发 FlowForge。设计细节参考 docs/ 下的专项文档。

## Commands

- Build (当前平台): `make dev` 或 `go build -trimpath -o bin/flowforge ./cmd/flowforge`
- Build (所有平台): `make build` 或 `./scripts/build.sh <version> all`
- Test: `go test ./internal/...`
- Lint: `golangci-lint run ./...`
- Release: `make release` 或 `./scripts/release.sh <version>`

## Boundaries

- ✅ **Always**: 变更前先读现有代码模式；变更后运行 `go test ./internal/...`；保持 SKILL 文件 < 200 tokens；错误必须显式处理（不用 `_` 忽略）
- ⚠️ **Ask first**: 添加新依赖、修改卡片 schema、变更 CLI 命令签名、修改 `internal/version/` 的注入逻辑
- 🚫 **Never**: 在 `assets/` 中放不部署的内容；使用 `panic()` 处理可恢复错误；直接编辑目标项目文件（必须通过 CLI）；在 `internal/` 中写公开 API（用 `pkg/` 或独立模块）

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
| [架构设计](docs/architecture.md) | 项目定位、核心设计决策 |
| [CLI 设计](docs/cli-design.md) | 命令体系、init/upgrade/uninstall |
| [知识卡片系统](docs/knowledge-system.md) | 卡片模型、ID 规范、目录结构、索引系统 |
| [v1 分析](docs/v1-analysis.md) | v1 问题诊断 |

## 语言偏好

使用中文进行对话和文档编写。

<!-- FLOWFORGE:START -->
## FlowForge

CLI is the only write path for cards. Never hand-write card files or frontmatter.

### CLI
- `--body -` heredoc `<<'EOF' ... EOF` for multi-line content
- `card batch - <<'EOF' ... EOF` for multi-card creation
- `-o json` for machine-readable output

### Skills
| When | Skill |
|------|-------|
| Design / decompose proposal | `flowforge-design` |
| Execute implementation task | `flowforge-implement` |
| Report bug / finding / gap | `flowforge-feedback` |
| Import docs / archive proposal | `flowforge-curate` |
<!-- FLOWFORGE:END -->