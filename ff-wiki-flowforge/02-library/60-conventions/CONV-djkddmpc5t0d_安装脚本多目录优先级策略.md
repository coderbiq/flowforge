---
id: CONV-djkddmpc5t0d
title: 安装脚本多目录优先级策略
type: convention
status: draft
importance: should
links:
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: DES-CR26062102-dji4escqb5oo
      relation: references
    - target: STR-djkddmozn3z9
      relation: indexes
created: 2026-06-28T03:41:42.028258708Z
updated: 2026-06-28T03:41:42.028258708Z
---

## Rule

安装脚本 `install.sh` 按优先级尝试将 CLI 二进制安装到已在 PATH 中的可写目录：`/usr/local/bin` → `/opt/homebrew/bin` → `$HOME/.local/bin` → `$HOME/.flowforge/bin`（自动配置 shell profile）。

## Rationale

优先使用系统 PATH 目录让用户安装后立即可用。Apple Silicon Mac 上 Homebrew 在 `/opt/homebrew/bin`，Intel Mac/Linux 在 `/usr/local/bin`。`$HOME/.local/bin` 是 freedesktop.org 用户 bin 标准。

## Applies When

编写或修改安装脚本时，必须按上述优先级选择安装目录。安装脚本支持 `--version` 指定版本、`--prefix` 自定义目录、SHA256 校验、GitHub Releases 下载。

## Links

- None

