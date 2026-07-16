---
id: FIND-安装版本检查与自动升级-djkdd6b80x0d
title: Windows 自更新文件替换锁机制
type: finding
status: draft
importance: should
links:
    - target: PROP-CR26062102_flowforge-安装版本检查与自动升级
      relation: belongs_to
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: LOG-CR26062102-dji5335tyyuf
      relation: references
    - target: STR-djkdd6aag0ep
      relation: indexes
created: 2026-06-28T03:41:06.34564472Z
updated: 2026-06-28T03:41:06.346271186Z
source: CR26062102_flowforge-安装版本检查与自动升级
---

## Finding

`github.com/minio/selfupdate` 库已原生支持 Windows 平台的文件替换锁处理。在 Windows 上，该库使用 `MoveFileEx` API 配合 `MOVEFILE_DELAY_UNTIL_REBOOT` 标志处理正在运行的可执行文件。若 exe 被占用，替换操作排入重启后执行。

## Impact

无需为 Windows 平台实现额外的文件锁定处理逻辑。直接使用 minio/selfupdate 即可覆盖 Linux、macOS、Windows 三平台。

## Links

### Outgoing

- `PROP-CR26062102_flowforge-安装版本检查与自动升级` [belongs_to]

