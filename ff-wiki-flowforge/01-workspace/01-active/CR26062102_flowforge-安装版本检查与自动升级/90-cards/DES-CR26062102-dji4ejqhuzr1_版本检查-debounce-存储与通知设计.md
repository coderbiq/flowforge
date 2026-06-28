---
id: DES-CR26062102-dji4ejqhuzr1
title: 版本检查 debounce 存储与通知设计
type: design
status: draft
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wtrqe8g
      relation: implements
created: 2026-06-25T12:14:24.029610782Z
updated: 2026-06-25T12:16:10.4953208Z
source: CR26062102
---

# 版本检查 debounce 存储与通知设计

## Goal

设计版本检查的 debounce 状态存储和用户通知机制。

## Decision

使用 sqlite runtime state store 存储 debounce 状态（最后检查时间和版本号），与 ConfigService 使用同一存储后端。通知策略为简单二元：有新版本则提示，无则静默。debounce 间隔设为 1 小时，在一个 CLI 使用 session 内最多检查一次，避免频繁请求 CDN。

## Rationale

- 复用 ConfigService 的 sqlite runtime state store，减少新增存储机制
- 二元通知策略简单可靠，无需区分 patch/minor/major（版本号比较已有 semver，用户可自行判断升级优先级）
- 1 小时间隔平衡了及时性和请求频率：用户在同一天内多次使用 CLI 时能较快获知新版本，同时避免每次命令都请求 CDN
- `--no-version-check` 和 `config set version_check false` 为用户提供细粒度控制

## Constraints

- 复用 `internal/config/runtime_state_store.go` 中的 sqlite 存储
- debounce 检查异步执行，不阻塞 CLI 命令
- 版本检查失败时静默忽略，不影响正常操作
- semver 版本比较使用 `github.com/Masterminds/semver/v3`

## Impact

- 新增 runtimeStateStore 中的 `version_check` 表（check_time, checked_version 字段）
- 新增 `internal/update/checker.go` 实现异步版本检查逻辑
- CLI 启动时注入 version checker，在每个命令执行前触发检查

## Verification

- 首次检查：无 debounce 记录，执行 HTTP 请求，存储结果
- 1 小时内再次检查：有 debounce 记录，跳过 HTTP 请求
- 1 小时后检查：debounce 过期，执行 HTTP 请求
- 检查失败：不更新 debounce 记录，不输出错误
- `--no-version-check` 和 `config set version_check false` 均跳过检查

## Follow-up Tasks

- 实现 version_check 表 schema
- 实现 checker.go 异步检查逻辑
- 实现 CLI 启动钩子注入

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [REQ-CR26062102-djeu2wtrqe8g](REQ-CR26062102-djeu2wtrqe8g_版本检查与更新通知.md) [requirement] - 版本检查与更新通知

### Incoming

- [TASK-CR26062102-i-dji5kxuvxgod](TASK-CR26062102-i-dji5kxuvxgod_实现版本检查-sqlite-schema-checkergo-cli.md) [task] - 实现版本检查 — sqlite schema + checker.go + CLI 钩子注入

