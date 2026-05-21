# Journal Entry

- Timestamp: 2026-05-21T12:22:04Z
- Actor: Jon.Bi

## What changed

完成了对已安装 FlowForge 升级边界的初步判断：核心安装产物是可重新生成的，`config.json` 是用户态配置，`state/` 是运行态数据。

## Evidence

- `scripts/install.sh`
- `workflow/guides/configuration.md`
- `workflow/guides/adapter-contract.md`
- `docs/GETTING-STARTED.md`

## New questions

- 是否需要单独的 `upgrade` 命令。
- 是否需要记录安装版本或 payload 版本。
