# 开发指南

> 版本：v2.0.0-alpha | 最后更新：2026-06-12

## 开发环境

### 前置条件

- Go 1.24+
- Git

### 快速开始

```bash
# 克隆仓库
git clone <repo-url>
cd flowforge

# 当前平台构建
make dev

# 运行
./bin/flowforge version

# 运行测试
make test

# 本地开发（链接到 PATH）
go install ./cmd/flowforge
flowforge --version
```

### 环境变量

```bash
# Go 代理（中国大陆）
export GOPROXY=https://goproxy.cn,direct
```

## 代码规范

### Go 编码规范

1. **错误处理** — 所有错误必须显式处理，禁止 `_ = someFunc()` 忽略错误
2. **错误包装** — 使用 `fmt.Errorf("context: %w", err)` 保留错误链
3. **包结构** — `cmd/` 只放 `package main`（< 50 行），逻辑全在 `internal/`
4. **依赖管理** — 优先标准库，新依赖需审查

### 文件组织

- 每个命令一个文件：`internal/command/init.go`
- 业务逻辑与 CLI 路由分离：`internal/core/`
- 卡片操作通过 `CardStore` 统一接口

### 提交规范

```
add: 新功能
fix: Bug 修复
refactor: 重构（不改变功能）
update: 更新依赖或文档
remove: 删除功能
```

## 测试

```bash
# 运行全部测试
make test
# 或
go test ./internal/...

# 运行单个测试文件
go test ./internal/command/ -run TestInit

# 运行并显示覆盖率
go test -cover ./internal/...
```

## 构建

### 当前平台（开发用）

```bash
make dev
# 或
go build -trimpath -o bin/flowforge ./cmd/flowforge
```

### 所有平台（发布用）

```bash
make build
# 或
./scripts/build.sh v0.1.0 all
```

输出在 `dist/v0.1.0/`：

```
dist/v0.1.0/
├── flowforge-x86_64-unknown-linux-gnu.tar.gz
├── flowforge-x86_64-unknown-linux-gnu.tar.gz.sha256
├── flowforge-aarch64-apple-darwin.tar.gz
├── flowforge-aarch64-apple-darwin.tar.gz.sha256
├── flowforge-x86_64-pc-windows-msvc.zip
├── flowforge-x86_64-pc-windows-msvc.zip.sha256
└── checksums.txt
```

## 发布流程

### 手动发布（当前方案）

```bash
# 1. 打 tag
git tag v0.1.0

# 2. 交叉编译所有平台
./scripts/build.sh v0.1.0 all

# 3. 打包 + 签名 + 上传
./scripts/release.sh v0.1.0

# 产物在 dist/v0.1.0/，手动上传到 CDN
# 然后更新 release-latest.txt
```

### 发布脚本说明

| 脚本 | 功能 |
|------|------|
| `scripts/build.sh` | 交叉编译、打包、生成 checksum |
| `scripts/release.sh` | 生成 manifest.json、Ed25519 签名、上传到 OSS |
| `scripts/install.sh` | 用户安装脚本（macOS/Linux） |
| `scripts/install.ps1` | 用户安装脚本（Windows） |

## 用户安装

### macOS / Linux

```bash
curl -fsSL https://get.flowforge.dev | sh
```

### Windows

```powershell
irm https://get.flowforge.dev/install.ps1 | iex
```

### 安装脚本环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `FLOWFORGE_INSTALL` | 安装目录 | `$HOME/.flowforge` |
| `FLOWFORGE_CDN` | CDN 地址 | `https://cdn.flowforge.dev` |

## 项目结构

```
flowforge/
├── cmd/flowforge/     # CLI 入口
├── internal/          # 私有业务逻辑
│   ├── command/       # Cobra 命令
│   ├── config/        # 配置加载
│   ├── core/          # 核心业务
│   ├── update/        # 自更新引擎
│   ├── daemon/        # 守护进程（未来）
│   └── version/       # 版本注入
├── assets/            # 部署制品
├── docs/              # 开发文档
├── scripts/           # 构建/分发脚本
├── go.mod
└── Makefile
```

## 调试技巧

### 本地测试自更新

```bash
# 启动本地 HTTP 服务器模拟 CDN
cd dist/v0.1.0
python3 -m http.server 8080

# 设置环境变量指向本地服务器
export FLOWFORGE_CDN=http://localhost:8080

# 运行 CLI
./bin/flowforge upgrade
```

### 查看版本注入

```bash
# 查看编译时注入的版本信息
go build -ldflags="-X flowforge/internal/version.injected=v1.2.3" \
    -trimpath -o bin/flowforge ./cmd/flowforge
./bin/flowforge version
# 输出: flowforge v1.2.3
```
