---
id: DES-CR26062102-dji4escqb5oo
title: 安装脚本增强设计
type: design
status: draft
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wt6uknc
      relation: implements
created: 2026-06-25T12:14:42.788317088Z
updated: 2026-06-25T12:16:12.656410283Z
source: CR26062102
---

# 安装脚本增强设计

## Goal

在现有 `scripts/install.sh` 和 `scripts/install.ps1` 基础上增强，支持版本选择、SHA256 校验、自定义安装目录和 GitHub Releases 备用下载。

## Decision

保留现有安装脚本结构，增加以下功能：`--version <ver>` 参数指定安装版本，`--prefix <dir>` 自定义安装目录，下载后 SHA256 校验，CDN 失败时自动 fallback 到 GitHub Releases。AGENTS.md 采用区块包裹策略部署（见 DES-CR26062102-dji5*-agents-md），安装脚本本身不直接处理 AGENTS.md，由 `flowforge init` 处理。

## Rationale

- 现有 `install.sh` 已有平台检测、架构检测、CDN 下载、安装目录创建等基础功能
- SHA256 校验防止下载损坏或中间人篡改
- GitHub Releases 作为 fallback 提高可用性
- 安装脚本职责单一：下载并安装 CLI 二进制；init 负责项目制品部署

## Constraints

- 安装脚本必须单一文件，支持 `curl | sh` 管道执行
- 不依赖除 curl/mktemp/chmod/uname 之外的命令
- SHA256 checksum 从 manifest.json 获取
- Windows install.ps1 同样支持这些参数

## Impact

- 修改 `scripts/install.sh` 增加参数解析、SHA256 校验、fallback 逻辑
- 修改 `scripts/install.ps1` 同步增加功能
- manifest.json 需要包含各平台的 SHA256 值（已有）

## Verification

- 默认安装最新版本：`curl ... | sh` → 安装最新 → `flowforge --version` 验证
- 指定版本安装：`curl ... | sh -s -- --version v0.1.0` → 安装指定版本
- 自定义目录：`curl ... | sh -s -- --prefix ~/bin` → 安装到自定义路径
- SHA256 校验失败：输出错误并清理
- CDN 不可用：自动 fallback 到 GitHub Releases

## Follow-up Tasks

- 增强 install.sh 参数解析和校验
- 增强 install.ps1 参数解析和校验
- 实现 GitHub Releases fallback 逻辑

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [REQ-CR26062102-djeu2wt6uknc](REQ-CR26062102-djeu2wt6uknc_一键安装脚本支持多平台.md) [requirement] - 一键安装脚本支持多平台

### Incoming

#### implements
- [TASK-CR26062102-i-dji5kdaztlcj](TASK-CR26062102-i-dji5kdaztlcj_增强-installsh-参数解析sha256-校验git-hub.md) [task] - 增强 install.sh — 参数解析、SHA256 校验、GitHub Releases fallback
- [TASK-CR26062102-i-dji5kgdglii9](TASK-CR26062102-i-dji5kgdglii9_增强-installps1-参数解析sha256-校验git.md) [task] - 增强 install.ps1 — 参数解析、SHA256 校验、GitHub Releases fallback

