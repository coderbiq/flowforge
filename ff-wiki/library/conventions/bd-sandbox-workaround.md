---
doc_type: convention
title: bd 写操作超时时使用 --sandbox 绕过 auto-sync
status: superseded
convention_status: superseded
enforcement: should
domain:
  scope: system
  type: convention
created: 2026-06-06
updated: 2026-06-06
---

> ⚠️ **此惯例已被废弃 (CR26060602)**。`--sandbox` 应在 beads.js 后端自动处理，不应暴露给 Agent。
> 替代方案见 `library/architecture/sandbox-leak-analysis.md`。

# bd 写操作超时时使用 --sandbox 绕过 auto-sync

## 规则

1. **应该**在 `bd create/update/close` 写操作超时时，添加 `--sandbox` 标志禁用 auto-sync 后重试
2. **应该**在会话收尾时运行 `bd dolt push` 手动同步远端——仅在网络可达时执行
3. **不应该**在网络不可达时反复重试无 `--sandbox` 的 bd 写操作

## 适用场景

- `flowforge task` CLI 内部调 `bd` 操作超时（30-60s）
- `bd` 写操作报 `ETIMEDOUT` 错误
- 外网无法访问 GitLab SSH 端口 7001

## 原因

`bd` 每次写操作后触发 `dolt auto-push` 到远端 GitLab（`git+ssh://git@gitlab.bytesforce.com:7001`）。SSH 端口 7001 在外网不可达时，push 挂起直到超时。**本地数据库写入本身是成功的**，仅远端同步受阻。

## 反例

```bash
# ❌ 错误：在网络不可达时反复使用无 --sandbox 的命令
flowforge task claim --proposal CR-id flowforge-cbe.3.1  # 超时 60s，即使 claim 已生效
flowforge task done --proposal CR-id flowforge-cbe.3.1 --summary "..."  # 同样超时

# ✅ 正确：代理操作使用 --sandbox
bd --sandbox update flowforge-cbe.3.1 --claim
bd --sandbox close flowforge-cbe.3.1 --reason "完成"

# 会话收尾时手动同步（仅在网络可达时）
bd dolt push
```

**为什么不对**：在网络不可达时，每次 bd 写操作会浪费 30-60 秒等待 auto-push 超时，严重影响开发效率。使用 `--sandbox` 可以立即完成本地写入。
