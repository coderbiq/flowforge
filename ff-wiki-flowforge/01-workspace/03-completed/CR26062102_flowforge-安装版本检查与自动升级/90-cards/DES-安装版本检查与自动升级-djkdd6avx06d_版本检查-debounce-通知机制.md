---
id: DES-安装版本检查与自动升级-djkdd6avx06d
title: 版本检查 debounce 通知机制
type: design
status: draft
importance: should
links:
    - target: PROP-CR26062102_flowforge-安装版本检查与自动升级
      relation: belongs_to
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: DES-CR26062102-dji4ejqhuzr1
      relation: references
    - target: STR-djkdd6aag0ep
      relation: indexes
created: 2026-06-28T03:41:06.325215279Z
updated: 2026-06-28T03:41:06.325933391Z
source: CR26062102_flowforge-安装版本检查与自动升级
---

## Goal

每次 CLI 执行时后台异步检查新版本，1 小时间隔 debounce，有新版本时提示用户升级。

## Decision

使用 sqlite runtime state store 存储 debounce 状态（`version_check` 表记录 check_time 和 checked_version），与 ConfigService 共用存储后端。通知策略为二元：新版本可用则 stderr 提示，否则静默。

## Constraints

- 1 小时间隔 debounce，同一版本号不重复检查
- 异步 goroutine 执行，不阻塞当前命令
- HTTP 请求失败静默忽略
- `--no-version-check` flag 跳过单次检查
- `config set version_check false` 全局禁用
- semver 比较使用手动实现（避免外部依赖）

## Links

### Outgoing

- `PROP-CR26062102_flowforge-安装版本检查与自动升级` [belongs_to]

