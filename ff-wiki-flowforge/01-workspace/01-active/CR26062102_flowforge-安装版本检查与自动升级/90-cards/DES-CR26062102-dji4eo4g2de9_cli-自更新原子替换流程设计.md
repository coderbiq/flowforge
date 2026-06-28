---
id: DES-CR26062102-dji4eo4g2de9
title: CLI 自更新原子替换流程设计
type: design
status: draft
importance: should
links:
    - target: DES-djdothhisojr
      relation: related
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wuool60
      relation: implements
created: 2026-06-25T12:14:33.580252136Z
updated: 2026-06-25T12:16:11.461390459Z
source: CR26062102
---

# CLI 自更新原子替换流程设计

## Goal

设计 `flowforge upgrade` 命令的自更新流程，实现从 CDN 下载、签名验证、原子替换到失败回滚的完整链路。

## Decision

使用 `internal/update/` 包封装自更新逻辑，集成 `github.com/minio/selfupdate` 库实现原子替换。流程顺序：版本检查 → 下载二进制 → Ed25519 签名验证 → SHA256 校验 → 原子替换（备份 .old → 替换 → 成功则删除 .old / 失败则恢复 .old）。

## Rationale

- minio/selfupdate 是成熟库，支持 Linux/macOS/Windows 原子替换
- Ed25519 + SHA256 双重校验保证二进制完整性和来源可信
- 备份-替换-恢复-回滚模式保证升级失败不破坏现有运行能力

## Constraints

- 必须使用 `CGO_ENABLED=0` 静态编译，确保 selfupdate 可替换自身
- Windows 平台：如 selfupdate 不支持文件替换锁，需实现 `movefileex` 延迟重命名
- 签名验证失败的二进制不写入磁盘
- 下载的二进制放在临时目录，验证通过后再执行替换

## Impact

- 新增 `internal/update/upgrade.go` 实现自更新逻辑
- 新增 `internal/update/signature.go` 实现 Ed25519 签名验证
- 新增 `internal/command/upgrade.go` 实现 CLI 命令
- 新增 `internal/core/manifest.go` 实现 manifest.json 解析
- 依赖新增：`github.com/minio/selfupdate`、`golang.org/x/crypto/ed25519`

## Verification

- 正常升级：下载最新版本 → 验证通过 → 替换成功 → 显示新版本号
- 签名失败：下载文件 → 签名不匹配 → 清理临时文件 → 输出错误
- SHA256 不匹配：下载文件 → hash 不一致 → 清理临时文件 → 输出错误
- 回滚：下载成功 → 替换失败 → 自动恢复 .old 文件 → 输出错误
- Windows：exe 文件被占用时能正确处理

## Follow-up Tasks

- 实现 upgrade.go 自更新流程
- 实现 signature.go 签名验证
- 实现 manifest.go 解析
- 实现 CLI upgrade 命令

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [REQ-CR26062102-djeu2wuool60](REQ-CR26062102-djeu2wuool60_cli-自更新原子替换与回滚.md) [requirement] - CLI 自更新（原子替换与回滚）
- [DES-djdothhisojr](../../../../02-library/30-designs/DES-djdothhisojr_flowforge-upgrade-自更新流程.md) [design] - flowforge upgrade 自更新流程

### Incoming

#### implements
- [TASK-CR26062102-i-dji5l5lrfdkq](TASK-CR26062102-i-dji5l5lrfdkq_实现-manifestjson-解析与-ed25519-签名验证.md) [task] - 实现 manifest.json 解析与 Ed25519 签名验证
- [TASK-CR26062102-i-dji5la2j9llm](TASK-CR26062102-i-dji5la2j9llm_实现-flowforge-upgrade-命令.md) [task] - 实现 flowforge upgrade 命令 — 自更新下载、验证、原子替换

