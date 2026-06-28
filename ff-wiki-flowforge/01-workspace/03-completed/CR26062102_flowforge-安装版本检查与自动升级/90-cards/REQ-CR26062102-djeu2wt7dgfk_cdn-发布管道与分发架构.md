---
id: REQ-CR26062102-djeu2wt7dgfk
title: CDN 发布管道与分发架构
type: requirement
status: draft
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
created: 2026-06-21T15:31:01.152517Z
updated: 2026-06-21T15:31:01.152519Z
source: CR26062102
---

# CDN 发布管道与分发架构

## Summary

通过 GoReleaser 编译多平台二进制，上传到 GitHub Releases 作为分发源，提供 manifest.json 作为版本元数据入口。无需自建 CDN 或域名备案。

## Source

设计卡 DES-CR26062102-dji4f30jvnes 已定义发布流程。

## Acceptance

- GoReleaser 编译 6 平台二进制（linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64）
- 发布产物自动上传到 GitHub Releases
- 每个 release 包含 `manifest.json`（各平台 URL、SHA256、Ed25519 签名引用）
- 发布流程：`make release` 或 git tag 触发 CI 自动发布
- GitHub Releases 作为唯一分发源，全球可达，无需 CDN 配置

## Scope

发布管道包括：版本号注入、GoReleaser 配置、manifest 生成和签名、GitHub Releases 上传

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

#### analyzes
- [TASK-CR26062102-a-dji4ebxyuc7p](TASK-CR26062102-a-dji4ebxyuc7p_分析-cdn-签名密钥管理方案.md) [task] - 分析 CDN 签名密钥管理方案
- [TASK-CR26062102-a-dji4fg0ekv55](TASK-CR26062102-a-dji4fg0ekv55_分析-cdn-缓存策略与签名文件分发.md) [task] - 分析 CDN 缓存策略与签名文件分发
- [DES-CR26062102-dji4f30jvnes](DES-CR26062102-dji4f30jvnes_cdn-发布管道设计.md) [design] - CDN 发布管道设计
#### records
- [LOG-CR26062102-dji534c6oz39](LOG-CR26062102-dji534c6oz39_分析结论-cdn-签名密钥管理方案.md) [log] - 分析结论: CDN 签名密钥管理方案
- [LOG-CR26062102-dji5378kq0uy](LOG-CR26062102-dji5378kq0uy_分析结论-cdn-缓存策略与签名文件分发.md) [log] - 分析结论: CDN 缓存策略与签名文件分发
- [LOG-CR26062102-dji5m5xwuu8m](LOG-CR26062102-dji5m5xwuu8m_split-turn-拆解-11-个实现任务.md) [log] - split turn: 拆解 11 个实现任务
- [LOG-安装版本检查与自动升级-dji4dluy95gj](LOG-安装版本检查与自动升级-dji4dluy95gj_design-turn-需求卡内容填充与重复卡片清理.md) [log] - design turn: 需求卡内容填充与重复卡片清理
#### satisfies
- [TASK-CR26062102-i-dji5kn5ajp5i](TASK-CR26062102-i-dji5kn5ajp5i_增强-releasesh-ed25519-签名生成与-cdn.md) [task] - 增强 release.sh — Ed25519 签名生成与 CDN 上传
- [TASK-CR26062102-i-dji5kqtak6ew](TASK-CR26062102-i-dji5kqtak6ew_创建-goreleaseryml-与-git-hub-actions-release.md) [task] - 创建 .goreleaser.yml 与 GitHub Actions release workflow

## Open Questions

- None（签名密钥方案见 LOG-CR26062102-dji534c6oz39，分发见 DES-CR26062102-dji4f30jvnes）

