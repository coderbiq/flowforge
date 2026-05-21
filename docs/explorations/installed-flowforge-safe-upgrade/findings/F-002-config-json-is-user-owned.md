# F-002 `config.json` 已经形成项目级持久配置边界

- Status: validated
- Source: `scripts/install.sh`, `workflow/guides/configuration.md`, `workflow/guides/adapter-contract.md`

## Statement

`install.sh` 只有在 `.flowforge/config.json` 不存在时才创建默认配置；配置文档也把它定义为项目级配置入口。这说明 `config.json` 不是升级时应覆盖的模板文件，而是用户态数据。

## Why it matters

安全升级必须默认保留项目级配置，否则会破坏 workspace 定义、任务后端设置和记忆提供器配置。

## References

- [Configuration guide](../../../workflow/guides/configuration.md)
- [Adapter contract](../../../workflow/guides/adapter-contract.md)
- [Installation script](../../../scripts/install.sh)
