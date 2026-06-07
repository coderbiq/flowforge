# Library Index

> 摘要看板 — 35 篇文档 | 更新于 2026-06-07 10:22

## 📊 概况

| 目录 | 文档数 |
|------|--------|
| architecture/ | 30 |
| conventions/ | 4 |
| decisions/ | 1 |

| 类型 | 文档数 |
|------|--------|
| architecture | 24 |
| finding | 6 |
| convention | 4 |
| decision | 1 |

| 成熟度 | 文档数 |
|--------|--------|
| growing | 27 |
| seed | 7 |
| deprecated | 1 |

## ⚠️ 铁律

- [Context 脚本应扫描目录而非检查 meta.status 来发现 proposal](conventions/proposal-discovery-by-directory.md)

## 🗑️ 待清理（已废弃）

- [bd 写操作超时时使用 --sandbox 绕过 auto-sync](conventions/bd-sandbox-workaround.md)

## 🔍 查找

```bash
flowforge library search "keyword"    # 全文搜索
flowforge library list --type design  # 按类型过滤
flowforge library list --module X     # 按模块过滤
flowforge library check --staleness   # 过期检测
```
