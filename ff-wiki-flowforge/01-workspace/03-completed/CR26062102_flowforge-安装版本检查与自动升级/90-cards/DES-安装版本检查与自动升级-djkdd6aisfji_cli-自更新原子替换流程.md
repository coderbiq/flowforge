---
id: DES-安装版本检查与自动升级-djkdd6aisfji
title: CLI 自更新原子替换流程
type: design
status: draft
importance: should
links:
    - target: PROP-CR26062102_flowforge-安装版本检查与自动升级
      relation: belongs_to
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: DES-CR26062102-dji4eo4g2de9
      relation: references
    - target: STR-djkdd6aag0ep
      relation: indexes
created: 2026-06-28T03:41:06.303080942Z
updated: 2026-06-28T03:41:06.303884944Z
source: CR26062102_flowforge-安装版本检查与自动升级
---

## Goal

实现 `flowforge upgrade` 命令，从 GitHub Releases 下载新版本二进制，经 Ed25519 签名和 SHA256 双重校验后原子替换，失败自动回滚。

## Decision

使用 `github.com/minio/selfupdate` 库实现跨平台原子替换。流程顺序：版本检查 → 下载二进制 → Ed25519 签名验证 → SHA256 校验 → selfupdate.Apply 原子替换。替换前备份当前二进制为 `<binary>.old`，成功删除备份，失败自动恢复。

## Constraints

- CGO_ENABLED=0 静态编译确保 selfupdate 可替换自身
- minio/selfupdate 原生支持 Windows MoveFileEx 处理文件锁定
- 签名或 SHA256 验证失败不写入磁盘
- 支持 `--version` 指定版本（允许降级）和 `--dry-run` 预览

## Links

### Outgoing

- `PROP-CR26062102_flowforge-安装版本检查与自动升级` [belongs_to]

